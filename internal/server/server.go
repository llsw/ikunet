package server

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	actor "github.com/asynkron/protoactor-go/actor"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	kerrors "github.com/cloudwego/kitex/pkg/kerrors"
	ksvc "github.com/cloudwego/kitex/server"
	transport "github.com/llsw/ikunet/internal/kitex_gen/transport"
	transportSvc "github.com/llsw/ikunet/internal/kitex_gen/transport/transportservice"
)

var (
	asys *actor.ActorSystem
	_    Server = &server{}
)

func init() {
	asys = actor.NewActorSystem()
}

func GetActorSystem() *actor.ActorSystem {
	return asys
}

// TransportServiceImpl implements the last service interface defined in the IDL.
type TransportServiceImpl struct {
	eps Endpoint
}

func NewTransportServiceImpl(eps Endpoint) *TransportServiceImpl {
	return &TransportServiceImpl{
		eps: eps,
	}
}

// Call implements the TransportServiceImpl interface.
func (s *TransportServiceImpl) Call(ctx context.Context, req *transport.Transport) (resp *transport.Transport, err error) {
	resp = &transport.Transport{}
	err = s.eps(ctx, req, resp)
	return
}

func richMWsWithBuilder(ctx context.Context, mwBs []MiddlewareBuilder, s *server) []Middleware {
	for i := range mwBs {
		s.mws = append(s.mws, mwBs[i](ctx))
	}
	return s.mws
}

// newErrorHandleMW provides a hook point for server error handling.
func newErrorHandleMW(errHandle func(context.Context, error) error) Middleware {
	return func(next Endpoint) Endpoint {
		return func(ctx context.Context, request, response *transport.Transport) error {
			err := next(ctx, request, response)
			if err == nil {
				return nil
			}
			if errHandle != nil {
				return errHandle(ctx, err)
			}
			return err
		}
	}
}

type Transport struct {
	Addr    string
	Session int64
	Meta    *transport.Meta
	Cmd     string
	Msg     any
}

type Message struct {
	ctx      context.Context
	request  any
	response any
}

func GetTraceId(ctx context.Context) string {
	traceId, ok := ctx.Value("traceId").(string)
	if !ok {
		traceId = ""
	}
	return traceId
}

func TracesToBytes(cluster, svc, cmd string) ([]byte, error) {
	return []byte{0, 0, 0, 0, 0, 0}, nil
}

func IdToTraces(cluster, svc, cmd int) []string {
	return []string{
		string(cluster),
		string(svc),
		string(cmd),
	}
}

func BytesToTraces(traces []byte) ([]string, error) {
	l := len(traces)
	if l < 6 || l%6 != 0 {
		return nil, kerrors.ErrInternalException.WithCause(errors.New("traces length error"))
	}
	res := make([]string, l/2)
	for i := 0; i < l; i += 6 {
		cluter := int(traces[i])*256 + int(traces[i+1])
		svc := int(traces[i+2])*256 + int(traces[i+3])
		cmd := int(traces[i+4])*256 + int(traces[i+5])
		trs := IdToTraces(cluter, svc, cmd)
		idx := i / 2
		res[idx] = trs[0]
		res[idx+1] = trs[1]
		res[idx+2] = trs[2]
	}
	return res, nil
}

func GetTrace(ctx context.Context, request *transport.Transport) string {
	res, err := BytesToTraces(request.Traces)
	var str string
	if err != nil {
		str = res[0]
		for i := 1; i < len(res); i++ {
			str = " " + res[i]
		}
	}
	return fmt.Sprintf("%s %s", GetTraceId(ctx), str)
}

func SetTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, "traceId", traceId)
}

func logRpcErr(ctx context.Context, request *transport.Transport, err error) {
	hlog.Errorf(
		fmt.Sprintf(
			"rpc error: %s %s %d %s \n%s",
			request.Addr,
			request.Cmd,
			request.Session,
			err.Error(),
			GetTrace(ctx, request),
		),
	)
}

func newTraceMW(cluster string) Middleware {
	return func(next Endpoint) Endpoint {
		return func(ctx context.Context, request, response *transport.Transport) error {
			traceId := GetTraceId(ctx)
			if traceId == "" {
				SetTraceId(ctx, fmt.Sprintf("%s-%d", request.Meta.Uuid, request.Session))
			}
			trs, err := TracesToBytes(cluster, request.Addr, request.Cmd)
			if err == nil {
				request.Traces = append(request.Traces, trs...)
			}
			return next(ctx, request, response)
		}
	}
}

func newActorMW() Middleware {
	return func(next Endpoint) Endpoint {
		return func(ctx context.Context, request, response *transport.Transport) error {
			_, ok := asys.ProcessRegistry.GetLocal(request.Addr)
			if ok {
				pid := actor.NewPID(asys.ProcessRegistry.Address, request.Addr)
				_, err := asys.Root.RequestFuture(pid, &Message{
					ctx:     ctx,
					request: request,
				}, time.Second*30).Result()

				if err != nil {
					if errors.Is(err, actor.ErrTimeout) {
						err = kerrors.ErrRPCTimeout.WithCause(err)
					}
					logRpcErr(ctx, request, err)
					return err
				}
			} else {
				return kerrors.ErrNoDestService
			}
			return next(ctx, request, response)
		}
	}
}

type Server interface {
	Run() error
	Stop() error
	GetServerInfo() *ServerInfo
}

type server struct {
	opt Options
	svc ksvc.Server
	mws []Middleware
}

func NewServer(opts ...Option) Server {
	s := &server{
		opt: *NewOptions(opts),
	}
	s.init()
	return s
}

func (s *server) init() {
	s.mws = richMWsWithBuilder(context.Background(), s.opt.MWBs, s)
	emw := newErrorHandleMW(s.opt.ErrHandle)
	tmw := newTraceMW(s.opt.Name)
	amw := newActorMW()

	//  error handler first
	s.mws = slices.Insert(s.mws, 0, emw)
	s.mws = append(s.mws, tmw, amw)
	eps := Chain(s.mws...)(nilEndpoint)
	s.svc = transportSvc.NewServer(NewTransportServiceImpl(eps))
}

func (s *server) Run() (err error) {
	err = s.svc.Run()
	if err != nil {
		return
	}
	err = s.register()
	return
}

func (s *server) Stop() (err error) {
	err = s.unregister()
	if err != nil {
		return
	}
	err = s.svc.Stop()
	return
}

func (s *server) GetServerInfo() *ServerInfo {
	return &ServerInfo{
		Name:    s.opt.Name,
		Address: s.opt.Address,
		Version: s.opt.Version,
	}
}

func (s *server) register() (err error) {
	if s.opt.Register != nil {
		err = s.opt.Register(s.GetServerInfo())
	}
	return
}

func (s *server) unregister() (err error) {
	if s.opt.UnRegister != nil {
		err = s.opt.UnRegister(s.GetServerInfo())
	}
	return
}
