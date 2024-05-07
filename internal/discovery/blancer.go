package discovery

import (
	"context"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	transport "github.com/llsw/ikunet/internal/kitex_gen/transport"
)

var (
	_ loadbalance.Loadbalancer = &balancer{}
	_ loadbalance.Picker       = &picker{}
)

type picker struct {
	dr *discovery.Result
}

func (p *picker) Next(ctx context.Context, request interface{}) discovery.Instance {
	req := request.(*transport.Transport)
	cmd := req.GetCmd()
	if chcmd, ok := p.dr.Instances[0].Tag("chcmd"); ok && chcmd == cmd {
		return p.dr.Instances[0]
	}
	// TODO: 负载均衡
	return nil
}

type balancer struct {
}

func (b *balancer) GetPicker(dr discovery.Result) loadbalance.Picker {
	return &picker{
		dr: &dr,
	}
}

func (b *balancer) Name() string {
	return "ikunet_balancer"
}
