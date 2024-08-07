package discovery

import (
	disc "github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/registry"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/kitex-contrib/registry-etcd/retry"
	balance "github.com/llsw/ikunet/internal/knet/balance"
)

func NewEtcdRegistryWithRetry(endpoints []string, retryConfig *retry.Config) (registry.Registry, error) {
	return etcd.NewEtcdRegistryWithRetry(endpoints, retryConfig)
}

func NewEtcdResolver(endpoints []string, opts ...etcd.Option) (disc.Resolver, error) {
	return etcd.NewEtcdResolver(endpoints, opts...)
}

func NewBalancer() *balance.Balancer {
	return &balance.Balancer{}
}
