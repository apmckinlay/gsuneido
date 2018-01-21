package global

import (
	"testing"

	"github.com/apmckinlay/gsuneido/base"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestGlobals(t *testing.T) {
	Assert(t).That(Num("foo"), Equals(1))
	Assert(t).That(Num("foo"), Equals(1))
	Assert(t).That(Add("bar", nil), Equals(2))
	Assert(t).That(Num("bar"), Equals(2))
	Assert(t).That(func() { Add("foo", nil) }, Panics("duplicate"))
	Assert(t).That(Name(1), Equals("foo"))
	Assert(t).That(Name(2), Equals("bar"))
}

var V base.Value

func BenchmarkBuffer(b *testing.B) {
	Num("foo")
	Num("bar")
	for n := 0; n < b.N; n++ {
		V = Get(1)
		V = Get(2)
	}
}
