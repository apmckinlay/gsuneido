package tuple

import (
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/ints"
	v "github.com/apmckinlay/gsuneido/value"
)

func TestTupleM(t *testing.T) {
	tm := TupleM{}
	Assert(t).That(tm.Size(), Equals(0))
	tm.Add(v.SuInt(123))
	tm.Add(v.SuStr("foobar"))
	Assert(t).That(tm.Size(), Equals(2))
	Assert(t).That(tm.Get(0), Equals(v.SuInt(123)))
	Assert(t).That(tm.Get(1), Equals(v.SuStr("foobar")))

	tb := tm.ToTupleB()
	Assert(t).That(tb.Size(), Equals(2))
	Assert(t).That(tb.Get(0), Equals(v.SuInt(123)))
	Assert(t).That(tb.Get(1), Equals(v.SuStr("foobar")))
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

func TestCompare(t *testing.T) {
	data := []Tuple{
		record(), record("one"), record("one", "three"),
		record("one", "two"), record("three"), record("two")}
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data); j++ {
			Assert(t).That(data[i].Compare(data[j]), Equals(ints.Compare(i, j)))
		}
	}
}

func record(strs ...string) Tuple {
	t := TupleM{}
	for _, s := range strs {
		t.Add(v.SuStr(s))
	}
	return t
}
