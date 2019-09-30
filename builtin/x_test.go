package builtin

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestNumberPat(t *testing.T) {
	Assert(t).True(numberPat.Matches("0"))
	Assert(t).True(numberPat.Matches("123"))
	Assert(t).True(numberPat.Matches("+123"))
	Assert(t).True(numberPat.Matches("-123"))
	Assert(t).True(numberPat.Matches(".123"))
	Assert(t).True(numberPat.Matches("123.465"))
	Assert(t).True(numberPat.Matches("-.5"))
	Assert(t).True(numberPat.Matches("-1.5"))
	Assert(t).True(numberPat.Matches("-1.5e2"))
	Assert(t).True(numberPat.Matches("1.5e-23"))

	Assert(t).False(numberPat.Matches(""))
	Assert(t).False(numberPat.Matches("."))
	Assert(t).False(numberPat.Matches("+"))
	Assert(t).False(numberPat.Matches("-"))
	Assert(t).False(numberPat.Matches("-."))
	Assert(t).False(numberPat.Matches("+-."))
	Assert(t).False(numberPat.Matches("1.2.3"))
}

func TestHeap(*testing.T) {
	p := alloc(64)
	*(*[64]byte)(p) = [64]byte{}
	q := alloc(8)
	*(*int64)(q) = 0
	free(q, 8)
	free(p, 64)
}
