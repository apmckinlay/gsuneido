// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package check

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSet(t *testing.T) {
	assert := assert.T(t)
	var x set
	x = x.with("a")
	assert.This(x).Is(set{"a"})
	x = x.with("b")
	assert.This(x).Is(set{"a", "b"})
	x = x.with("c")
	assert.This(x).Is(set{"a", "b", "c"})
	for _, s := range x {
		x = x.with(s)
		assert.This(x).Is(set{"a", "b", "c"})
		assert.That(x.has(s))
	}
	assert.That(!x.has("d"))

	y := x.union(set{"c", "d", "b"})
	assert.This(y).Is(set{"a", "b", "c", "d"})
	assert.This(x).Is(set{"a", "b", "c"})

	x = x.intersect(set{"c", "d", "b"})
	assert.This(x).Is(set{"b", "c"})

	x = set{"a"}
	y = x.with("b").with("c")
	z := x.with("c").with("d")
	x = x.unionIntersect(y, z)
	assert.This(x).Is(set{"a", "c"})
}
