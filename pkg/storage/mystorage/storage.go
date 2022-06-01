package mystorage

import (
	"sync"
	"time"
)

const (
	deleteExpiredDataInterval = 5 * time.Second
)

type Storage interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

type container struct {
	data      interface{}
	createdAt time.Time
	ttl       time.Duration
}

type dataStore struct {
	data map[string]container
	mux  sync.RWMutex
}

type myStorage struct {
	dataStore dataStore
}

func NewMyStorage() Storage {
	data := make(map[string]container)

	s := &myStorage{
		dataStore: dataStore{
			data: data,
			mux:  sync.RWMutex{},
		},
	}

	go s.clearCache(time.NewTicker(deleteExpiredDataInterval))

	return s
}

func (s *myStorage) Set(key string, value interface{}, ttl time.Duration) {
	s.dataStore.mux.Lock()
	defer s.dataStore.mux.Unlock()

	cont := container{
		data:      value,
		createdAt: time.Now(),
		ttl:       ttl,
	}
	s.dataStore.data[key] = cont
}

func (s *myStorage) Get(key string) (interface{}, bool) {
	s.dataStore.mux.RLock()
	defer s.dataStore.mux.RUnlock()

	cont, ok := s.dataStore.data[key]
	if !ok {
		return nil, ok
	}
	if cont.ttl != 0 && time.Now().Before(cont.createdAt.Add(cont.ttl)) {
		return nil, ok
	}
	return cont.data, ok
}

func (s *myStorage) Delete(key string) {
	s.dataStore.mux.Lock()
	defer s.dataStore.mux.Unlock()

	delete(s.dataStore.data, key)
}

func (s *myStorage) clearCache(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			expiredKeys := s.getExpiredDataKeys()
			s.deleteDataByKeys(expiredKeys)
		}
	}
}

func (s *myStorage) getExpiredDataKeys() []string {
	s.dataStore.mux.RLock()
	defer s.dataStore.mux.RUnlock()

	expiredKeys := make([]string, 0)
	for key, cont := range s.dataStore.data {
		if cont.ttl != 0 && time.Now().Before(cont.createdAt.Add(cont.ttl)) {
			expiredKeys = append(expiredKeys, key)
		}
	}
	return expiredKeys
}

func (s *myStorage) deleteDataByKeys(keys []string) {
	if len(keys) == 0 {
		return
	}
	s.dataStore.mux.Lock()
	defer s.dataStore.mux.Unlock()

	for i := 0; i < len(keys); i++ {
		delete(s.dataStore.data, keys[i])
	}
}
