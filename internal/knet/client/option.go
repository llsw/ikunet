package client

import (
	"context"
	"net"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/utils"
	midw "github.com/llsw/ikunet/internal/knet/middleware"
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
	Cluster   string
	Name      string
	Version   string
	Address   net.Addr
	ErrHandle func(context.Context, error) error
	DebugInfo utils.Slice
	Resolver  discovery.Resolver
	MWBs      []midw.MiddlewareBuilder
}

// NewOptions creates a default options.
func NewOptions(opts []Option) *Options {
	o := &Options{}
	ApplyOptions(opts, o)
	return o
}

// ApplyOptions applies the given options.
func ApplyOptions(opts []Option, o *Options) {
	for _, op := range opts {
		op.F(o, &o.DebugInfo)
	}
}
