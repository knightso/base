package utils

import (
	"testing"
)

func TestHiraganaToKatakana(t *testing.T) {
	assert(t, "ゴーファーヲabcded#$%", Hiragana2Katakana("ごーふぁーをabcded#$%"))
	assert(t, "ごーふぁーをabcded#$%", Katakana2Hiragana("ゴーファーヲabcded#$%"))
}

func assert(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Errorf("expected %s, but was %s\n", expected, actual)
	}
}
