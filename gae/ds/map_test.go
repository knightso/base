package ds

import (
	"appengine/aetest"
	"appengine/datastore"
	"testing"
)

func Test(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	m := NewSyncMap()

	key1 := datastore.NewKey(c, "Test", "", 999, nil)
	m.Put(key1, "test-IntID")

	key2 := datastore.NewKey(c, "Test", "123", 0, nil)
	m.Put(key2, "test-StringID")

	key3 := datastore.NewKey(c, "Test", "456", 0, key1)
	m.Put(key3, "test-ancestor1")

	key4 := datastore.NewKey(c, "Test", "456", 0, key2)
	m.Put(key4, "test-ancestor2")

	var v interface{}
	var ok bool

	v, ok = m.Get(key1)
	assertTrue(t, ok)
	assert(t, "test-IntID", v)

	v, ok = m.Get(key2)
	assertTrue(t, ok)
	assert(t, "test-StringID", v)

	v, ok = m.Get(key3)
	assertTrue(t, ok)
	assert(t, "test-ancestor1", v)

	v, ok = m.Get(key4)
	assertTrue(t, ok)
	assert(t, "test-ancestor2", v)
}

func assertTrue(t *testing.T, b bool) {
	if !b {
		t.Errorf("expected true, but was false\n")
	}
}

func assert(t *testing.T, expected string, actual interface{}) {
	s, ok := actual.(string)
	if !ok {
		t.Errorf("expected actual is string, but was %T\n", expected, actual)
	}
	if expected != s {
		t.Errorf("expected %s, but was %s\n", expected, s)
	}
}
