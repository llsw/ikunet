package server

import (
	"context"
	"errors"
	"time"

	actor "github.com/asynkron/protoactor-go/actor"
	kerrors "github.com/cloudwego/kitex/pkg/kerrors"
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

func richMWsWithBuilder(ctx context.Context, mwBs []MiddlewareBuilder, s *server) []Middleware {
	for i := range mwBs {
		s.mws = append(s.mws, mwBs[i](ctx))
	}
	return s.mws
}

func nilEndpoint(ctx context.Context, req, resp *transport.Transport) error {
	return nil
}

// newErrorHandleMW provides a hook point for server error handling.
func newErrorHandleMW(s *server) Middleware {
	return func(next Endpoint) Endpoint {
		return func(ctx context.Context, request, response *transport.Transport) error {
			err := next(ctx, request, response)
			if err == nil {
				s.logRpcErr(ctx, request, err)
				return nil
			}
			if s.opt.ErrHandle != nil {
				return s.opt.ErrHandle(ctx, err)
			}
			return err
		}
	}
}

func newTraceMW(s *server) Middleware {
	return func(next Endpoint) Endpoint {
		return func(ctx context.Context, request, response *transport.Transport) error {
			traceId := s.opt.GetTraceId(ctx)
			if traceId == "" {
				ctx = s.opt.SetTraceId(ctx, request)
			}
			trs := s.opt.SetTrace(s.opt.Name, request.Addr, request.Cmd)
			if len(trs) > 0 {
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
					return err
				}
			} else {
				return kerrors.ErrNoDestService
			}
			return next(ctx, request, response)
		}
	}
}
