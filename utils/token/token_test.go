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

	c, _ := der.(*config)
	fmt.Printf("%x\n", uint64(c.maskToMask))

	for i := 1; i < 100; i++ {
		fmt.Printf("%d:%d\n", i, der.Encode(int64(i)))
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
