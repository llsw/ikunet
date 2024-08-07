package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/gofunc"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/utils"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/kitex-contrib/registry-etcd/retry"
	kdisc "github.com/llsw/ikunet/internal/knet/discovery"
	midw "github.com/llsw/ikunet/internal/knet/middleware"
	trace "github.com/llsw/ikunet/internal/knet/trace"
)

func init() {
}

type Info struct {
	Cluster string
	Name    string
	Version string
	Address net.Addr
	Type    string
	Id      string
}

// Option is the only way to config a server.
type Option struct {
	F func(o *Options, di *utils.Slice)
}

// Options is used to initialize the server.
type Options struct {
	Cluster       string
	Name          string
	Version       string
	Address       net.Addr
	Type          string
	Id            string
	ErrHandle     func(context.Context, error) error
	ExitSignal    func() <-chan error
	DebugInfo     utils.Slice
	Register      func(info *Info) error
	UnRegister    func(info *Info) error
	MWBs          []midw.MiddlewareBuilder
	GetTraceId    trace.GetTraceId
	SetTraceId    trace.SetTraceId
	GetTrace      trace.BytesToTraces
	SetTrace      trace.TracesToBytes
	Retry         registry.Registry
	Resolver      discovery.Resolver
	BalancerCalls []string
}

// NewOptions creates a default options.
func NewOptions(opts []Option) *Options {
	o := &Options{
		ExitSignal: DefaultSysExitSignal,
		GetTraceId: trace.DefaultGetTraceId,
		SetTraceId: trace.DefaultSetTraceId,
		GetTrace:   trace.DefaultGetTrace,
		SetTrace:   trace.DefaultSetTrace,
	}
	ApplyOptions(opts, o)
	return o
}

// ApplyOptions applies the given options.
func ApplyOptions(opts []Option, o *Options) {
	for _, op := range opts {
		op.F(o, &o.DebugInfo)
	}
}

func DefaultSysExitSignal() <-chan error {
	errCh := make(chan error, 1)
	gofunc.GoFunc(context.Background(), func() {
		sig := SysExitSignal()
		defer signal.Stop(sig)
		<-sig
		errCh <- nil
	})
	return errCh
}

func SysExitSignal() chan os.Signal {
	signals := make(chan os.Signal, 1)
	notifications := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	if !signal.Ignored(syscall.SIGHUP) {
		notifications = append(notifications, syscall.SIGHUP)
	}
	signal.Notify(signals, notifications...)
	return signals
}

func WithBalancerCalls(calls []string) Option {
	return Option{
		F: func(o *Options, di *utils.Slice) {
			o.BalancerCalls = calls
		},
	}
}

func WithEtcdRetry(endpoints []string, retryConfig *retry.Config, opts ...etcd.Option) Option {
	r, err := kdisc.NewEtcdRegistryWithRetry(endpoints, retryConfig)
	if err != nil {
		hlog.Fatal(err)
	}
	return Option{
		F: func(o *Options, di *utils.Slice) {
			o.Retry = r
			r, err := kdisc.NewEtcdResolver(endpoints, opts...)
			if err != nil {
				hlog.Fatal(err)
				return
			}
			o.Resolver = r
		},
	}
}
