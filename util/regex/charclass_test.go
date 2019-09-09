package regex

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestCharClass(t *testing.T) {
	test := func(in inst, s string, expected int) {
		t.Helper()
		pat := Pattern([]inst{in})
		var r Result
		Assert(t).That(pat.match(s, 0, 0, &r), Equals(expected))
	}
	test(digit, "x", -1)
	test(digit, "0", 0)
	test(digit, "5", 0)
	test(digit, "9", 0)
	test(notWord, "x", -1)
	test(notWord, "_", -1)
	test(notWord, "5", -1)
	test(notWord, " ", 0)
	test(notWord, "+", 0)
	test(space, " ", 0)
	test(space, "\t", 0)
	test(space, "\r", 0)
	test(space, "\n", 0)
	test(space, "x", -1)
	test(space, "0", -1)
}
