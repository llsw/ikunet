package main

import (
	"context"

	transport "github.com/llsw/ikunet/internal/kitex_gen/transport"
)

// TransportServiceImpl implements the last service interface defined in the IDL.
type TransportServiceImpl struct{}

// Call implements the TransportServiceImpl interface.
func (s *TransportServiceImpl) Call(ctx context.Context, req *transport.Transport) (resp *transport.Transport, err error) {
	// TODO: Your code here...
	return
}
