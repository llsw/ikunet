package balance

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	kmuxer "github.com/llsw/ikunet/internal/knet/muxer/tcp"
	kretry "github.com/llsw/ikunet/internal/knet/registry"
	kc "github.com/llsw/ikunet/pkg/common/cache"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	etcdPrefixTpl = "knet/rule/%v"
	EXPIRE        = 1 * time.Minute
)

var (
	rulecache = kc.New[*kmuxer.Muxer](6*time.Minute, 6*time.Minute)
)

func ruleKeyPrefix(serviceName string) string {
	return fmt.Sprintf(etcdPrefixTpl, serviceName)
}

// serviceKey generates the key used to stored in etcd.
func ruleKey(serviceName string) string {
	return ruleKeyPrefix(serviceName)
}

type RuleResolver interface {
	GetRules(serviceName string) (*kmuxer.Muxer, bool)
	SetRules(serviceName string, rules []string) error
}

type etcdRuleResolver struct {
	etcdClient *clientv3.Client
}

// SetRules implements RuleResolver.
func (e *etcdRuleResolver) SetRules(serviceName string, rules []string) error {
	key := ruleKey(serviceName)

	val, err := json.Marshal(&rules)
	if err != nil {
		hlog.Errorf("set rules  marshal fail error:%s", err.Error())
		return err
	}

	var muxer *kmuxer.Muxer
	muxer, err = kmuxer.NewMuxer()

	if err != nil {
		hlog.Errorf("set rules new muxer fail, key:%s error:%s", key, err.Error())
		return err
	}

	for _, v := range rules {
		err = muxer.AddRoute(v)
		hlog.Errorf("add rule to muxer fail, key:%s error:%s", key, err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = e.etcdClient.Put(ctx, key, string(val))
	if err != nil {
		hlog.Errorf("set rules put ectd fail error:%s", err.Error())
		return err
	}

	rulecache.Set(serviceName, muxer, EXPIRE)

	return nil
}

func genMuxer(key string, kvs []*mvccpb.KeyValue) (val *kmuxer.Muxer, err error) {
	if len(kvs) > 0 {
		val, err = kmuxer.NewMuxer()
		if err != nil {
			hlog.Errorf("get rules new muxer fail, key:%s error:%s", key, err.Error())
			return
		}
		count := 0
		for _, kv := range kvs {
			var info []string
			err = json.Unmarshal(kv.Value, &info)
			if err != nil {
				hlog.Warnf("fail to unmarshal with err: %v, ignore key: %v", err, string(kv.Key))
				continue
			}
			if len(info) == 0 {
				hlog.Warnf("value is empty, ignore key: %v", string(kv.Key))
				continue
			}
			for _, rule := range info {
				count = count + 1
				err = val.AddRoute(rule)
				if err != nil {
					hlog.Errorf("add rule to muxer fail, key:%s error:%s", key, err.Error())
					return nil, err
				}
			}
		}
		if count == 0 {
			val = nil
			err = fmt.Errorf("not exit")
		}
	} else {
		err = fmt.Errorf("not exit")
	}
	return
}

func (e *etcdRuleResolver) getRules(key string) (val *kmuxer.Muxer, err error) {
	key = ruleKey(key)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := e.etcdClient.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	return genMuxer(key, resp.Kvs)
}

func (e *etcdRuleResolver) GetRules(serviceName string) (*kmuxer.Muxer, bool) {
	muxer, ok := rulecache.Get(serviceName, func(key string) (val *kmuxer.Muxer, err error) {
		return e.getRules(serviceName)
	}, EXPIRE)
	if ok {
		return muxer, true
	}
	return nil, false
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
