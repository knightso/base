package ngram

import (
	"testing"
)

func TestToUnigrams(t *testing.T) {
	result := ToUnigrams("abc dあいbCh")
	if len(result) != 7 {
		t.Errorf("len(result) exected:%d, but was:%d\n", 7, len(result))
	}

	if !result['a'] {
		t.Errorf("Unigram notfound. 'a'")
	}
	if !result['b'] {
		t.Errorf("Unigrbm notfound. 'a'")
	}
	if !result['c'] {
		t.Errorf("Unigrbm notfound. 'c'")
	}
	if !result['d'] {
		t.Errorf("Unigrbm notfound. 'd'")
	}
	if !result['あ'] {
		t.Errorf("Unigrbm notfound. 'あ'")
	}
	if !result['い'] {
		t.Errorf("Unigrbm notfound. 'い'")
	}
	if !result['h'] {
		t.Errorf("Unigrbm notfound. 'h'")
	}
}

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
