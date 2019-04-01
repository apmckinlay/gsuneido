package tuple

import (
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/runtime"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestTupleBuilder(t *testing.T) {
	var tb TupleBuilder
	tup := tb.Build()
	Assert(t).That([]byte(tup), Equals([]byte{type8 << 6, 0, 3}))
	tb.AddRaw("one")
	tup = tb.Build()
	Assert(t).That([]byte(tup), Equals([]byte{type8 << 6, 1, 7, 4, 'o', 'n', 'e'}))
	Assert(t).That(tup.GetRaw(0), Equals("one"))

	tb = TupleBuilder{}
	tb.Add(SuInt(123))
	tb.Add(SuStr("foobar"))

	tup = tb.Build()
	Assert(t).That(tup.mode(), Equals(type8))
	Assert(t).That(tup.Count(), Equals(2))
	Assert(t).That(tup.GetVal(0), Equals(SuInt(123)))
	Assert(t).That(tup.GetVal(1), Equals(SuStr("foobar")))

	s := strings.Repeat("helloworld", 30)
	tb.AddRaw(s)
	tup = tb.Build()
	Assert(t).That(tup.mode(), Equals(type16))
	Assert(t).That(tup.GetRaw(2), Equals(s))

}

func TestLength(t *testing.T) {
	Assert(t).That(tblength(0, 0), Equals(3))
	Assert(t).That(tblength(1, 1), Equals(5))
	Assert(t).That(tblength(1, 200), Equals(204))
	Assert(t).That(tblength(1, 248), Equals(252))

	Assert(t).That(tblength(1, 252), Equals(258))
	Assert(t).That(tblength(1, 300), Equals(306))

	Assert(t).That(tblength(1, 0x10000), Equals(0x1000a))
}
