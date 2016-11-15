package errors

import (
	"errors"
	"testing"
)

func TestErrorf(t *testing.T) {
	err := Errorf("test %d %s", 123, "hoge")

	expected := "test 123 hoge"
	if err.Error() != expected {
		t.Fatalf("expected: %s, but was: %s", expected, err.Error())
	}
}

func TestSyncMultiError(t *testing.T) {
	var sme SyncMultiError
	sme.Append(errors.New("aaa"))
	sme.Append(errors.New("bbb"))
	sme.Append(errors.New("ccc"))

	expected := "aaa\nbbb\nccc\n"
	if sme.Error() != expected {
		t.Fatalf("sme.Error() expected:%s, but was:%s", expected, sme.Error())
	}

	expectedLen := 3
	if sme.Len() != expectedLen {
		t.Fatalf("sme.Len() expected:%d, but was:%d", expectedLen, sme.Len())
	}
}
