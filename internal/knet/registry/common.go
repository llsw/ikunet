package registry

import "fmt"

const (
	etcdPrefixTpl = "kitex/registry-etcd/%v/"
)

func serviceKeyPrefix(serviceName string) string {
	return fmt.Sprintf(etcdPrefixTpl, serviceName)
}

// serviceKey generates the key used to stored in etcd.
func serviceKey(serviceName, addr string) string {
	return serviceKeyPrefix(serviceName) + addr
}

// instanceInfo used to stored service basic info in etcd.
type instanceInfo struct {
	Network string            `json:"network"`
	Address string            `json:"address"`
	Weight  int               `json:"weight"`
	Tags    map[string]string `json:"tags"`
}
