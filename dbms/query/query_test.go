// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestKeys(t *testing.T) {
	test := func(query, expected string) {
		t.Helper()
		q := ParseQuery(query, testTran{})
		assert.T(t).This(idxsToS(q.Keys())).Is(expected)
	}
	test("tables", "table, tablename")
	test("columns", "table+column")
	test("columns rename column to col", "table+col")
	test("tables extend x=1,b=2", "table, tablename")
	test("tables intersect tables", "table, tablename")
	test("abc intersect bcd", "b, c")
	test("hist project item, cost", "item+cost")
	test("hist2 project date, item", "date")
	test("abc times inven", "b+item, c+item")
	test("abc union abc", "a+b+c") // not disjoint
	test("(bcd where b is 1) union (bcd where b is 2)", "b")
}

func TestByContainsKey(t *testing.T) {
	test := func(by string, keys string, expected bool) {
		t.Helper()
		result := containsKey(strings.Fields(by), sToIdxs(keys))
		assert.T(t).This(result).Is(expected)
	}
	test("", "a b, c", false)
	test("a b", "a", true)
	test("a b", "b", true)
	test("a b", "x, a+b, y", true)
}

func TestProjectIndexes(t *testing.T) {
	test := func(idxs string, cols string, expected string) {
		t.Helper()
		result := projectIndexes(sToIdxs(idxs), strings.Fields(cols))
		assert.T(t).This(idxsToS(result)).Is(expected)
	}
	test("a, b+c, d+e+f", "a b c d", "a, b+c")
	test("a, b+c, d+e+f", "c e f", "")
}

// sToIdxs splits strings like: "a+b, c, d+e+f"
func sToIdxs(s string) [][]string {
	var idxs [][]string
	for _, ix := range strings.Split(s, ", ") {
		idxs = append(idxs, strings.Split(ix, "+"))
	}
	return idxs
}

// idxsToS converts [][]string to a string like: "a+b, c, d+e+f"
func idxsToS(idxs [][]string) string {
	tmp := make([]string, len(idxs))
	for i, ix := range idxs {
		tmp[i] = strings.Join(ix, "+")
	}
	return strings.Join(tmp, ", ")
}

func TestForeignKeys(*testing.T) {
	store := stor.HeapStor(8192)
	db, err := db19.CreateDb(store)
	ck(err)
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *rt.SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(ut, act)
		assert.This(n).Is(1)
	}
	DoAdmin(db, "create hdr (a,b) key(a)")
	act("insert { a: 1, b: 2 } into hdr")
	act("insert { a: 3, b: 4 } into hdr")
	DoAdmin(db, "create lin (a,c) key(c) index(a) in hdr")
	act("insert { a: 1, c: 5 } into lin")

	assert.This(func() { act("delete hdr where a = 1") }).
		Panics("blocked by foreign key")
	act("delete hdr where a = 3") // no lin so ok

	assert.This(func() { act("insert { a: 9, c: 6 } into lin") }).
		Panics("blocked by foreign key")
	act("insert { a: '', c: 6 } into lin") // '' allowed

	act("insert { a: '', b: 22 } into hdr")
	act("delete hdr where a = ''")

	assert.This(func() { act("update lin set a = 9") }).
		Panics("blocked by foreign key")
	assert.This(func() { act("update hdr set a = 9") }).
		Panics("blocked by foreign key")
	act("update lin where a = 1 set a = ''") // '' allowed
}
