package client

import (
	"context"
	"fmt"
	"slices"

	midw "github.com/llsw/ikunet/internal/knet/middleware"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	kclient "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	transport "github.com/llsw/ikunet/internal/kitex_gen/transport"
	transportSvc "github.com/llsw/ikunet/internal/kitex_gen/transport/transportservice"
)

var (
	_ Client = &client{}
)

type Client interface {
	Call(ctx context.Context, request *transport.Transport) (response *transport.Transport, err error)
}

type client struct {
	opt    *Options
	client transportSvc.Client
	mws    []midw.Middleware
	eps    midw.Endpoint
}

func NewClient(opts ...Option) (Client, error) {
	c := &client{
		opt: NewOptions(opts),
	}
	err := c.init()
	return c, err
}

func (c *client) init() (err error) {
	c.client, err = transportSvc.NewClient(c.opt.Name, kclient.WithResolver(c.opt.Resolver))
	if err != nil {
		return
	}
	c.mws = midw.RichMWsWithBuilder(context.Background(), c.opt.MWBs, c.mws)
	emw := newErrorHandleMW(c)
	// tmw := newTraceMW(s)
	amw := newCallMW(c)
	//  error handler first
	c.mws = slices.Insert(c.mws, 0, emw)
	c.mws = append(c.mws, amw)
	c.eps = midw.Chain(c.mws...)(midw.NilEndpoint)
	return
}

func (c *client) Call(ctx context.Context, request *transport.Transport) (response *transport.Transport, err error) {
	response = &transport.Transport{}
	// err = kerrors.ErrRPCTimeout.WithCause(errors.New("unimplement"))
	err = c.eps(ctx, request, response)
	return
}

func (c *client) GetResolver() discovery.Resolver {
	return c.opt.Resolver
}

func (c *client) logRpcErr(ctx context.Context, request *transport.Transport, err error) {
	cluster, addr, cmd := c.opt.GetTrace(request.Traces)
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

// client, err := echo.NewClient("echo", client.WithResolver(p2p.NewP2PResolver("tcp", ":8888")))
// if err != nil {
// 	log.Fatal(err)
// }
// for {
// 	req := &api.Request{Message: "my request"}
// 	resp, err := client.Echo(context.Background(), req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println(resp)
// 	time.Sleep(time.Second)
// }
