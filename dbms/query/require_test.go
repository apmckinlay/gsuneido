// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestMergeReq(t *testing.T) {
	assert := assert.T(t)
	test := func(req1 Use, cols1 []string, req2 Use, cols2 []string, expReq Use, expCols []string) {
		req, cols := MergeReq(req1, cols1, req2, cols2)
		assert.This(req).Is(expReq)
		assert.This(cols).Is(expCols)
		req, cols = MergeReq(req2, cols2, req1, cols1)
		assert.This(req).Is(expReq)
		assert.This(cols).Is(expCols)
	}
	// unordered + anything = the other
	test(ReqUnordered, nil, ReqOrdered, []string{"a"}, ReqOrdered, []string{"a"})
	test(ReqOrdered, []string{"a"}, ReqUnordered, nil, ReqOrdered, []string{"a"})

	// ordered + ordered: prefix match returns longer of the two
	test(ReqOrdered, []string{"a"}, ReqOrdered, []string{"a"}, ReqOrdered, []string{"a"})
	test(ReqOrdered, []string{"a"}, ReqOrdered, []string{"a", "b"}, ReqOrdered, []string{"a", "b"})
	test(ReqOrdered, []string{"a", "b"}, ReqOrdered, []string{"a"}, ReqOrdered, []string{"a", "b"})
	test(ReqOrdered, []string{"a"}, ReqOrdered, []string{"b"}, ReqConflict, nil)

	// grouped + grouped: equal works, otherwise conflict
	test(ReqGrouped, []string{"a"}, ReqGrouped, []string{"a"}, ReqGrouped, []string{"a"})
	test(ReqGrouped, []string{"a"}, ReqGrouped, []string{"b"}, ReqConflict, nil)

	// ordered + grouped
	test(ReqOrdered, []string{"a"}, ReqGrouped, []string{"a"}, ReqOrdered, []string{"a"})
	test(ReqOrdered, []string{"a", "b"}, ReqGrouped, []string{"a"}, ReqOrdered, []string{"a", "b"})
	test(ReqOrdered, []string{"a"}, ReqGrouped, []string{"a", "b"}, ReqOrdered, []string{"a", "b"})
	test(ReqOrdered, []string{"a"}, ReqGrouped, []string{"b"}, ReqConflict, nil)
	test(ReqOrdered, []string{"a"}, ReqGrouped, []string{"b", "c"}, ReqConflict, nil)

	// ordered + lookup
	test(ReqOrdered, []string{"a"}, ReqLookup, []string{"a"}, ReqOrdered, []string{"a"})
	test(ReqOrdered, []string{"a"}, ReqLookup, []string{"b"}, ReqConflict, nil)

	// grouped + lookup
	test(ReqGrouped, []string{"a"}, ReqLookup, []string{"a", "b"}, ReqLookup, []string{"a", "b"})
	test(ReqGrouped, []string{"a"}, ReqLookup, []string{"a"}, ReqLookup, []string{"a"})
	test(ReqGrouped, []string{"a"}, ReqLookup, []string{"b"}, ReqConflict, nil)

	// lookup + lookup
	test(ReqLookup, []string{"a"}, ReqLookup, []string{"a"}, ReqLookup, []string{"a"})
	test(ReqLookup, []string{"a"}, ReqLookup, []string{"b"}, ReqConflict, nil)
}
