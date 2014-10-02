package ngram

import (
	"fmt"
	"strings"
)

type Bigram struct {
	A, B rune
}

func (bigram *Bigram) String() string {
	return fmt.Sprintf("%c%c", bigram.A, bigram.B)
}

func ToBigrams(value string) map[Bigram]bool {
	result := make(map[Bigram]bool)
	var prev rune
	for i, r := range strings.ToLower(value) {
		if i > 0 && prev != ' ' && r != ' ' {
			result[Bigram{prev, r}] = true
		}
		prev = r
	}
	return result
}
