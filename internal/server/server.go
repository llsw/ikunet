package server

import (
	"context"

	ksvc "github.com/cloudwego/kitex/server"
	transport "github.com/llsw/ikunet/internal/kitex_gen/transport"
	transportSvc "github.com/llsw/ikunet/internal/kitex_gen/transport/transportservice"
)

var (
	_ Server = &server{}
)

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
			return nil
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
	s.mws = richMWsWithBuilder(context.Background(), s.opt.MWBs, s)
	emw := newErrorHandleMW(s.opt.ErrHandle)
	s.mws = append(s.mws, emw)
	eps := Chain(s.mws...)(func(ctx context.Context, req, resp *transport.Transport) error {
		return nil
	})
	s.svc = transportSvc.NewServer(NewTransportServiceImpl(eps))
	return s
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
