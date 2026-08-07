// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/bits"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
	"github.com/apmckinlay/gsuneido/util/str"
)

func init() {
	sortForTest = true
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
}

const nfuzz = 200

var (
	reqNoneCount   atomic.Int64
	reqOrderCount  atomic.Int64
	reqGroupCount  atomic.Int64
	reqUniqueCount atomic.Int64
)

type fuzzRunner struct {
	build func(*FT) Query
}

func (fr fuzzRunner) Run(t *testing.T, seed1, seed2 uint64) {
	defer func(jr int) { joinRev = jr }(joinRev)
	joinRev = impossible // keep joins in the order created
	defer func(ti int) { ticostAdj = ti }(ticostAdj)
	ticostAdj = 9999999 // discourage temp indexes unless impossible without
	ft := newFT(seed1, seed2)
	defer ft.db.Close()
	defer func(rb *rand.Rand) { randomBest = rb }(randomBest)
	randomBest = ft.rnd // choose best randomly (with the seeded generator) to exercise more possibilities
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
// go test -run '^$' -fuzz=FuzzRandom ./dbms/query

var fuzzRandomRunner = fuzzRunner{build: fuzzRandom}

func FuzzRandom(f *testing.F) {
	fuzzRandomRunner.Fuzz(f)
}

func TestFuzzRandomDebug(t *testing.T) {
	fuzzRandomRunner.Run(t, 12, 226)
}

func TestFuzzRandom(t *testing.T) {
	startNone := reqNoneCount.Load()
	startOrder := reqOrderCount.Load()
	startGroup := reqGroupCount.Load()
	startUnique := reqUniqueCount.Load()

	fuzzRandomRunner.Test(t)

	fmt.Printf("Require uses: none=%d order=%d group=%d unique=%d\n",
		reqNoneCount.Load()-startNone,
		reqOrderCount.Load()-startOrder,
		reqGroupCount.Load()-startGroup,
		reqUniqueCount.Load()-startUnique)
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
// go test -run '^$' -fuzz=FuzzTable ./dbms/query

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
// go test -run '^$' -fuzz=FuzzProject ./dbms/query

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
	deps := ft.ruleDeps
	if ext, ok := qs.(*Extend); ok && len(ext.fwd) > 0 {
		deps = make(map[string]string, len(ft.ruleDeps)+len(ext.fwd))
		for k, v := range ft.ruleDeps {
			deps[k] = v
		}
		for i, fwd := range ext.fwd {
			target := fwd[1:]
			if dep, ok := ft.ruleDeps[target]; ok {
				deps[ext.cols[i]] = dep
			}
		}
	}
	projCols := randomProjectCols(ft.rnd, qs.Columns(), qs.Indexes(), deps)
	return NewProject(qs, projCols)
}

func randomProjectCols(rnd *rand.Rand, srcCols []string, indexes [][]string,
	ruleDeps map[string]string) []string {
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
	cols := make([]string, 0, n)
	for i := range n {
		c := srcCols[perm[i]]
		// if c is a rule, check if its dep is available
		if dep, ok := ruleDeps[c]; ok {
			if !slices.Contains(srcCols, dep) {
				continue // skip this rule col since dep not available
			}
			cols = append(cols, c, dep)
		} else {
			cols = append(cols, c)
		}
	}
	if len(cols) == 0 {
		// fallback: at least one column
		return []string{srcCols[0]}
	}
	return cols
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzRename ./dbms/query

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
	from, to := randomRename(ft.rnd, qs.Columns(), qs.Header().Rules(), ft.ruleRefs)
	return NewRename(qs, from, to)
}

func randomRename(rnd *rand.Rand, srcCols, ruleCols, ruleRefs []string) (from, to []string) {
	if len(srcCols) == 0 {
		return nil, nil
	}

	// exclude rule columns (derived, not stored) and columns referenced by rules from renaming
	excluded := make([]string, 0, len(ruleCols)+len(ruleRefs))
	excluded = append(excluded, ruleCols...)
	excluded = append(excluded, ruleRefs...)
	nonRuleCols := make([]string, 0, len(srcCols))
	for _, col := range srcCols {
		if !slices.Contains(excluded, col) {
			nonRuleCols = append(nonRuleCols, col)
		}
	}
	if len(nonRuleCols) == 0 {
		return nil, nil
	}

	// Determine how many columns to rename (1 to 3)
	n := 1 + rnd.IntN(min(3, len(nonRuleCols)))

	// Choose random columns to rename
	perm := rnd.Perm(len(nonRuleCols))
	from = make([]string, n)
	for i := range n {
		from[i] = nonRuleCols[perm[i]]
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
// go test -run '^$' -fuzz=FuzzSummarize ./dbms/query

var fuzzSummarizeRunner = fuzzRunner{build: fuzzSummarize}

func FuzzSummarize(f *testing.F) {
	fuzzSummarizeRunner.Fuzz(f)
}

func TestFuzzSummarizeDebug(t *testing.T) {
	fuzzSummarizeRunner.Run(t, 16602272769697815038, 9374908162863325624)
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
		// avoid 'on' columns that conflict with 'by' (now disallowed)
		nonByCols := set.Difference(srcCols, by)
		// build a pool of distinct (op, on) pairs whose default output
		// column names do not conflict with 'by' columns
		type pair struct{ op, on string }
		bySet := make(map[string]struct{}, len(by))
		for _, c := range by {
			bySet[c] = struct{}{}
		}
		pool := []pair(nil)
		add := func(op, on string) {
			if _, ok := bySet[defaultColName(op, on)]; !ok {
				pool = append(pool, pair{op, on})
			}
		}
		add("count", "")
		for _, op := range sumOps {
			if op == "count" {
				continue
			}
			for _, on := range nonByCols {
				add(op, on)
			}
		}
		// select pairs so that no 'on' collides with a generated output name
		selected := make([]pair, 0, nops)
		outputNames := make(map[string]struct{}, nops)
		selectedOns := make(map[string]struct{}, nops)
		rnd.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
		for _, p := range pool {
			if len(selected) >= nops {
				break
			}
			out := defaultColName(p.op, p.on)
			if _, dup := outputNames[out]; dup {
				continue
			}
			if _, conflict := selectedOns[out]; conflict {
				continue
			}
			if p.on != "" {
				if _, conflict := outputNames[p.on]; conflict {
					continue
				}
			}
			selected = append(selected, p)
			outputNames[out] = struct{}{}
			selectedOns[p.on] = struct{}{}
		}
		nops = len(selected)
		cols = cols[:nops]
		ops = ops[:nops]
		ons = ons[:nops]
		for i := range nops {
			ops[i] = selected[i].op
			ons[i] = selected[i].on
			cols[i] = "" // use default name
		}
	}
	return
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzMinus ./dbms/query

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

func TestFuzzMinusDebug(t *testing.T) {
	fuzzMinusRunner.Run(t, 192, 1)
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzIntersect ./dbms/query

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
	fuzzIntersectRunner.Run(t, 8687611748617673601, 15585953995166667714)
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzUnion ./dbms/query

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

func TestFuzzUnionDebug(t *testing.T) {
	fuzzUnionRunner.Run(t, 192, 1)
}

func fuzzUnion(ft *FT) Query {
	q1, q2 := newCompatibleQS(ft)
	return NewUnion(q1, q2)
}

//-------------------------------------------------------------------

// newCompatibleQS creates QuerySources for Union, Intersect, Minus
func newCompatibleQS(ft *FT) (Query, Query) {
	rnd := ft.rnd
	if rnd.IntN(5) == 3 {
		return ft.NewFuzzTable(), ft.NewFuzzTable()
	}

	b := ft.newFT().Sizes(73, 5, 5).construct()
	b1 := *b
	b2 := *b

	if rnd.IntN(5) == 2 {
		return b1.finish(), b2.finish()
	}

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

	// 10% of the time, set key to all columns (like non-disjoint Union result)
	switch rnd.IntN(19) {
	case 5:
		makeAllColsKey(&b1)
	case 13:
		makeAllColsKey(&b2)
	case 17:
		makeAllColsKey(&b1)
		makeAllColsKey(&b2)
	}

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
	addExtraData(rnd, b1.data, b.columns, b1.columns)

	b2.columns = slices.Clip(b.columns)
	i = len(b.columns)
	for range rnd.IntN(7) {
		col := "c" + strconv.Itoa(i)
		b2.columns = append(b2.columns, col)
		i++
	}
	addExtraData(rnd, b2.data, b.columns, b2.columns)

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

// addExtraData populates data values for columns added after data generation.
// This ensures extra columns have varied non-empty values so that Lookup
// with sels beyond the index can produce mismatches
// (catching bugs like Union.Lookup not verifying the full sels).
func addExtraData(rnd *rand.Rand, data [][]string, baseCols, allCols []string) {
	nbase := len(baseCols)
	nextra := len(allCols) - nbase
	if nextra == 0 {
		return
	}
	for i := range data {
		extra := make([]string, nextra)
		for j := range nextra {
			col := allCols[nbase+j]
			extra[j] = col + "_" + strconv.Itoa(rnd.IntN(50))
		}
		data[i] = append(data[i], extra...)
	}
}

func makeAllColsKey(b *buildFT) {
	// exclude rule columns since they can't be part of a key (they're derived)
	cols := make([]string, 0, len(b.columns))
	for _, col := range b.columns {
		if !slices.Contains(b.ruleCols, col) {
			cols = append(cols, col)
		}
	}
	b.keys = [][]string{cols}
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
// go test -run '^$' -fuzz=FuzzTimes ./dbms/query

var fuzzTimesRunner = fuzzRunner{build: fuzzTimes}

func FuzzTimes(f *testing.F) {
	fuzzTimesRunner.Fuzz(f)
}

func TestFuzzTimes(t *testing.T) {
	fuzzTimesRunner.Test(t)
}

func fuzzTimes(ft *FT) Query {
	q1, q2 := newTimesQS(ft)
	return NewTimes(q1, q2)
}

func newTimesQS(ft *FT) (Query, Query) {
	q1 := ft.newFT().Sizes(20, 3, 3).Build()
	q2 := ft.newFT().Sizes(20, 3, 3).Prefix("d").Build()
	return q1, q2
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzJoin ./dbms/query

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
	// Sometimes add a shared rule column to exercise rule-in-by paths.
	// The column name is capitalized (Jr) so it's derived (not stored),
	// but Table.Columns() returns it uncapitalized (jr), so by uses "jr".
	if rnd.IntN(4) == 0 {
		ruleCol := "Jr"
		b1.columns = append(b1.columns, ruleCol)
		b2.columns = append(b2.columns, ruleCol)
		b1.ruleCols = append(b1.ruleCols, ruleCol)
		b2.ruleCols = append(b2.ruleCols, ruleCol)
		by = append(by, "jr")
	}
	q1, q2 := b1.finish(), b2.finish()
	if slices.Contains(by, "jr") {
		Global.TestDef("Rule_jr", compile.Constant("function() { return 123 }"))
	}
	return q1, q2, by
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
// go test -run '^$' -fuzz=FuzzLeftJoin ./dbms/query

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
// go test -run '^$' -fuzz=FuzzSemiJoin ./dbms/query

var fuzzSemiJoinRunner = fuzzRunner{build: fuzzSemiJoin}

func FuzzSemiJoin(f *testing.F) {
	fuzzSemiJoinRunner.Fuzz(f)
}

func TestFuzzSemiJoin(t *testing.T) {
	start11Count := semiJoin11Count.Load()
	start1nCount := semiJoin1nCount.Load()
	startn1Count := semiJoinn1Count.Load()
	startnnCount := semiJoinnnCount.Load()

	fuzzSemiJoinRunner.Test(t)

	fmt.Println("11:", semiJoin11Count.Load()-start11Count,
		"1n:", semiJoin1nCount.Load()-start1nCount,
		"n1:", semiJoinn1Count.Load()-startn1Count,
		"nn:", semiJoinnnCount.Load()-startnnCount)
	assert.T(t).This(semiJoin11Count.Load() - start11Count).Isnt(0)
	assert.T(t).This(semiJoin1nCount.Load() - start1nCount).Isnt(0)
	assert.T(t).This(semiJoinn1Count.Load() - startn1Count).Isnt(0)
	assert.T(t).This(semiJoinnnCount.Load() - startnnCount).Isnt(0)
	fmt.Println("no results", noResults, "/", fuzzCount)
}

func TestFuzzSemiJoinDebug(t *testing.T) {
	fuzzSemiJoinRunner.Run(t, 18275177146248160504, 16758281424391211506)
}

func fuzzSemiJoin(ft *FT) Query {
	q1, q2, to := newFuzzJoin(ft)
	return NewSemiJoin(q1, q2, to, ft.rt)
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzWhere ./dbms/query

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
// go test -run '^$' -fuzz=FuzzExtend ./dbms/query

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

var fuzzCount = 0
var noResults = 0

func fuzzQuery(t *testing.T, q Query, ft *FT) {
	before := String(q) // before Transform
	defer func() {
		if r := recover(); r != nil || t.Failed() {
			PrintTableInfo(q)
			fmt.Println("original:", before)
			fmt.Println("optimized:", String(q))
			if r != nil {
				panic(r)
			}
		}
	}()
	use := random([]Use{ReqNone, ReqOrder, ReqGroup, ReqUnique}, ft.rnd)
	indexes := q.Indexes()
	var index []string
	if len(indexes) == 0 || isEmptyKey(indexes) {
		use = ReqNone
	} else if use == ReqUnique {
		keyIdxs := keyIndexes(q)
		if len(keyIdxs) > 0 {
			index = random(keyIdxs, ft.rnd)
		} else {
			// fallback to all columns since this must contain a key
			index = q.Columns()
		}
	} else {
		index = random(indexes, ft.rnd)
	}
	var req Require
	switch use {
	case ReqNone:
		req = NoneReq(1)
	case ReqOrder:
		// randomly shorten (truncate) index keeping at least one column
		index = index[:1+ft.rnd.IntN(len(index))]
		req = OrderReq(index, 1)
	case ReqGroup:
		nseeks := int32(1 + ft.rnd.IntN(10))
		frac := float32(1) / float32(1+ft.rnd.IntN(4))
		req = GroupReq(slices.Clone(index), frac, nseeks)
	case ReqUnique:
		// add extra columns (to exercise Lookup with sels beyond the index)
		index = slices.Clip(index)
		nextra := ft.rnd.IntN(len(index))
		for range nextra {
			index = set.AddUnique(index, random(q.Columns(), ft.rnd))
		}
		nseeks := int32(1 + ft.rnd.IntN(10))
		req = UniqueReq(index, nseeks)
	}
	q = q.Transform()
	fixcost, varcost := Optimize(q, ReadMode, req)
	if fixcost+varcost >= impossible {
		// fall back to an unordered read
		// fmt.Println("IMPOSSIBLE:", String(q))
		use = ReqNone
		req = NoneReq(1)
		index = nil
		fixcost, varcost = Optimize(q, ReadMode, req)
		if fixcost+varcost >= impossible {
			t.Fatal("impossible\n", format(0, q, 0))
		}
	}
	switch req.use {
	case ReqNone:
		reqNoneCount.Add(1)
	case ReqOrder:
		reqOrderCount.Add(1)
	case ReqGroup:
		reqGroupCount.Add(1)
	case ReqUnique:
		reqUniqueCount.Add(1)
	}
	fuzzCount++
	// fmt.Println(String(q))
	q = SetApproach(q, req, ft.rt)
	q.SetTran(ft.rt)

	hdr := q.Header()
	th := &Thread{}
	expected := q.Simple(th)
	// fmt.Println("Simple", len(expected))
	if len(expected) == 0 {
		noResults++
	}

	qh := NewQueryHasher(hdr).CheckDups()
	for _, row := range expected {
		qh.Row(row)
	}

	// match implicit contract, see require.go
	fc := &fuzzCtx{t: t, ft: ft, rnd: ft.rnd, q: q, hdr: hdr, th: th}
	fc.testRandomGet(qh, nil)
	if use == ReqUnique {
		if len(index) > 0 {
			fc.testRandomLookups(index, expected)
		}
	}
	if use == ReqOrder || use == ReqGroup {
		if len(index) > 0 {
			fc.testRandomSelects(index, expected)
		}
	}
}

type fuzzCtx struct {
	t   *testing.T
	ft  *FT
	rnd *rand.Rand
	q   Query
	hdr *Header
	th  *Thread
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

func (fc *fuzzCtx) testRandomGet(qh *QueryHash, sels Sels) {
	// Get all rows using Next first to establish correct iteration order
	fc.q.Rewind()
	nextRows := getAllRows(fc.q, Next, fc.th)
	if !rowSetsEqual(nextRows, qh, fc.hdr) {
		PrintTableInfo(fc.q)
		for k, v := range fc.ft.rules {
			fmt.Println(k, "=", v)
		}
		fmt.Println("QUERY:", String(fc.q))
		fmt.Println("=== Optimized Get rows (actual) ===")
		for i, row := range nextRows {
			if i < 60 {
				fmt.Printf("  row %d: %s\n", i, RowStr(fc.hdr, row))
				fmt.Println(row)
			}
		}
		fmt.Printf("... total: %d (actual) vs %d (expected)\n\n", len(nextRows), qh.nrows)
		fc.q.Rewind()
		expectedRows := fc.q.Simple(fc.th)
		fmt.Println("=== Simple rows (expected) ===")
		for i, row := range expectedRows {
			if i < 60 {
				fmt.Printf("  row %d: %s\n", i, RowStr(fc.hdr, row))
				fmt.Println(row)
			}
		}
		fmt.Printf("... total: %d\n", len(expectedRows))
		fc.t.Fatalf("Next iteration returned %d rows, expected %d", len(nextRows), qh.nrows)
	}

	// Run deterministic cursor pattern checks before random walk
	fc.testCursorPatterns(nextRows)

	data := NewDataSource(nextRows)

	// Redo the Select after getAllRows to reset indexed state for projMap
	fc.q.Select(sels)

	// Do a random walk with Next/Prev using nextRows as expected
	history := ""
	nsteps := min(100, len(nextRows)*3)
	for i := range nsteps {
		// Occasionally add a Select to reset indexed flag for projMap
		if fc.rnd.IntN(20) == 0 { // 5% chance
			if sels == nil {
				fc.q.Select(nil) // this also rewinds
			} else {
				fc.q.Rewind()
			}
			data.rewind()
		}

		pos := data.pos
		if data.pos == dsEof {
			history += "r"
			fc.q.Rewind()
			data.rewind()
		}
		dir := random([]Dir{Next, Prev}, fc.rnd)
		history += string(dir)
		expectedRow := data.get(dir)
		row := fc.q.Get(fc.th, dir)

		if expectedRow == nil && row != nil {
			fc.t.Fatalf("random walk step %d: %c from %v: expected nil, got row\nhistory %s",
				i, dir, pos, history)
		} else if expectedRow != nil && row == nil {
			fc.t.Log(fc.q)
			fc.t.Fatalf("random walk step %d: %c from %v: expected row, got nil\nhistory %s",
				i, dir, pos, history)
		} else if expectedRow != nil && row != nil {
			if !fc.hdr.EqualRows(row, expectedRow, nil, nil) {
				fc.t.Fatalf("random walk step %d: %c from %v: row mismatch\nhistory %s",
					i, dir, pos, history)
			}
		}
	}

	// Get all rows using Prev
	fc.q.Rewind()
	prevRows := getAllRows(fc.q, Prev, fc.th)
	if !rowSetsEqual(prevRows, qh, fc.hdr) {
		fc.t.Fatalf("Prev iteration returned %d rows, expected %d", len(prevRows), qh.nrows)
	}
}

func getAllRows(q Query, dir Dir, th *Thread) []Row {
	q.Rewind()
	var rows []Row
	for {
		row := q.Get(th, dir)
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
func (fc *fuzzCtx) testCursorPatterns(nextRows []Row) {
	n := len(nextRows)

	check := func(name string, row, expected Row) {
		t := fc.t
		t.Helper()
		if expected == nil && row != nil {
			t.Fatalf("%s: expected nil, got row", name)
		} else if expected != nil && row == nil {
			t.Fatalf("%s: expected row, got nil", name)
		} else if expected != nil && row != nil {
			if !fc.hdr.EqualRows(row, expected, nil, nil) {
				t.Fatalf("%s: row mismatch", name)
			}
		}
	}

	// Pattern 0: Rewind, Prev - should go to last row
	fc.q.Rewind()
	row := fc.q.Get(fc.th, Prev) // last row
	if n > 0 {
		check("Rewind, Prev", row, nextRows[n-1])
	}

	// Pattern 1: Rewind, Next, Prev - after first Next, Prev should return nil
	fc.q.Rewind()
	row = fc.q.Get(fc.th, Next) // first row or nil if empty
	if n > 0 {
		check("Next,Prev: N", row, nextRows[0])
	} else {
		check("Next,Prev: N (empty)", row, nil)
	}
	row = fc.q.Get(fc.th, Prev) // should be nil - nothing before first
	check("Next,Prev: P", row, nil)

	// Pattern 2: Rewind, Prev, Next - Prev from rewind goes to last, then Next should be nil
	if n > 0 {
		fc.q.Rewind()
		row = fc.q.Get(fc.th, Prev) // last row
		check("Prev,Next: P", row, nextRows[n-1])
		row = fc.q.Get(fc.th, Next) // should be nil - nothing after last
		check("Prev,Next: N", row, nil)
	}
	// Pattern 3: Rewind, Prev, Prev, Next, Next
	if n >= 2 {
		fc.q.Rewind()
		row = fc.q.Get(fc.th, Prev) // last row (n-1)
		check("PPNN: P1", row, nextRows[n-1])
		row = fc.q.Get(fc.th, Prev) // second to last (n-2)
		check("PPNN: P2", row, nextRows[n-2])
		row = fc.q.Get(fc.th, Next) // back to last (n-1)
		check("PPNN: N1", row, nextRows[n-1])
		row = fc.q.Get(fc.th, Next) // should be nil
		check("PPNN: N2", row, nil)
	}

	// Pattern 4: Rewind, Next, Next, Prev, Prev
	if n >= 2 {
		fc.q.Rewind()
		row = fc.q.Get(fc.th, Next) // first row (0)
		check("NNPP: N1", row, nextRows[0])
		row = fc.q.Get(fc.th, Next) // second row (1)
		check("NNPP: N2", row, nextRows[1])
		row = fc.q.Get(fc.th, Prev) // back to first (0)
		check("NNPP: P1", row, nextRows[0])
		row = fc.q.Get(fc.th, Prev) // should be nil
		check("NNPP: P2", row, nil)
	}

	// Pattern 5: Next to end, past end (nil), then Prev
	// plain stick at eof: Prev should also return nil
	if n > 0 {
		fc.q.Rewind()
		for i := range n {
			row = fc.q.Get(fc.th, Next)
			check("ToEnd: N"+strconv.Itoa(i), row, nextRows[i])
		}
		row = fc.q.Get(fc.th, Next) // past end
		check("ToEnd: N-past", row, nil)
		row = fc.q.Get(fc.th, Prev) // plain stick: should be nil
		check("ToEnd: P", row, nil)
	}

	// Pattern 6: Prev to beginning, past beginning (nil), then Next
	// plain stick at eof: Next should also return nil
	if n > 0 {
		fc.q.Rewind()
		for i := n - 1; i >= 0; i-- {
			row = fc.q.Get(fc.th, Prev)
			check("ToBegin: P"+strconv.Itoa(n-1-i), row, nextRows[i])
		}
		row = fc.q.Get(fc.th, Prev) // past beginning
		check("ToBegin: P-past", row, nil)
		row = fc.q.Get(fc.th, Next) // plain stick: should be nil
		check("ToBegin: N", row, nil)
	}

	// Pattern 7: Rewind, Next, Prev, Next - plain stick at eof after Prev
	if n > 0 {
		fc.q.Rewind()
		row = fc.q.Get(fc.th, Next) // first
		check("NPN: N1", row, nextRows[0])
		row = fc.q.Get(fc.th, Prev) // nil
		check("NPN: P", row, nil)
		row = fc.q.Get(fc.th, Next) // plain stick: nil
		check("NPN: N2", row, nil)
	}

	// Pattern 8: Rewind, Prev, Next, Prev - plain stick at eof after Next
	if n > 0 {
		fc.q.Rewind()
		row = fc.q.Get(fc.th, Prev) // last
		check("PNP: P1", row, nextRows[n-1])
		row = fc.q.Get(fc.th, Next) // nil
		check("PNP: N", row, nil)
		row = fc.q.Get(fc.th, Prev) // plain stick: nil
		check("PNP: P2", row, nil)
	}

	// Reset for subsequent tests
	fc.q.Rewind()
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

func (fc *fuzzCtx) testRandomSelects(index []string, allRows []Row) {
	fc.t.Helper()
	hdr := fc.q.Header()
	fc.testExistentSelect(allRows, hdr, index)
	fc.testNonExistentSelect(allRows, hdr, index)
}

func (fc *fuzzCtx) testExistentSelect(allRows []Row, hdr *Header, index []string) {
	if len(allRows) == 0 {
		return
	}
	for range 10 {
		srcRow := random(allRows, fc.rnd)
		sels := fc.indexSelectCriteria(srcRow, hdr, index)
		fc.q.Select(sels)

		qh := NewQueryHasher(hdr)
		for _, row := range allRows {
			if selMatchIndex(hdr, row, sels, index) {
				qh.Row(row)
			}
		}

		fc.testRandomGet(qh, sels)

		fc.q.Select(nil) // clear select
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

// indexSelectCriteria uses all columns of the index for select criteria.
func (fc *fuzzCtx) indexSelectCriteria(row Row, hdr *Header, index []string) Sels {
	selCols := slices.Clone(index)
	fc.rnd.Shuffle(len(selCols), func(i, j int) {
		selCols[i], selCols[j] = selCols[j], selCols[i]
	})
	return makeSels(hdr, row, selCols, fc.th, nil)
}

func (fc *fuzzCtx) testNonExistentSelect(allRows []Row, hdr *Header, index []string) {
	for range 10 {
		// If there are no rows, use a dummy row sized to match hdr.Fields.
		// This avoids panics when hdr.Fields references derived records (e.g. Extend).
		srcRow := make(Row, len(hdr.Fields))
		if len(allRows) > 0 {
			srcRow = random(allRows, fc.rnd)
		}
		sels := fc.indexSelectCriteria(srcRow, hdr, index)
		sels[fc.rnd.IntN(len(sels))].val = "nonexistent"
		fc.q.Select(sels)
		if fc.q.Get(fc.th, Next) != nil {
			fc.t.Fatal("non-existent select returned a row")
		}
		fc.q.Select(nil) // clear select
	}
}

//-------------------------------------------------------------------

func (fc *fuzzCtx) testRandomLookups(index []string, allRows []Row) {
	fc.t.Helper()
	if len(allRows) == 0 {
		return
	}
	lookupCols := slices.Clone(index)
	slc.Shuffle(fc.rnd, lookupCols)
	hdr := fc.q.Header()
	cols := hdr.Columns
	for range min(20, len(allRows)) {
		srcRow := random(allRows, fc.rnd)
		sels := makeSels(hdr, srcRow, lookupCols, fc.th, nil)
		if fc.rnd.IntN(2) == 0 {
			result := lookup(fc.q, sels, fc.th, nil)
			if result == nil {
				PrintTableInfo(fc.q)
				for k, v := range fc.ft.rules {
					fmt.Println(k, "=", v)
				}
				fmt.Println("srcRow", RowStr(hdr, srcRow))
				fmt.Println("sels", sels)
				fc.t.Fatal("lookup returned nil for existing key")
			}
			assert.That(rowsEqual(result, srcRow, hdr, cols))
		} else {
			// set one of the keyVals to a non-existent value
			r := fc.rnd.IntN(len(lookupCols))
			if srcRow.GetRaw(hdr, lookupCols[r]) == "" {
				sels[r].val = sels[r].col + "_" + strconv.Itoa(fc.rnd.IntN(100))
			} else if fc.rnd.IntN(2) == 1 {
				sels[r].val = "nonexistent"
			} else {
				sels[r].val = ""
			}
			result := lookup(fc.q, sels, fc.th, nil)
			if result != nil {
				if exists(hdr, allRows, fc.th, sels) {
					continue
				}
				fc.t.Fatal("lookup returned row for non-existent key")
			}
		}
	}
}

func exists(hdr *Header, allRows []Row, th *Thread, sels Sels) bool {
	for _, row := range allRows {
		if lookupFilter(hdr, row, sels, th, nil) != nil {
			return true
		}
	}
	return false
}

func random[E any](list []E, rnd *rand.Rand) E {
	return list[rnd.IntN(len(list))]
}

//-------------------------------------------------------------------
// go test -run '^$' -fuzz=FuzzSplitShare ./dbms/query

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

//-------------------------------------------------------------------

// PrintTableInfo prints the table schema details for each leaf Table in the query
func PrintTableInfo(q Query) {
	switch q := q.(type) {
	case q2i:
		PrintTableInfo(q.Source())
		PrintTableInfo(q.Source2())
	case q1i:
		PrintTableInfo(q.Source())
	case *Table:
		printTableSchema(q)
	case *ProjectNone:
		PrintTableInfo(q.source)
	}
}

func printTableSchema(tbl *Table) {
	fmt.Println(tbl.Name())
	fmt.Printf("  columns: %s\n", str.Join(",", tbl.Columns()))

	fmt.Print("  indexes:")
	for _, ix := range tbl.Indexes() {
		fmt.Printf(" %s", str.Join("(,)", ix))
	}

	fmt.Print("\n  keys:")
	for _, key := range tbl.Keys() {
		fmt.Printf(" %s", str.Join("(,)", key))
	}
	fmt.Println()
}

func PrintQueryData(q Query, th *Thread) {
	switch q := q.(type) {
	case q2i:
		if u, ok := q.(*Union); ok {
			fmt.Println("--- union source1 simple ---")
			printRows(u.Source(), th)
			fmt.Println("--- union source2 simple ---")
			printRows(u.Source2(), th)
		}
		PrintQueryData(q.Source(), th)
		PrintQueryData(q.Source2(), th)
	case q1i:
		PrintQueryData(q.Source(), th)
	case *Table:
		fmt.Println("--- table", q.Name(), "---")
		printRows(q, th)
	}
}

func printRows(q Query, th *Thread) {
	hdr := q.Header()
	q.Rewind()
	rows := q.Simple(th)
	for i, row := range rows {
		fmt.Printf("  row %d: %s\n", i, RowStr(hdr, row))
	}
	if len(rows) == 0 {
		fmt.Println("  (no rows)")
	}
}
