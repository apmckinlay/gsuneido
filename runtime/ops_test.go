package runtime

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/dnum"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestDiv(t *testing.T) {
	q := Div(SuInt(999), SuInt(3))
	xi, xok := SmiToInt(q)
	Assert(t).That(xok, Equals(true))
	Assert(t).That(xi, Equals(333))
	q = Div(SuInt(1), SuInt(3))
	_ = q.(SuDnum)
}

func TestBool(t *testing.T) {
	Assert(t).True(SuBool(true) == True)
	Assert(t).True(SuBool(false) == False)
}
func TestIndex(t *testing.T) {
	_, ok := Index2(False)
	Assert(t).False(ok)
	_, ok = Index2(EmptyStr)
	Assert(t).False(ok)
	Assert(t).That(Index(SuInt(123)), Equals(123))
	Assert(t).That(Index(SuDnum{Dnum: dnum.FromInt(123)}), Equals(123))
}
