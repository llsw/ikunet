package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudwego/kitex/pkg/gofunc"
	"github.com/cloudwego/kitex/pkg/utils"
)

func init() {
}

type ServerInfo struct {
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
	Name       string
	Version    string
	Address    net.Addr
	ErrHandle  func(context.Context, error) error
	ExitSignal func() <-chan error
	DebugInfo  utils.Slice
	Register   func(info *ServerInfo) error
	UnRegister func(info *ServerInfo) error
	MWBs       []MiddlewareBuilder
}

// NewOptions creates a default options.
func NewOptions(opts []Option) *Options {
	o := &Options{
		ExitSignal: DefaultSysExitSignal,
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
