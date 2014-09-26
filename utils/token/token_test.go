package token

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestNew(t *testing.T) {
	der := New([]int64{
		5577006791947779410,
		8674665223082153551,
		6129484611666145821,
		4037200794235010051,
		3916589616287113937,
		6334824724549167320,
		605394647632969758,
		1443635317331776148,
		894385949183117216,
		2775422040480279449,
		4751997750760398084,
		7504504064263669287,
		1976235410884491574,
		3510942875414458836,
		2933568871211445515,
		4324745483838182873,
	}, 2610529275472644968)

	assertToken(t, der, 1, "3DTo3UjhU29")
	assertToken(t, der, 2, "mmK8MTuTyV7")
	assertToken(t, der, 3, "JVYJWN0CLt1")
	assertToken(t, der, 4, "gGAzJCrdnI9")
	assertToken(t, der, 5, "tESoSWuf2b5")
	assertToken(t, der, 10000000, "ICgQmf4I6D9")
	assertToken(t, der, 20000000, "CjC3idfp6e4")
	assertToken(t, der, 30000000, "eriiRP54Ne1")
}

func assertToken(t *testing.T, der EncodeDecoder, n int64, token string) {
	encoded := der.Encode(n)
	if encoded != token {
		t.Errorf("Encode error. expected:%s, but was:%s", token, encoded)
	}
	decoded := der.Decode(encoded)
	if decoded != n {
		t.Errorf("Decode error. expected:%d, but was:%d", n, decoded)
	}
}

func _TestGenerateMasks(t *testing.T) {

	fmt.Println("=== masks")
	for i := 0; i < 16; i++ {
		fmt.Printf("%d\n", rand.Int63())
	}
	fmt.Println("=== zeroMask")
	fmt.Printf("%d\n", rand.Int63())
}
