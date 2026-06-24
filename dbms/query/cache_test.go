// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCache(t *testing.T) {
	var c cache
	test := func(use Use, idx []string, fixcost, varcost int) {
		t.Helper()
		fc, vc, _ := c.cacheGet(use, idx, 0)
		assert.T(t).This(fc).Is(fixcost)
		assert.T(t).This(vc).Is(varcost)
	}
	test(ReqOrdered, nil, -1, -1)
	ix1 := []string{"one"}
	test(ReqOrdered, ix1, -1, -1)
	c.cacheAdd(ReqOrdered, ix1, 0, 12, 34, nil)
	test(ReqOrdered, ix1, 12, 34)
	ix2 := []string{"one", "two"}
	test(ReqOrdered, ix2, -1, -1)
	c.cacheAdd(ReqOrdered, ix2, 0, 45, 67, nil)
	test(ReqOrdered, ix1, 12, 34)
	test(ReqOrdered, ix2, 45, 67)
	// use should differentiate
	test(ReqUnordered, ix1, -1, -1)
	c.cacheAdd(ReqUnordered, ix1, 0, 99, 99, nil)
	test(ReqUnordered, ix1, 99, 99)
	test(ReqOrdered, ix1, 12, 34) // still returns ReqOrdered entry
}
