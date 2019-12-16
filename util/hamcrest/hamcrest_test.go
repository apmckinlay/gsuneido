// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hamcrest

import gotesting "testing"

func TestHamcrest(t *gotesting.T) {
	assert := Assert(t)
	assert.That(func() { panic("error") }, Panics("error"))
}
