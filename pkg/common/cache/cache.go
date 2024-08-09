package cache

import (
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	fc "github.com/patrickmn/go-cache"
)

type FromSource[T any] func(key string) (val T, err error)

type Cahce[T any] interface {
	Set(key string, val T, expire time.Duration)
	Get(key string, fromSource FromSource[T], expire time.Duration) (val T, ok bool)
	Del(key string)
}

var (
	_ Cahce[any] = &cache[any]{}
)

type cache[T any] struct {
	fc *fc.Cache
}

func New[T any](defaultExpiration, cleanupInterval time.Duration) Cahce[T] {
	c :=  &cache[T]{
		fc: fc.New(defaultExpiration, cleanupInterval),
	}
	return c
}

// Del implements Cahce.
func (c *cache[T]) Del(key string) {
	c.fc.Delete(key)
}

// Get implements Cahce.
func (c *cache[T]) Get(key string, fromSource FromSource[T], expire time.Duration) (val T, ok bool) {
	if v, o := c.fc.Get(key); o {
		if tv, o := v.(T); o {
			val = tv
			ok = o
		}
	} else {
		if v, err := fromSource(key); err == nil {
			val = v
			c.fc.Set(key, val, expire)
			ok = true
		} else {
			hlog.Errorf("get key:%s from source error:%s", key ,err.Error())
		}
	}
	return
}

// Set implements Cahce.
func (c *cache[T]) Set(key string, val T, expire time.Duration) {
	c.fc.Set(key, val, expire)
}
