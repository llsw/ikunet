package discovery

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/utils"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/kitex-contrib/registry-etcd/retry"
	balance "github.com/llsw/ikunet/internal/knet/balance"
	kclient "github.com/llsw/ikunet/internal/knet/client"
	kserver "github.com/llsw/ikunet/internal/knet/server"
)

func WithEtcdRetry(endpoints []string, retryConfig *retry.Config, opts ...etcd.Option) kserver.Option {
	r, err := etcd.NewEtcdRegistryWithRetry(endpoints, retryConfig)
	if err != nil {
		hlog.Fatal(err)
	}
	return kserver.Option{
		F: func(o *kserver.Options, di *utils.Slice) {
			o.Retry = r
			r, err := etcd.NewEtcdResolver(endpoints, opts...)
			if err != nil {
				hlog.Fatal(err)
				return
			}
			o.Resolver = r
		},
	}
}

func WithEtcdResolver(endpoints []string, opts ...etcd.Option) kclient.Option {
	return kclient.Option{
		F: func(o *kclient.Options, di *utils.Slice) {
			r, err := etcd.NewEtcdResolver(endpoints, opts...)
			if err != nil {
				hlog.Fatal(err)
				return
			}
			o.Resolver = r
		},
	}
}

func WithBalancer() kclient.Option {
	return kclient.Option{
		F: func(o *kclient.Options, di *utils.Slice) {
			o.Balancer = &balance.Balancer{}
		},
	}
}

func WithBalancerCalls(calls []string) kserver.Option {
	return kserver.Option{
		F: func(o *kserver.Options, di *utils.Slice) {
			o.BalancerCalls = calls
		},
	}
}
