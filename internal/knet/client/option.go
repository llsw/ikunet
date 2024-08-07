package client

import (
	"context"
	"net"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/utils"
	etcd "github.com/kitex-contrib/registry-etcd"
	kdisc "github.com/llsw/ikunet/internal/knet/discovery"
	midw "github.com/llsw/ikunet/internal/knet/middleware"
	"github.com/llsw/ikunet/internal/knet/trace"
)

func init() {
}

type Info struct {
	Cluster string
	Name    string
	Version string
	Address net.Addr
}

// Option is the only way to config a server.
type Option struct {
	F func(o *Options, di *utils.Slice)
}

// Options is used to initialize the server.
type Options struct {
	Cluster    string
	Name       string
	Version    string
	Address    net.Addr
	ErrHandle  func(context.Context, error) error
	DebugInfo  utils.Slice
	Resolver   discovery.Resolver
	Balancer   loadbalance.Loadbalancer
	MWBs       []midw.MiddlewareBuilder
	GetTraceId trace.GetTraceId
	SetTraceId trace.SetTraceId
	GetTrace   trace.BytesToTraces
	SetTrace   trace.TracesToBytes
}

// NewOptions creates a default options.
func NewOptions(opts []Option) *Options {
	o := &Options{
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

func WithEtcdResolver(endpoints []string, opts ...etcd.Option) Option {
	return Option{
		F: func(o *Options, di *utils.Slice) {
			r, err := kdisc.NewEtcdResolver(endpoints, opts...)
			if err != nil {
				hlog.Fatal(err)
				return
			}
			o.Resolver = r
		},
	}
}

func WithBalancer() Option {
	return Option{
		F: func(o *Options, di *utils.Slice) {
			o.Balancer = kdisc.NewBalancer()
		},
	}
}
