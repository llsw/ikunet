package middleware

import (
	"context"

	transport "github.com/llsw/ikunet/internal/kitex_gen/transport"
)

type Endpoint func(ctx context.Context, req, resp *transport.Transport) (err error)

type Middleware func(Endpoint) Endpoint

type MiddlewareBuilder func(ctx context.Context) Middleware

// Chain connect middlewares into one middleware.
func Chain(mws ...Middleware) Middleware {
	return func(next Endpoint) Endpoint {
		for i := len(mws) - 1; i >= 0; i-- {
			next = mws[i](next)
		}
		return next
	}
}

func RichMWsWithBuilder(ctx context.Context, mwBs []MiddlewareBuilder, mws []Middleware) []Middleware {
	for i := range mwBs {
		mws = append(mws, mwBs[i](ctx))
	}
	return mws
}
