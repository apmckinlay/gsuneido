package runtime

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestGlobals(t *testing.T) {
	foo := Global.Num("foo")
	Assert(t).That(Global.Num("foo"), Equals(foo))
	Assert(t).That(Global.Add("bar", nil), Equals(foo+1))
	Assert(t).That(Global.Num("bar"), Equals(foo+1))
	Global.Add("baz", True)
	Assert(t).That(func() { Global.Add("baz", False) }, Panics("duplicate"))
	Assert(t).That(Global.Name(foo), Equals("foo"))
	Assert(t).That(Global.Name(foo+1), Equals("bar"))
}

var V Value

func BenchmarkBuffer(b *testing.B) {
	f := Global.Num("foo")
	g := Global.Num("bar")
	for n := 0; n < b.N; n++ {
		V = Global.Get(f)
		V = Global.Get(g)
	}
}
