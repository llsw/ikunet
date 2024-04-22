// Code generated by Kitex v0.9.1. DO NOT EDIT.
package transportservice

import (
	server "github.com/cloudwego/kitex/server"
	transport "github.com/llsw/ikunet/kitex_gen/transport"
)

// NewServer creates a server.Server with the given handler and options.
func NewServer(handler transport.TransportService, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}

func RegisterService(svr server.Server, handler transport.TransportService, opts ...server.RegisterOption) error {
	return svr.RegisterService(serviceInfo(), handler, opts...)
}