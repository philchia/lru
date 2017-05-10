package lru

import (
	"container/list"
	"sync"
	"time"
)

// Cache represent a lru cache
type Cache interface {
	// Get value form cache with key, return nil if not exists
	Get(key interface{}) interface{}
	// Set value with key, and optional expires time
	Set(key, value interface{}, expires ...time.Time)
	// Del key from cache
	Del(key interface{})
}

type cache struct {
	lru   *list.List
	items map[interface{}]*list.Element
	size  int
}

type lockCache struct {
	sync.RWMutex
	cache
}

type item struct {
	k, v    interface{}
	expires time.Time
}

// New create a Cache instance
func New(size int) Cache {
	if size <= 0 {
		panic("lur: must provide a positive num")
	}

	return &cache{
		lru:   list.New(),
		items: make(map[interface{}]*list.Element),
		size:  size,
	}
}

func (c *cache) Get(key interface{}) interface{} {
	if v, ok := c.items[key]; ok {
		item := v.Value.(*item)
		if item.expires.IsZero() || item.expires.After(time.Now()) {
			return item.v
		}
		c.removeItem(v)
	}
	return nil
}

func (c *cache) Set(key, value interface{}, expires ...time.Time) {
	if ele, ok := c.items[key]; ok {
		c.lru.MoveToFront(ele)
		item := ele.Value.(*item)
		item.v = value
		if len(expires) > 0 {
			item.expires = expires[0]
		}
		return
	}

	item := &item{
		k: key,
		v: value,
	}
	if len(expires) > 0 {
		item.expires = expires[0]
	}

	c.items[key] = c.lru.PushFront(item)

	if len(c.items) > c.size {
		item := c.lru.Back()
		c.removeItem(item)
	}
}

func (c *cache) Del(key interface{}) {
	if ele, ok := c.items[key]; ok {
		c.removeItem(ele)
	}
}

func (c *cache) removeItem(ele *list.Element) {
	c.lru.Remove(ele)
	delete(c.items, ele.Value.(*item).k)
}

// NewLockCache create a thread safe Cache
func NewLockCache(size int) Cache {
	if size <= 0 {
		panic("lur: must provide a positive num")
	}

	return &lockCache{
		cache: cache{
			lru:   list.New(),
			items: make(map[interface{}]*list.Element),
			size:  size,
		},
	}
}

func (c *lockCache) Get(key interface{}) interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.cache.Get(key)
}

func (c *lockCache) Set(key, value interface{}, expires ...time.Time) {
	c.Lock()
	defer c.Unlock()
	c.cache.Set(key, value, expires...)
}

func (c *lockCache) Del(key interface{}) {
	c.Lock()
	defer c.Unlock()
	c.cache.Del(key)
}
