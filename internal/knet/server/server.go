package server

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	actor "github.com/asynkron/protoactor-go/actor"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	ksvc "github.com/cloudwego/kitex/server"
	transport "github.com/llsw/ikunet/internal/kitex_gen/transport"
	transportSvc "github.com/llsw/ikunet/internal/kitex_gen/transport/transportservice"
	cv "github.com/llsw/ikunet/internal/knet/const"
	midw "github.com/llsw/ikunet/internal/knet/middleware"
	kretry "github.com/llsw/ikunet/internal/knet/registry"
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
	eps midw.Endpoint
}

func NewTransportServiceImpl(eps midw.Endpoint) *TransportServiceImpl {
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
	GetServerInfo() *Info
}

type server struct {
	opt *Options
	svc ksvc.Server
	mws []midw.Middleware
	sync.Mutex
}

func NewServer(opts ...Option) Server {
	s := &server{
		opt: NewOptions(opts),
	}
	s.init()
	return s
}

func (s *server) serverTags(info *Info) map[string]string {
	tags := make(map[string]string)
	tags[cv.TAG_VERSION] = info.Version
	tags[cv.TAG_CLUSTER] = info.Cluster
	tags[cv.TAG_TYPE] = info.Type
	tags[cv.TAG_ID] = info.Id

	if info.Maintain {
		tags[cv.TAG_MAINTAIN] = cv.TAG_MAINTAIN_ON
	} else {
		tags[cv.TAG_MAINTAIN] = cv.TAG_MAINTAIN_OFF
	}

	return tags
}

func (s *server) init() {
	s.mws = midw.RichMWsWithBuilder(context.Background(), s.opt.MWBs, s.mws)
	emw := newErrorHandleMW(s)
	tmw := newTraceMW(s)
	amw := newActorMW(s)

	//  error handler first
	s.mws = slices.Insert(s.mws, 0, emw)
	s.mws = append(s.mws, tmw, amw)
	eps := midw.Chain(s.mws...)(midw.NilEndpoint)
	info := s.GetServerInfo()

	s.svc = transportSvc.NewServer(
		NewTransportServiceImpl(eps),
		ksvc.WithServiceAddr(s.opt.Address),
		// ksvc.WithRegistry(s.opt.Retry),
		ksvc.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: info.Name,
				// 增加tag，tag可包含服务版本
				Tags: s.serverTags(info),
			},
		),
	)
}

func (s *server) buildRegistryInfo() *registry.Info {
	svc := s.GetServerInfo()
	return &registry.Info{
		ServiceName: svc.Name,
		Addr:        svc.Address,
		Tags:        s.serverTags(svc),
		StartTime:   time.Now(),
		Weight:      1,
	}
}

func (s *server) waitExit() error {
	exitSignal := s.opt.ExitSignal()
	// service may not be available as soon as startup.
	delayRegister := time.After(1 * time.Second)
	for {
		select {
		case err := <-exitSignal:
			return err
		case <-delayRegister:
			s.Lock()
			if s.opt.Registry != nil {
				s.opt.RegistryInfo = s.buildRegistryInfo()
				if err := s.opt.Registry.Register(s.opt.RegistryInfo); err != nil {
					s.Unlock()
					return err
				}
			}
			s.Unlock()
		}
	}
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

func (s *server) GetServerInfo() *Info {
	return &Info{
		Cluster:  s.opt.Cluster,
		Name:     s.opt.Name,
		Address:  s.opt.Address,
		Version:  s.opt.Version,
		Maintain: s.opt.Maintain,
	}
}

func (s *server) GetResolver() discovery.Resolver {
	return s.opt.Resolver
}

func (s *server) register() (err error) {
	if s.opt.Register != nil {
		err = s.opt.Register(s.GetServerInfo())
	}

	if err == nil {
		s.waitExit()
	}
	return
}

func (s *server) unregister() (err error) {
	if s.opt.UnRegister != nil {
		err = s.opt.UnRegister(s.GetServerInfo())
	}
	if err == nil {
		if s.opt.Registry != nil {
			err = s.opt.Registry.Deregister(s.opt.RegistryInfo)
		}
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

func (s *server) modify() (err error) {
	s.opt.RegistryInfo = s.buildRegistryInfo()
	if s.opt.Registry != nil {
		if retry, ok := s.opt.Registry.(*kretry.EtcdRegistry); ok {
			err = retry.Modify(s.opt.RegistryInfo)
		}
	}
	return
}

func (s *server) Open() error {
	s.opt.Maintain = false
	return s.modify()
}

func (s *server) Maintain() error {
	s.opt.Maintain = true
	return s.modify()
}
