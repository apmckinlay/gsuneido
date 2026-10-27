// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package set

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestNPerms(t *testing.T) {
	assert.T(t).This(NPerms(5, 0)).Is(1)
	assert.T(t).This(NPerms(5, 3)).Is(60)
	assert.T(t).This(NPerms(30, 4)).Is(657720)
	assert.T(t).This(func() { NPerms(0, 4) }).Panics("NPerms invalid arguments")
	assert.T(t).This(func() { NPerms(55, 44) }).Panics("NPerms overflow")
}

func TestSubPerm(t *testing.T) {
	elements := []string{"a", "b", "c", "d", "e"}

	assert.T(t).This(NthPerm(elements, 2, 0)).Is([]string{"a", "b"})
	assert.T(t).This(NthPerm(elements, 0, 0)).Is([]string{})

	n := NPerms(len(elements), 3)
	seen := make(map[string]bool)
	for i := range n {
		perm := NthPerm(elements, 3, i)
		assert.That(len(perm) == 3)
		s := strings.Join(perm, ",")
		assert.That(!seen[s])
		seen[s] = true
	}
}
