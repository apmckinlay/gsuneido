// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCache(t *testing.T) {
	var c cache
	test := func (idx []string, expected int) {
		t.Helper()
		cost,_ := c.cacheGet(idx)
		assert.T(t).This(cost).Is(expected)
	}
	test(nil, -1)
	ix1 := []string{"one"}
	test(ix1, -1)
	c.cacheAdd(ix1, 123, nil)
	test(ix1, 123)
	ix2 := []string{"one", "two"}
	test(ix2, -1)
	c.cacheAdd(ix2, 456, nil)
	test(ix1, 123)
	test(ix2, 456)
}
