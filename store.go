package main

import (
	"errors"
	"sync"
)

// Storage to initialize our stateless storage as a global variable.
var Storage Store = newMemoryStore()

// Store is the main interface of storage in this app.
type Store interface {
	// Put to store data, will return error if the UniqueID exists.
	Put(uniqueID UniqueID) error
}

// MemoryStore is a dummy in memory storage. It can be replaced with any storage.
// For example if we want quorom in consensus system we can consider other options.
type MemoryStore struct {
	// lock for kv.
	lock sync.Mutex
	kv   map[UniqueID]struct{}
}

func newMemoryStore() Store {
	return &MemoryStore{kv: make(map[UniqueID]struct{})}
}

func (m *MemoryStore) Put(uniqueID UniqueID) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, exist := m.kv[uniqueID]; exist {
		return errors.New("duplicated data")
	}
	// Save dummy data.
	m.kv[uniqueID] = struct{}{}
	return nil
}
