package ds

import (
	"appengine/datastore"
	"sync"
)

type SyncMap struct {
	sync.RWMutex
	M map[string]interface{}
}

func NewSyncMap() *SyncMap {
	return &SyncMap{M: make(map[string]interface{})}
}

func (s *SyncMap) Get(k *datastore.Key) (interface{}, bool) {
	s.RLock()
	r, ok := s.M[k.Encode()]
	s.RUnlock()
	return r, ok
}

func (s *SyncMap) Put(k *datastore.Key, v interface{}) {
	s.Lock()
	s.M[k.Encode()] = v
	s.Unlock()
}

func (s *SyncMap) Delete(k *datastore.Key) {
	s.RLock()
	delete(s.M, k.Encode())
	s.RUnlock()
}

func (s *SyncMap) ForEach(f func(k *datastore.Key, v interface{}) error) error {
	s.RLock()
	defer s.RUnlock()
	for k, v := range s.M {
		key, err := datastore.DecodeKey(k)
		if err != nil {
			return err
		}
		if err := f(key, v); err != nil {
			return err
		}
	}
	return nil
}
