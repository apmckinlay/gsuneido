// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package sample

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestFrom(t *testing.T) {
	seen := make(map[int]bool)
	for v := range From(10) {
		assert.T(t).This(seen[v]).Is(false)
		seen[v] = true
		assert.T(t).Msg(v).This(v >= 0 && v < 10).Is(true)
		if len(seen) == 5 {
			break
		}
	}
}

func TestTake(t *testing.T) {
	seen := make(map[int]bool)
	for v := range Take(5, 10) {
		assert.T(t).This(seen[v]).Is(false)
		seen[v] = true
		assert.T(t).Msg(v).This(v >= 0 && v < 10).Is(true)
	}
	assert.T(t).This(len(seen)).Is(5)
}
