package balance

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	transport "github.com/llsw/ikunet/internal/kitex_gen/transport"
	"github.com/llsw/ikunet/internal/knet/muxer/tcp"
	kretry "github.com/llsw/ikunet/internal/knet/registry"
	cmap "github.com/orcaman/concurrent-map"
	"golang.org/x/exp/rand"
)

var (
	_         loadbalance.Loadbalancer = &Balancer{}
	pic       loadbalance.Picker       = &picker{}
	cache     cmap.ConcurrentMap       = cmap.New()
	whiteList cmap.ConcurrentMap       = cmap.New()
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

type picker struct {
	dr           *discovery.Result
	ruleResolver RuleResolver
}

func (p *picker) Next(ctx context.Context, request interface{}) discovery.Instance {
	req := request.(*transport.Transport)
	var ins []discovery.Instance
	if muxer, ok := p.ruleResolver.GetRules(req.GetAddr()); ok {
		ins = make([]discovery.Instance, 0, len(p.dr.Instances))
		for _, v := range p.dr.Instances {
			match := muxer.Match(&tcp.Data{
				Req:      req,
				Instance: &v,
			})
			if match {
				ins = append(ins, v)
			}
		}
	} else {
		ins = p.dr.Instances
	}

	ll := len(ins)
	if ll > 0 {
		idx := rand.Intn(ll)
		return ins[idx]
	}
	return nil
}

type Balancer struct {
	ruleResolver RuleResolver
}

func NewBalancer(endpoints []string, opts ...kretry.Option) (*Balancer, error) {
	r, err := NewEtcdRuleResolver(endpoints, opts...)
	if err != nil {
		return nil, err
	}

	return &Balancer{
		ruleResolver: r,
	}, nil
}

func (b *Balancer) GetPicker(dr discovery.Result) loadbalance.Picker {
	// 每次都要赋值，防止服务发现结果改变了
	pic.(*picker).dr = &dr
	pic.(*picker).ruleResolver = b.ruleResolver
	return pic
}

func (b *Balancer) Name() string {
	return "kbalancer"
}

func GetBlCallKey(call string) string {
	return fmt.Sprintf("blcall-%s", call)
}
