package base

import (
	"testing"
)

func TestSuConcatEquals(t *testing.T) {
	data := []string{"", "a", "ab", "aba", "abc"}
	var str, cat [5]Value
	for i, s := range data {
		str[i] = SuStr(s)
		cat[i] = NewSuConcat().Add(s)
	}
	for i, s := range str {
		for j, c := range cat {
			expected := i == j
			if s.Equals(c) != expected || c.Equals(s) != expected {
				t.Error(s, "vs", c)
			}
		}
	}
}
