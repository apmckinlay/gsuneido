package interp

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestArgSpecString(t *testing.T) {
	Assert(t).That(ArgSpec{0, nil, nil}.String(),
		Equals("ArgSpec()"))
	Assert(t).That(ArgSpec{3, nil, nil}.String(),
		Equals("ArgSpec(?, ?, ?)"))
	Assert(t).That(ArgSpec{EACH, nil, nil}.String(),
		Equals("ArgSpec(@)"))
	Assert(t).That(ArgSpec{EACH1, nil, nil}.String(),
		Equals("ArgSpec(@+1)"))
	Assert(t).That(ArgSpec{0, []byte{2, 0, 1}, []string{"a", "b", "c"}}.String(),
		Equals("ArgSpec(c:, a:, b:)"))
	Assert(t).That(ArgSpec{2, []byte{2, 1}, []string{"a", "b", "c", "d"}}.String(),
		Equals("ArgSpec(?, ?, c:, b:)"))
}
