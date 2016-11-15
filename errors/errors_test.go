package errors

import (
	"errors"
	"testing"
)

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
