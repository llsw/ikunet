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

type GetTraceId func(context.Context) string
type SetTraceId func(ctx context.Context, request *transport.Transport) context.Context

type TracesToBytes func(cluster, svc, cmd string) ([]byte, error)
type BytesToTraces func([]byte) (cluster, svc, cmd string)

// func GetTraceId(ctx context.Context) string {
// 	traceId, ok := ctx.Value("traceId").(string)
// 	if !ok {
// 		traceId = ""
// 	}
// 	return traceId
// }

// func TracesToBytes(cluster, svc, cmd string) ([]byte, error) {
// 	return []byte{0, 0, 0, 0, 0, 0}, nil
// }

// func IdToTraces(cluster, svc, cmd int) []string {
// 	return []string{
// 		string(cluster),
// 		string(svc),
// 		string(cmd),
// 	}
// }

// func BytesToTraces(traces []byte) ([]string, error) {
// 	l := len(traces)
// 	if l < 6 || l%6 != 0 {
// 		return nil, kerrors.ErrInternalException.WithCause(errors.New("traces length error"))
// 	}
// 	res := make([]string, l/2)
// 	for i := 0; i < l; i += 6 {
// 		cluter := int(traces[i])*256 + int(traces[i+1])
// 		svc := int(traces[i+2])*256 + int(traces[i+3])
// 		cmd := int(traces[i+4])*256 + int(traces[i+5])
// 		trs := IdToTraces(cluter, svc, cmd)
// 		idx := i / 2
// 		res[idx] = trs[0]
// 		res[idx+1] = trs[1]
// 		res[idx+2] = trs[2]
// 	}
// 	return res, nil
// }

func newTraceMW(s *server) Middleware {
	return func(next Endpoint) Endpoint {
		return func(ctx context.Context, request, response *transport.Transport) error {
			traceId := s.opt.GetTraceId(ctx)
			if traceId == "" {
				ctx = s.opt.SetTraceId(ctx, request)
			}
			trs, err := s.opt.SetTrace(s.opt.Name, request.Addr, request.Cmd)
			if err == nil {
				request.Traces = append(request.Traces, trs...)
			}
			return next(ctx, request, response)
		}
	}
}

func newActorMW(s *server) Middleware {
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
					s.logRpcErr(ctx, request, err)
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
	tmw := newTraceMW(s)
	amw := newActorMW(s)

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

func (s *server) logRpcErr(ctx context.Context, request *transport.Transport, err error) {
	cluster, addr, cmd := s.opt.GetTrace(request.Traces)
	hlog.Errorf(
		fmt.Sprintf(
			"rpc error: %s %s %d %s \n%s %s %s",
			request.Addr,
			request.Cmd,
			request.Session,
			err.Error(),
			cluster,
			addr,
			cmd,
		),
	)
}
