package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudwego/kitex/pkg/gofunc"
	"github.com/cloudwego/kitex/pkg/utils"
	transport "github.com/llsw/ikunet/internal/kitex_gen/transport"
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
	GetTraceId GetTraceId
	SetTraceId SetTraceId
	GetTrace   BytesToTraces
	SetTrace   TracesToBytes
}

// NewOptions creates a default options.
func NewOptions(opts []Option) *Options {
	o := &Options{
		ExitSignal: DefaultSysExitSignal,
		GetTraceId: DefaultGetTraceId,
		SetTraceId: DefaultSetTraceId,
		GetTrace:   DefaultGetTrace,
		SetTrace:   DefaultSetTrace,
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

func DefaultGetTraceId(ctx context.Context) string {
	return ctx.Value(TRACEID_KEY).(string)
}

func DefaultSetTraceId(ctx context.Context, request *transport.Transport) context.Context {
	traceId := fmt.Sprintf("%s-%d", request.Meta.Uuid, request.Session)
	ctx = context.WithValue(ctx, TRACEID_KEY, traceId)
	return ctx
}

func DefaultSetTrace(cluster, svc, cmd string) []byte {
	return nil
}
func DefaultGetTrace([]byte) (cluster, svc, cmd string) {
	return "", "", ""
}
