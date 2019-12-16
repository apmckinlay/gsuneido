// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hash

import (
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"testing"
)

func TestHash(t *testing.T) {
	test := func(s string, expected uint32) {
		Assert(t).That(HashString(s), Equals(expected))
		Assert(t).That(HashBytes([]byte(s)), Equals(expected))
	}
	test("", 0x811c9dc5)
	test("foobar", 0xbf9cf968)
}
