package interp

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestGlobals(t *testing.T) {
	Assert(t).That(NameNumG("foo"), Equals(1))
	Assert(t).That(NameNumG("foo"), Equals(1))
	Assert(t).That(AddG("bar", nil), Equals(2))
	Assert(t).That(NameNumG("bar"), Equals(2))
	Assert(t).That(func() { AddG("foo", nil) }, Panics("duplicate"))
}

var V Value

func BenchmarkBuffer(b *testing.B) {
	NameNumG("foo")
	NameNumG("bar")
	for n := 0; n < b.N; n++ {
		V = GetG(1)
		V = GetG(2)
	}
}
