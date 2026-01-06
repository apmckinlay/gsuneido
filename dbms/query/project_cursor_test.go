// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// Tests to isolate Project cursor issues with bidirectional walking

// TestProjectSeqCursor tests projSeq strategy cursor behavior
func TestProjectSeqCursor(t *testing.T) {
	// Create source with 5 rows, no duplicates
	source := newTestSource([]string{"a", "b", "c", "d", "e"})
	proj := NewProject(source, []string{"val"})
	proj.projectApproach = projectApproach{strat: projSeq}
	
	// Get all rows with Next to establish expected order
	proj.Rewind()
	nextRows := getAllTestRows(t, proj, Next)
	assert.T(t).This(len(nextRows)).Is(5)

	// Test 1: All Next
	t.Run("AllNext", func(t *testing.T) {
		proj.Rewind()
		for i := 0; i < 5; i++ {
			row := proj.Get(nil, Next)
			assert.T(t).That(row != nil)
			assert.T(t).This(getValue(row, proj)).Is(nextRows[i])
		}
		// One more should return nil
		assert.T(t).That(proj.Get(nil, Next) == nil)
	})

	// Test 2: All Prev
	t.Run("AllPrev", func(t *testing.T) {
		proj.Rewind()
		for i := 4; i >= 0; i-- {
			row := proj.Get(nil, Prev)
			assert.T(t).That(row != nil)
			assert.T(t).This(getValue(row, proj)).Is(nextRows[i])
		}
		// One more should return nil
		assert.T(t).That(proj.Get(nil, Prev) == nil)
	})

	// Test 3: Next then Prev (at start)
	t.Run("NextThenPrev", func(t *testing.T) {
		proj.Rewind()
		// Next returns row[0]
		row := proj.Get(nil, Next)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])
		// Prev from row[0] should return nil (nothing before first)
		row = proj.Get(nil, Prev)
		assert.T(t).That(row == nil)
	})

	// Test 4: Multiple Nexts then Prev
	t.Run("MultipleNextsThenPrev", func(t *testing.T) {
		proj.Rewind()
		// Next returns row[0], row[1], row[2]
		proj.Get(nil, Next)        // row[0]
		proj.Get(nil, Next)        // row[1]
		row := proj.Get(nil, Next) // row[2]
		assert.T(t).This(getValue(row, proj)).Is(nextRows[2])

		// Prev should return row[1]
		row = proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1])

		// Prev should return row[0]
		row = proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])

		// Prev should return nil
		row = proj.Get(nil, Prev)
		assert.T(t).That(row == nil)
	})

	// Test 5: Prev then Next (from end)
	t.Run("PrevThenNext", func(t *testing.T) {
		proj.Rewind()
		// Prev from rewind returns row[4] (last)
		row := proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[4])

		// Next should return nil (nothing after last)
		row = proj.Get(nil, Next)
		assert.T(t).That(row == nil)
	})

	t.Run("NextPrevNext", func(t *testing.T) {
		proj.Rewind()
		// Next from rewind returns row[0] (first)
		row := proj.Get(nil, Next)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])

		// Prev should return nil (nothing before first)
		row = proj.Get(nil, Prev)
		assert.T(t).That(row == nil)

		// Next should row[0] (first)
		row = proj.Get(nil, Next)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])
	})

	// Test 6: Alternating Next/Prev in middle
	t.Run("AlternatingMiddle", func(t *testing.T) {
		proj.Rewind()
		// Go to row[2]
		proj.Get(nil, Next) // row[0]
		proj.Get(nil, Next) // row[1]
		proj.Get(nil, Next) // row[2]

		// Prev back to row[1]
		row := proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1])

		// Next to row[2]
		row = proj.Get(nil, Next)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[2])

		// Prev back to row[1]
		row = proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1])
	})

	// Test 7: Next to end, then Prev (critical edge case)
	t.Run("NextToEndThenPrev", func(t *testing.T) {
		proj.Rewind()
		// Go through all rows with Next
		for i := 0; i < 5; i++ {
			row := proj.Get(nil, Next)
			assert.T(t).This(getValue(row, proj)).Is(nextRows[i])
		}
		// One more Next returns nil
		assert.T(t).That(proj.Get(nil, Next) == nil)

		// Now Prev should return the last row
		row := proj.Get(nil, Prev)
		assert.T(t).That(row != nil)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[4])

		// Prev again should return row[3]
		row = proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[3])
	})

	// Test 8: Multiple direction changes

	t.Run("MultipleDirectionChanges", func(t *testing.T) {
		proj.Rewind()
		proj.Get(nil, Next) // row[0]
		proj.Get(nil, Next) // row[1]

		// Prev to row[0]
		row := proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])

		// Prev returns nil
		row = proj.Get(nil, Prev)
		assert.T(t).That(row == nil)

		// Next from before-first returns row[0]
		row = proj.Get(nil, Next)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])

		// Next to row[1]
		row = proj.Get(nil, Next)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1])
	})
}

// TestProjectMapCursor tests projMap strategy cursor behavior
func TestProjectMapCursor(t *testing.T) {
	// Create source with 5 rows, no duplicates
	source := newTestSource([]string{"a", "b", "c", "d", "e"})
	proj := NewProject(source, []string{"val"})
	proj.projectApproach = projectApproach{strat: projMap}

	// Get all rows with Next to establish expected order
	proj.Rewind()
	nextRows := getAllTestRows(t, proj, Next)
	assert.T(t).This(len(nextRows)).Is(5)

	// Test 1: All Next
	t.Run("AllNext", func(t *testing.T) {
		proj.Rewind()
		for i := 0; i < 5; i++ {
			row := proj.Get(nil, Next)
			assert.T(t).That(row != nil)
			assert.T(t).This(getValue(row, proj)).Is(nextRows[i])
		}
		assert.T(t).That(proj.Get(nil, Next) == nil)
	})

	// Test 2: All Prev
	t.Run("AllPrev", func(t *testing.T) {
		proj.Rewind()
		for i := 4; i >= 0; i-- {
			row := proj.Get(nil, Prev)
			assert.T(t).That(row != nil)
			assert.T(t).This(getValue(row, proj)).Is(nextRows[i])
		}
		assert.T(t).That(proj.Get(nil, Prev) == nil)
	})

	// Test 3: Next then Prev
	t.Run("NextThenPrev", func(t *testing.T) {
		proj.Rewind()
		row := proj.Get(nil, Next)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])
		// Prev from row[0] should return nil
		row = proj.Get(nil, Prev)
		assert.T(t).That(row == nil)
	})

	// Test 4: Multiple Nexts then Prev
	t.Run("MultipleNextsThenPrev", func(t *testing.T) {
		proj.Rewind()
		proj.Get(nil, Next)        // row[0]
		proj.Get(nil, Next)        // row[1]
		row := proj.Get(nil, Next) // row[2]
		assert.T(t).This(getValue(row, proj)).Is(nextRows[2])

		row = proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1])

		row = proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])

		row = proj.Get(nil, Prev)
		assert.T(t).That(row == nil)
	})

	// Test 5: Prev then Next
	t.Run("PrevThenNext", func(t *testing.T) {
		proj.Rewind()
		row := proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[4])

		row = proj.Get(nil, Next)
		assert.T(t).That(row == nil)
	})

	// Test 6: Alternating in middle
	t.Run("AlternatingMiddle", func(t *testing.T) {
		proj.Rewind()
		proj.Get(nil, Next) // row[0]
		proj.Get(nil, Next) // row[1]
		proj.Get(nil, Next) // row[2]

		row := proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1])

		row = proj.Get(nil, Next)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[2])

		row = proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1])
	})

	// Test 7: Next to end, then Prev (critical edge case)
	t.Run("NextToEndThenPrev", func(t *testing.T) {
		proj.Rewind()
		// Go through all rows with Next
		for i := 0; i < 5; i++ {
			row := proj.Get(nil, Next)
			assert.T(t).This(getValue(row, proj)).Is(nextRows[i])
		}
		// One more Next returns nil
		assert.T(t).That(proj.Get(nil, Next) == nil)

		// Now Prev should return the last row
		row := proj.Get(nil, Prev)
		assert.T(t).That(row != nil)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[4])

		// Prev again should return row[3]
		row = proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[3])
	})

	// Test 8: Multiple direction changes
	t.Run("MultipleDirectionChanges", func(t *testing.T) {
		proj.Rewind()
		proj.Get(nil, Next) // row[0]
		proj.Get(nil, Next) // row[1]

		// Prev to row[0]
		row := proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])

		// Prev returns nil
		row = proj.Get(nil, Prev)
		assert.T(t).That(row == nil)

		// Next from before-first returns row[0]
		row = proj.Get(nil, Next)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])

		// Next to row[1]
		row = proj.Get(nil, Next)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1])
	})

	// Test 11: Prev first (builds full index), then Next through all rows
	// This tests that when indexed=true and Next exhausts uniqueRows, we return nil
	t.Run("PrevThenNextToEnd", func(t *testing.T) {
		proj.Rewind()
		// Prev first - this triggers buildMap (indexed=true)
		for i := 4; i >= 0; i-- {
			row := proj.Get(nil, Prev)
			assert.T(t).This(getValue(row, proj)).Is(nextRows[i])
		}
		// One more Prev returns nil
		assert.T(t).That(proj.Get(nil, Prev) == nil)

		// Now Next from beginning should iterate through all rows
		row := proj.Get(nil, Next)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])

		// Continue Next through all rows
		for i := 1; i < 5; i++ {
			row = proj.Get(nil, Next)
			assert.T(t).This(getValue(row, proj)).Is(nextRows[i])
		}
		// One more Next returns nil (indexed=true, uniqueRows exhausted)
		assert.T(t).That(proj.Get(nil, Next) == nil)
	})

	// Test 12: Prev (indexed), Rewind, Next through all, Rewind, Prev through all
	t.Run("IndexedMixedDirections", func(t *testing.T) {
		proj.Rewind()
		// Prev first to build index
		proj.Get(nil, Prev) // last row

		// Rewind and Next through all
		proj.Rewind()
		for i := 0; i < 5; i++ {
			row := proj.Get(nil, Next)
			assert.T(t).This(getValue(row, proj)).Is(nextRows[i])
		}
		assert.T(t).That(proj.Get(nil, Next) == nil)

		// Rewind and Prev through all
		proj.Rewind()
		for i := 4; i >= 0; i-- {
			row := proj.Get(nil, Prev)
			assert.T(t).This(getValue(row, proj)).Is(nextRows[i])
		}
		assert.T(t).That(proj.Get(nil, Prev) == nil)
	})
}

// TestProjectCopyCursor tests projCopy strategy cursor behavior (should work correctly)
func TestProjectCopyCursor(t *testing.T) {
	// Create source with unique key so projCopy is used
	source := newTestSourceWithKey([]string{"a", "b", "c", "d", "e"})
	proj := NewProject(source, []string{"val"})

	assert.T(t).That(proj.unique) // Should use projCopy
	proj.projectApproach = projectApproach{strat: projCopy}

	proj.Rewind()
	nextRows := getAllTestRows(t, proj, Next)
	assert.T(t).This(len(nextRows)).Is(5)

	// Test: Alternating in middle
	t.Run("AlternatingMiddle", func(t *testing.T) {
		proj.Rewind()
		proj.Get(nil, Next) // row[0]
		proj.Get(nil, Next) // row[1]
		proj.Get(nil, Next) // row[2]

		row := proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1])

		row = proj.Get(nil, Next)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[2])

		row = proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1])
	})
}

// Helper: create a simple test source
func newTestSource(values []string) *testProjectSource {
	hdr := NewHeader([][]string{{"val"}}, []string{"val"})
	var rows []Row
	for i, v := range values {
		var rb RecordBuilder
		rb.Add(SuStr(v))
		rows = append(rows, Row{DbRec{Record: rb.Build(), Off: uint64(i)}})
	}
	return &testProjectSource{
		data:   *NewDataSource(rows),
		header: hdr,
	}
}

// Helper: create a test source with a key (so projCopy will be used)
func newTestSourceWithKey(values []string) *testProjectSource {
	src := newTestSource(values)
	src.keys = [][]string{{"val"}}
	return src
}

// testProjectSource is a simple Query source for testing
type testProjectSource struct {
	QueryMock
	data    dataSource
	header  *Header
	keys    [][]string
	selCols []string
	selVals []string
}

func (s *testProjectSource) Header() *Header {
	return s.header
}

func (s *testProjectSource) Columns() []string {
	return []string{"val"}
}

func (s *testProjectSource) Keys() [][]string {
	if s.keys != nil {
		return s.keys
	}
	return [][]string{} // empty by default - will use projSeq/projMap
}

func (s *testProjectSource) Indexes() [][]string {
	return [][]string{{"val"}}
}

func (s *testProjectSource) Fixed() []Fixed {
	return nil
}

func (s *testProjectSource) Nrows() (int, int) {
	return len(s.data.rows), len(s.data.rows)
}

func (s *testProjectSource) Rewind() {
	s.data.rewind()
}

func (s *testProjectSource) Get(_ *Thread, dir Dir) Row {
	for {
		row := s.data.get(dir)
		if row == nil {
			return nil
		}
		if s.matchesSelect(row) {
			return row
		}
	}
}

func (s *testProjectSource) matchesSelect(row Row) bool {
	if s.selCols == nil {
		return true
	}
	for i, col := range s.selCols {
		if row.GetRaw(s.header, col) != s.selVals[i] {
			return false
		}
	}
	return true
}

func (s *testProjectSource) Select(cols, vals []string) {
	s.selCols = cols
	s.selVals = vals
	s.Rewind()
}

func (s *testProjectSource) Simple(*Thread) []Row {
	return s.data.rows
}

func (s *testProjectSource) Transform() Query {
	return s
}

func (s *testProjectSource) setApproach([]string, float64, any, QueryTran) {
}

// Helper to get all rows
func getAllTestRows(t *testing.T, proj Query, dir Dir) []string {
	t.Helper()
	proj.Rewind()
	var result []string
	for {
		row := proj.Get(nil, dir)
		if row == nil {
			break
		}
		result = append(result, getValue(row, proj))
	}
	return result
}

// Helper to get value from row
func getValue(row Row, proj Query) string {
	return row.GetRaw(proj.Header(), "val")
}

// TestProjectMapWithDuplicates tests projMap with duplicate values in source
func TestProjectMapWithDuplicates(t *testing.T) {
	// Create source with duplicates: a, b, b, c, c, c, d
	source := newTestSource([]string{"a", "b", "b", "c", "c", "c", "d"})
	proj := NewProject(source, []string{"val"})
	proj.projectApproach = projectApproach{strat: projMap}

	// Get all unique rows with Next
	proj.Rewind()
	nextRows := getAllTestRows(t, proj, Next)
	assert.T(t).This(len(nextRows)).Is(4) // a, b, c, d

	t.Run("AllNext", func(t *testing.T) {
		proj.Rewind()
		for i := 0; i < 4; i++ {
			row := proj.Get(nil, Next)
			assert.T(t).This(getValue(row, proj)).Is(nextRows[i])
		}
		assert.T(t).That(proj.Get(nil, Next) == nil)
	})

	t.Run("AllPrev", func(t *testing.T) {
		proj.Rewind()
		for i := 3; i >= 0; i-- {
			row := proj.Get(nil, Prev)
			assert.T(t).This(getValue(row, proj)).Is(nextRows[i])
		}
		assert.T(t).That(proj.Get(nil, Prev) == nil)
	})

	t.Run("NextThenPrevThenNext", func(t *testing.T) {
		proj.Rewind()
		row := proj.Get(nil, Next) // row[0]
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])
		row = proj.Get(nil, Next) // row[1]
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1])
		row = proj.Get(nil, Prev) // row[0]
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])
		row = proj.Get(nil, Next) // row[1]
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1])
		row = proj.Get(nil, Next) // row[2]
		assert.T(t).This(getValue(row, proj)).Is(nextRows[2])
	})
}

// TestProjectNonCopyLookup tests the Lookup method for non-projCopy strategies.
// This code path is rarely/never reached via normal query execution because:
// - For non-unique projects, the only key is ALL columns
// - Joins only call Lookup when joinType.toOne() is true
// - toOne() requires the join's by columns to form a key in the project
// - This almost never happens for non-unique projects
// This test directly exercises that code path.
func TestProjectNonCopyLookup(t *testing.T) {
	test := func(t *testing.T, strat projectStrategy) {
		source := newTestSource([]string{"a", "b", "c", "d", "e"})
		source.keys = nil // no key - forces non-unique project
		proj := NewProject(source, []string{"val"})
		proj.projectApproach = projectApproach{strat: strat}
		assert.T(t).That(!proj.unique)
		assert.T(t).That(proj.strat != projCopy)

		// Lookup existing value
		row := proj.Lookup(nil, []string{"val"}, []string{Pack(SuStr("c"))})
		assert.T(t).That(row != nil)
		assert.T(t).This(row.GetRaw(proj.Header(), "val")).Is(Pack(SuStr("c")))

		// Lookup another existing value
		row = proj.Lookup(nil, []string{"val"}, []string{Pack(SuStr("a"))})
		assert.T(t).That(row != nil)
		assert.T(t).This(row.GetRaw(proj.Header(), "val")).Is(Pack(SuStr("a")))

		// Lookup non-existing value
		row = proj.Lookup(nil, []string{"val"}, []string{Pack(SuStr("z"))})
		assert.T(t).That(row == nil)
	}

	t.Run("projSeq", func(t *testing.T) { test(t, projSeq) })
	t.Run("projMap", func(t *testing.T) { test(t, projMap) })
}

// TestProjectMapEmpty tests projMap with empty result set
func TestProjectMapEmpty(t *testing.T) {
	source := newTestSource([]string{})
	proj := NewProject(source, []string{"val"})
	proj.projectApproach = projectApproach{strat: projMap}

	t.Run("NextReturnsNil", func(t *testing.T) {
		proj.Rewind()
		assert.T(t).That(proj.Get(nil, Next) == nil)
	})

	t.Run("PrevReturnsNil", func(t *testing.T) {
		proj.Rewind()
		assert.T(t).That(proj.Get(nil, Prev) == nil)
	})

	t.Run("NextThenPrevAtEmpty", func(t *testing.T) {
		proj.Rewind()
		assert.T(t).That(proj.Get(nil, Next) == nil)
		assert.T(t).That(proj.Get(nil, Prev) == nil)
	})

	t.Run("PrevThenNextAtEmpty", func(t *testing.T) {
		proj.Rewind()
		assert.T(t).That(proj.Get(nil, Prev) == nil)
		assert.T(t).That(proj.Get(nil, Next) == nil)
	})
}

// TestProjectSeqEmpty tests projSeq with empty result set
func TestProjectSeqEmpty(t *testing.T) {
	source := newTestSource([]string{})
	proj := NewProject(source, []string{"val"})
	proj.projectApproach = projectApproach{strat: projSeq}

	t.Run("NextReturnsNil", func(t *testing.T) {
		proj.Rewind()
		assert.T(t).That(proj.Get(nil, Next) == nil)
	})

	t.Run("PrevReturnsNil", func(t *testing.T) {
		proj.Rewind()
		assert.T(t).That(proj.Get(nil, Prev) == nil)
	})

	t.Run("NextThenPrevAtEmpty", func(t *testing.T) {
		proj.Rewind()
		assert.T(t).That(proj.Get(nil, Next) == nil)
		assert.T(t).That(proj.Get(nil, Prev) == nil)
	})

	t.Run("PrevThenNextAtEmpty", func(t *testing.T) {
		proj.Rewind()
		assert.T(t).That(proj.Get(nil, Prev) == nil)
		assert.T(t).That(proj.Get(nil, Next) == nil)
	})
}

// TestProjectSeqDirectionChangeAtBoundary tests the specific bug where
// after Prev reaches beginning (returns nil), Next returns row[0],
// then Prev should return nil (not row[0] again).
// This is the sequence: P→row[0], P→nil, N→row[0], P→nil (not row!)
func TestProjectSeqDirectionChangeAtBoundary(t *testing.T) {
	source := newTestSource([]string{"a", "b", "c"})
	proj := NewProject(source, []string{"val"})
	proj.projectApproach = projectApproach{strat: projSeq}

	proj.Rewind()
	nextRows := getAllTestRows(t, proj, Next)
	assert.T(t).This(len(nextRows)).Is(3)

	// The critical sequence that was failing
	t.Run("PrevToBeginningThenNextThenPrev", func(t *testing.T) {
		proj.Rewind()
		// Prev to row[0]
		row := proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[2]) // last row (c)
		row = proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1]) // b
		row = proj.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0]) // a (first)

		// Prev returns nil (at beginning)
		row = proj.Get(nil, Prev)
		assert.T(t).That(row == nil)

		// Next returns row[0] (first row)
		row = proj.Get(nil, Next)
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])

		// KEY TEST: Prev should return nil (we're at first row, nothing before)
		row = proj.Get(nil, Prev)
		assert.T(t).That(row == nil) // This was the bug - it returned row[0]
	})

	// Same test with single-row source
	t.Run("SingleRowPrevNilNextPrev", func(t *testing.T) {
		source1 := newTestSource([]string{"only"})
		proj1 := NewProject(source1, []string{"val"})
		proj1.projectApproach = projectApproach{strat: projSeq}

		// Get the expected row format
		proj1.Rewind()
		expectedRows := getAllTestRows(t, proj1, Next)
		assert.T(t).This(len(expectedRows)).Is(1)

		proj1.Rewind()
		// Prev returns "only"
		row := proj1.Get(nil, Prev)
		assert.T(t).This(getValue(row, proj1)).Is(expectedRows[0])

		// Prev returns nil (at beginning)
		row = proj1.Get(nil, Prev)
		assert.T(t).That(row == nil)

		// Next returns "only"
		row = proj1.Get(nil, Next)
		assert.T(t).This(getValue(row, proj1)).Is(expectedRows[0])

		// Prev should return nil (we're at first row)
		row = proj1.Get(nil, Prev)
		assert.T(t).That(row == nil)
	})
}

// TestProjectSeqWithDuplicates tests projSeq with duplicate values in source
func TestProjectSeqWithDuplicates(t *testing.T) {
	// Create source with duplicates: a, b, b, c, c, c, d
	source := newTestSource([]string{"a", "b", "b", "c", "c", "c", "d"})
	proj := NewProject(source, []string{"val"})
	proj.projectApproach = projectApproach{strat: projSeq}

	// Get all unique rows with Next
	proj.Rewind()
	nextRows := getAllTestRows(t, proj, Next)
	assert.T(t).This(len(nextRows)).Is(4) // a, b, c, d

	t.Run("AllNext", func(t *testing.T) {
		proj.Rewind()
		for i := 0; i < 4; i++ {
			row := proj.Get(nil, Next)
			assert.T(t).This(getValue(row, proj)).Is(nextRows[i])
		}
		assert.T(t).That(proj.Get(nil, Next) == nil)
	})

	t.Run("AllPrev", func(t *testing.T) {
		proj.Rewind()
		for i := 3; i >= 0; i-- {
			row := proj.Get(nil, Prev)
			assert.T(t).This(getValue(row, proj)).Is(nextRows[i])
		}
		assert.T(t).That(proj.Get(nil, Prev) == nil)
	})

	t.Run("NextThenPrevThenNext", func(t *testing.T) {
		proj.Rewind()
		row := proj.Get(nil, Next) // row[0]
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])
		row = proj.Get(nil, Next) // row[1]
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1])
		row = proj.Get(nil, Prev) // row[0]
		assert.T(t).This(getValue(row, proj)).Is(nextRows[0])
		row = proj.Get(nil, Next) // row[1]
		assert.T(t).This(getValue(row, proj)).Is(nextRows[1])
		row = proj.Get(nil, Next) // row[2]
		assert.T(t).This(getValue(row, proj)).Is(nextRows[2])
	})
}
