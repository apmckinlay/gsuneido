// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSetFields(t *testing.T) {
	test := func(table string, iIndex int, expected []string) {
		t.Helper()
		schema := testSchemas[table]
		setFields(schema)
		assert.T(t).This(schema.Indexes[iIndex].Fields).Is(expected)
	}
	// key - Fields equals Columns
	test("supplier", 0, []string{"supplier"})
	// non-key - Fields includes key columns
	test("supplier", 1, []string{"city", "supplier"})

	// key with multiple columns
	test("hist", 1, []string{"date", "item", "id"})
	// non-key - key columns not in index columns are added
	test("hist", 0, []string{"item", "date", "id"})

	// single column key
	test("hist2", 1, []string{"date"})
	// non-key with single column key
	test("hist2", 0, []string{"id", "date"})

	// all schemas have Fields set
	for table, schema := range testSchemas {
		setFields(schema)
		for i, ix := range schema.Indexes {
			if len(ix.Fields) == 0 {
				t.Errorf("%s index %d Fields not set", table, i)
			}
			if ix.Mode == 'k' {
				if !slices.Equal(ix.Fields, ix.Columns) {
					t.Errorf("%s key %d Fields should equal Columns", table, i)
				}
			}
		}
	}
}
