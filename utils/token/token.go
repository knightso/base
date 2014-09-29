package token

import (
	"bytes"
	"math"
	"strings"
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

func (c *config) Decode(token string) int64 {
	return c.applyMasks(c.toInt64(token))
}

// TODO:concern negative and zero
func (c *config) toBase62(v int64) string {
	var buf bytes.Buffer
	for v > 0 {
		buf.WriteByte(BASE62CHARS[v%62])
		v /= 62
	}
	return buf.String()
}

func (c *config) toInt64(token string) int64 {
	var value int64
	for i := len(token) - 1; i >= 0; i-- {
		value += int64(strings.IndexByte(BASE62CHARS, token[i])) * c.pow62(i)
	}
	return value
}

func (c *config) pow62(k int) int64 {
	value := int64(1)
	for i := 0; i < k; i++ {
		value *= 62
	}
	return value
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
