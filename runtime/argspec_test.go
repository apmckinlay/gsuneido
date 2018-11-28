package runtime

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestArgSpecString(t *testing.T) {
	test := func(as *ArgSpec, expected string) {
		Assert(t).That(as.String(), Equals(expected))
	}
	test(&ArgSpec{0, nil, nil}, "ArgSpec()")
	test(&ArgSpec{3, nil, nil}, "ArgSpec(?, ?, ?)")
	test(&ArgSpec{EACH, nil, nil}, "ArgSpec(@)")
	test(&ArgSpec{EACH1, nil, nil}, "ArgSpec(@+1)")
	test(&ArgSpec{0, []byte{2, 0, 1}, vals("a", "b", "c")},
		"ArgSpec(c:, a:, b:)")
	test(&ArgSpec{2, []byte{2, 1}, vals("a", "b", "c", "d")},
		"ArgSpec(?, ?, c:, b:)")
}

func TestArgSpecEqual(t *testing.T) {
	as := []*ArgSpec{
		&ArgSpec{Unnamed: 0},
		&ArgSpec{Unnamed: 4},
		&ArgSpec{Unnamed: EACH},
		&ArgSpec{Unnamed: EACH1},
		&ArgSpec{Spec: []byte{0}, Names: []Value{SuStr("foo")}},
		&ArgSpec{Spec: []byte{0}, Names: []Value{SuStr("block")}},
		&ArgSpec{Spec: []byte{0,1}, Names: []Value{SuStr("foo"), SuStr("bar")}},
	}
	for i, x := range as {
		for j, y := range as {
			Assert(t).That(x.Equal(y), Equals(i == j))
		}
	}
}
