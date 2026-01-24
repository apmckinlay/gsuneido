// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/core"
)

func init() {
	sortForTest = true
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzRandom ./dbms/query/

func FuzzRandom(f *testing.F) {
	f.Add(uint64(123), uint64(456))
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzRandom(t, rnd)
	})
}

func TestFuzzRandom(t *testing.T) {
	for range 1000 {
		seed1, seed2 := rand.Uint64(), rand.Uint64()
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzRandom(t, rnd)
	}
}

func fuzzRandom(t *testing.T, rnd *rand.Rand) {
	fuzzers := []func(*testing.T, *rand.Rand){
		fuzzQuerySource,
		fuzzTempIndex,
		fuzzRename,
		fuzzExtend,
		fuzzProject,
		fuzzSummarize,
		fuzzMinus,
		fuzzIntersect,
		fuzzUnion,
	}
	f := random(fuzzers, rnd)
	f(t, rnd)
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzQuerySource ./dbms/query/

func FuzzQuerySource(f *testing.F) {
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzQuerySource(t, rnd)
	})
}

func TestFuzzQuerySource(t *testing.T) {
	for range 1000 {
		seed1, seed2 := rand.Uint64(), rand.Uint64()
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzQuerySource(t, rnd)
	}
}

func fuzzQuerySource(t *testing.T, rnd *rand.Rand) {
	qs := NewQuerySource(rnd)
	index := chooseIndex(rnd, qs)
	qs.setApproach(index, 1, nil, nil)
	fuzzQuery(t, qs, rnd, index)
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
	for range 1000 {
		seed1, seed2 := rand.Uint64(), rand.Uint64()
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzProject(t, rnd)
	}
}

func fuzzProject(t *testing.T, rnd *rand.Rand) {
	qs := NewQuerySource(rnd)
	projCols := randomProjectCols(rnd, qs.ColumnsResult, qs.IndexesResult)
	q := NewProject(qs, projCols)
	index := chooseIndex(rnd, q)
	fuzzQuery(t, q, rnd, index)
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

func chooseIndex(rnd *rand.Rand, source Query) []string {
	if isEmptyKey(source.Keys()) {
		// choose random columns
		cols := source.Header().Columns
		if len(cols) == 0 {
			return nil
		}
		n := 1 + rnd.IntN(len(cols)) // 1 to all columns
		perm := rnd.Perm(len(cols))
		randomCols := make([]string, n)
		for i := range n {
			randomCols[i] = cols[perm[i]]
		}
		return randomCols
	}
	indexes := source.Indexes()
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
	for range 1000 {
		seed1, seed2 := rand.Uint64(), rand.Uint64()
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzRename(t, rnd)
	}
}

func fuzzRename(t *testing.T, rnd *rand.Rand) {
	qs := NewQuerySource(rnd)
	from, to := randomRename(rnd, qs.ColumnsResult)
	q := NewRename(qs, from, to)
	index := chooseIndex(rnd, q)
	fuzzQuery(t, q, rnd, index)
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

func FuzzSummarize(f *testing.F) {
	f.Add(uint64(122), uint64(334))
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzSummarize(t, rnd)
	})
}

func TestFuzzSummarize(t *testing.T) {
	rnd := rand.New(rand.NewPCG(1091395294133611146, 8719992948325563695))
	fuzzSummarize(t, rnd)
	for range 1000 {
		seed1, seed2 := rand.Uint64(), rand.Uint64()
		// fmt.Printf("%d, %d\n", seed1, seed2)
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzSummarize(t, rnd)
	}
}

var sumOps = []string{"count", "total", "average", "min", "max", "list"}

func fuzzSummarize(t *testing.T, rnd *rand.Rand) {
	qs := NewQuerySource(rnd)
	by, cols, ops, ons := randomSummarize(rnd, qs.ColumnsResult, qs.IndexesResult)
	q := NewSummarize(qs, "", by, cols, ops, ons)
	index := chooseIndex(rnd, q)
	fuzzQuery(t, q, rnd, index)
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
			ops[i] = random(sumOps, rnd)
			if ops[i] == "count" {
				ons[i] = ""
			} else {
				ons[i] = random(srcCols, rnd)
			}
			cols[i] = "" // use default name
		}
	}
	return
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzMinus ./dbms/query/

func FuzzMinus(f *testing.F) {
	f.Add(uint64(122), uint64(334))
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzMinus(t, rnd)
	})
}

func TestFuzzMinus(t *testing.T) {
	for range 1000 {
		seed1, seed2 := rand.Uint64(), rand.Uint64()
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzMinus(t, rnd)
	}
}

func fuzzMinus(t *testing.T, rnd *rand.Rand) {
	qs1, qs2 := NewCompatibleQS(rnd)
	// fmt.Printf("minus %d %d = ", len(qs1.rows), len(qs2.rows))
	q := NewMinus(qs1, qs2)
	index := chooseIndex(rnd, q)
	fuzzQuery(t, q, rnd, index)
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzIntersect ./dbms/query/

func FuzzIntersect(f *testing.F) {
	f.Add(uint64(122), uint64(334))
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzIntersect(t, rnd)
	})
}

func TestFuzzIntersect(t *testing.T) {
	for range 1000 {
		seed1, seed2 := rand.Uint64(), rand.Uint64()
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzIntersect(t, rnd)
	}
}

func fuzzIntersect(t *testing.T, rnd *rand.Rand) {
	qs1, qs2 := NewCompatibleQS(rnd)
	// fmt.Print("intersect ", len(qs1.rows), " ", len(qs2.rows), " = ")
	q := NewIntersect(qs1, qs2)
	index := chooseIndex(rnd, q)
	fixcost, varcost := Optimize(q, ReadMode, index, 1)
	if fixcost+varcost >= impossible {
		return
	}
	fuzzQuery(t, q, rnd, index)
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzUnion ./dbms/query/

func FuzzUnion(f *testing.F) {
	f.Add(uint64(122), uint64(334))
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzUnion(t, rnd)
	})
}

func TestFuzzUnion(t *testing.T) {
	for range 1000 {
		seed1, seed2 := rand.Uint64(), rand.Uint64()
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzUnion(t, rnd)
	}
}

func fuzzUnion(t *testing.T, rnd *rand.Rand) {
	qs1, qs2 := NewCompatibleQS(rnd)
	q := NewUnion(qs1, qs2)
	index := chooseIndex(rnd, q)
	fuzzQuery(t, q, rnd, index)
}

//-------------------------------------------------------------------

func NewCompatibleQS(rnd *rand.Rand) (*QuerySource, *QuerySource) {
	qs := newQuerySource(rnd, 199, 5, 5)
	// temporarily remove keys from indexes (restore them after changes)
	qs.IndexesResult = qs.IndexesResult[:len(qs.KeysResult)]
	qs1 := *qs
	qs2 := *qs

	qs1.rows, qs2.rows = splitShare(rnd, qs.rows)
	qs2.rows = slices.Clone(qs2.rows) // so they don't share
	if len(qs1.rows) > 100 {
		qs1.rows = qs1.rows[:100]
	}
	qs1.NrowsN = len(qs1.rows)
	qs1.NrowsP = len(qs1.rows)
	if len(qs2.rows) > 100 {
		qs2.rows = qs2.rows[:100]
	}
	qs2.NrowsN = len(qs2.rows)
	qs2.NrowsP = len(qs2.rows)

	qs1.IndexesResult, qs2.IndexesResult = splitShare(rnd, qs.IndexesResult)

	qs1.KeysResult, qs2.KeysResult = splitShare(rnd, qs.KeysResult)
	// ensure at least one key in each
	if len(qs1.KeysResult) == 0 {
		qs1.KeysResult = append(qs1.KeysResult, random(qs.KeysResult, rnd))
	}
	if len(qs2.KeysResult) == 0 {
		qs2.KeysResult = append(qs2.KeysResult, random(qs.KeysResult, rnd))
	}

	// keep the original columns on both to ensure indexes are valid
	// and add some new ones
	qs1.ColumnsResult = slices.Clip(qs.ColumnsResult)
	i := len(qs.ColumnsResult)
	for range rnd.IntN(7) {
		col := "c" + strconv.Itoa(i)
		qs1.ColumnsResult = append(qs1.ColumnsResult, col)
		i++
	}
	qs1.HeaderResult = SimpleHeader(qs1.ColumnsResult)

	qs2.ColumnsResult = slices.Clip(qs.ColumnsResult)
	i = len(qs.ColumnsResult)
	for range rnd.IntN(7) {
		col := "c" + strconv.Itoa(i)
		qs2.ColumnsResult = append(qs2.ColumnsResult, col)
		i++
	}
	qs2.HeaderResult = SimpleHeader(qs2.ColumnsResult)

	// add the keys back to the indexes
	qs1.IndexesResult = append(qs1.IndexesResult, qs1.KeysResult...)
	qs2.IndexesResult = append(qs2.IndexesResult, qs2.KeysResult...)

	return &qs1, &qs2
}

// splitShare splits a slice into three parts and returns two slices
// one contains part 1 and 2, the other contains part 2 and 3
func splitShare[E any](rnd *rand.Rand, s []E) ([]E, []E) {
	n := len(s)
	if n < 3 {
		return s, s
	}
	a := rnd.IntN(n + 1)
	b := rnd.IntN(n + 1)
	if a > b {
		a, b = b, a
	}
	return slices.Clip(s[:b]), slices.Clip(s[a:])
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzExtend ./dbms/query/

func FuzzExtend(f *testing.F) {
	f.Add(uint64(122), uint64(334))
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzExtend(t, rnd)
	})
}

func TestFuzzExtend(t *testing.T) {
	for range 1000 {
		seed1, seed2 := rand.Uint64(), rand.Uint64()
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzExtend(t, rnd)
	}
}

func fuzzExtend(t *testing.T, rnd *rand.Rand) {
	qs := NewQuerySource(rnd)
	cols, exprs := randomExtend(rnd, qs.ColumnsResult)
	ext := NewExtend(qs, cols, exprs)

	index := chooseIndex(rnd, ext)
	fixcost, varcost := Optimize(ext, ReadMode, index, 1)
	if fixcost+varcost >= impossible {
		t.Fatal("impossible")
	}

	q := SetApproach(ext, index, 1, nil)

	fuzzQuery(t, q, rnd, index)
}

func randomExtend(rnd *rand.Rand, srcCols []string) (cols []string, exprs []ast.Expr) {
	if len(srcCols) == 0 {
		return nil, nil
	}

	// Keep this simple: we are fuzzing cursor behavior and query plumbing,
	// not expression evaluation.
	n := rnd.IntN(min(3, len(srcCols)) + 1) // 0-3 (or fewer if few cols)
	if n == 0 {
		return nil, nil
	}

	cols = make([]string, n)
	exprs = make([]ast.Expr, n)
	for i := range n {
		// unique new column name, not in source and not duplicate within extend
		for {
			c := "x" + strconv.Itoa(rnd.IntN(1000))
			if !slices.Contains(srcCols, c) && !slices.Contains(cols[:i], c) {
				cols[i] = c
				break
			}
		}
		if rnd.IntN(2) == 0 {
			exprs[i] = &ast.Constant{Val: IntVal(rnd.IntN(1000))}
		} else {
			exprs[i] = &ast.Ident{Name: random(srcCols, rnd)}
		}
	}
	return cols, exprs
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzTempIndex ./dbms/query/

func FuzzTempIndex(f *testing.F) {
	f.Add(uint64(122), uint64(334))
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzTempIndex(t, rnd)
	})
}

func TestFuzzTempIndex(t *testing.T) {
	for range 1000 {
		seed1, seed2 := rand.Uint64(), rand.Uint64()
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzTempIndex(t, rnd)
	}
}

func fuzzTempIndex(t *testing.T, rnd *rand.Rand) {
	qs := NewQuerySource(rnd)
	if len(qs.ColumnsResult) == 0 {
		return
	}
	order := randomOrder(rnd, qs.ColumnsResult)
	if len(order) == 0 {
		return
	}
	ti := NewTempIndex(qs, order, nil)
	fuzzQuery(t, ti, rnd, order)
}

func randomOrder(rnd *rand.Rand, cols []string) []string {
	n := 1 + rnd.IntN(min(3, len(cols)))
	perm := rnd.Perm(len(cols))
	order := make([]string, n)
	for i := range n {
		order[i] = cols[perm[i]]
	}
	return order
}

//-------------------------------------------------------------------

func fuzzQuery(t *testing.T, q Query, rnd *rand.Rand, index []string) {
	fixcost, varcost := Optimize(q, ReadMode, index, 1)
	if fixcost+varcost >= impossible {
		t.Fatal("impossible\n", format(0, q, 0))
	}
	q = SetApproach(q, index, 1, nil)
	// fmt.Println(format(0, q, 0))
	// fmt.Println("index", index)

	hdr := q.Header()
	cols := hdr.Columns
	expected := q.Simple(nil)
	// fmt.Println(len(expected))

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
			dir = random([]Dir{Next, Prev}, rnd)
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
		srcRow := random(allRows, rnd)
		selCols, selVals := indexSelectCriteria(rnd, srcRow, hdr, index, cols, false)
		q.Select(selCols, selVals)

		// Select only filters by index columns, not extra columns
		qh := NewQueryHasher(hdr)
		for _, row := range allRows {
			if selMatchIndex(hdr, row, selCols, selVals, index) {
				qh.Row(row)
			}
		}

		testRandomGet(t, rnd, q, qh, hdr, selCols, selVals)

		q.Select(nil, nil) // clear select
	}
}

// selMatchIndex checks only the index columns, not extra columns.
// It iterates the index in order and stops at the first missing column,
// matching the behavior of selKeys/TempIndex.makeKey.
func selMatchIndex(hdr *Header, row Row, selCols, selVals, index []string) bool {
	for _, col := range index {
		i := slices.Index(selCols, col)
		if i == -1 {
			break // stop at first missing column (matches selKeys behavior)
		}
		if row.GetRaw(hdr, col) != selVals[i] {
			return false
		}
	}
	return true
}

// indexSelectCriteria picks a random prefix of the index for select criteria.
func indexSelectCriteria(rnd *rand.Rand, row Row, hdr *Header, index, cols []string, nonexist bool) ([]string, []string) {
	n := 1 + rnd.IntN(len(index))
	selCols := slices.Clone(index[:n])
	_ = cols
	// Add extra columns, but avoid duplicates
	// perm := rnd.Perm(len(cols))
	// for _, i := range perm {
	// 	col := cols[i]
	// 	if !slices.Contains(selCols, col) {
	// 		selCols = append(selCols, col)
	// 		if len(selCols) >= n+2 {
	// 			break
	// 		}
	// 	}
	// }
	rnd.Shuffle(len(selCols), func(i, j int) {
		selCols[i], selCols[j] = selCols[j], selCols[i]
	})

	selVals := make([]string, len(selCols))
	for i, col := range selCols {
		selVals[i] = row.GetRaw(hdr, col)
	}
	if nonexist {
		// Must set nonexistent on an index column, not an extra column
		col := index[rnd.IntN(n)]
		i := slices.Index(selCols, col)
		selVals[i] = "nonexistent"
	}
	return selCols, selVals
}

func testNonExistentSelect(t *testing.T, allRows []Row, rnd *rand.Rand, hdr *Header, index []string, cols []string, q Query) {
	for range 10 {
		// If there are no rows, use a dummy row sized to match hdr.Fields.
		// This avoids panics when hdr.Fields references derived records (e.g. Extend).
		srcRow := make(Row, len(hdr.Fields))
		if len(allRows) > 0 {
			srcRow = random(allRows, rnd)
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
		srcRow := random(allRows, rnd)

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

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzSplitShare ./dbms/query/

func FuzzSplitShare(f *testing.F) {
	f.Add(uint64(123), uint64(456))
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		fuzzSplitShare(t, rnd)
	})
}

func TestFuzzSplitShare(t *testing.T) {
	stats := struct {
		part1Empty int
		part2Empty int
		part3Empty int
		total      int
	}{}

	for range 10000 {
		seed1, seed2 := rand.Uint64(), rand.Uint64()
		rnd := rand.New(rand.NewPCG(seed1, seed2))
		part1Empty, part2Empty, part3Empty := fuzzSplitShare(t, rnd)
		stats.total++
		if part1Empty {
			stats.part1Empty++
		}
		if part2Empty {
			stats.part2Empty++
		}
		if part3Empty {
			stats.part3Empty++
		}
	}

	t.Logf("splitShare stats (n=%d): part1 empty=%d (%.1f%%), part2 empty=%d (%.1f%%), part3 empty=%d (%.1f%%)",
		stats.total,
		stats.part1Empty, float64(stats.part1Empty)/float64(stats.total)*100,
		stats.part2Empty, float64(stats.part2Empty)/float64(stats.total)*100,
		stats.part3Empty, float64(stats.part3Empty)/float64(stats.total)*100)

	if stats.part1Empty == 0 {
		t.Error("part1 was never empty")
	}
	if stats.part2Empty == 0 {
		t.Error("part2 was never empty")
	}
	if stats.part3Empty == 0 {
		t.Error("part3 was never empty")
	}
}

func fuzzSplitShare(t *testing.T, rnd *rand.Rand) (part1Empty, part2Empty, part3Empty bool) {
	t.Helper()

	n := rnd.IntN(101)
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}

	result1, result2 := splitShare(rnd, s)

	len1 := len(result1)
	len2 := len(result2)

	part2len := (len1 + len2) - n
	part2Empty = part2len == 0
	part1Empty = len2 == n
	part3Empty = len1 == n

	return
}
