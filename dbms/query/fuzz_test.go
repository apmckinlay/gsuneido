// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/bits"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
)

func init() {
	sortForTest = true
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
}

const nfuzz = 200

type fuzzRunner struct {
	build func(*FT) Query
}

func (fr fuzzRunner) Run(t *testing.T, seed1, seed2 uint64) {
	defer func(jr int) { joinRev = jr }(joinRev)
	joinRev = impossible
	defer func(ti int) { ticostAdj = ti }(ticostAdj)
	ticostAdj = 9999999
	ft := newFT(seed1, seed2)
	defer ft.db.Close()
	q := fr.build(ft)
	fuzzQuery(t, q, ft)
}

func (fr fuzzRunner) Fuzz(f *testing.F) {
	f.Add(uint64(122), uint64(334))
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		fr.Run(t, seed1, seed2)
	})
}

func (fr fuzzRunner) Test(t *testing.T) {
	var seed1, seed2 uint64
	defer func() {
		if r := recover(); r != nil || t.Failed() {
			fmt.Printf("failing seed: %d, %d\n", seed1, seed2)
			if r != nil {
				panic(r)
			}
		}
	}()
	for range nfuzz {
		seed1, seed2 = rand.Uint64(), rand.Uint64()
		fr.Run(t, seed1, seed2)
	}
	fmt.Println("tempindex", tempIndexCount.Load())
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzRandom ./dbms/query/

var fuzzRandomRunner = fuzzRunner{build: fuzzRandom}

func FuzzRandom(f *testing.F) {
	fuzzRandomRunner.Fuzz(f)
}

func TestFuzzRandomDebug(t *testing.T) {
	fuzzRandomRunner.Run(t, 120, 291)
}

func TestFuzzRandom(t *testing.T) {
	fuzzRandomRunner.Test(t)
}

func fuzzRandom(ft *FT) Query {
	builders := []func(*FT) Query{
		fuzzTable,
		fuzzProject,
		fuzzRename,
		fuzzExtend,
		fuzzSummarize,
		fuzzWhere,
		fuzzMinus,
		fuzzIntersect,
		fuzzUnion,
		fuzzTimes,
		fuzzJoin,
		fuzzLeftJoin,
		fuzzSemiJoin,
	}
	composers := []func(*FT, Query) Query{
		composeFuzzProject,
		composeFuzzRename,
		composeFuzzExtend,
		composeFuzzSummarize,
		composeFuzzWhere,
	}
	if ft.rnd.IntN(3) == 0 {
		return random(builders, ft.rnd)(ft)
	}
	inner := random(builders, ft.rnd)
	outer := random(composers, ft.rnd)
	return outer(ft, inner(ft))
}

//-------------------------------------------------------------------

func TestFuzzNothing(t *testing.T) {
	ft := testFT()
	defer ft.db.Close()
	q := &Nothing{table: "nothing"}
	q.header = SimpleHeader([]string{"a", "b", "c"})
	for range nfuzz {
		fuzzQuery(t, q, ft)
	}
}

func TestFuzzProjectNone(t *testing.T) {
	ft := testFT()
	defer ft.db.Close()
	empty := &Nothing{table: "nothing"}
	empty.header = SimpleHeader([]string{})
	q := &ProjectNone{source: empty}
	for range nfuzz {
		fuzzQuery(t, q, ft)
	}
	tbl := ft.NewFuzzTable()
	q = &ProjectNone{source: tbl}
	for range nfuzz {
		fuzzQuery(t, q, ft)
	}
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzTable ./dbms/query/

var fuzzTableRunner = fuzzRunner{build: fuzzTable}

func fuzzTable(ft *FT) Query {
	return ft.NewFuzzTable()
}

func FuzzTable(f *testing.F) {
	fuzzTableRunner.Fuzz(f)
}

func TestFuzzTable(t *testing.T) {
	fuzzTableRunner.Test(t)
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzProject ./dbms/query/

var fuzzProjectRunner = fuzzRunner{build: fuzzProject}

func FuzzProject(f *testing.F) {
	fuzzProjectRunner.Fuzz(f)
}

func TestFuzzProjectDebug(t *testing.T) {
	fuzzProjectRunner.Run(t, 4886698708123789290, 16491253703327079940)
}

func TestFuzzProject(t *testing.T) {
	startCopy := projCopyCount.Load()
	startSeq := projSeqCount.Load()
	startMap := projMapCount.Load()

	fuzzProjectRunner.Test(t)

	deltaCopy := projCopyCount.Load() - startCopy
	deltaSeq := projSeqCount.Load() - startSeq
	deltaMap := projMapCount.Load() - startMap
	fmt.Printf("Project strategies: copy=%d seq=%d map=%d\n", deltaCopy, deltaSeq, deltaMap)
	if nfuzz >= 1000 {
		if deltaCopy == 0 {
			t.Error("projCopy strategy not used")
		}
		if deltaSeq == 0 {
			t.Error("projSeq strategy not used")
		}
		if deltaMap == 0 {
			t.Error("projMap strategy not used")
		}
	}
}

func fuzzProject(ft *FT) Query {
	return composeFuzzProject(ft, ft.NewFuzzTable())
}

func composeFuzzProject(ft *FT, qs Query) Query {
	if len(qs.Columns()) == 0 {
		return qs
	}
	projCols := randomProjectCols(ft.rnd, qs.Columns(), qs.Indexes())
	return NewProject(qs, projCols)
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

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzRename ./dbms/query/

var fuzzRenameRunner = fuzzRunner{build: fuzzRename}

func FuzzRename(f *testing.F) {
	fuzzRenameRunner.Fuzz(f)
}

func TestFuzzRenameDebug(t *testing.T) {
	fuzzRenameRunner.Run(t, 2736498751574507473, 11100617320412793980)
}

func TestFuzzRename(t *testing.T) {
	fuzzRenameRunner.Test(t)
}

func fuzzRename(ft *FT) Query {
	return composeFuzzRename(ft, ft.NewFuzzTable())
}

func composeFuzzRename(ft *FT, qs Query) Query {
	from, to := randomRename(ft.rnd, qs.Columns())
	return NewRename(qs, from, to)
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
			newName := "renamed_" + strconv.Itoa(rnd.IntN(nfuzz))
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

var fuzzSummarizeRunner = fuzzRunner{build: fuzzSummarize}

func FuzzSummarize(f *testing.F) {
	fuzzSummarizeRunner.Fuzz(f)
}

func TestFuzzSummarizeDebug(t *testing.T) {
	fuzzSummarizeRunner.Run(t, 18388908088726648744, 13779681344780478899)
}

func TestFuzzSummarize(t *testing.T) {
	startSeq := sumSeqCount.Load()
	startMap := sumMapCount.Load()
	startIdx := sumIdxCount.Load()
	startTbl := sumTblCount.Load()
	startUnique := sumUniqueCount.Load()
	startWholeRow := sumWholeRowCount.Load()

	fuzzSummarizeRunner.Test(t)

	deltaSeq := sumSeqCount.Load() - startSeq
	deltaMap := sumMapCount.Load() - startMap
	deltaIdx := sumIdxCount.Load() - startIdx
	deltaTbl := sumTblCount.Load() - startTbl
	deltaUnique := sumUniqueCount.Load() - startUnique
	deltaWholeRow := sumWholeRowCount.Load() - startWholeRow
	fmt.Printf("Summarize strategies: seq=%d map=%d idx=%d tbl=%d unique=%d wholerow=%d\n",
		deltaSeq, deltaMap, deltaIdx, deltaTbl, deltaUnique, deltaWholeRow)
	if nfuzz >= 1000 {
		if deltaSeq == 0 {
			t.Error("sumSeq strategy not used")
		}
		if deltaMap == 0 {
			t.Error("sumMap strategy not used")
		}
		if deltaIdx == 0 {
			t.Error("sumIdx strategy not used")
		}
		if deltaTbl == 0 {
			t.Error("sumTbl strategy not used")
		}
		if deltaUnique == 0 {
			t.Error("sumUnique variation not used")
		}
		if deltaWholeRow == 0 {
			t.Error("sumWholeRow variation not used")
		}
	}
}
func fuzzSummarize(ft *FT) Query {
	return composeFuzzSummarize(ft, ft.NewFuzzTable())
}

func composeFuzzSummarize(ft *FT, qs Query) Query {
	if len(qs.Columns()) == 0 {
		return NewSummarize(qs, "", nil, []string{""}, []string{"count"}, []string{""})
	}
	by, cols, ops, ons := randomSummarize(ft.rnd, qs.Columns(), qs.Indexes())
	return NewSummarize(qs, "", by, cols, ops, ons)
}

var sumOps = []string{"count", "total", "average", "min", "max"}

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

var fuzzMinusRunner = fuzzRunner{build: fuzzMinus}

func fuzzMinus(ft *FT) Query {
	q1, q2 := newCompatibleQS(ft)
	return NewMinus(q1, q2, ft.rt)
}

func FuzzMinus(f *testing.F) {
	fuzzMinusRunner.Fuzz(f)
}

func TestFuzzMinus(t *testing.T) {
	fuzzMinusRunner.Test(t)
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzIntersect ./dbms/query/

var fuzzIntersectRunner = fuzzRunner{build: fuzzIntersect}

func fuzzIntersect(ft *FT) Query {
	q1, q2 := newCompatibleQS(ft)
	return NewIntersect(q1, q2, ft.rt)
}

func FuzzIntersect(f *testing.F) {
	fuzzIntersectRunner.Fuzz(f)
}

func TestFuzzIntersect(t *testing.T) {
	fuzzIntersectRunner.Test(t)
}

func TestFuzzIntersectDebug(t *testing.T) {
	fuzzIntersectRunner.Run(t, 15551282355907782167, 17075134520393906833)
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzUnion ./dbms/query/

var fuzzUnionRunner = fuzzRunner{build: fuzzUnion}

func FuzzUnion(f *testing.F) {
	fuzzUnionRunner.Fuzz(f)
}

func TestFuzzUnion(t *testing.T) {
	startMerge := unionMergeCount.Load()
	startLookup := unionLookupCount.Load()
	startDisjoint := unionDisjointCount.Load()
	startMergeDisjoint := unionMergeDisjoint.Load()

	fuzzUnionRunner.Test(t)

	deltaMerge := unionMergeCount.Load() - startMerge
	deltaLookup := unionLookupCount.Load() - startLookup
	deltaDisjoint := unionDisjointCount.Load() - startDisjoint
	deltaMergeDisjoint := unionMergeDisjoint.Load() - startMergeDisjoint
	fmt.Printf("Union strategies: merge=%d lookup=%d disjoint=%d merge-disjoint=%d\n",
		deltaMerge, deltaLookup, deltaDisjoint, deltaMergeDisjoint)
	if nfuzz >= 1000 {
		if deltaMerge == 0 {
			t.Error("unionMerge strategy not used")
		}
		if deltaLookup == 0 {
			t.Error("unionLookup strategy not used")
		}
		if deltaDisjoint == 0 {
			t.Error("unionDisjoint variation not used")
		}
		if deltaMerge+deltaLookup-deltaDisjoint == 0 {
			t.Error("union non-disjoint variation not used")
		}
	}
}

func fuzzUnion(ft *FT) Query {
	q1, q2 := newCompatibleQS(ft)
	return NewUnion(q1, q2)
}

//-------------------------------------------------------------------

// newCompatibleQS creates QuerySources for Union, Intersect, Minus
func newCompatibleQS(ft *FT) (Query, Query) {
	b := ft.newFT().Sizes(73, 5, 5).construct()
	b1 := *b
	b2 := *b
	rnd := ft.rnd

	b1.data, b2.data = splitShare(rnd, b.data)
	if len(b1.data) > 100 {
		b1.data = b1.data[:100]
	}
	if len(b2.data) > 100 {
		b2.data = b2.data[:100]
	}

	b2.data = slices.Clone(b2.data) // so they don't share
	if len(b1.data) > 100 {
		b1.data = b1.data[:100]
	}
	if len(b2.data) > 100 {
		b2.data = b2.data[:100]
	}

	b1.indexes, b2.indexes = splitShare(rnd, b.indexes)

	b1.keys, b2.keys = splitShare(rnd, b.keys)

	// 10% of the time, force empty keys
	switch rnd.IntN(19) {
	case 7:
		makeEmptyKey(rnd, &b1)
	case 11:
		makeEmptyKey(rnd, &b2)
	}

	// ensure at least one key in each
	if len(b1.keys) == 0 {
		b1.keys = append(b1.keys, random(b.keys, rnd))
	}
	if len(b2.keys) == 0 {
		b2.keys = append(b2.keys, random(b.keys, rnd))
	}

	// keep the original columns on both to ensure indexes are valid
	// and add some new ones
	b1.columns = slices.Clip(b.columns)
	i := len(b.columns)
	for range rnd.IntN(7) {
		col := "c" + strconv.Itoa(i)
		b1.columns = append(b1.columns, col)
		i++
	}

	b2.columns = slices.Clip(b.columns)
	i = len(b.columns)
	for range rnd.IntN(7) {
		col := "c" + strconv.Itoa(i)
		b2.columns = append(b2.columns, col)
		i++
	}

	q1, q2 := b1.finish(), b2.finish()

	// make the tables disjoint (by fixed) half the time
	switch rnd.IntN(8) {
	case 0:
		q1 = composeFuzzExtend(ft, q1)
	case 1:
		q2 = composeFuzzExtend(ft, q2)
	case 2:
		q1 = composeFuzzWhere(ft, q1)
	case 3:
		q2 = composeFuzzWhere(ft, q2)
	}
	return q1, q2
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

func makeEmptyKey(rnd *rand.Rand, qs *buildFT) {
	qs.keys = emptyKey
	qs.indexes = nil
	if len(qs.data) > 1 {
		qs.data = qs.data[:1]
		if rnd.IntN(2) == 1 {
			qs.data = nil
		}
	} else {
		qs.data = nil
	}
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzTimes ./dbms/query/

var fuzzTimesRunner = fuzzRunner{build: fuzzTimes}

func FuzzTimes(f *testing.F) {
	fuzzTimesRunner.Fuzz(f)
}

func TestFuzzTimes(t *testing.T) {
	fuzzTimesRunner.Test(t)
}

func fuzzTimes(ft *FT) Query {
	q1, q2 := newDisjointQS(ft)
	return NewTimes(q1, q2)
}

func newDisjointQS(ft *FT) (Query, Query) {
	q1 := ft.newFT().Sizes(20, 3, 3).Build()
	q2 := ft.newFT().Sizes(20, 3, 3).Prefix("d").Build()
	return q1, q2
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzJoin ./dbms/query/

var fuzzJoinRunner = fuzzRunner{build: fuzzJoin}

func FuzzJoin(f *testing.F) {
	fuzzJoinRunner.Fuzz(f)
}

func TestFuzzJoin(t *testing.T) {
	start11Count := join11Count.Load()
	start1nCount := join1nCount.Load()
	startn1Count := joinn1Count.Load()
	startnnCount := joinnnCount.Load()

	fuzzJoinRunner.Test(t)

	fmt.Println("11:", join11Count.Load()-start11Count,
		"1n:", join1nCount.Load()-start1nCount,
		"n1:", joinn1Count.Load()-startn1Count,
		"nn:", joinnnCount.Load()-startnnCount)
	assert.T(t).This(join11Count.Load() - start11Count).Isnt(0)
	assert.T(t).This(join1nCount.Load() - start1nCount).Isnt(0)
	assert.T(t).This(joinn1Count.Load() - startn1Count).Isnt(0)
	assert.T(t).This(joinnnCount.Load() - startnnCount).Isnt(0)
	fmt.Println("no results", noResults, "/", fuzzCount)
}

func TestFuzzJoinDebug(t *testing.T) {
	fuzzJoinRunner.Run(t, 10854391646124096407, 353583731168819573)
}

func fuzzJoin(ft *FT) Query {
	q1, q2, to := newFuzzJoin(ft)
	return NewJoin(q1, q2, to, ft.rt)
}

func newFuzzJoin(ft *FT) (Query, Query, []string) {
	b1 := ft.newFT().NoEmptyKey().construct()
	b2 := ft.newFT().NoEmptyKey().Prefix("d").construct()
	rnd := ft.rnd
	var by []string
	switch rnd.IntN(4) {
	case 0: // 1:1
		b1nc := len(b1.columns)
		key := joinBy(rnd, b1, b2)
		by = key
		addKey(rnd, b1, key)
		// join data on b2
		perm := rnd.Perm(len(b1.data))
		for i := range b2.data {
			if len(perm) == 0 || rnd.IntN(2) == 0 {
				for range key {
					b2.data[i] = append(b2.data[i], "J"+strconv.Itoa(i))
				}
			} else {
				row := b1.data[perm[0]]
				perm = perm[1:]
				b2.data[i] = append(b2.data[i], row[b1nc:]...)
			}
		}
		b2.keys = append(b2.keys, key)
	case 1, 2: // 1:n or n:1
		if len(b1.data) < len(b2.data) {
			b1, b2 = b2, b1
		}
		b1nc := len(b1.columns)
		by = joinBy(rnd, b1, b2)
		addKey(rnd, b1, by)

		span := calcSpan(len(by), b1, b2)
		for i := range b2.data {
			if rnd.IntN(2) == 0 || len(b1.data) == 0 {
				for range by {
					b2.data[i] = append(b2.data[i], "j"+strconv.Itoa(span))
				}
			} else {
				row := random(b1.data, rnd)
				b2.data[i] = append(b2.data[i], row[b1nc:]...)
			}
		}
		b2.indexes = append(b2.indexes, by)
		if rnd.IntN(2) == 1 {
			b1, b2 = b2, b1
		}
	case 3: // n:n
		by = joinBy(rnd, b1, b2)
		ncols := len(by)
		span := calcSpan(ncols, b1, b2)
		for i := range b1.data {
			for range ncols {
				b1.data[i] = append(b1.data[i], "j"+strconv.Itoa(rnd.IntN(span)))
			}
		}
		for i := range b2.data {
			for range ncols {
				b2.data[i] = append(b2.data[i], "j"+strconv.Itoa(rnd.IntN(span)))
			}
		}
		b1.indexes = append(b1.indexes, by)
		b2.indexes = append(b2.indexes, by)
	}
	return b1.finish(), b2.finish(), by
}

func calcSpan(ncols int, b1, b2 *buildFT) int {
	switch ncols {
	case 1:
		return len(b1.data) + len(b2.data)
	case 2:
		return 15
	case 3:
		return 7
	default:
		panic(assert.ShouldNotReachHere())
	}
}

// joinBy adds join columns to both sources
func joinBy(rnd *rand.Rand, b1 *buildFT, b2 *buildFT) []string {
	ncols := 1 + rnd.IntN(3)
	cols := make([]string, ncols)
	for i := range cols {
		cols[i] = "j" + strconv.Itoa(i)
	}
	b1.columns = append(b1.columns, cols...)
	b2.columns = append(b2.columns, cols...)
	return cols
}

// addKey adds unique key data to a source, and creates a key index
func addKey(rnd *rand.Rand, b *buildFT, key []string) {
	x := uint16(rnd.Int())
	for i := range b.data {
		x = bits.Shuffle16(x) // shuffle ensures unique key values
		n := x
		for k := range key {
			// split n (a unique value) over the columns of the key
			var v uint16
			if k < len(key)-1 {
				// 4 bits = 0 - 15 gives chance of duplicates
				v = n & 0b1111
				n >>= 4
			} else {
				// last column gets the rest
				v = n
			}
			b.data[i] = append(b.data[i], "j"+strconv.Itoa(int(v)))
		}
	}
	b.keys = append(b.keys, key)
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzLeftJoin ./dbms/query/

var fuzzLeftJoinRunner = fuzzRunner{build: fuzzLeftJoin}

func FuzzLeftJoin(f *testing.F) {
	fuzzLeftJoinRunner.Fuzz(f)
}

func TestFuzzLeftJoin(t *testing.T) {
	start11Count := leftJoin11Count.Load()
	start1nCount := leftJoin1nCount.Load()
	startn1Count := leftJoinn1Count.Load()
	startnnCount := leftJoinnnCount.Load()

	fuzzLeftJoinRunner.Test(t)

	fmt.Println("11:", leftJoin11Count.Load()-start11Count,
		"1n:", leftJoin1nCount.Load()-start1nCount,
		"n1:", leftJoinn1Count.Load()-startn1Count,
		"nn:", leftJoinnnCount.Load()-startnnCount)
	assert.T(t).This(leftJoin11Count.Load() - start11Count).Isnt(0)
	assert.T(t).This(leftJoin1nCount.Load() - start1nCount).Isnt(0)
	assert.T(t).This(leftJoinn1Count.Load() - startn1Count).Isnt(0)
	assert.T(t).This(leftJoinnnCount.Load() - startnnCount).Isnt(0)
	fmt.Println("no results", noResults, "/", fuzzCount)
}

func fuzzLeftJoin(ft *FT) Query {
	q1, q2, to := newFuzzJoin(ft)
	return NewLeftJoin(q1, q2, to, ft.rt)
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzSemiJoin ./dbms/query/

var fuzzSemiJoinRunner = fuzzRunner{build: fuzzSemiJoin}

func FuzzSemiJoin(f *testing.F) {
	fuzzSemiJoinRunner.Fuzz(f)
}

func TestFuzzSemiJoin(t *testing.T) {
	fuzzSemiJoinRunner.Test(t)
	fmt.Println("no results", noResults, "/", fuzzCount)
}

func fuzzSemiJoin(ft *FT) Query {
	q1, q2, to := newFuzzJoin(ft)
	return NewSemiJoin(q1, q2, to, ft.rt)
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzWhere ./dbms/query/

var fuzzWhereRunner = fuzzRunner{build: fuzzWhere}

func TestFuzzWhereDebug(t *testing.T) {
	fuzzWhereRunner.Run(t, 12582666410114420314, 13574490830499976766)
}

func FuzzWhere(f *testing.F) {
	fuzzWhereRunner.Fuzz(f)
}

func TestFuzzWhere(t *testing.T) {
	startSingleton := whereSingletonCount.Load()
	fuzzWhereRunner.Test(t)
	deltaSingleton := whereSingletonCount.Load() - startSingleton
	fmt.Println("Where strategies: singleton", deltaSingleton, "/", fuzzCount)
}

func fuzzWhere(ft *FT) Query {
	// Use richer index topologies 75% of the time to exercise skip scan,
	// but keep some coverage of empty keys and plain QuerySource.
	var q Query
	if ft.rnd.IntN(4) == 0 {
		q = ft.NewFuzzTable()
	} else {
		q = ft.newFT().NoEmptyKey().Sizes(151, 6, 8).Build()
	}
	return composeFuzzWhere(ft, q)
}

func composeFuzzWhere(ft *FT, q Query) Query {
	expr := randomWhereExpr(ft.rnd, q.Columns(), q.Keys(), q.Indexes())
	return NewWhere(q, expr, ft.rt)
}

func randomWhereExpr(rnd *rand.Rand, cols []string, keys [][]string, indexes [][]string) ast.Expr {
	if len(keys) > 0 && rnd.IntN(10) == 0 {
		key := random(keys, rnd)
		if len(key) > 0 {
			exprs := make([]ast.Expr, len(key))
			for i, col := range key {
				val := SuStr(col + "_" + strconv.Itoa(rnd.IntN(16)))
				exprs[i] = &ast.Binary{Tok: tok.Is, Lhs: &ast.Ident{Name: col}, Rhs: &ast.Constant{Val: val}}
			}
			if len(exprs) == 1 {
				return exprs[0]
			}
			return &ast.Nary{Tok: tok.And, Exprs: exprs}
		}
	}

	if len(cols) == 0 {
		return &ast.Constant{Val: True}
	}
	n := 1 + rnd.IntN(4)
	exprs := make([]ast.Expr, n)
	var ix []string
	if len(indexes) > 0 && rnd.IntN(2) == 0 {
		ix = random(indexes, rnd)
	}
	for i := range n {
		col := random(cols, rnd)
		if len(ix) > 0 && rnd.IntN(2) == 0 {
			if len(ix) > 1 && rnd.IntN(3) != 0 {
				col = ix[1+rnd.IntN(len(ix)-1)]
			} else {
				col = random(ix, rnd)
			}
		}
		val := SuStr(col + "_" + strconv.Itoa(rnd.IntN(16)))
		switch rnd.IntN(5) {
		case 0: // =
			exprs[i] = &ast.Binary{Tok: tok.Is, Lhs: &ast.Ident{Name: col}, Rhs: &ast.Constant{Val: val}}
		case 1: // <
			exprs[i] = &ast.Binary{Tok: tok.Lt, Lhs: &ast.Ident{Name: col}, Rhs: &ast.Constant{Val: val}}
		case 2: // >
			exprs[i] = &ast.Binary{Tok: tok.Gt, Lhs: &ast.Ident{Name: col}, Rhs: &ast.Constant{Val: val}}
		case 3: // in
			nvals := 1 + rnd.IntN(3)
			vals := make([]ast.Expr, nvals)
			for j := range nvals {
				vals[j] = &ast.Constant{Val: IntVal(rnd.IntN(10))}
			}
			exprs[i] = &ast.In{E: &ast.Ident{Name: col}, Exprs: vals}
		case 4: // col = col (not a btree range; works on any type)
			col2 := random(cols, rnd)
			exprs[i] = &ast.Binary{Tok: tok.Is, Lhs: &ast.Ident{Name: col}, Rhs: &ast.Ident{Name: col2}}
		}
	}
	if n == 1 {
		return exprs[0]
	}
	return &ast.Nary{Tok: tok.And, Exprs: exprs}
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzExtend ./dbms/query/

var fuzzExtendRunner = fuzzRunner{build: fuzzExtend}

func FuzzExtend(f *testing.F) {
	fuzzExtendRunner.Fuzz(f)
}

func TestFuzzExtend(t *testing.T) {
	fuzzExtendRunner.Test(t)
}

func TestFuzzExtendDebug(t *testing.T) {
	fuzzExtendRunner.Run(t, 493, 913)
}

func fuzzExtend(ft *FT) Query {
	return composeFuzzExtend(ft, ft.NewFuzzTable())
}

func composeFuzzExtend(ft *FT, qs Query) Query {
	n := 1 + ft.rnd.IntN(5)
	cols := make([]string, n)
	exprs := make([]ast.Expr, n)
	qcols := qs.Columns()
	for i := range n {
		for j := 0; ; j++ {
			name := "x" + strconv.Itoa(i)
			if j > 0 {
				name += "_" + strconv.Itoa(j)
			}
			if !slices.Contains(qcols, name) {
				cols[i] = name
				break
			}
		}
		if ft.rnd.IntN(2) == 0 && len(qcols) > 0 {
			exprs[i] = &ast.Ident{Name: random(qcols, ft.rnd)}
		} else {
			exprs[i] = &ast.Constant{Val: IntVal(ft.rnd.IntN(1000))}
		}
	}
	return NewExtend(qs, cols, exprs)
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzTempIndex ./dbms/query/

var fuzzTempIndexRunner = fuzzRunner{build: fuzzTempIndex}

func FuzzTempIndex(f *testing.F) {
	fuzzTempIndexRunner.Fuzz(f)
}

func TestFuzzTempIndexDebug(t *testing.T) {
	fuzzTempIndexRunner.Run(t, 493, 913)
}

func TestFuzzTempIndex(t *testing.T) {
	fuzzTempIndexRunner.Test(t)
}

func fuzzTempIndex(ft *FT) Query {
	q := ft.newFT().NoEmptyKey().Build()
	cols := q.Columns()
	n := 1 + ft.rnd.IntN(min(3, len(cols)))
	order := set.RandPerm(ft.rnd, cols, n)
	return NewTempIndex(q, order, ft.rt)
}

//-------------------------------------------------------------------

var fuzzCount = 0
var noResults = 0

func fuzzQuery(t *testing.T, q Query, ft *FT) {
	before := String(q) // before Transform
	defer func() {
		if r := recover(); r != nil || t.Failed() {
			fmt.Println(before)
			fmt.Println(String(q))
			if r != nil {
				panic(r)
			}
		}
	}()
	which := random([]string{"lookup", "select"}, ft.rnd)
	if isEmptyKey(q.Indexes()) {
		which = "get"
	}
	var index []string
	if which == "lookup" {
		ki := keyIndexes(q)
		if len(ki) > 0 {
			index = random(keyIndexes(q), ft.rnd)
		} else {
			which = "select"
		}
	}
	if which == "select" {
		indexes := q.Indexes()
		if len(indexes) > 0 {
			index = random(q.Indexes(), ft.rnd)
		} else {
			which = "get"
		}
	}
	q = q.Transform()
	req := OrderedReq(index, 1)
	fixcost, varcost := Optimize2(q, ReadMode, req)
	fuzzCount++
	if fixcost+varcost >= impossible {
		t.Fatal("impossible\n", format(0, q, 0))
	}
	// fmt.Println(String(q))
	q = SetApproach2(q, req, ft.rt)
	q.SetTran(ft.rt)

	hdr := q.Header()
	expected := q.Simple(nil)
	// fmt.Println("Simple", len(expected))
	if len(expected) == 0 {
		noResults++
	}

	qh := NewQueryHasher(hdr).CheckDups()
	for _, row := range expected {
		qh.Row(row)
	}
	testRandomGet(t, ft.rnd, q, qh, hdr, nil)

	switch which {
	case "lookup":
		cols := hdr.Columns
		testRandomLookups(t, ft.rnd, q, index, cols, expected)
	case "select":
		testRandomSelects(t, ft.rnd, q, index, expected)
	}
}

func keyIndexes(q Query) [][]string {
	var keyIndexes [][]string
	for _, index := range q.Indexes() {
		for _, key := range q.Keys() {
			if set.Equal(index, key) {
				keyIndexes = append(keyIndexes, index)
			}
		}
	}
	return keyIndexes
}

func testRandomGet(t *testing.T, rnd *rand.Rand, q Query, qh *QueryHash, hdr *Header, sels Sels) {
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
	q.Select(sels)

	// Do a random walk with Next/Prev using nextRows as expected
	history := ""
	nsteps := min(100, len(nextRows)*3)
	for i := range nsteps {
		// Occasionally add a Select to reset indexed flag for projMap
		if rnd.IntN(20) == 0 { // 5% chance
			if sels == nil {
				q.Select(nil) // this also rewinds
			} else {
				q.Rewind()
			}
			data.rewind()
		}

		pos := data.pos
		if data.pos == dsEof {
			history += "r"
			q.Rewind()
			data.rewind()
		}
		dir := random([]Dir{Next, Prev}, rnd)
		history += string(dir)
		expectedRow := data.get(dir)
		row := q.Get(nil, dir)

		if expectedRow == nil && row != nil {
			t.Fatalf("random walk step %d: %c from %v: expected nil, got row\nhistory %s",
				i, dir, pos, history)
		} else if expectedRow != nil && row == nil {
			t.Log(q)
			t.Fatalf("random walk step %d: %c from %v: expected row, got nil\nhistory %s",
				i, dir, pos, history)
		} else if expectedRow != nil && row != nil {
			if !hdr.EqualRows(row, expectedRow, nil, nil) {
				t.Fatalf("random walk step %d: %c from %v: row mismatch\nhistory %s",
					i, dir, pos, history)
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

	// Pattern 0: Rewind, Prev - should go to last row
	q.Rewind()
	row := q.Get(nil, Prev) // last row
	if n > 0 {
		check("Rewind, Prev", row, nextRows[n-1])
	}

	// Pattern 1: Rewind, Next, Prev - after first Next, Prev should return nil
	q.Rewind()
	row = q.Get(nil, Next) // first row or nil if empty
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
	// plain stick at eof: Prev should also return nil
	if n > 0 {
		q.Rewind()
		for i := range n {
			row = q.Get(nil, Next)
			check("ToEnd: N"+strconv.Itoa(i), row, nextRows[i])
		}
		row = q.Get(nil, Next) // past end
		check("ToEnd: N-past", row, nil)
		row = q.Get(nil, Prev) // plain stick: should be nil
		check("ToEnd: P", row, nil)
	}

	// Pattern 6: Prev to beginning, past beginning (nil), then Next
	// plain stick at eof: Next should also return nil
	if n > 0 {
		q.Rewind()
		for i := n - 1; i >= 0; i-- {
			row = q.Get(nil, Prev)
			check("ToBegin: P"+strconv.Itoa(n-1-i), row, nextRows[i])
		}
		row = q.Get(nil, Prev) // past beginning
		check("ToBegin: P-past", row, nil)
		row = q.Get(nil, Next) // plain stick: should be nil
		check("ToBegin: N", row, nil)
	}

	// Pattern 7: Rewind, Next, Prev, Next - plain stick at eof after Prev
	if n > 0 {
		q.Rewind()
		row = q.Get(nil, Next) // first
		check("NPN: N1", row, nextRows[0])
		row = q.Get(nil, Prev) // nil
		check("NPN: P", row, nil)
		row = q.Get(nil, Next) // plain stick: nil
		check("NPN: N2", row, nil)
	}

	// Pattern 8: Rewind, Prev, Next, Prev - plain stick at eof after Next
	if n > 0 {
		q.Rewind()
		row = q.Get(nil, Prev) // last
		check("PNP: P1", row, nextRows[n-1])
		row = q.Get(nil, Next) // nil
		check("PNP: N", row, nil)
		row = q.Get(nil, Prev) // plain stick: nil
		check("PNP: P2", row, nil)
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

func testRandomSelects(t *testing.T, rnd *rand.Rand, q Query, index []string, allRows []Row) {
	t.Helper()
	hdr := q.Header()
	testExistentSelect(t, allRows, rnd, hdr, index, q)
	testNonExistentSelect(t, allRows, rnd, hdr, index, q)
}

func testExistentSelect(t *testing.T, allRows []Row, rnd *rand.Rand, hdr *Header, index []string, q Query) {
	if len(allRows) == 0 {
		return
	}
	for range 10 {
		srcRow := random(allRows, rnd)
		sels := indexSelectCriteria(rnd, srcRow, hdr, index)
		q.Select(sels)

		qh := NewQueryHasher(hdr)
		for _, row := range allRows {
			if selMatchIndex(hdr, row, sels, index) {
				qh.Row(row)
			}
		}

		testRandomGet(t, rnd, q, qh, hdr, sels)

		q.Select(nil) // clear select
	}
}

// selMatchIndex checks only the index columns, not extra columns.
// It iterates the index in order and stops at the first missing column,
// matching the behavior of selKeys/TempIndex.makeKey.
func selMatchIndex(hdr *Header, row Row, sels Sels, index []string) bool {
	for _, col := range index {
		val, ok := sels.Get(col)
		if !ok {
			break // stop at first missing column (matches selKeys behavior)
		}
		if row.GetRaw(hdr, col) != val {
			return false
		}
	}
	return true
}

// indexSelectCriteria picks a random prefix of the index for select criteria.
func indexSelectCriteria(rnd *rand.Rand, row Row, hdr *Header, index []string) Sels {
	n := 1 + rnd.IntN(len(index))
	selCols := slices.Clone(index[:n])
	rnd.Shuffle(len(selCols), func(i, j int) {
		selCols[i], selCols[j] = selCols[j], selCols[i]
	})

	sels := make(Sels, len(selCols))
	for i, col := range selCols {
		sels[i] = Sel{col: col, val: row.GetRaw(hdr, col)}
	}
	return sels
}

func testNonExistentSelect(t *testing.T, allRows []Row, rnd *rand.Rand, hdr *Header, index []string, q Query) {
	for range 10 {
		// If there are no rows, use a dummy row sized to match hdr.Fields.
		// This avoids panics when hdr.Fields references derived records (e.g. Extend).
		srcRow := make(Row, len(hdr.Fields))
		if len(allRows) > 0 {
			srcRow = random(allRows, rnd)
		}
		sels := indexSelectCriteria(rnd, srcRow, hdr, index)
		sels[rnd.IntN(len(sels))].val = "nonexistent"
		q.Select(sels)
		if q.Get(nil, Next) != nil {
			t.Fatal("non-existent select returned a row")
		}
		q.Select(nil) // clear select
	}
}

//-------------------------------------------------------------------

func testRandomLookups(t *testing.T, rnd *rand.Rand, q Query, index, cols []string, allRows []Row) {
	lookupCols := slices.Clone(index)
	slc.Shuffle(rnd, lookupCols)

	testExistentLookup(t, allRows, rnd, lookupCols, q, cols)
	testNonExistentLookup(t, rnd, q, lookupCols)
}

// canLookup checks if a lookup with the given index is valid.
// The index columns must form a key (or be a subset of a key with
// remaining key columns being fixed).
func canLookup(keys [][]string, fixed Fixed, index []string) bool {
	for _, key := range keys {
		// Check if index is subset of key+fixed (no extra columns)
		subset := true
		for _, col := range index {
			if !slices.Contains(key, col) && !fixed.Single(col) {
				subset = false
				break
			}
		}

		if subset {
			// Check if key is subset of index+fixed (complete key)
			superset := true
			for _, col := range key {
				if !slices.Contains(index, col) && !fixed.Single(col) {
					superset = false
					break
				}
			}
			if superset {
				return true
			}
		}
	}
	return false
}

func testExistentLookup(t *testing.T, allRows []Row, rnd *rand.Rand, lookupCols []string, q Query, cols []string) {
	t.Helper()
	if len(allRows) == 0 {
		return
	}
	hdr := q.Header()
	for range min(10, len(allRows)) {
		srcRow := random(allRows, rnd)

		sels := make(Sels, len(lookupCols))
		for i, col := range lookupCols {
			sels[i] = Sel{col: col, val: srcRow.GetRaw(hdr, col)}
		}

		result := q.Lookup(nil, sels)

		if result == nil {
			t.Fatal("lookup returned nil for existing key")
		}

		for i, col := range lookupCols {
			if result.GetRaw(hdr, col) != sels[i].val {
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
		sels := make(Sels, len(lookupCols))
		// set one of the keyVals to a non-existent value
		// the others to possibly existing values
		r := rnd.IntN(len(lookupCols))
		for i, col := range lookupCols {
			sels[i].col = col
			if i == r {
				sels[i].val = "nonexistent"
			} else {
				sels[i].val = col + "_" + strconv.Itoa(rnd.IntN(100))
			}
		}
		result := q.Lookup(nil, sels)
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

	for range nfuzz {
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

	fmt.Printf("splitShare stats: empty1=%d (%.1f%%), empty2=%d (%.1f%%), empty3=%d (%.1f%%)\n",
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

func TestCanLookup(t *testing.T) {
	// Case 1: Index is subset of key (should be false, currently true)
	keys := [][]string{{"a", "b"}}
	fixed := Fixed{}
	index := []string{"a"}
	assert.T(t).This(canLookup(keys, fixed, index)).Is(false)

	// Case 2: Index contains complete key but has extra column (should be false due to Where.Lookup restriction)
	keys = [][]string{{"a"}}
	fixed = Fixed{}
	index = []string{"a", "b"}
	assert.T(t).This(canLookup(keys, fixed, index)).Is(false)

	// Case 3: Index matches key exactly (should be true)
	keys = [][]string{{"a", "b"}}
	fixed = Fixed{}
	index = []string{"b", "a"}
	assert.T(t).This(canLookup(keys, fixed, index)).Is(true)

	// Case 4: Key part is fixed (should be true)
	keys = [][]string{{"a", "b"}}
	fixed = Fixed{{col: "b", values: []string{"1"}}}
	index = []string{"a"}
	assert.T(t).This(canLookup(keys, fixed, index)).Is(true)

	// Case 5: Key part is fixed (should be true)
	keys = [][]string{{"a"}}
	fixed = Fixed{{col: "b", values: []string{"1"}}}
	index = []string{"a", "b"}
	assert.T(t).This(canLookup(keys, fixed, index)).Is(true)
}
