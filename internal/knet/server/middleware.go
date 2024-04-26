package server

import (
	"context"
	"errors"
	"time"

	actor "github.com/asynkron/protoactor-go/actor"
	"github.com/cloudwego/kitex/pkg/kerrors"
	transport "github.com/llsw/ikunet/internal/kitex_gen/transport"
	knet "github.com/llsw/ikunet/internal/knet"
	midw "github.com/llsw/ikunet/internal/knet/middleware"
)

// newErrorHandleMW provides a hook point for server error handling.
func newErrorHandleMW(s *server) midw.Middleware {
	return func(next midw.Endpoint) midw.Endpoint {
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

func newTraceMW(s *server) midw.Middleware {
	return func(next midw.Endpoint) midw.Endpoint {
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

func newActorMW(s *server) midw.Middleware {
	return func(next midw.Endpoint) midw.Endpoint {
		return func(ctx context.Context, request, response *transport.Transport) error {
			_, ok := asys.ProcessRegistry.GetLocal(request.Addr)
			if ok {
				pid := actor.NewPID(asys.ProcessRegistry.Address, request.Addr)
				_, err := asys.Root.RequestFuture(pid, &Message{
					ctx:     ctx,
					request: request,
				}, time.Second*knet.TIME_OUT).Result()

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
