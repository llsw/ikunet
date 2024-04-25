package server

import (
	"context"
	"fmt"
	"slices"

	actor "github.com/asynkron/protoactor-go/actor"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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
	emw := newErrorHandleMW(s)
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
