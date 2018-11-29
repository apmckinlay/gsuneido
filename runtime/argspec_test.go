package runtime

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestArgSpecString(t *testing.T) {
	test := func(as *ArgSpec, expected string) {
		t.Helper()
		Assert(t).That(as.String(), Equals(expected))
	}
	test(ArgSpec0, "ArgSpec()")
	test(ArgSpec3, "ArgSpec(?, ?, ?)")
	test(ArgSpecEach, "ArgSpec(@)")
	test(ArgSpecEach1, "ArgSpec(@+1)")
	test(ArgSpecBlock, "ArgSpec(block:)")
	test(&ArgSpec{Nargs: 0, Spec: []byte{2, 0, 1}, Names: vals("a", "b", "c")},
		"ArgSpec(c:, a:, b:)")
	test(&ArgSpec{Nargs: 4, Spec: []byte{2, 1}, Names: vals("a", "b", "c", "d")},
		"ArgSpec(?, ?, c:, b:)")
}

func TestArgSpecEqual(t *testing.T) {
	as := []*ArgSpec{
		ArgSpec0,
		ArgSpec4,
		ArgSpecEach,
		ArgSpecEach1,
		ArgSpecBlock,
		&ArgSpec{Nargs: 2, Spec: []byte{0,1}, Names: []Value{SuStr("foo"), SuStr("bar")}},
		&ArgSpec{Nargs: 2, Spec: []byte{0,1}, Names: []Value{SuStr("foo"), SuStr("baz")}},
	}
	for i, x := range as {
		for j, y := range as {
			Assert(t).That(x.Equal(y), Equals(i == j))
		}
	}
}
