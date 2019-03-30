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
	Assert(t).That([]byte(tup), Equals([]byte{'c', 0, 0, 0, 5}))
	tb.Add("one")
	tup = tb.Build()
	Assert(t).That([]byte(tup), Equals([]byte{'c', 0, 1, 0, 9, 6, 'o', 'n', 'e'}))

	tb = TupleBuilder{}
	tb.AddVal(SuInt(123))
	tb.AddVal(SuStr("foobar"))

	tup = tb.Build()
	Assert(t).That(tup.mode(), Equals('c'))
	Assert(t).That(tup.Count(), Equals(2))
	Assert(t).That(tup.GetVal(0), Equals(SuInt(123)))
	Assert(t).That(tup.GetVal(1), Equals(SuStr("foobar")))

	s := strings.Repeat("helloworld", 30)
	tb.Add(s)
	tup = tb.Build()
	Assert(t).That(tup.mode(), Equals('s'))
	Assert(t).That(tup.Get(2), Equals(s))

}

func TestLength(t *testing.T) {
	Assert(t).That(tblength(0, 0), Equals(5))
	Assert(t).That(tblength(1, 1), Equals(7))
	Assert(t).That(tblength(1, 200), Equals(206))
	Assert(t).That(tblength(1, 248), Equals(254))

	Assert(t).That(tblength(1, 250), Equals(258))
	Assert(t).That(tblength(1, 300), Equals(308))

	Assert(t).That(tblength(1, 0x10000), Equals(0x1000c))
}
