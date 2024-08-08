package balance

import (
	kretry "github.com/llsw/ikunet/internal/knet/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type RuleResolver interface {
	GetRules(dstServiceName string) []string
}

type etcdRuleResolver struct {
	etcdClient *clientv3.Client
}

func (e *etcdRuleResolver) GetRules(dstServiceName string) []string {
	return nil
}

func NewEtcdRuleResolver(endpoints []string, opts ...kretry.Option) (RuleResolver, error) {
	cfg := clientv3.Config{
		Endpoints: endpoints,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	return &etcdRuleResolver{
		etcdClient: etcdClient,
	}, nil
}
