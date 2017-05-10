package lru

import (
	"container/list"
	"time"
)

// Cache represent a lru cache
type Cache interface {
	Get(k interface{}) interface{}
	Set(k, v interface{}, expires ...time.Time)
	Del(k interface{})
}

type cache struct {
	lru   *list.List
	items map[interface{}]*list.Element
	size  int
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

func (c *cache) Get(k interface{}) interface{} {
	if v, ok := c.items[k]; ok {
		item := v.Value.(*item)
		if item.expires.IsZero() || item.expires.After(time.Now()) {
			return item.v
		}
		c.removeItem(v)
	}
	return nil
}

func (c *cache) Set(k, v interface{}, expires ...time.Time) {
	if ele, ok := c.items[k]; ok {
		c.lru.MoveToFront(ele)
		item := ele.Value.(*item)
		item.v = v
		if len(expires) > 0 {
			item.expires = expires[0]
		}
		return
	}

	item := &item{
		k: k,
		v: v,
	}
	if len(expires) > 0 {
		item.expires = expires[0]
	}

	c.items[k] = c.lru.PushFront(item)

	if len(c.items) > c.size {
		item := c.lru.Back()
		c.removeItem(item)
	}
}

func (c *cache) Del(k interface{}) {
	if ele, ok := c.items[k]; ok {
		c.removeItem(ele)
	}
}

func (c *cache) removeItem(ele *list.Element) {
	c.lru.Remove(ele)
	delete(c.items, ele.Value.(*item).k)
}
