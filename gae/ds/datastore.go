package ds

import (
	"context"
	"crypto/md5"
	goerrors "errors"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/knightso/base/errors"
	"github.com/qedus/nds"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var OptimisticLockError = goerrors.New("Optimistic lock failure.")
var UniqueConstraintError = goerrors.New("Already exists.")

var DefaultCache = false
var CacheKinds = make(map[string]bool)

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
	Version   int            `datastore:",noindex" json:"version"`
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
	if m.CreatedAt.IsZero() {
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

func Put(c context.Context, e HasKey) error {
	if ht, ok := e.(HasTime); ok {
		ht.SetMetaTime()
	}
	if hv, ok := e.(HasVersion); ok {
		hv.IncrementVersion()
	}

	f := getPutFunc(e.GetKey().Kind())
	key, err := f(c, e.GetKey(), e)
	if err != nil {
		return errors.WrapOr(err)
	}

	e.SetKey(key)

	return nil
}

func Get(c context.Context, key *datastore.Key, dst interface{}) error {

	f := getGetFunc(key.Kind())

	if err := f(c, key, dst); err != nil {
		return errors.WrapOr(err)
	}

	if hk, ok := dst.(HasKey); ok {
		hk.SetKey(key)
	}

	return nil
}

func GetWithVersion(c context.Context, key *datastore.Key, version int, dst interface{}) error {

	f := getGetFunc(key.Kind())

	if err := f(c, key, dst); err != nil {
		return errors.WrapOr(err)
	}

	if hv, ok := dst.(HasVersion); ok {
		if hv.GetVersion() != version {
			return errors.WrapOr(OptimisticLockError)
		}
	}

	if hk, ok := dst.(HasKey); ok {
		hk.SetKey(key)
	}

	return nil
}

// keys only query is not supported(use datastore package)
// TODO:check keysonly
func ExecuteQuery(c context.Context, q *datastore.Query, dst interface{}) error {

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

func GetMulti(c context.Context, keys []*datastore.Key, dst interface{}) error {

	f := datastore.GetMulti
	if len(keys) > 0 {
		f = getGetMultiFunc(keys[0].Kind())
	}

	err := f(c, keys, dst)
	if err != nil {
		if _, ok := err.(appengine.MultiError); !ok {
			return errors.WrapOr(err)
		}
	}

	forEach(dst, func(index int, elem interface{}) error {
		if hk, ok := elem.(HasKey); ok {
			hk.SetKey(keys[index])
		}
		return nil
	})

	if err != nil {
		return errors.WrapOr(err)
	}

	return nil
}

func GetMultiWithVersion(c context.Context, keys []*datastore.Key, versions []int, dst interface{}) error {

	f := datastore.GetMulti
	if len(keys) > 0 {
		f = getGetMultiFunc(keys[0].Kind())
	}

	err := f(c, keys, dst)
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
				return errors.WrapOr(OptimisticLockError)
			}
		}
		return nil
	})

	if err != nil {
		return errors.WrapOr(err)
	}

	return nil
}

// TODO: consider interface (how do you treat key?)
func PutMulti(c context.Context, keys []*datastore.Key, dst interface{}) error {

	f := datastore.PutMulti
	if len(keys) > 0 {
		f = getPutMultiFunc(keys[0].Kind())
	}

	forEach(dst, func(index int, elem interface{}) error {
		if ht, ok := elem.(HasTime); ok {
			ht.SetMetaTime()
		}
		if hv, ok := elem.(HasVersion); ok {
			hv.IncrementVersion()
		}
		return nil
	})

	keysput, err := f(c, keys, dst)
	if err != nil {
		if _, ok := err.(appengine.MultiError); !ok {
			return errors.WrapOr(err)
		}
	}

	forEach(dst, func(index int, elem interface{}) error {
		if hk, ok := elem.(HasKey); ok {
			hk.SetKey(keysput[index])
		}
		return nil
	})

	if err != nil {
		return errors.WrapOr(err)
	}

	return nil
}

func Delete(c context.Context, key *datastore.Key) error {
	f := getDeleteFunc(key.Kind())
	return f(c, key)
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

func FindIndexedKey(c context.Context, kind, id string) (*datastore.Key, error) {

	f := getGetFunc(kind)

	idxKey := datastore.NewKey(c, kind+INDEX_PREFIX, id, 0, nil)

	var index GenericIndex
	if err := f(c, idxKey, &index); err != nil {
		return nil, err
	}

	return index.Key, nil
}

// call in Tx
func PutKeyToIndex(c context.Context, key *datastore.Key, id, oldId string) error {
	kind := key.Kind()
	getF := getGetFunc(kind)
	putF := getPutFunc(kind)
	deleteF := getDeleteFunc(kind)

	if id == oldId {
		// do nothing
		return nil
	}

	idxKey := datastore.NewKey(c, key.Kind()+INDEX_PREFIX, id, 0, nil)

	var index GenericIndex
	err := getF(c, idxKey, &index)
	if err == nil {
		return errors.WrapOr(UniqueConstraintError)
	} else if err != datastore.ErrNoSuchEntity {
		return errors.WrapOr(err)
	}

	idx := GenericIndex{Key: key}
	if _, err := putF(c, idxKey, &idx); err != nil {
		return errors.WrapOr(err)
	}

	if oldId != "" {
		oldIdxKey := datastore.NewKey(c, key.Kind()+INDEX_PREFIX, oldId, 0, nil)
		var oldIndex GenericIndex
		if err := getF(c, oldIdxKey, &oldIndex); err != nil {
			return errors.WrapOr(err)
		} else {
			//remove old index.
			if err := deleteF(c, oldIdxKey); err != nil {
				return errors.WrapOr(err)
			}
		}
		if !oldIndex.Key.Equal(key) {
			return errors.New("data conflict error")
		}
	}

	return nil
}

func RemoveKeyFromIndex(c context.Context, key *datastore.Key, id string) error {
	deleteF := getDeleteFunc(key.Kind())

	idxKey := datastore.NewKey(c, key.Kind()+INDEX_PREFIX, id, 0, nil)
	if err := deleteF(c, idxKey); err != nil {
		return errors.WrapOr(err)
	}
	return nil
}

// to avoid tab split
func AddHashPrefix(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x-%s", h.Sum(nil)[:3], s)
}

func needCache(kind string) bool {
	nc, ok := CacheKinds[kind]
	if ok {
		return nc
	} else {
		return DefaultCache
	}
}

func getGetFunc(kind string) func(context.Context, *datastore.Key, interface{}) error {
	if needCache(kind) {
		return nds.Get
	} else {
		return datastore.Get
	}
}

func getGetMultiFunc(kind string) func(context.Context, []*datastore.Key, interface{}) error {
	if needCache(kind) {
		return nds.GetMulti
	} else {
		return datastore.GetMulti
	}
}

func getPutFunc(kind string) func(context.Context, *datastore.Key, interface{}) (*datastore.Key, error) {
	if needCache(kind) {
		return nds.Put
	} else {
		return datastore.Put
	}
}

func getPutMultiFunc(kind string) func(context.Context, []*datastore.Key, interface{}) ([]*datastore.Key, error) {
	if needCache(kind) {
		return nds.PutMulti
	} else {
		return datastore.PutMulti
	}
}

func getDeleteFunc(kind string) func(context.Context, *datastore.Key) error {
	if needCache(kind) {
		return nds.Delete
	} else {
		return datastore.Delete
	}
}

type Sequence struct {
	LastID int64
}

func GenerateID(c context.Context, kind string) (int64, error) {
	key := datastore.NewKey(c, "Sequence", kind, 0, nil)
	var sequence Sequence
	var err error
	err = datastore.RunInTransaction(c, func(c context.Context) error {
		err = datastore.Get(c, key, &sequence)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		sequence.LastID++
		_, err = datastore.Put(c, key, &sequence)
		if err != nil {
			return err
		}
		return nil
	}, nil)
	if err != nil {
		return 0, err
	}

	return sequence.LastID, nil
}

// RunInTransaction wraps nds's RunInTransaction
func RunInTransaction(c context.Context, f func(tc context.Context) error,
	opts *datastore.TransactionOptions) error {

	return nds.RunInTransaction(c, f, opts)
}

/* under construction
func FilterStartsWith(q *datastore.Query, filterStr string, value string) *datastore.Query {
	return q.Filter(filterStr + " >=", value).Filter(filterStr + " <", value + "\ufffd")
}
*/
