package storage

import (
	"errors"
	"log"
	"sync"
)

type Store struct {
	mu    sync.RWMutex
	kvLog (map[string]string)
}

var store = Store{kvLog: make(map[string]string)}

var ErrorKeyNotFound = errors.New("key not found")

func NewStore() error {
	return nil
}

func Put(key string, value string) error {
	store.mu.Lock()
	log.Printf("key %s, value %s\n", key, value)
	store.kvLog[key] = value
	store.mu.Unlock()
	return nil
}

func Get(key string) (string, error) {
	store.mu.RLock()
	value, ok := store.kvLog[key]
	store.mu.RUnlock()
	if !ok {
		return "", ErrorKeyNotFound
	}
	return value, nil
}

func Delete(key string) error {
	store.mu.Lock()
	delete(store.kvLog, key)
	store.mu.Unlock()
	return nil
}
