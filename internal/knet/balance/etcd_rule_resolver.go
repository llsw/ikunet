package balance

import (
	"fmt"
	"time"

	kretry "github.com/llsw/ikunet/internal/knet/registry"
	kc "github.com/llsw/ikunet/pkg/common/cache"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	etcdPrefixTpl = "knet/rule/%v/"
)

var (
	rulecache = kc.New[[]string](6*time.Minute, 6*time.Minute)
)

func ruleKeyPrefix(serviceName string) string {
	return fmt.Sprintf(etcdPrefixTpl, serviceName)
}

// serviceKey generates the key used to stored in etcd.
func ruleKey(serviceName, addr string) string {
	return ruleKeyPrefix(serviceName) + addr
}

type RuleResolver interface {
	GetRules(serviceName string) []string
	SetRules(serviceName string, rules []string)
}

type etcdRuleResolver struct {
	etcdClient *clientv3.Client
}

// SetRules implements RuleResolver.
func (e *etcdRuleResolver) SetRules(serviceName string, rules []string) {
	panic("unimplemented")
}

func (e *etcdRuleResolver) GetRules(serviceName string) []string {
	rules, ok := rulecache.Get(serviceName, func(key string) (val []string, err error) {
		return nil , fmt.Errorf("unimplemented")
	})

	if ok {
		return rules
	}
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
