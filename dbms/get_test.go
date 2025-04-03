// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestFindIndex(t *testing.T) {
	test := func(indexes [][]string, fields []string, expectedIndex []string, expectedLen int) {
		t.Helper()
		idx, idxlen := findIndex(indexes, fields)
		assert.T(t).This(idx).Is(expectedIndex)
		assert.T(t).This(idxlen).Is(expectedLen)
	}

	// Test empty indexes
	test([][]string{},
		[]string{"a", "b"}, nil, 0)

	// Test empty fields
	test([][]string{{"a", "b"}, {"c", "d"}},
		[]string{}, nil, 0)

	// Test exact match
	test([][]string{{"a", "b"}, {"c", "d"}},
		[]string{"a", "b"}, []string{"a", "b"}, 2)

	// Test partial match at beginning
	test([][]string{{"a", "b", "c"}, {"d", "e"}},
		[]string{"a", "b"}, []string{"a", "b", "c"}, 2)

	// Test partial match in middle (should not match)
	test([][]string{{"a", "b", "c"}, {"d", "e"}},
		[]string{"b", "c"}, nil, 0)

	// Test no match
	test([][]string{{"a", "b"}, {"c", "d"}},
		[]string{"x", "y"}, nil, 0)

	// Test multiple matches - first longest
	test([][]string{{"a", "b", "c"}, {"a", "b"}, {"a"}},
		[]string{"a", "b", "c", "d"}, []string{"a", "b", "c"}, 3)

	// Test multiple matches - later longest
	test([][]string{{"a"}, {"a", "b"}, {"a", "b", "c"}},
		[]string{"a", "b", "c", "d"}, []string{"a", "b", "c"}, 3)

	// Test fields in different order
	test([][]string{{"a", "b", "c"}},
		[]string{"c", "b", "a"}, []string{"a", "b", "c"}, 3)

	// Test fields subset in different order
	test([][]string{{"a", "b", "c", "d"}},
		[]string{"c", "a", "d"}, []string{"a", "b", "c", "d"}, 1)
}
