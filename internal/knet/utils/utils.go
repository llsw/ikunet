package utils

import (
	"github.com/cloudwego/kitex/pkg/discovery"
)

func GetTagVal(ins discovery.Instance, tag string) string {
	if val, ok := ins.Tag(tag); ok {
		return val
	}
	return ""
}
