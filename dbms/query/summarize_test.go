// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSummarizeSelectFilter(t *testing.T) {
	// Test that Select and Lookup properly split and filter on summarized columns
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }

	doAdmin(db, "create test (a, b, c) key(a)")
	act(db, "insert { a: 1, b: 1, c: 10 } into test")
	act(db, "insert { a: 2, b: 1, c: 20 } into test")
	act(db, "insert { a: 3, b: 2, c: 30 } into test")
	act(db, "insert { a: 4, b: 2, c: 25 } into test")

	tran := sizeTran{db.NewReadTran()}
	q := ParseQuery("test summarize a, max c", tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)

	// Test Select with summarized column (max_c)
	// splitSelect should pass a to source, keep max_c for filtering
	// Should filter to only row where a=2 and max_c = 20
	cols := []string{"a", "max_c"}
	vals := []string{Pack(SuInt(2)), Pack(SuInt(20))}
	q.Select(cols, vals)
	row := q.Get(nil, Next)
	assert.T(t).That(row != nil)
	hdr := q.Header()
	assert.T(t).This(row.GetVal(hdr, "a", nil, nil)).Is(SuInt(2))
	assert.T(t).This(row.GetVal(hdr, "max_c", nil, nil)).Is(SuInt(20))

	// Should not get another row
	row = q.Get(nil, Next)
	assert.T(t).This(row).Is(nil)

	// Test with wrong max_c - should return nil
	q.Select(nil, nil) // clear
	vals2 := []string{Pack(SuInt(2)), Pack(SuInt(999))}
	q.Select(cols, vals2)
	row = q.Get(nil, Next)
	assert.T(t).This(row).Is(nil)
}

func TestSummarize_Keys(t *testing.T) {
	srckeys := [][]string{{"a", "b"}, {"c"}}
	keys := projectKeys(srckeys, []string{"x", "y"})
	assert.T(t).This(keys).Is([][]string{{"x", "y"}})
	keys = projectKeys(srckeys, []string{})
	assert.T(t).This(keys).Is([][]string{{}})
}

func TestSummarize_Indexes(t *testing.T) {
	srcidxs := [][]string{{"a", "b"}, {"c"}}
	idxs := projectIndexes(srcidxs, []string{"x", "y"})
	assert.T(t).This(idxs).Is(nil)
	idxs = projectIndexes(srcidxs, []string{})
	assert.T(t).This(idxs).Is(nil)
}
