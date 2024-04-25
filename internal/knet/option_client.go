package knet

import (
	"context"
	"net"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/utils"
)

func init() {
}

type ClientInfo struct {
	Cluster string
	Name    string
	Version string
	Address net.Addr
}

// ClientOption is the only way to config a server.
type ClientOption struct {
	F func(o *ClientOptions, di *utils.Slice)
}

// ClientOptions is used to initialize the server.
type ClientOptions struct {
	Cluster   string
	Name      string
	Version   string
	Address   net.Addr
	ErrHandle func(context.Context, error) error
	DebugInfo utils.Slice
	Resolver  discovery.Resolver
	MWBs      []MiddlewareBuilder
}

// NewClientOptions creates a default options.
func NewClientOptions(opts []ClientOption) *ClientOptions {
	o := &ClientOptions{}
	ApplyClientOptions(opts, o)
	return o
}

// ApplyClientOptions applies the given options.
func ApplyClientOptions(opts []ClientOption, o *ClientOptions) {
	for _, op := range opts {
		op.F(o, &o.DebugInfo)
	}
}
