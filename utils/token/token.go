package token

import (
	"bytes"
	"math"
	"unicode/utf8"
)

const BASE62CHARS = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type EncodeDecoder interface {
	Encode(int64) string
	Decode(string) int64
}

type config struct {
	masks      []int64
	zeroMask   int64
	maskToMask int64
}

func (c *config) Encode(v int64) string {
	return c.toBase62(c.applyMasks(v))
}

func (c *config) Decode(s string) int64 {
	return 0
}

// TODO:concern negative and zero
func (c *config) toBase62(v int64) string {
	var buf bytes.Buffer
	for v > 0 {
		r, _ := utf8.DecodeRune([]byte(BASE62CHARS)[v%62:])
		buf.WriteRune(r)
		v /= 62
	}
	return buf.String()
}

func (c *config) applyMasks(v int64) int64 {
	v ^= c.zeroMask & c.maskToMask
	for i, m := range c.masks {
		if (v & (int64(1) << uint(i))) != 0 {
			v ^= m & c.maskToMask
		}
	}
	return v
}

func New(masks []int64, zeroMask int64) EncodeDecoder {
	c := &config{masks, zeroMask, 0}
	c.maskToMask = ^((int64(1) << uint(len(masks))) - 1) & math.MaxInt64
	return c
}
