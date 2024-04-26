package trace

import (
	"context"
	"fmt"

	transport "github.com/llsw/ikunet/internal/kitex_gen/transport"
	knet "github.com/llsw/ikunet/internal/knet"
)

type GetTraceId func(context.Context) string
type SetTraceId func(ctx context.Context, request *transport.Transport) context.Context

type TracesToBytes func(cluster, svc, cmd string) []byte
type BytesToTraces func([]byte) (cluster, svc, cmd string)

func DefaultGetTraceId(ctx context.Context) string {
	return ctx.Value(knet.TRACEID_KEY).(string)
}

func DefaultSetTraceId(ctx context.Context, request *transport.Transport) context.Context {
	traceId := fmt.Sprintf("%s-%d", request.Meta.Uuid, request.Session)
	ctx = context.WithValue(ctx, knet.TRACEID_KEY, traceId)
	return ctx
}

func DefaultSetTrace(cluster, svc, cmd string) []byte {
	return nil
}
func DefaultGetTrace([]byte) (cluster, svc, cmd string) {
	return "", "", ""
}

// func GetTraceId(ctx context.Context) string {
// 	traceId, ok := ctx.Value("traceId").(string)
// 	if !ok {
// 		traceId = ""
// 	}
// 	return traceId
// }

// func TracesToBytes(cluster, svc, cmd string) ([]byte, error) {
// 	return []byte{0, 0, 0, 0, 0, 0}, nil
// }

// func IdToTraces(cluster, svc, cmd int) []string {
// 	return []string{
// 		string(cluster),
// 		string(svc),
// 		string(cmd),
// 	}
// }

// func BytesToTraces(traces []byte) ([]string, error) {
// 	l := len(traces)
// 	if l < 6 || l%6 != 0 {
// 		return nil, kerrors.ErrInternalException.WithCause(errors.New("traces length error"))
// 	}
// 	res := make([]string, l/2)
// 	for i := 0; i < l; i += 6 {
// 		cluter := int(traces[i])*256 + int(traces[i+1])
// 		svc := int(traces[i+2])*256 + int(traces[i+3])
// 		cmd := int(traces[i+4])*256 + int(traces[i+5])
// 		trs := IdToTraces(cluter, svc, cmd)
// 		idx := i / 2
// 		res[idx] = trs[0]
// 		res[idx+1] = trs[1]
// 		res[idx+2] = trs[2]
// 	}
// 	return res, nil
// }
