// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestBuffer(t *testing.T) {
	var b buffer

	test1 := func(n int) {
		Assert(t).That(b.put1(n).get1(), Equals(n))
	}
	test1(0)
	test1(1<<8 - 1)

	test2 := func(n int) {
		Assert(t).That(b.put2(n).get2(), Equals(n))
	}
	test2(0)
	test2(1<<16 - 1)

	test3 := func(n int) {
		Assert(t).That(b.put3(n).get3(), Equals(n))
	}
	test3(0)
	test3(1<<24 - 1)

	test4 := func(n int) {
		Assert(t).That(b.put4(n).get4(), Equals(n))
	}
	test4(0)
	test4(1<<31 - 1)

	test5 := func(n uint64) {
		Assert(t).That(b.put5(n).get5(), Equals(n))
	}
	test5(0)
	test5(1<<40 - 1)
}
