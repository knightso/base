package model

import (
	"appengine"
	"appengine/datastore"
	"base/errors"
	"github.com/qedus/nds"
	"reflect"
	"time"
)

type OptimisticLockError struct {
	*errors.BaseError
}

type AlreadyIndexedError struct {
	*errors.BaseError
}

type HasKey interface {
	GetKey() *datastore.Key
	SetKey(*datastore.Key)
}

type HasTime interface {
	SetMetaTime()
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

type HasVersion interface {
	GetVersion() int
	IncrementVersion()
}

type Meta struct {
	Key       *datastore.Key `datastore:"-" json:"-"`
	Version   int            `json:"version"`
	Deleted   bool           `json:"deleted"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

func (m *Meta) GetKey() *datastore.Key {
	return m.Key
}

func (m *Meta) SetKey(key *datastore.Key) {
	m.Key = key
}

func (m *Meta) SetMetaTime() {
	if m.CreatedAt == (time.Time{}) {
		m.CreatedAt = time.Now()
	}
	m.UpdatedAt = time.Now()
}

func (m *Meta) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *Meta) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

func (m *Meta) GetVersion() int {
	return m.Version
}

func (m *Meta) IncrementVersion() {
	m.Version++
}

func Put(c appengine.Context, e HasKey) error {
	if ht, ok := e.(HasTime); ok {
		ht.SetMetaTime()
	}
	if hv, ok := e.(HasVersion); ok {
		hv.IncrementVersion()
	}

	key, err := nds.Put(c, e.GetKey(), e)
	if err != nil {
		return errors.WrapOr(err)
	}

	e.SetKey(key)

	return nil
}

func Get(c appengine.Context, key *datastore.Key, dst interface{}) error {

	if err := nds.Get(c, key, dst); err != nil {
		return errors.WrapOr(err)
	}

	if hk, ok := dst.(HasKey); ok {
		hk.SetKey(key)
	}

	return nil
}

func GetWithVersion(c appengine.Context, key *datastore.Key, version int, dst interface{}) error {

	if err := nds.Get(c, key, dst); err != nil {
		return errors.WrapOr(err)
	}

	if hv, ok := dst.(HasVersion); ok {
		if hv.GetVersion() != version {
			return OptimisticLockError{
				errors.New("optimistic lock failure. somebody may have updated it."),
			}
		}
	}

	if hk, ok := dst.(HasKey); ok {
		hk.SetKey(key)
	}

	return nil
}

func ExecuteQuery(c appengine.Context, q *datastore.Query, dst interface{}) error {

	keys, err := q.GetAll(c, dst)
	if err != nil {
		return errors.WrapOr(err)
	}

	return forEach(dst, func(index int, elem interface{}) error {
		if hk, ok := elem.(HasKey); ok {
			hk.SetKey(keys[index])
		}
		return nil
	})
}

func GetMulti(c appengine.Context, keys []*datastore.Key, dst interface{}) error {

	err := nds.GetMulti(c, keys, dst)
	if err != nil {
		if _, ok := err.(appengine.MultiError); !ok {
			return errors.WrapOr(err)
		}
	}

	return forEach(dst, func(index int, elem interface{}) error {
		if hk, ok := elem.(HasKey); ok {
			hk.SetKey(keys[index])
		}
		return nil
	})
}

func GetMultiWithVersion(c appengine.Context, keys []*datastore.Key, versions []int, dst interface{}) error {

	err := nds.GetMulti(c, keys, dst)
	if err != nil {
		if _, ok := err.(appengine.MultiError); !ok {
			return errors.WrapOr(err)
		}
	}

	return forEach(dst, func(index int, elem interface{}) error {
		if hk, ok := elem.(HasKey); ok {
			hk.SetKey(keys[index])
		}
		if hv, ok := elem.(HasVersion); ok {
			if hv.GetVersion() != versions[index] {
				return OptimisticLockError{
					errors.New("optimistic lock failure. somebody may have updated it."),
				}
			}
		}
		return nil
	})
}

// TODO: consider interface (how do you treat key?)
func PutMulti(c appengine.Context, keys []*datastore.Key, dst interface{}) error {

	forEach(dst, func(index int, elem interface{}) error {
		if ht, ok := elem.(HasTime); ok {
			ht.SetMetaTime()
		}
		if hv, ok := elem.(HasVersion); ok {
			hv.IncrementVersion()
		}
		return nil
	})

	keysput, err := nds.PutMulti(c, keys, dst)
	if err != nil {
		if _, ok := err.(appengine.MultiError); !ok {
			return errors.WrapOr(err)
		}
	}

	return forEach(dst, func(index int, elem interface{}) error {
		if hk, ok := elem.(HasKey); ok {
			hk.SetKey(keysput[index])
		}
		return nil
	})
}

func forEach(list interface{}, f func(index int, elem interface{}) error) error {

	v := reflect.ValueOf(list)
	v = reflect.Indirect(v)
	elemKind := v.Type().Elem().Kind()

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)

		if elemKind == reflect.Struct {
			elem = elem.Addr()
		} else if elem.IsNil() {
			continue
		}

		ei := elem.Interface()

		if err := f(i, ei); err != nil {
			return errors.WrapOr(err)
		}
	}
	return nil
}

type GenericIndex struct {
	Key *datastore.Key
}

const INDEX_PREFIX = "Index"

func FindIndexedKey(c appengine.Context, kind, id string) (*datastore.Key, error) {
	idxKey := datastore.NewKey(c, kind+INDEX_PREFIX, id, 0, nil)

	var index GenericIndex
	if err := nds.Get(c, idxKey, &index); err != nil {
		return nil, err
	}

	return index.Key, nil
}

// call in Tx
func PutKeyToIndex(c appengine.Context, key *datastore.Key, id, oldId string) error {
	if id == oldId {
		// do nothing
		return nil
	}

	idxKey := datastore.NewKey(c, key.Kind()+INDEX_PREFIX, id, 0, nil)

	var index GenericIndex
	err := nds.Get(c, idxKey, &index)
	if err == nil {
		return AlreadyIndexedError{
			errors.New("ID already exists."),
		}
	} else if err != datastore.ErrNoSuchEntity {
		return errors.WrapOr(err)
	}

	idx := GenericIndex{Key: key}
	if _, err := nds.Put(c, idxKey, &idx); err != nil {
		return errors.WrapOr(err)
	}

	if oldId != "" {
		oldIdxKey := datastore.NewKey(c, key.Kind()+INDEX_PREFIX, oldId, 0, nil)
		var oldIndex GenericIndex
		if err := nds.Get(c, oldIdxKey, &oldIndex); err != nil {
			return errors.WrapOr(err)
		}
		if !oldIndex.Key.Equal(key) {
			return errors.New("data conflict error")
		}
	}

	return nil
}
