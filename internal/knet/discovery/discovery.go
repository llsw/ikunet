package discovery

import (
	disc "github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/kitex-contrib/registry-etcd/retry"
	balance "github.com/llsw/ikunet/internal/knet/balance"
	kretry "github.com/llsw/ikunet/internal/knet/registry"
)

func NewEtcdRegistryWithRetry(endpoints []string, retryConfig *retry.Config) (registry.Registry, error) {
	return kretry.NewEtcdRegistryWithRetry(endpoints, retryConfig)
}

func NewEtcdResolver(endpoints []string, opts ...kretry.Option) (disc.Resolver, error) {
	return kretry.NewEtcdResolver(endpoints, opts...)
}

func NewBalancer() *balance.Balancer {
	return &balance.Balancer{}
}
