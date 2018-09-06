package base

import (
	"testing"

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
