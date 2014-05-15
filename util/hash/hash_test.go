package hash

import (
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestHash(t *testing.T) {
	test := func(s string, expected uint32) {
		Assert(t).That(HashString(s), Equals(expected))
	}
	test("", 0x811c9dc5)
	test("foobar", 0xbf9cf968)
}
