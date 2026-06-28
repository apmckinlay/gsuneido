// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCache(t *testing.T) {
	var c cache
	test := func(req Require, fixcost, varcost int) {
		t.Helper()
		fc, vc, _ := c.cacheGet(req)
		assert.T(t).This(fc).Is(fixcost)
		assert.T(t).This(vc).Is(varcost)
	}
	test(NewRequire(nil, 0, 0), -1, -1)
	r1 := NewRequire([]string{"one"}, 0, 0)
	test(r1, -1, -1)
	c.cacheAdd(r1, 12, 34, nil)
	test(r1, 12, 34)
	r2 := NewRequire([]string{"one", "two"}, 0, 0)
	test(r2, -1, -1)
	c.cacheAdd(r2, 45, 67, nil)
	test(r1, 12, 34)
	test(r2, 45, 67)
	// frac should differentiate
	r3 := NewRequire([]string{"one"}, .5, 0)
	test(r3, -1, -1)
	c.cacheAdd(r3, 99, 99, nil)
	test(r3, 99, 99)
	test(r1, 12, 34) // still returns original entry
	c.cacheClear()
	test(r1, -1, -1)
	test(r3, -1, -1)
}
