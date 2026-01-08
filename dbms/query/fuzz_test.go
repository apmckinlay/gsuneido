// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
)

// go test -run '^$' -fuzz=FuzzQuerySource ./dbms/query/

func init() {
	sortForTest = true
}

func FuzzQuerySource(f *testing.F) {
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		qs := NewQuerySource(rnd)
		fuzzQuery(t, qs, rnd, chooseIndex(rnd, qs.Indexes()))
	})
}

func TestFuzzQuerySource(t *testing.T) {
	rnd := rand.New(rand.NewPCG(12351, 68081))
	for range 1000 {
		qs := NewQuerySource(rnd)
		fuzzQuery(t, qs, rnd, chooseIndex(rnd, qs.Indexes()))
	}
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzProject ./dbms/query/

func FuzzProject(f *testing.F) {
	f.Add(uint64(122), uint64(334))
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzProject(t, rnd)
	})
	// fmt.Printf("Strategy: copy=%d, seq=%d, map=%d\n",
	// 	strategyCounts[projCopy], strategyCounts[projSeq], strategyCounts[projMap])
}

func TestFuzzProject(t *testing.T) {
	rnd := rand.New(rand.NewPCG(12351, 68081))
	for range 10000 {
		fuzzProject(t, rnd)
	}
	// fmt.Printf("Strategy: copy=%d, seq=%d, map=%d\n",
	// 	strategyCounts[projCopy], strategyCounts[projSeq], strategyCounts[projMap])
}

var (
	strategyCounts = map[projectStrategy]int{
		projCopy: 0,
		projSeq:  0,
		projMap:  0,
	}
	uniqueCount = 0
)

func fuzzProject(t *testing.T, rnd *rand.Rand) {
	qs := NewQuerySource(rnd)
	projCols := randomProjectCols(rnd, qs.ColumnsResult, qs.IndexesResult)
	proj := NewProject(qs, projCols)

	index := chooseIndex(rnd, proj.Indexes())
	// CursorMode to prevent temp indexes
	fixcost, varcost, approach := proj.optimize(ReadMode, index, 1)
	if fixcost+varcost >= impossible {
		t.Fatal("impossible")
	}
	proj.cacheAdd(index, 1, fixcost, varcost, approach)

	SetApproach(proj, index, 1, nil)

	if approach != nil {
		strat := approach.(*projectApproach).strat
		strategyCounts[strat]++

		if proj.unique {
			uniqueCount++
		}
	}

	fuzzQuery(t, proj, rnd, index)
}

func randomProjectCols(rnd *rand.Rand, srcCols []string, indexes [][]string) []string {
	// 20% of the time, choose columns that allow projSeq by selecting a prefix of an index
	if len(srcCols) > 0 && len(indexes) > 0 && rnd.IntN(5) == 0 { // 20% chance
		// Choose a random index
		index := random(indexes, rnd)
		// Skip empty indexes
		if len(index) > 0 {
			// Choose a prefix of this index (1 to full length)
			prefixLen := 1 + rnd.IntN(len(index))
			return index[:prefixLen]
		}
		// Fall through to original random selection
	}

	// 80% of the time, use original random selection (or when index is empty)
	n := 1 + rnd.IntN(len(srcCols)) // 1 to all columns
	perm := rnd.Perm(len(srcCols))
	cols := make([]string, n)
	for i := range n {
		cols[i] = srcCols[perm[i]]
	}
	return cols
}

func chooseIndex(rnd *rand.Rand, indexes [][]string) []string {
	if len(indexes) == 0 {
		return nil
	}
	index := random(indexes, rnd)
	if len(index) == 0 {
		return nil
	}
	n := rnd.IntN(len(index))
	if n == 0 {
		return nil
	}
	return index[:n]
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzRename ./dbms/query/

func FuzzRename(f *testing.F) {
	f.Add(uint64(122), uint64(334))
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzRename(t, rnd)
	})
}

func TestFuzzRename(t *testing.T) {
	rnd := rand.New(rand.NewPCG(12351, 68081))
	for range 1000 {
		fuzzRename(t, rnd)
	}
}

func fuzzRename(t *testing.T, rnd *rand.Rand) {
	qs := NewQuerySource(rnd)
	from, to := randomRename(rnd, qs.ColumnsResult)
	ren := NewRename(qs, from, to)

	index := chooseIndex(rnd, ren.Indexes())
	fixcost, varcost, approach := ren.optimize(ReadMode, index, 1)
	if fixcost+varcost >= impossible {
		t.Fatal("impossible")
	}
	ren.cacheAdd(index, 1, fixcost, varcost, approach)

	SetApproach(ren, index, 1, nil)

	fuzzQuery(t, ren, rnd, index)
}

func randomRename(rnd *rand.Rand, srcCols []string) (from, to []string) {
	if len(srcCols) == 0 {
		return nil, nil
	}

	// Determine how many columns to rename (1 to 3)
	n := 1 + rnd.IntN(min(3, len(srcCols)))

	// Choose random columns to rename
	perm := rnd.Perm(len(srcCols))
	from = make([]string, n)
	for i := range n {
		from[i] = srcCols[perm[i]]
	}

	// Generate new names for the columns
	to = make([]string, n)
	for i := range n {
		// Generate a unique new name that doesn't conflict with existing columns
		for {
			newName := "renamed_" + strconv.Itoa(rnd.IntN(1000))
			if !slices.Contains(srcCols, newName) && !slices.Contains(to[:i], newName) {
				to[i] = newName
				break
			}
		}
	}

	return from, to
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzSummarize ./dbms/query/

var sumStrategyCounts = map[sumStrategy]int{
	sumSeq: 0,
	sumMap: 0,
	sumIdx: 0,
	sumTbl: 0,
}

func FuzzSummarize(f *testing.F) {
	f.Add(uint64(122), uint64(334))
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzSummarize(t, rnd)
	})
}

func TestFuzzSummarize(t *testing.T) {
	rnd := rand.New(rand.NewPCG(12351, 68081))
	for range 1000 {
		fuzzSummarize(t, rnd)
	}
	fmt.Printf("Strategy: seq=%d, map=%d, idx=%d, tbl=%d\n",
		sumStrategyCounts[sumSeq], sumStrategyCounts[sumMap],
		sumStrategyCounts[sumIdx], sumStrategyCounts[sumTbl])
}

var sumOps = []string{"count", "total", "average", "min", "max", "list"}

func fuzzSummarize(t *testing.T, rnd *rand.Rand) {
	qs := NewQuerySource(rnd)
	by, cols, ops, ons := randomSummarize(rnd, qs.ColumnsResult, qs.IndexesResult)
	sum := NewSummarize(qs, "", by, cols, ops, ons)

	index, fixcost, varcost, approach := chooseSummarizeIndex(rnd, sum)
	if fixcost+varcost >= impossible {
		return // skip if no valid index found
	}
	sum.cacheAdd(index, 1, fixcost, varcost, approach)

	SetApproach(sum, index, 1, nil)

	if approach != nil {
		strat := approach.(*summarizeApproach).strat
		sumStrategyCounts[strat]++
	}

	fuzzQuery(t, sum, rnd, index)
}

func chooseSummarizeIndex(rnd *rand.Rand, sum *Summarize) ([]string, Cost, Cost, any) {
	indexes := sum.Indexes()
	candidates := make([][]string, 0, len(indexes)+1)
	candidates = append(candidates, indexes...)
	candidates = append(candidates, nil) // fall back to nil index
	perm := rnd.Perm(len(candidates))
	for _, i := range perm {
		index := candidates[i]
		if len(index) > 0 {
			index = index[:1+rnd.IntN(len(index))]
		}
		fixcost, varcost, approach := sum.optimize(ReadMode, index, 1)
		if fixcost+varcost < impossible {
			return index, fixcost, varcost, approach
		}
	}
	return nil, impossible, impossible, nil
}

func randomSummarize(rnd *rand.Rand, srcCols []string, indexes [][]string) (by, cols, ops, ons []string) {
	// 20% of the time, choose 'by' columns that allow sumSeq
	if len(srcCols) > 0 && len(indexes) > 0 && rnd.IntN(5) == 0 {
		index := random(indexes, rnd)
		if len(index) > 0 {
			prefixLen := 1 + rnd.IntN(len(index))
			by = slices.Clone(index[:prefixLen])
		}
	}
	if by == nil {
		n := rnd.IntN(len(srcCols) + 1) // 0 to all columns
		if n > 0 {
			perm := rnd.Perm(len(srcCols))
			by = make([]string, n)
			for i := range n {
				by[i] = srcCols[perm[i]]
			}
		}
	}

	nops := 1 + rnd.IntN(3)
	cols = make([]string, nops)
	ops = make([]string, nops)
	ons = make([]string, nops)

	// 10% of the time, create conditions for sumIdx (single min/max with no 'by' columns)
	// sumIdx requires the 'on' column to be an index
	if len(by) == 0 && len(indexes) > 0 && rnd.IntN(10) == 0 {
		// Find the first column of an index for sumIdx
		for _, idx := range indexes {
			if len(idx) > 0 {
				cols = make([]string, 1)
				ops = make([]string, 1)
				ons = make([]string, 1)
				if rnd.IntN(2) == 0 {
					ops[0] = "min"
				} else {
					ops[0] = "max"
				}
				ons[0] = idx[0]
				cols[0] = "" // use default name
				return
			}
		}
	}
	if len(by) == 0 && rnd.IntN(7) == 0 {
		// 10% of the time, create conditions for sumTbl (single count with no 'by' columns)
		cols = make([]string, 1)
		ops = make([]string, 1)
		ons = make([]string, 1)
		ops[0] = "count"
		ons[0] = ""
		cols[0] = "" // use default name
	} else {
		for i := range nops {
			ops[i] = sumOps[rnd.IntN(len(sumOps))]
			if ops[i] == "count" {
				ons[i] = ""
			} else {
				ons[i] = srcCols[rnd.IntN(len(srcCols))]
			}
			cols[i] = "" // use default name
		}
	}
	return
}

//-------------------------------------------------------------------

func fuzzQuery(t *testing.T, q Query, rnd *rand.Rand, index []string) {
	t.Helper()
	hdr := q.Header()
	cols := hdr.Columns
	expected := q.Simple(nil)

	verifyNoDuplicates(t, expected, hdr, cols)

	qh := NewQueryHasher(hdr)
	for _, row := range expected {
		qh.Row(row)
	}
	testRandomGet(t, rnd, q, qh, hdr, nil, nil)

	if len(index) != 0 {
		testRandomSelects(t, rnd, q, index, cols, expected)
		testRandomLookups(t, rnd, q, index, cols, expected)
	}
}

func verifyNoDuplicates(t *testing.T, rows []Row, hdr *Header, cols []string) {
	t.Helper()
	seen := make(map[string]bool)
	for _, row := range rows {
		key := rowKey(row, hdr, cols)
		if seen[key] {
			t.Fatal("duplicate row in project result")
		}
		seen[key] = true
	}
}

func rowKey(row Row, hdr *Header, cols []string) string {
	var key string
	for _, col := range cols {
		key += row.GetRaw(hdr, col) + "|"
	}
	return key
}

func testRandomGet(t *testing.T, rnd *rand.Rand, q Query, qh *QueryHash, hdr *Header,
	selCols, selVals []string) {
	t.Helper()

	// Get all rows using Next first to establish correct iteration order
	q.Rewind()
	nextRows := getAllRows(q, Next)
	if !rowSetsEqual(nextRows, qh, hdr) {
		t.Fatalf("Next iteration returned %d rows, expected %d", len(nextRows), qh.nrows)
	}

	// Run deterministic cursor pattern checks before random walk
	testCursorPatterns(t, q, hdr, nextRows)

	data := NewDataSource(nextRows)

	// Redo the Select after getAllRows to reset indexed state for projMap
	q.Select(selCols, selVals)

	// Do a random walk with Next/Prev using nextRows as expected
	nsteps := min(100, len(nextRows)*3)
	for range nsteps {
		// Occasionally add a Select to reset indexed flag for projMap
		if rnd.IntN(20) == 0 { // 5% chance
			if selCols == nil {
				q.Select(nil, nil) // this also rewinds
			} else {
				q.Rewind()
			}
			data.rewind()
		}

		pos := data.pos
		var dir Dir
		switch data.pos {
		case dsAtOrg:
			dir = Next
		case dsAtEnd:
			dir = Prev
		default: // rewound or with -> random dir
			dir = []Dir{Next, Prev}[rnd.IntN(2)]
		}
		expectedRow := data.get(dir)
		row := q.Get(nil, dir)

		if expectedRow == nil && row != nil {
			t.Fatalf("random walk %c from %v: expected nil, got row\nsteps so far: %d\nexpected rows count: %d",
				dir, pos, nsteps-len(nextRows)*3+1, len(nextRows))
		} else if expectedRow != nil && row == nil {
			t.Log(q)
			t.Fatalf("random walk %c from %v: expected row, got nil\nsteps so far: %d\nexpected rows count: %d",
				dir, pos, nsteps-len(nextRows)*3+1, len(nextRows))
		} else if expectedRow != nil && row != nil {
			if !hdr.EqualRows(row, expectedRow, nil, nil) {
				t.Fatalf("random walk %c from %v: row mismatch\nsteps so far: %d\nexpected rows count: %d",
					dir, pos, nsteps-len(nextRows)*3+1, len(nextRows))
			}
		}
	}

	// Get all rows using Prev
	q.Rewind()
	prevRows := getAllRows(q, Prev)
	if !rowSetsEqual(prevRows, qh, hdr) {
		t.Fatalf("Prev iteration returned %d rows, expected %d", len(prevRows), qh.nrows)
	}
}

func getAllRows(q Query, dir Dir) []Row {
	q.Rewind()
	var rows []Row
	for {
		row := q.Get(nil, dir)
		if row == nil {
			break
		}
		rows = append(rows, row)
	}
	return rows
}

// testCursorPatterns runs deterministic cursor navigation patterns.
// These are run before the random walk because failures are clearer -
// they test specific edge cases with known expected behavior.
func testCursorPatterns(t *testing.T, q Query, hdr *Header, nextRows []Row) {
	t.Helper()
	n := len(nextRows)

	check := func(name string, row, expected Row) {
		t.Helper()
		if expected == nil && row != nil {
			t.Fatalf("%s: expected nil, got row", name)
		} else if expected != nil && row == nil {
			t.Fatalf("%s: expected row, got nil", name)
		} else if expected != nil && row != nil {
			if !hdr.EqualRows(row, expected, nil, nil) {
				t.Fatalf("%s: row mismatch", name)
			}
		}
	}

	// Pattern 1: Rewind, Next, Prev - after first Next, Prev should return nil
	q.Rewind()
	row := q.Get(nil, Next) // first row or nil if empty
	if n > 0 {
		check("Next,Prev: N", row, nextRows[0])
	} else {
		check("Next,Prev: N (empty)", row, nil)
	}
	row = q.Get(nil, Prev) // should be nil - nothing before first
	check("Next,Prev: P", row, nil)

	// Pattern 2: Rewind, Prev, Next - Prev from rewind goes to last, then Next should be nil
	if n > 0 {
		q.Rewind()
		row = q.Get(nil, Prev) // last row
		check("Prev,Next: P", row, nextRows[n-1])
		row = q.Get(nil, Next) // should be nil - nothing after last
		check("Prev,Next: N", row, nil)
	}

	// Pattern 3: Rewind, Prev, Prev, Next, Next
	if n >= 2 {
		q.Rewind()
		row = q.Get(nil, Prev) // last row (n-1)
		check("PPNN: P1", row, nextRows[n-1])
		row = q.Get(nil, Prev) // second to last (n-2)
		check("PPNN: P2", row, nextRows[n-2])
		row = q.Get(nil, Next) // back to last (n-1)
		check("PPNN: N1", row, nextRows[n-1])
		row = q.Get(nil, Next) // should be nil
		check("PPNN: N2", row, nil)
	}

	// Pattern 4: Rewind, Next, Next, Prev, Prev
	if n >= 2 {
		q.Rewind()
		row = q.Get(nil, Next) // first row (0)
		check("NNPP: N1", row, nextRows[0])
		row = q.Get(nil, Next) // second row (1)
		check("NNPP: N2", row, nextRows[1])
		row = q.Get(nil, Prev) // back to first (0)
		check("NNPP: P1", row, nextRows[0])
		row = q.Get(nil, Prev) // should be nil
		check("NNPP: P2", row, nil)
	}

	// Pattern 5: Next to end, past end (nil), then Prev
	if n > 0 {
		q.Rewind()
		for i := 0; i < n; i++ {
			row = q.Get(nil, Next)
			check("ToEnd: N"+strconv.Itoa(i), row, nextRows[i])
		}
		row = q.Get(nil, Next) // past end
		check("ToEnd: N-past", row, nil)
		row = q.Get(nil, Prev) // should return last row
		check("ToEnd: P", row, nextRows[n-1])
	}

	// Pattern 6: Prev to beginning, past beginning (nil), then Next
	if n > 0 {
		q.Rewind()
		for i := n - 1; i >= 0; i-- {
			row = q.Get(nil, Prev)
			check("ToBegin: P"+strconv.Itoa(n-1-i), row, nextRows[i])
		}
		row = q.Get(nil, Prev) // past beginning
		check("ToBegin: P-past", row, nil)
		row = q.Get(nil, Next) // should return first row
		check("ToBegin: N", row, nextRows[0])
	}

	// Pattern 7: Rewind, Next, Prev, Next - should get first row twice
	if n > 0 {
		q.Rewind()
		row = q.Get(nil, Next) // first
		check("NPN: N1", row, nextRows[0])
		row = q.Get(nil, Prev) // nil
		check("NPN: P", row, nil)
		row = q.Get(nil, Next) // first again
		check("NPN: N2", row, nextRows[0])
	}

	// Pattern 8: Rewind, Prev, Next, Prev - should get last row twice
	if n > 0 {
		q.Rewind()
		row = q.Get(nil, Prev) // last
		check("PNP: P1", row, nextRows[n-1])
		row = q.Get(nil, Next) // nil
		check("PNP: N", row, nil)
		row = q.Get(nil, Prev) // last again
		check("PNP: P2", row, nextRows[n-1])
	}

	// Reset for subsequent tests
	q.Rewind()
}

func rowSetsEqual(a []Row, qh *QueryHash, hdr *Header) bool {
	if len(a) != qh.nrows {
		return false
	}

	// Use QueryHash for efficient comparison instead of O(N^2) approach
	qh2 := NewQueryHasher(hdr)
	for _, row := range a {
		qh2.Row(row)
	}

	// Compare the final hash values
	return qh2.hash == qh.hash && qh2.nrows == qh.nrows
}

func rowsEqual(a, b Row, hdr *Header, cols []string) bool {
	for _, col := range cols {
		if a.GetRaw(hdr, col) != b.GetRaw(hdr, col) {
			return false
		}
	}
	return true
}

//-------------------------------------------------------------------

func testRandomSelects(t *testing.T, rnd *rand.Rand, q Query, index, cols []string, allRows []Row) {
	t.Helper()
	hdr := q.Header()
	testExistentSelect(t, allRows, rnd, hdr, index, cols, q)
	testNonExistentSelect(t, allRows, rnd, hdr, index, cols, q)
}

func testExistentSelect(t *testing.T, allRows []Row, rnd *rand.Rand, hdr *Header, index []string, cols []string, q Query) {
	if len(allRows) == 0 {
		return
	}
	for range 10 {
		srcRow := allRows[rnd.IntN(len(allRows))]
		selCols, selVals := indexSelectCriteria(rnd, srcRow, hdr, index, cols, false)
		q.Select(selCols, selVals)

		qh := NewQueryHasher(hdr)
		for _, row := range allRows {
			if selMatch(hdr, row, selCols, selVals) {
				qh.Row(row)
			}
		}

		testRandomGet(t, rnd, q, qh, hdr, selCols, selVals)

		q.Select(nil, nil) // clear select
	}
}

func selMatch(hdr *Header, row Row, selCols, selVals []string) bool {
	for i, col := range selCols {
		if row.GetRaw(hdr, col) != selVals[i] {
			return false
		}
	}
	return true
}

// indexSelectCriteria picks a random prefix of the index for select criteria.
func indexSelectCriteria(rnd *rand.Rand, row Row, hdr *Header, index, cols []string, nonexist bool) ([]string, []string) {
	selCols := slices.Clone(index)
	n := 1 + rnd.IntN(len(selCols))
	selCols = selCols[:n]
	// Add extra columns, but avoid duplicates
	perm := rnd.Perm(len(cols))
	for _, i := range perm {
		col := cols[i]
		if !slices.Contains(selCols, col) {
			selCols = append(selCols, col)
			if len(selCols) >= n+2 {
				break
			}
		}
	}
	rnd.Shuffle(len(selCols), func(i, j int) {
		selCols[i], selCols[j] = selCols[j], selCols[i]
	})

	selVals := make([]string, len(selCols))
	for i, col := range selCols {
		selVals[i] = row.GetRaw(hdr, col)
	}
	if nonexist {
		selVals[rnd.IntN(len(selVals))] = "nonexistent"
	}
	return selCols, selVals
}

func testNonExistentSelect(t *testing.T, allRows []Row, rnd *rand.Rand, hdr *Header, index []string, cols []string, q Query) {
	for range 10 {
		var srcRow = Row{DbRec{Record: Record("\x00")}}
		if len(allRows) > 0 {
			srcRow = allRows[rnd.IntN(len(allRows))]
		}
		selCols, selVals := indexSelectCriteria(rnd, srcRow, hdr, index, cols, true)
		q.Select(selCols, selVals)
		if q.Get(nil, Next) != nil {
			t.Fatal("non-existent select returned a row")
		}
		q.Select(nil, nil) // clear select
	}
}

//-------------------------------------------------------------------

func testRandomLookups(t *testing.T, rnd *rand.Rand, q Query, index, cols []string, allRows []Row) {
	t.Helper()

	// Only test lookups if the index is one of the query's keys
	if !slices.ContainsFunc(q.Keys(),
		func(k []string) bool { return slices.Equal(k, index) }) {
		return
	}

	lookupCols := slices.Clone(index)
	// Add extra columns, but avoid duplicates
	perm := rnd.Perm(len(cols))
	for _, i := range perm {
		col := cols[i]
		if !slices.Contains(lookupCols, col) {
			lookupCols = append(lookupCols, col)
			if len(lookupCols) >= len(index)+2 {
				break
			}
		}
	}
	rnd.Shuffle(len(lookupCols), func(i, j int) {
		lookupCols[i], lookupCols[j] = lookupCols[j], lookupCols[i]
	})

	testExistentLookup(t, allRows, rnd, lookupCols, q, cols)
	testNonExistentLookup(t, rnd, q, lookupCols)
}

func testExistentLookup(t *testing.T, allRows []Row, rnd *rand.Rand, lookupCols []string, q Query, cols []string) {
	t.Helper()
	if len(allRows) == 0 {
		return
	}
	hdr := q.Header()
	for range min(10, len(allRows)) {
		srcRow := allRows[rnd.IntN(len(allRows))]

		lookupVals := make([]string, len(lookupCols))
		for i, col := range lookupCols {
			lookupVals[i] = srcRow.GetRaw(hdr, col)
		}

		result := q.Lookup(nil, lookupCols, lookupVals)

		if result == nil {
			t.Fatal("lookup returned nil for existing key")
		}

		for i, col := range lookupCols {
			if result.GetRaw(hdr, col) != lookupVals[i] {
				t.Fatalf("lookup result doesn't match key: col=%s", col)
			}
		}

		// Verify result is one of allRows
		found := false
		for _, ar := range allRows {
			if rowsEqual(result, ar, hdr, cols) {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("lookup result not in all rows")
		}
	}
}

func testNonExistentLookup(t *testing.T, rnd *rand.Rand, q Query, lookupCols []string) {
	t.Helper()
	for range 10 {
		keyVals := make([]string, len(lookupCols))
		// set one of the keyVals to a non-existent value
		// the others to possibly existing values
		r := rnd.IntN(len(lookupCols))
		for i, col := range lookupCols {
			if i == r {
				keyVals[i] = "nonexistent"
			} else {
				keyVals[i] = col + "_" + strconv.Itoa(rnd.IntN(100))
			}
		}
		result := q.Lookup(nil, lookupCols, keyVals)
		if result != nil {
			t.Fatal("lookup returned row for non-existent key")
		}
	}
}

func random[E any](list []E, rnd *rand.Rand) E {
	return list[rnd.IntN(len(list))]
}
