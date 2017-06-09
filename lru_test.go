package lru

import (
	"container/list"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		size int
	}

	tests := []struct {
		name      string
		args      args
		wantNil   bool
		wantPanic bool
	}{
		{
			name: "case1",
			args: args{
				size: 10,
			},
			wantNil:   false,
			wantPanic: false,
		},
		{
			name: "case2",
			args: args{
				size: 0,
			},
			wantNil:   true,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			func() {
				defer func() {
					if err := recover(); err != nil {
						if !tt.wantPanic {
							t.Errorf("New() want panic %v", tt.wantPanic)
						}
					}
				}()
				if got := New(tt.args.size); (got == nil) != tt.wantNil {
					t.Errorf("New() = %v, want nil %v", got, tt.wantNil)
				}
			}()

		})
	}
}

func Test_cache_Get(t *testing.T) {
	type fields struct {
		lru   *list.List
		items map[interface{}]*list.Element
		size  int
	}
	type args struct {
		k interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(*cache)
		want   interface{}
	}{
		{
			name: "case1",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  10,
			},
			args: args{
				k: "key",
			},
			setup: func(c *cache) {},
			want:  nil,
		},
		{
			name: "case2",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  10,
			},
			args: args{
				k: "key",
			},
			setup: func(c *cache) {
				c.Set("key", "value")
			},
			want: "value",
		},
		{
			name: "case3",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  10,
			},
			args: args{
				k: "key",
			},
			setup: func(c *cache) {
				c.Set("key", "value", time.Now())
				time.Sleep(time.Millisecond * 10)
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				lru:   tt.fields.lru,
				items: tt.fields.items,
				size:  tt.fields.size,
			}
			tt.setup(c)
			if got := c.Get(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cache.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_Set(t *testing.T) {
	type fields struct {
		lru   *list.List
		items map[interface{}]*list.Element
		size  int
	}
	type args struct {
		k       interface{}
		v       interface{}
		expires []time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(*cache)
		wantF  func(*cache) bool
	}{
		{
			name: "case1",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  10,
			},
			args: args{
				k: "key",
				v: "value",
			},
			setup: func(c *cache) {},
			wantF: func(c *cache) bool {
				return len(c.items) == 1
			},
		},
		{
			name: "case2",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  1,
			},
			args: args{
				k: "key",
				v: "value",
			},
			setup: func(c *cache) {
				c.Set("k", "v")
			},
			wantF: func(c *cache) bool {
				return len(c.items) == 1
			},
		},
		{
			name: "case3",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  2,
			},
			args: args{
				k: "key",
				v: "value",
			},
			setup: func(c *cache) {
				c.Set("key", "val")
				c.Set("k", "v")
			},
			wantF: func(c *cache) bool {
				return c.lru.Front().Value.(*item).k.(string) == "key"
			},
		},
		{
			name: "case4",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  2,
			},
			args: args{
				k:       "key",
				v:       "value",
				expires: []time.Time{time.Now()},
			},
			setup: func(c *cache) {},
			wantF: func(c *cache) bool {
				return !c.lru.Front().Value.(*item).expires.IsZero()
			},
		},
		{
			name: "case5",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  2,
			},
			args: args{
				k:       "key",
				v:       "value",
				expires: []time.Time{time.Now()},
			},
			setup: func(c *cache) {
				c.Set("key", "value")
			},
			wantF: func(c *cache) bool {
				return !c.lru.Front().Value.(*item).expires.IsZero()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				lru:   tt.fields.lru,
				items: tt.fields.items,
				size:  tt.fields.size,
			}
			tt.setup(c)
			c.Set(tt.args.k, tt.args.v, tt.args.expires...)
			if !tt.wantF(c) {
				t.Fail()
			}
		})
	}
}

func Test_cache_Del(t *testing.T) {
	type fields struct {
		lru   *list.List
		items map[interface{}]*list.Element
		size  int
	}
	type args struct {
		k interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(*cache)
		wantF  func(*cache) bool
	}{
		{
			name: "case1",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  10,
			},
			args: args{
				k: "key",
			},
			setup: func(c *cache) {},
			wantF: func(c *cache) bool {
				return len(c.items) == 0
			},
		},
		{
			name: "case2",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  10,
			},
			args: args{
				k: "key",
			},
			setup: func(c *cache) {
				c.Set("key", "value")
			},
			wantF: func(c *cache) bool {
				return len(c.items) == 0
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				lru:   tt.fields.lru,
				items: tt.fields.items,
				size:  tt.fields.size,
			}
			tt.setup(c)
			c.Del(tt.args.k)
			if !tt.wantF(c) {
				t.Fail()
			}
		})
	}
}

func TestNewLockCache(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name      string
		args      args
		wantNil   bool
		wantPanic bool
	}{
		{
			name: "case1",
			args: args{
				size: 10,
			},
			wantNil:   false,
			wantPanic: false,
		},
		{
			name: "case2",
			args: args{
				size: 0,
			},
			wantNil:   true,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			func() {
				defer func() {
					if err := recover(); err != nil {
						if !tt.wantPanic {
							t.Errorf("New() want panic %v", tt.wantPanic)
						}
					}
				}()
				if got := NewLockCache(tt.args.size); (got == nil) != tt.wantNil {
					t.Errorf("New() = %v, want nil %v", got, tt.wantNil)
				}
			}()

		})
	}
}

func Test_lockCache_Get(t *testing.T) {
	type fields struct {
		lru   *list.List
		items map[interface{}]*list.Element
		size  int
	}
	type args struct {
		k interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(Cache)
		want   interface{}
	}{
		{
			name: "case1",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  10,
			},
			args: args{
				k: "key",
			},
			setup: func(c Cache) {},
			want:  nil,
		},
		{
			name: "case2",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  10,
			},
			args: args{
				k: "key",
			},
			setup: func(c Cache) {
				c.Set("key", "value")
			},
			want: "value",
		},
		{
			name: "case3",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  10,
			},
			args: args{
				k: "key",
			},
			setup: func(c Cache) {
				c.Set("key", "value", time.Now())
				time.Sleep(time.Millisecond * 10)
			},
			want: nil,
		},
	}
	for _, tt := range tests {

		c := &lockCache{
			cache: cache{
				lru:   tt.fields.lru,
				items: tt.fields.items,
				size:  tt.fields.size,
			},
		}
		tt.setup(c)
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Get(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cache.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lockCache_Set(t *testing.T) {
	type fields struct {
		lru   *list.List
		items map[interface{}]*list.Element
		size  int
	}
	type args struct {
		k       interface{}
		v       interface{}
		expires []time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(Cache)
		wantF  func(*lockCache) bool
	}{
		{
			name: "case1",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  10,
			},
			args: args{
				k: "key",
				v: "value",
			},
			setup: func(c Cache) {},
			wantF: func(c *lockCache) bool {
				return len(c.items) == 1
			},
		},
		{
			name: "case2",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  1,
			},
			args: args{
				k: "key",
				v: "value",
			},
			setup: func(c Cache) {
				c.Set("k", "v")
			},
			wantF: func(c *lockCache) bool {
				return len(c.items) == 1
			},
		},
		{
			name: "case3",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  2,
			},
			args: args{
				k: "key",
				v: "value",
			},
			setup: func(c Cache) {
				c.Set("key", "val")
				c.Set("k", "v")
			},
			wantF: func(c *lockCache) bool {
				return c.lru.Front().Value.(*item).k.(string) == "key"
			},
		},
		{
			name: "case4",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  2,
			},
			args: args{
				k:       "key",
				v:       "value",
				expires: []time.Time{time.Now()},
			},
			setup: func(c Cache) {},
			wantF: func(c *lockCache) bool {
				return !c.lru.Front().Value.(*item).expires.IsZero()
			},
		},
		{
			name: "case5",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  2,
			},
			args: args{
				k:       "key",
				v:       "value",
				expires: []time.Time{time.Now()},
			},
			setup: func(c Cache) {
				c.Set("key", "value")
			},
			wantF: func(c *lockCache) bool {
				return !c.lru.Front().Value.(*item).expires.IsZero()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &lockCache{
				cache: cache{
					lru:   tt.fields.lru,
					items: tt.fields.items,
					size:  tt.fields.size,
				},
			}
			tt.setup(c)
			c.Set(tt.args.k, tt.args.v, tt.args.expires...)
			if !tt.wantF(c) {
				t.Fail()
			}
		})
	}
}

func Test_lockCache_Del(t *testing.T) {
	type fields struct {
		lru   *list.List
		items map[interface{}]*list.Element
		size  int
	}
	type args struct {
		k interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(Cache)
		wantF  func(*lockCache) bool
	}{
		{
			name: "case1",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  10,
			},
			args: args{
				k: "key",
			},
			setup: func(c Cache) {},
			wantF: func(c *lockCache) bool {
				return len(c.items) == 0
			},
		},
		{
			name: "case2",
			fields: fields{
				lru:   list.New(),
				items: make(map[interface{}]*list.Element),
				size:  10,
			},
			args: args{
				k: "key",
			},
			setup: func(c Cache) {
				c.Set("key", "value")
			},
			wantF: func(c *lockCache) bool {
				return len(c.items) == 0
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &lockCache{
				cache: cache{
					lru:   tt.fields.lru,
					items: tt.fields.items,
					size:  tt.fields.size,
				},
			}
			tt.setup(c)
			c.Del(tt.args.k)
			if !tt.wantF(c) {
				t.Fail()
			}
		})
	}
}
