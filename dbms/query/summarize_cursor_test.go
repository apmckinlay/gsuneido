// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// TestSummarizeSeqCursor tests sumSeq strategy cursor behavior
func TestSummarizeSeqCursor(t *testing.T) {
	// Create source with data that produces 3 groups when summarized by "grp"
	source := newSummarizeTestSource()
	sum := NewSummarize(source, "", []string{"grp"}, []string{"count"}, []string{"count"}, []string{""})
	sum.summarizeApproach = summarizeApproach{strat: sumSeq}
	seqT := sumSeqT{}
	sum.get = seqT.getSeq
	sum.rewound = true

	// Get all rows with Next to establish expected order
	sum.Rewind()
	nextRows := getSummarizeTestRows(t, sum, Next)
	assert.T(t).This(len(nextRows)).Is(3) // a:2, b:3, c:1

	// Test 1: All Next
	t.Run("AllNext", func(t *testing.T) {
		sum.Rewind()
		for i := 0; i < 3; i++ {
			row := sum.Get(nil, Next)
			assert.T(t).That(row != nil)
			assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[i])
		}
		assert.T(t).That(sum.Get(nil, Next) == nil)
	})

	// Test 2: All Prev
	t.Run("AllPrev", func(t *testing.T) {
		sum.Rewind()
		for i := 2; i >= 0; i-- {
			row := sum.Get(nil, Prev)
			assert.T(t).That(row != nil)
			assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[i])
		}
		assert.T(t).That(sum.Get(nil, Prev) == nil)
	})

	// Test 3: Next then Prev (at start)
	t.Run("NextThenPrev", func(t *testing.T) {
		sum.Rewind()
		row := sum.Get(nil, Next)
		assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[0])
		row = sum.Get(nil, Prev)
		assert.T(t).That(row == nil)
	})

	// Test 4: Multiple Nexts then Prev
	t.Run("MultipleNextsThenPrev", func(t *testing.T) {
		sum.Rewind()
		sum.Get(nil, Next)        // row[0]
		sum.Get(nil, Next)        // row[1]
		row := sum.Get(nil, Next) // row[2]
		assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[2])

		row = sum.Get(nil, Prev)
		assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[1])

		row = sum.Get(nil, Prev)
		assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[0])

		row = sum.Get(nil, Prev)
		assert.T(t).That(row == nil)
	})

	// Test 5: Prev then Next (from end)
	t.Run("PrevThenNext", func(t *testing.T) {
		sum.Rewind()
		row := sum.Get(nil, Prev)
		assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[2])

		row = sum.Get(nil, Next)
		assert.T(t).That(row == nil)
	})

	// Test 6: Next to end, then Prev (critical edge case)
	t.Run("NextToEndThenPrev", func(t *testing.T) {
		sum.Rewind()
		for i := 0; i < 3; i++ {
			row := sum.Get(nil, Next)
			assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[i])
		}
		assert.T(t).That(sum.Get(nil, Next) == nil)

		row := sum.Get(nil, Prev)
		assert.T(t).That(row != nil)
		assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[2])

		row = sum.Get(nil, Prev)
		assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[1])
	})

	// Test 7: NextPrevNext
	t.Run("NextPrevNext", func(t *testing.T) {
		sum.Rewind()
		row := sum.Get(nil, Next)
		assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[0])

		row = sum.Get(nil, Prev)
		assert.T(t).That(row == nil)

		row = sum.Get(nil, Next)
		assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[0])
	})

	// Test 8: Alternating in middle
	t.Run("AlternatingMiddle", func(t *testing.T) {
		sum.Rewind()
		sum.Get(nil, Next) // row[0]
		sum.Get(nil, Next) // row[1]

		row := sum.Get(nil, Prev)
		assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[0])

		row = sum.Get(nil, Next)
		assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[1])

		row = sum.Get(nil, Prev)
		assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[0])
	})
}

// TestSummarizeMapCursor tests sumMap strategy cursor behavior
func TestSummarizeMapCursor(t *testing.T) {
	source := newSummarizeTestSource()
	sum := NewSummarize(source, sumSmall, []string{"grp"}, []string{"count"}, []string{"count"}, []string{""})
	sum.summarizeApproach = summarizeApproach{strat: sumMap}
	mapT := sumMapT{}
	sum.get = mapT.getMap
	sum.rewound = true

	sum.Rewind()
	nextRows := getSummarizeTestRows(t, sum, Next)
	assert.T(t).This(len(nextRows)).Is(3)

	// Test: All Next
	t.Run("AllNext", func(t *testing.T) {
		sum.Rewind()
		for i := 0; i < 3; i++ {
			row := sum.Get(nil, Next)
			assert.T(t).That(row != nil)
			assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[i])
		}
		assert.T(t).That(sum.Get(nil, Next) == nil)
	})

	// Test: All Prev
	t.Run("AllPrev", func(t *testing.T) {
		sum.Rewind()
		for i := 2; i >= 0; i-- {
			row := sum.Get(nil, Prev)
			assert.T(t).That(row != nil)
			assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[i])
		}
		assert.T(t).That(sum.Get(nil, Prev) == nil)
	})

	// Test: Next to end, then Prev
	t.Run("NextToEndThenPrev", func(t *testing.T) {
		sum.Rewind()
		for i := 0; i < 3; i++ {
			sum.Get(nil, Next)
		}
		assert.T(t).That(sum.Get(nil, Next) == nil)

		row := sum.Get(nil, Prev)
		assert.T(t).That(row != nil)
		assert.T(t).This(getGrpValue(row, sum)).Is(nextRows[2])
	})
}

// Helper: create a test source with grouped data
// Returns source with rows: a,a,b,b,b,c (groups: a:2, b:3, c:1)
func newSummarizeTestSource() *testSummarizeSource {
	hdr := NewHeader([][]string{{"grp", "val"}}, []string{"grp", "val"})
	var rows []Row
	data := []struct{ grp, val string }{
		{"a", "1"}, {"a", "2"},
		{"b", "3"}, {"b", "4"}, {"b", "5"},
		{"c", "6"},
	}
	for i, d := range data {
		var rb RecordBuilder
		rb.Add(SuStr(d.grp))
		rb.Add(SuStr(d.val))
		rows = append(rows, Row{DbRec{Record: rb.Build(), Off: uint64(i)}})
	}
	return &testSummarizeSource{
		data:   *NewDataSource(rows),
		header: hdr,
	}
}

type testSummarizeSource struct {
	QueryMock
	data   dataSource
	header *Header
}

func (s *testSummarizeSource) Header() *Header     { return s.header }
func (s *testSummarizeSource) Columns() []string   { return []string{"grp", "val"} }
func (s *testSummarizeSource) Keys() [][]string    { return [][]string{} }
func (s *testSummarizeSource) Indexes() [][]string { return [][]string{{"grp"}} }
func (s *testSummarizeSource) Fixed() []Fixed      { return nil }
func (s *testSummarizeSource) Nrows() (int, int)   { return len(s.data.rows), len(s.data.rows) }
func (s *testSummarizeSource) Rewind()             { s.data.rewind() }
func (s *testSummarizeSource) Get(_ *Thread, dir Dir) Row {
	return s.data.get(dir)
}
func (s *testSummarizeSource) Select(cols, vals []string) { s.Rewind() }
func (s *testSummarizeSource) Simple(*Thread) []Row       { return s.data.rows }
func (s *testSummarizeSource) Transform() Query           { return s }
func (s *testSummarizeSource) setApproach([]string, float64, any, QueryTran) {
}

func getSummarizeTestRows(t *testing.T, sum *Summarize, dir Dir) []string {
	t.Helper()
	sum.Rewind()
	var result []string
	for {
		row := sum.Get(nil, dir)
		if row == nil {
			break
		}
		result = append(result, getGrpValue(row, sum))
	}
	return result
}

func getGrpValue(row Row, sum *Summarize) string {
	return row.GetRaw(sum.Header(), "grp")
}
