package utils

import (
	"testing"
)

func TestHaversine(t *testing.T) {
	v := Haversine(35.6544, 139.74477, 21.4225, 39.8261)
	expected := 9491.280549
	hundredth := expected / 100
	if v < expected-hundredth || v > expected+hundredth {
		t.Errorf("expected: %f < value < %f, but in fact: value = %f", expected-hundredth, expected+hundredth, v)
	}
}
