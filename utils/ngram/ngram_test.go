package ngram

import (
	"testing"
)

func TestToBigrams(t *testing.T) {
	result := ToBigrams("abc debch iJあdeN")
	if len(result) != 9 {
		t.Errorf("len(result) exected:%d, but was:%d\n", 9, len(result))
	}
	assertBigram(t, result, Bigram{'a', 'b'})
	assertBigram(t, result, Bigram{'b', 'c'})
	assertBigram(t, result, Bigram{'d', 'e'})
	assertBigram(t, result, Bigram{'e', 'b'})
	assertBigram(t, result, Bigram{'c', 'h'})
	assertBigram(t, result, Bigram{'i', 'j'})
	assertBigram(t, result, Bigram{'j', 'あ'})
	assertBigram(t, result, Bigram{'あ', 'd'})
	assertBigram(t, result, Bigram{'e', 'n'})
}

func TestString(t *testing.T) {
	b := Bigram{'a', 'あ'}
	if b.String() != "aあ" {
		t.Errorf("exected:%s, but was:%s\n", "aあ", b.String())
	}
}

func assertBigram(t *testing.T, set map[Bigram]bool, bigram Bigram) {
	if !set[bigram] {
		t.Errorf("Bigram notfound. %v\n", bigram)
	}
}
