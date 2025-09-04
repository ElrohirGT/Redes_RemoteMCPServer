package lib

import (
	"sync"
	"time"
)

type Cache[T any] struct {
	mutex   *sync.Mutex
	data    *T
	expires time.Time
}

func NewCache[T any](data *T, expires time.Time) *Cache[T] {
	return &Cache[T]{
		data:    data,
		expires: expires,
		mutex:   &sync.Mutex{},
	}
}

// Returns the data if valid, nil otherwise
func (c *Cache[T]) GetData() (*T, bool) {
	if time.Now().After(c.expires) {
		return nil, false
	} else {
		return c.data, true
	}
}

func (c *Cache[T]) Update(data *T, expires time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = data
	c.expires = expires
}
