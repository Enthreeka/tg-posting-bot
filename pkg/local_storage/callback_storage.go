package store

import "sync"

type CallbackStorage struct {
	storage map[string]struct{}

	mu sync.RWMutex
}

func NewCallbackStorage() *CallbackStorage {
	return &CallbackStorage{
		storage: make(map[string]struct{}, 25),
	}
}

func (c *CallbackStorage) GetStorage() map[string]struct{} {
	return c.storage
}

func (c *CallbackStorage) AppendStorage(name string) {
	c.mu.Lock()
	c.storage[name] = struct{}{}
	c.mu.Unlock()
}
