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

func NewBalancer(endpoints []string, opts ...kretry.Option) (*balance.Balancer, error) {
	return balance.NewBalancer(endpoints, opts...)
}
