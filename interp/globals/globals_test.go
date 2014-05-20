package globals

import (
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/value"
)

func TestGlobals(t *testing.T) {
	Assert(t).That(NameNum("foo"), Equals(1))
	Assert(t).That(NameNum("foo"), Equals(1))
	Assert(t).That(Add("bar", nil), Equals(2))
	Assert(t).That(NameNum("bar"), Equals(2))
	Assert(t).That(func() { Add("foo", nil) }, Panics("duplicate"))
}

var V value.Value

func BenchmarkBuffer(b *testing.B) {
	NameNum("foo")
	NameNum("bar")
	for n := 0; n < b.N; n++ {
		V = Get(1)
		V = Get(2)
	}
}
