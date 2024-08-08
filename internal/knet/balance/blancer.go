package balance

import (
	"context"
	"fmt"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	transport "github.com/llsw/ikunet/internal/kitex_gen/transport"
	cmap "github.com/orcaman/concurrent-map"
)

var (
	_         loadbalance.Loadbalancer = &Balancer{}
	pic       loadbalance.Picker       = &picker{}
	cache     cmap.ConcurrentMap       = cmap.New()
	whiteList cmap.ConcurrentMap       = cmap.New()
)

const (
	TAG_CLUSTER      = "cluster"
	TAG_VERSION      = "version"
	TAG_ID           = "id"
	TAG_TYPE         = "type"
	TAG_MAINTAIN     = "maintain"
	TYPE_STATEFUL    = "1"
	TAG_MAINTAIN_ON  = "1"
	TAG_MAINTAIN_OFF = "0"
)

type picker struct {
	dr *discovery.Result
}

func getNewVer(uuid string, dr *discovery.Result) string {
	var isWl bool
	if mv, ok := whiteList.Get(dr.CacheKey); ok {
		wl := mv.(map[string]bool)
		if _, ok := wl[uuid]; ok {
			isWl = true
		}
	}
	var ver string
	if vt, ok := dr.Instances[0].Tag(TAG_VERSION); ok {
		ver = vt
	}

	for k := 1; k < len(dr.Instances); k++ {
		v := dr.Instances[k]
		if _, ok := v.Tag(TAG_MAINTAIN); ok {
			if !isWl {
				continue
			}
		}

		if vt, ok := v.Tag(TAG_VERSION); ok {
			if vt > ver {
				ver = vt
			}
		}
	}
	return ver
}

func isStateful(ins discovery.Instance) bool {
	if v, ok := ins.Tag(TAG_TYPE); ok {
		if v == TYPE_STATEFUL {
			return true
		}
	}
	return false
}

func getVer(uuid string, cmd string, dr *discovery.Result) (string, bool) {
	var (
		ver   string
		isNew bool
	)
	vck := fmt.Sprintf("%s_%s_ver", uuid, dr.CacheKey)
	if cv, ok := cache.Get(vck); !ok {
		ver = getNewVer(uuid, dr)
		isNew = true
	} else {
		ver = cv.(string)
	}

	if _, ok := dr.Instances[0].Tag(GetBlCallKey(cmd)); ok {
		newVer := getNewVer(uuid, dr)
		if ver != newVer {
			isNew = true
			cache.Set(vck, newVer)
			ver = newVer
		}
	}
	return ver, isNew
}

func (p *picker) Next(ctx context.Context, request interface{}) discovery.Instance {
	req := request.(*transport.Transport)
	cmd := req.GetCmd()
	uuid := req.GetMeta().GetUuid()
	ver, isNew := getVer(uuid, cmd, p.dr)

	isSt := isStateful(p.dr.Instances[0])
	// 不是新的又是有状态的服务，那就走原来的服务
	if !isNew && isSt {
		ick := fmt.Sprintf("%s_%s_stf", uuid, p.dr.CacheKey)
		if ins, ok := cache.Get(ick); ok {
			return ins.(discovery.Instance)
		}
	}

	for _, v := range p.dr.Instances {
		// TODO: 负载均衡
		if iv, ok := v.Tag(TAG_MAINTAIN); ok {
			if iv == ver {
				return v
			}
		}
	}

	return nil
}

type Balancer struct {
}

func (b *Balancer) GetPicker(dr discovery.Result) loadbalance.Picker {
	// 每次都要赋值，防止服务发现结果改变了
	pic.(*picker).dr = &dr
	return pic
}

func (b *Balancer) Name() string {
	return "ikunet_balancer"
}

func GetBlCallKey(call string) string {
	return fmt.Sprintf("blcall-%s", call)
}

func GetTagVal(ins *discovery.Instance, tag string) string {
	if val, ok := (*ins).Tag(tag); ok {
		return val
	}
	return ""
}
