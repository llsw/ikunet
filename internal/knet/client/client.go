package client

import (
	"context"
	"errors"

	midw "github.com/llsw/ikunet/internal/knet/middleware"

	kclient "github.com/cloudwego/kitex/client"
	kerrors "github.com/cloudwego/kitex/pkg/kerrors"
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
	return
}

func (c *client) Call(ctx context.Context, request *transport.Transport) (response *transport.Transport, err error) {
	response = &transport.Transport{}
	err = kerrors.ErrRPCTimeout.WithCause(errors.New("unimplement"))
	return
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
