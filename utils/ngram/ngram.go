package ngram

import (
	"fmt"
	"strings"
)

type Bigram struct {
	a, b rune
}

func (bigram *Bigram) String() string {
	return fmt.Sprintf("%s%s", bigram.a, bigram.b)
}

func ToBigrams(value string) map[Bigram]struct{} {
	result := make(map[Bigram]struct{})
	var prev rune
	for i, r := range strings.ToLower(value) {
		if i > 0 && prev != ' ' && r != ' ' {
			result[Bigram{prev, r}] = struct{}{}
		}
		prev = r
	}
	return result
}
