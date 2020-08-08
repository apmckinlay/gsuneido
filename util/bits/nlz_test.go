// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package bits

import "testing"
import . "github.com/apmckinlay/gsuneido/util/hamcrest"

func Test_Nlz(t *testing.T) {
	assert := Assert(t)
	test := func(n uint64, expected int) {
		assert.That(Nlz(n), Is(expected))
	}
	test(0, 64)
	test(0xff, 56)
	test(0xffff, 48)
	test(^uint64(0), 0)
}
