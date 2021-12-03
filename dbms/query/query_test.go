// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
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

	DoAdmin(db, "create hdr1 (a,b) key(a)")
	act("insert { a: 1, b: 2 } into hdr1")
	act("insert { a: 3, b: 4 } into hdr1")
	DoAdmin(db, "create lin1 (a,c) key(c) index(a) in hdr1")
	act("insert { a: 1, c: 5 } into lin1")

	assert.This(func() { act("delete hdr1 where a = 1") }).
		Panics("blocked by foreign key")
	act("delete hdr1 where a = 3") // no lin1 so ok

	assert.This(func() { act("insert { a: 9, c: 6 } into lin1") }).
		Panics("blocked by foreign key")
	act("insert { a: '', c: 6 } into lin1") // '' allowed

	act("insert { a: '', b: 22 } into hdr1")
	act("delete hdr1 where a = ''")

	assert.This(func() { act("update lin1 set a = 9") }).
		Panics("blocked by foreign key")
	assert.This(func() { act("update hdr1 set a = 9") }).
		Panics("blocked by foreign key")
	act("update lin1 where a = 1 set a = ''") // '' allowed

	DoAdmin(db, "create hdr2 (m) key(m)")
	act("insert { m: 1 } into hdr2")
	DoAdmin(db, "create lin2 (m, d) key(d) index(m) in hdr2 cascade")
	act("insert { m: 1, d: 10 } into lin2")
	act("insert { m: 1, d: 11 } into lin2")
	act("insert { m: 1, d: 12 } into lin2")
	assert.This(db.GetState().Meta.GetRoInfo("lin2").Nrows).Is(3)
	act("delete hdr2") // cascade
	assert.This(db.GetState().Meta.GetRoInfo("lin2").Nrows).Is(0)

	DoAdmin(db, "create hdr3 (m) key(m)")
	act("insert { m: 1 } into hdr3")
	DoAdmin(db, "create lin3 (m, d) key(d) index(m) in hdr3 cascade")
	act("insert { m: 1, d: 10 } into lin3")
	act("insert { m: 1, d: 11 } into lin3")
	act("insert { m: 1, d: 12 } into lin3")
	assert.This(queryAll(db, "lin3")).
		Is("m=1 d=10 | m=1 d=11 | m=1 d=12")
	act("update hdr3 set m = 2")
	assert.This(queryAll(db, "lin3")).
		Is("m=2 d=10 | m=2 d=11 | m=2 d=12")

	DoAdmin(db, "create hdr4 (a) key(a)")
	DoAdmin(db, "create lin4 (b,a) key(b)")
	act("insert { b: 1, a: 1 } into lin4")
	assert.This(func() { DoAdmin(db, "alter lin4 create index(a) in hdr4") }).
		Panics("blocked by foreign key")

	DoAdmin(db, "ensure hdr5 (a, b, c) key(a)")
	DoAdmin(db, "ensure lin5 (a, d, e) key(e) index(a) in hdr5 cascade")
	act("insert { a: 'a1', b: 'b1', c: 'c1'  } into hdr5")
	act("insert { a: 'a1', d: 'd1', e: 'e1' } into lin5")
	act("delete hdr5 where a is 'a1'")
	assert.This(queryAll(db, "lin5")).Is("")

	// requires encode
	DoAdmin(db, "create hdr6 (a, b, c) key(a)")
	DoAdmin(db, "create lin6 (a, d, e) key(e) index(a) in hdr6 cascade")
	act("insert { a: #20211110.132155918, b: 'b1', c: 'c1'  } into hdr6")
	act("insert { a: #20211110.132155918, d: 'd1', e: 'e1' } into lin6")
	act("update hdr6 where a is #20211110.132155918 set a = #20211110.132155919")
	act("delete hdr6 where a is #20211110.132155919")
	assert.This(queryAll(db, "lin6")).Is("")

	// requires rangeEnd in fkeyDeleteCascade
	DoAdmin(db, "ensure hdr7 (a, b, c) key(a,b)")
	DoAdmin(db, "ensure lin7 (a, b, d, e) key(e) index(a,b) in hdr7 cascade")
	act("insert { a: 'a1', c: 'c1'  } into hdr7")
	act("insert { a: 'a1', b: 'b2', c: 'c2'  } into hdr7")
	act("insert { a: 'a1', d: 'd1', e: 'e1' } into lin7")
	act("insert { a: 'a1', b: 'b2', d: 'd2', e: 'e2' } into lin7")
	act("delete hdr7 where a is 'a1' and b is ''")
	assert.This(queryAll(db, "lin7")).Is("a=a1 b=b2 d=d2 e=e2")

	// requires rangeEnd in fkeyUpdateCascade
	DoAdmin(db, "ensure hdr8 (a, b, c) key(a,b)")
	DoAdmin(db, "ensure lin8 (a, b, d, e) key(e) index(a,b) in hdr8 cascade")
	act("insert { a: 'a1', c: 'c1'  } into hdr8")
	act("insert { a: 'a1', b: 'b2', c: 'c2'  } into hdr8")
	act("insert { a: 'a1', d: 'd1', e: 'e1' } into lin8")
	act("insert { a: 'a1', b: 'b2', d: 'd2', e: 'e2' } into lin8")
	act("update hdr8 where a is 'a1' and b is '' set a = 'a0'")
	assert.This(queryAll(db, "lin8")).Is("a=a0 b= d=d1 e=e1 | a=a1 b=b2 d=d2 e=e2")

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
			fmt.Fprint(&sb, sep2, col, "=", rt.AsStr(val))
			sep2 = " "
		}
		sep = " | "
	}
	return sb.String()
}

func TestSelKeys(t *testing.T) {
	sep := ixkey.Sep
	max := ixkey.Max
	encode := false
	dstCols := []string{"one", "two"}
	srcCols := []string{"two", "one"}
	vals := []string{"2", "1"}
	test := func(org, end string) {
		t.Helper()
		o, e := selKeys(encode, dstCols, srcCols, vals)
		assert.Msg("org").This(o).Is(org)
		assert.Msg("end").This(e).Is(end)
	}
	test("1", "1\x00")
	encode = true
	test("1"+sep+"2", "1"+sep+"2"+sep+max)
	dstCols = []string{"a", "b", "c"}
	srcCols = []string{"a"}
	vals = []string{"1"}
	test("1", "1"+sep+sep+sep+max)
}

func TestQueryBug(*testing.T) {
	db, err := db19.CreateDb(stor.HeapStor(8192))
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
	DoAdmin(db, "create tmp (a,b) key(a)")
	act("insert { a: 1 } into tmp")
	assert.This(queryAll(db, "tmp where b > 0")).Is("")
}

func TestExtendAllRules(*testing.T) {
	MakeSuTran = func(qt QueryTran) *rt.SuTran { return nil }
	db := testDb()
	defer db.Close()
	tran := db.NewReadTran()
	q := ParseQuery("cus extend Foo, n=1, Bar", tran)
	q, _ = Setup(q, ReadMode, tran)
	assert.That(!q.SingleTable())
	assert.This(len(q.Header().Fields)).Is(2)
	q = ParseQuery("cus extend Foo, Bar", tran)
	q, _ = Setup(q, ReadMode, tran)
	assert.That(q.SingleTable())
	assert.This(len(q.Header().Fields)).Is(1)
}

func TestDuplicateKey(*testing.T) {
	db, err := db19.CreateDb(stor.HeapStor(8192))
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
	DoAdmin(db, "create tmp (k,u,i) key(k) index unique(u) index(i)")
	act("insert { k: 1, u: 2, i: 3 } into tmp")
	act("insert { k: 11, u: 22, i: 3 } into tmp")
	assert.This(func(){ act("insert { k: 11, u: 0, i: 0 } into tmp") }).
		Panics("duplicate key")
	assert.This(func(){ act("insert { k: 0, u: 22, i: 0 } into tmp") }).
		Panics("duplicate key")
	act("insert { k: 111, u: 222, i: 3 } into tmp")
	act("insert { k: 1111, u: '' } into tmp")
	act("insert { k: 11111, u: '' } into tmp")
}
