package runtime

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestGlobals(t *testing.T) {
	Assert(t).That(Global.Num("foo"), Equals(1))
	Assert(t).That(Global.Num("foo"), Equals(1))
	Assert(t).That(Global.Add("bar", nil), Equals(2))
	Assert(t).That(Global.Num("bar"), Equals(2))
	Global.Add("baz", True)
	Assert(t).That(func() { Global.Add("baz", False) }, Panics("duplicate"))
	Assert(t).That(Global.Name(1), Equals("foo"))
	Assert(t).That(Global.Name(2), Equals("bar"))
}

var V Value

func BenchmarkBuffer(b *testing.B) {
	Global.Num("foo")
	Global.Num("bar")
	for n := 0; n < b.N; n++ {
		V = Global.Get(1)
		V = Global.Get(2)
	}
}
