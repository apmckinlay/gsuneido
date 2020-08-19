// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package bits

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestNextPow2(t *testing.T) {
	assert := assert.T(t).This
	assert(NextPow2(0)).Is(0)
	assert(NextPow2(1)).Is(1)
	assert(NextPow2(2)).Is(2)
	assert(NextPow2(3)).Is(4)
	assert(NextPow2(123)).Is(128)
	assert(NextPow2(65536)).Is(65536)
}
