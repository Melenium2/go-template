package cache

import (
	"cmp"
	"sync"
)

const defaultKeysLimit = 500

type ExtendedConstraint interface {
	cmp.Ordered | ~[16]byte // uuid.UUID ([16]byte) type
}

type Cache[K ExtendedConstraint, V any] struct {
	limit int

	mutex sync.RWMutex
	c     map[K]V
	list  []K
	curr  int
}

func NewCache[K ExtendedConstraint, V any](keysLimit ...int) *Cache[K, V] {
	limit := defaultKeysLimit

	if len(keysLimit) > 0 && keysLimit[0] > 0 {
		limit = keysLimit[0]
	}

	return &Cache[K, V]{
		limit: limit,
		c:     make(map[K]V, limit+1),
		list:  make([]K, limit+1),
		curr:  0,
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.c[key] = value

	// not efficient way to purge map tail, but in this case
	// we don't care.
	k := c.list[c.curr]

	delete(c.c, k)

	c.list[c.curr] = key

	c.curr = (c.curr + 1) % c.limit
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()

	value, ok := c.c[key]

	c.mutex.RUnlock()

	return value, ok
}
