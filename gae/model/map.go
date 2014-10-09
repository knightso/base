package model

import (
	"appengine/datastore"
	"sync"
)

type SyncMap struct {
	sync.RWMutex
	M map[datastore.Key]interface{}
}

func NewSyncMap() *SyncMap {
	return &SyncMap{M: make(map[datastore.Key]interface{})}
}

func (s *SyncMap) Get(k *datastore.Key) (interface{}, bool) {
	s.RLock()
	r, ok := s.M[*k]
	s.RUnlock()
	return r, ok
}

func (s *SyncMap) Put(k *datastore.Key, v interface{}) {
	s.Lock()
	s.M[*k] = v
	s.Unlock()
}

func (s *SyncMap) ForEach(f func(k *datastore.Key, v interface{})) {
	s.RLock()
	defer s.RUnlock()
	for k, v := range s.M {
		f(&k, v)
	}
}
