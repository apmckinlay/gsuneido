// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
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

	DoAdmin(db, "create master (m) key(m)")
	act("insert { m: 1 } into master")
	DoAdmin(db, "create detail (m, d) key(d) index(m) in master cascade")
	act("insert { m: 1, d: 10 } into detail")
	act("insert { m: 1, d: 11 } into detail")
	act("insert { m: 1, d: 12 } into detail")
	assert.This(db.GetState().Meta.GetRoInfo("detail").Nrows).Is(3)
	act("delete master") // cascade
	assert.This(db.GetState().Meta.GetRoInfo("detail").Nrows).Is(0)

	DoAdmin(db, "create header (m) key(m)")
	act("insert { m: 1 } into header")
	DoAdmin(db, "create lines (m, d) key(d) index(m) in header cascade")
	act("insert { m: 1, d: 10 } into lines")
	act("insert { m: 1, d: 11 } into lines")
	act("insert { m: 1, d: 12 } into lines")
	assert.This(queryAll(db, "lines")).
		Is("m=1 d=10 | m=1 d=11 | m=1 d=12")
	act("update header set m = 2")
	assert.This(queryAll(db, "lines")).
		Is("m=2 d=10 | m=2 d=11 | m=2 d=12")

	DoAdmin(db, "create one (a) key(a)")
	DoAdmin(db, "create two (b,a) key(b)")
	act("insert { b: 1, a: 1 } into two")
	assert.This(func() { DoAdmin(db, "alter two create index(a) in one") }).
		Panics("blocked by foreign key")

	DoAdmin(db, "ensure test_table1 (a, b, c) key(a)")
	DoAdmin(db, "ensure test_table2 (a, d, e) key(e) index(a) in test_table1 cascade")
	act("insert { a: 'a1', b: 'b1', c: 'c1'  } into test_table1")
	act("insert { a: 'a1', d: 'd1', e: 'e1' } into test_table2")
	act("delete test_table1 where a is 'a1'")
	assert.This(queryAll(db, "test_table2")).Is("")

	db.Check()
}

func queryAll(db *db19.Database, query string) string {
	tran := sizeTran{db.NewReadTran()}
	q := ParseQuery(query, tran)
	q, _ = Setup(q, ReadMode, tran)
	hdr := q.Header()
	sep := ""
	var sb strings.Builder
	for row := q.Get(rt.Next); row != nil; row = q.Get(rt.Next) {
		sep2 := ""
		sb.WriteString(sep)
		for _, col := range hdr.Columns {
			val := row.GetVal(hdr, col, nil, nil)
			fmt.Fprint(&sb, sep2, col, "=", val.String())
			sep2 = " "
		}
		sep = " | "
	}
	return sb.String()
}
