// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ord

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestMin(t *testing.T) {
	assert.This(Min(1, 2)).Is(1)
}
