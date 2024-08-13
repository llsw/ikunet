package balance

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	transport "github.com/llsw/ikunet/internal/kitex_gen/transport"
	"github.com/llsw/ikunet/internal/knet/muxer/tcp"
	kretry "github.com/llsw/ikunet/internal/knet/registry"
	cmap "github.com/orcaman/concurrent-map"
	"golang.org/x/exp/rand"
	"golang.org/x/sync/singleflight"
)

var (
	_ loadbalance.Loadbalancer = &Balancer{}
	_ loadbalance.Picker       = &picker{}

	cache cmap.ConcurrentMap = cmap.New()

	pickerPool = sync.Pool{
		New: newPicker,
	}
	sfg singleflight.Group
)

const (
	RULE_CACHE_KEY    = "RCK"
	MATCH_INS_CAP     = 10
	MATCH_RESET_CHECK = 200
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

type picker struct {
	dr           *discovery.Result
	ruleResolver RuleResolver
	matchIns     []discovery.Instance
	useCount     int
}

func newMatchIns() []discovery.Instance {
	return make([]discovery.Instance, MATCH_INS_CAP)
}

// Recycle implements internal.Reusable.
func (p *picker) Recycle() {
	p.dr = nil
	p.ruleResolver = nil
	p.useCount = p.useCount + 1

	if p.useCount == MATCH_RESET_CHECK {
		p.useCount = 0
		p.matchIns = newMatchIns()
	}

	pickerPool.Put(p)
}

func newPicker() any {
	return &picker{
		matchIns: newMatchIns(),
	}
}

func (p *picker) Next(ctx context.Context, request interface{}) discovery.Instance {
	defer p.Recycle()
	req := request.(*transport.Transport)
	var ins []discovery.Instance
	num := 0
	// 服务路由规则
	if muxer, ok := p.ruleResolver.GetRules(req.GetAddr()); ok {
		// 复用这段内存空间, 不够再增加内存空间
		driLen := len(p.dr.Instances)
		if len(p.matchIns) < driLen {
			p.matchIns = make([]discovery.Instance, driLen)
		}
		ins = p.matchIns
		for k := range p.dr.Instances {
			match := muxer.Match(&tcp.Data{
				Req:      req,
				Instance: p.dr.Instances[k],
			})
			if match {
				ins[num] = p.dr.Instances[k]
				num = num + 1
			}
		}
	} else {
		ins = p.dr.Instances
		num = len(ins)
	}

	if num > 0 {
		idx := rand.Intn(num)
		return ins[idx]
	}
	return nil
}

type Balancer struct {
	ruleResolver RuleResolver
}

func NewBalancer(endpoints []string, opts ...kretry.Option) (*Balancer, error) {
	var (
		r  interface{}
		ok bool
	)
	if r, ok = cache.Get(RULE_CACHE_KEY); !ok {
		var err error
		r, err, _ = sfg.Do(RULE_CACHE_KEY, func() (interface{}, error) {
			nr, e := NewEtcdRuleResolver(endpoints, opts...)
			if e != nil {
				return nil, e
			}
			cache.Set(RULE_CACHE_KEY, nr)
			return nr, nil
		})

		if err != nil {
			return nil, err
		}
	}

	return &Balancer{
		ruleResolver: r.(RuleResolver),
	}, nil
}

// 获取服务发现的拾取器
func (b *Balancer) GetPicker(dr discovery.Result) loadbalance.Picker {
	// 需要池化，否则会有并发问题
	pk := pickerPool.Get().(*picker)
	pk.dr = &dr
	pk.ruleResolver = b.ruleResolver
	return pk
}

func (b *Balancer) Name() string {
	return "kbalancer"
}

func GetBlCallKey(call string) string {
	return fmt.Sprintf("blcall-%s", call)
}
