package server

import (
	"context"

	ksvc "github.com/cloudwego/kitex/server"
	transport "github.com/llsw/ikunnet/kitex_gen/transport"
	transportSvc "github.com/llsw/ikunnet/kitex_gen/transport/transportservice"
)

var (
	_ Server = &server{}
)

// TransportServiceImpl implements the last service interface defined in the IDL.
type TransportServiceImpl struct{}

// Call implements the TransportServiceImpl interface.
func (s *TransportServiceImpl) Call(ctx context.Context, req *transport.Transport) (resp *transport.Transport, err error) {
	// TODO: Your code here...
	return
}

type Server interface {
	Run() error
	Stop() error
	GetServerInfo() *ServerInfo
}

type server struct {
	opt Options
	svc ksvc.Server
}

func NewServer(opts ...Option) Server {
	s := &server{
		opt: *NewOptions(opts),
	}
	s.svc = transportSvc.NewServer(new(TransportServiceImpl))
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
