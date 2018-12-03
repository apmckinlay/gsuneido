package runtime

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestGlobals(t *testing.T) {
	Assert(t).That(GlobalNum("foo"), Equals(1))
	Assert(t).That(GlobalNum("foo"), Equals(1))
	Assert(t).That(AddGlobal("bar", nil), Equals(2))
	Assert(t).That(GlobalNum("bar"), Equals(2))
	AddGlobal("baz", True)
	Assert(t).That(func() { AddGlobal("baz", False) }, Panics("duplicate"))
	Assert(t).That(GlobalName(1), Equals("foo"))
	Assert(t).That(GlobalName(2), Equals("bar"))
}

var V Value

func BenchmarkBuffer(b *testing.B) {
	GlobalNum("foo")
	GlobalNum("bar")
	for n := 0; n < b.N; n++ {
		V = GetGlobal(1)
		V = GetGlobal(2)
	}
}
