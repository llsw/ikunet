package client

import (
	"context"
	"errors"
	"time"

	knet "github.com/llsw/ikunet/internal/knet"

	kerrors "github.com/cloudwego/kitex/pkg/kerrors"
	transport "github.com/llsw/ikunet/internal/kitex_gen/transport"
	midw "github.com/llsw/ikunet/internal/knet/middleware"
)

func newCallMW(c *client) midw.Middleware {
	return func(next midw.Endpoint) midw.Endpoint {
		return func(ctx context.Context, request, response *transport.Transport) (err error) {
			ctx, cancel := context.WithTimeout(ctx, time.Second*knet.TIME_OUT)
			defer cancel()
			response, err = c.client.Call(ctx, request)

			if errors.Is(err, context.DeadlineExceeded) {
				err = kerrors.ErrRPCTimeout.WithCause(err)
				return
			}
			return next(ctx, request, response)
		}
	}
}
