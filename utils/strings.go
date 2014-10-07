package utils

import (
	"bytes"
)

// TODO support hankaku-kana

const (
	hstart rune = 'ぁ'
	hend rune = 'ん'
	kstart rune = 'ァ'
	kend rune = 'ン'
	gap = kstart - hstart
)

func Hiragana2Katakana(h string) string {
	var buf bytes.Buffer
	for _, r := range h {
		if r >= hstart && r <= hend {
			buf.WriteRune(r + gap)
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func Katakana2Hiragana(k string) string {
	var buf bytes.Buffer
	for _, r := range k {
		if r >= kstart && r <= kend {
			buf.WriteRune(r - gap)
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}
