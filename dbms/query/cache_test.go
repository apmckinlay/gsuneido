// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCache(t *testing.T) {
	var c cache
	test := func(idx []string, fixcost, varcost int) {
		t.Helper()
		fc, vc, _ := c.cacheGet(idx, 0)
		assert.T(t).This(fc).Is(fixcost)
		assert.T(t).This(vc).Is(varcost)
	}
	test(nil, -1, -1)
	ix1 := []string{"one"}
	test(ix1, -1, -1)
	c.cacheAdd(ix1, 0, 12, 34, nil)
	test(ix1, 12, 34)
	ix2 := []string{"one", "two"}
	test(ix2, -1, -1)
	c.cacheAdd(ix2, 0, 45, 67, nil)
	test(ix1, 12, 34)
	test(ix2, 45, 67)
}
