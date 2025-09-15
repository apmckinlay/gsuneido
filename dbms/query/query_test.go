// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
)

func TestKeys(t *testing.T) {
	test := func(query, expected string) {
		t.Helper()
		q := ParseQuery(query, testTran{}, nil)
		assert.T(t).This(idxsToS(q.Keys())).Is(expected)
	}
	test("tables", "table")
	test("columns", "table+column")
	test("columns rename column to col", "table+col")
	test("tables extend x=1,b=2", "table")
	test("tables intersect tables", "table")
	test("abc intersect bcd", "b, c")
	test("hist project item, cost", "item+cost")
	test("hist2 project date, item", "date")
	test("abc times inven", "b+item, c+item")
	test("abc union abc", "a+b+c") // not disjoint
	test("(bcd where b is 1) union (bcd where b is 2)", "b")
}

func TestForeignKeys(t *testing.T) {
	//	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act)
		assert.This(n).Is(1)
	}
	test := func(key, index string) {
		t.Helper()
		// note that lin fields are different column names
		// and in a different order relative to the indexes
		doAdmin(db, "create hdr (a,b) key("+key+")")
		defer doAdmin(db, "drop hdr")
		act("insert { a: '1\x00', b: 2 } into hdr")
		act("insert { a: '3\x00', b: 4 } into hdr")
		doAdmin(db, "create lin (f,e,d) key(e) "+index+" in hdr("+key+")")
		defer doAdmin(db, "drop lin")

		// output (lin) block
		assert.This(func() { act("insert { d: 9, e: 9 } into lin") }).
			Panics("blocked by foreign key")

		// output (lin) NO block
		act("insert { d: '1\x00', e: 2, f: 5 } into lin")

		// update (lin) block
		assert.This(func() { act("update lin set d = 9") }).
			Panics("blocked by foreign key")

		// update (lin) NO block
		act("update lin set d = '3\x00', e = 4")

		// delete (hdr) block
		assert.This(func() { act("delete hdr where a = '3\x00' and b = 4") }).
			Panics("blocked by foreign key")

		// delete (hdr) NO block
		act("delete hdr where a = '1\x00' and b = 2")

		doAdmin(db, "drop lin")
		doAdmin(db, "create lin (f,e,d) key(e) "+index+" in hdr("+key+") cascade")
		act("insert { d: '3\x00', e: 4, f: 5 } into lin")
		
		// update (hdr) cascade
		act("update hdr set a = 33, b = 44")
		assert.This(queryAll(db, "hdr")).Is("a=33 b=44")
		data := queryAll(db, "lin")
		assert.Msg(data).That(strings.Contains(data, "d=33"))

		// delete (hdr) cascade
		act("delete hdr")
		assert.This(queryAll(db, "lin")).Is("")
	}
	test("a", "key(d)")
	test("a", "index(d)")
	test("a", "index(d,e)")
	test("a,b", "index(d,e)")
	test("a,b", "index(d,e,f)")
	db.MustCheck()
}

func TestForeignKeyDeleteBlock(t *testing.T) {
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act)
		assert.This(n).Is(1)
	}
	doAdmin(db, "create hdr (a, b, c) key(a)")
	doAdmin(db, "create lin (d, e, f) key(e) index(d) in hdr(a)")
	act("insert { a: '\x00', b: 'b1', c: 'c1'  } into hdr")
	act("insert { d: '\x00', e: 'e1', f: 'f1' } into lin")
	assert.This(func() { act("delete hdr where a = '\x00'") }).
		Panics("blocked by foreign key")
}

func TestForeignKeyDeleteCascade(t *testing.T) {
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act)
		assert.This(n).Is(1)
	}
	doAdmin(db, "create hdr (a, b, c) key(a)")
	doAdmin(db, "create lin (d, e, f) key(e) index(d,f) in hdr(a) cascade")
	act("insert { a: '\x00', b: 'b1', c: 'c1'  } into hdr")
	act("insert { d: '\x00', e: 'e1', f: 'f1' } into lin")
	act("delete hdr where a is '\x00'")
	// cascade should delete lin
	assert.T(t).This(queryAll(db, "lin")).Is("")
}

func queryAll(db *db19.Database, query string) string {
	tran := sizeTran{db.NewReadTran()}
	q := ParseQuery(query, tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)
	return queryAll2(q)
}

func queryAll2(q Query) string {
	hdr := q.Header()
	sep := ""
	var sb strings.Builder
	th := &Thread{}
	for row := q.Get(th, Next); row != nil; row = q.Get(th, Next) {
		sb.WriteString(sep)
		sb.WriteString(row2str(hdr, row))
		sep = " | "
	}
	return sb.String()
}

func row2str(hdr *Header, row Row) string {
	if row == nil {
		return "nil"
	}
	var sb strings.Builder
	sep := ""
	for _, col := range hdr.Columns {
		val := row.GetVal(hdr, col, nil, nil)
		if val != EmptyStr {
			fmt.Fprint(&sb, sep, col, "=", AsStr(val))
			sep = " "
		}
	}
	return sb.String()
}

func TestSelKeys(t *testing.T) {
	sep := ixkey.Sep
	max := ixkey.Max
	encode := false
	dstCols := []string{"one"}
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
	dstCols = []string{"one", "two"}
	test("1"+sep+"2", "1"+sep+"2"+sep+max)
	dstCols = []string{"a", "b", "c"}
	srcCols = []string{"a"}
	vals = []string{"1"}
	test("1", "1"+sep+max)
}

func TestQueryBug(*testing.T) {
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act)
		assert.This(n).Is(1)
	}
	doAdmin(db, "create tmp (a,b) key(a)")
	act("insert { a: 1 } into tmp")
	assert.This(queryAll(db, "tmp where b > 0")).Is("")
	// inconsistent with the language, but how it has worked historically
}

func TestExtendAllRules(*testing.T) {
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	db := testDb()
	defer db.Close()
	tran := db.NewReadTran()
	q := ParseQuery("cus extend Foo, n=1, Bar", tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)
	assert.That(!q.SingleTable())
	assert.This(len(q.Header().Fields)).Is(2)
	q = ParseQuery("cus extend Foo, Bar", tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)
	assert.That(q.SingleTable())
	assert.This(len(q.Header().Fields)).Is(1)
}

func TestDuplicateKey(*testing.T) {
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act)
		assert.This(n).Is(1)
	}
	doAdmin(db, "create tmp (k,u,i) key(k) index unique(u) index(i)")
	act("insert { k: 1, u: 2, i: 3 } into tmp")
	act("insert { k: 11, u: 22, i: 3 } into tmp")
	assert.This(func() { act("insert { k: 11, u: 0, i: 0 } into tmp") }).
		Panics("duplicate key")
	assert.This(func() { act("insert { k: 0, u: 22, i: 0 } into tmp") }).
		Panics("duplicate key")
	act("insert { k: 111, u: 222, i: 3 } into tmp")
	act("insert { k: 1111, u: '' } into tmp")
	act("insert { k: 11111, u: '' } into tmp")
}

func TestWhereSelectBug(t *testing.T) {
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act)
		assert.This(n).Is(1)
	}
	doAdmin(db, "create t2 (d) key (d)")
	doAdmin(db, "create t1 (a, b, d) key(a) index(b) index(d)")
	act("insert {d: '1'} into t2")
	act("insert {d: '1', a: '2', b: '8'} into t1")
	act("insert {d: '1', a: '3', b: '7'} into t1")
	act("insert {d: '1', a: '4', b: '8'} into t1")
	act("insert {d: '1', a: '5', b: '7'} into t1")
	query := "t1 join t2 where d is '1' and b < 'z'"
	assert.T(t).This(queryAll(db, query)).
		Is("d=1 a=3 b=7 | d=1 a=5 b=7 | d=1 a=2 b=8 | d=1 a=4 b=8")

	tran := sizeTran{db.NewReadTran()}
	q := ParseQuery("t1 where d is '1' and b < 'z'", tran, nil)
	idx := []string{"d"}
	q = q.Transform()
	_, _, app := q.optimize(ReadMode, idx, 1)
	q.setApproach(idx, 1, app, tran)
	assert.T(t).This(String(q)).Is("t1^(b) where d is '1' and b < 'z'")
	vals := []string{Pack(SuStr("1"))}
	q.Select(idx, vals)
	assert.T(t).This(queryAll2(q)).
		Is("a=3 b=7 d=1 | a=5 b=7 d=1 | a=2 b=8 d=1 | a=4 b=8 d=1")
}

func TestJoinBug(t *testing.T) {
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act)
		assert.This(n).Is(1)
	}
	doAdmin(db, "create t1 (a) key(a)")
	doAdmin(db, "create t2 (a, b) key(a,b)")
	act("insert {a: '1'} into t1")
	act("insert {a: '1', b: '2'} into t2")
	assert.T(t).This(queryAll(db, "t1 join t2")).Is("a=1 b=2")
}

func TestSelectOnSingleton(t *testing.T) {
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act)
		assert.This(n).Is(1)
	}
	doAdmin(db, "create t1 (a) key(a)")
	doAdmin(db, "create t2 (a, b) key()")
	act("insert {a: '1'} into t1")
	act("insert {a: '1', b: '2'} into t2")
	assert.T(t).This(queryAll(db, "t1 leftjoin t2")).Is("a=1 b=2")
}

func TestSingleton(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act)
		assert.This(n).Is(1)
	}
	doAdmin(db, "create tmp (a,b) key(a) key(b)")
	act("insert { a: 1, b: 2 } into tmp")
	act("insert { a: 3, b: 4 } into tmp")
	tran := sizeTran{db.NewReadTran()}
	q := ParseQuery("tmp where a = 3", tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)
	assert.This(String(q)).Is("tmp^(a) where*1 a is 3") // singleton
	// reading by a, but singleton so we can Select/Lookup on b
	bcols := []string{"b"}
	bvals := []string{Pack(SuInt(4))}
	q.Select(bcols, bvals)
	assert.This(queryAll2(q)).Is("a=3 b=4")
	hdr := q.Header()
	assert.This(row2str(hdr, q.Lookup(nil, bcols, bvals))).Is("a=3 b=4")

	bvals = []string{Pack(SuInt(2))}
	q.Select(bcols, bvals)
	assert.This(queryAll2(q)).Is("")
	assert.This(q.Lookup(nil, bcols, bvals)).Is(nil)
}

func TestWithoutDupsOrSupersets(t *testing.T) {
	test := func(keys, expected [][]string) {
		result := withoutDupsOrSupersets(keys)
		assert.T(t).This(result).Is(expected)
	}
	test([][]string{}, [][]string{})
	test([][]string{{"a"}}, [][]string{{"a"}})
	test([][]string{{"a"}, {"b", "c"}}, [][]string{{"a"}, {"b", "c"}})
	test([][]string{{"a", "b"}, {"b", "a"}}, [][]string{{"a", "b"}})
	test([][]string{{"a", "b"}, {"b", "a", "c"}}, [][]string{{"a", "b"}})
}

func TestWhereSplitBug(t *testing.T) {
	Global.TestDef("Rule_hx",
		compile.Constant("function() { return .b }"))
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	doAdmin(db, "create tmp1 (a, b, hx) key (a)")
	doAdmin(db, "create tmp2 (a, b) key (a)")

	act(db, "insert { a: 1, b: 2 } into tmp2")
	assert.T(t).This(queryAll(db, "tmp1 union (tmp2 extend hx)")).
		Is("a=1 b=2")
	assert.T(t).This(queryAll(db, "(tmp1 union (tmp2 extend hx)) where hx = 2")).
		Is("a=1 b=2")

	act(db, "insert { a: 1, b: 2, hx: 2 } into tmp1")
	assert.T(t).This(queryAll(db, "tmp1 join (tmp2 extend hx)")).
		Is("a=1 b=2 hx=2")
	assert.T(t).This(queryAll(db, "(tmp1 join (tmp2 extend hx)) where hx = 2")).
		Is("a=1 b=2 hx=2")
}

var result [][]string

func BenchmarkNoOptMod(b *testing.B) {
	orig := [][]string{{"a"}, {"b"}, {"c"}, {"d"}, {"e"}, {"f"}}
	for b.Loop() {
		result = make([][]string, len(orig))
		for _, o := range orig { //nolint:gosimple
			result = append(result, o)
		}
	}
}

func BenchmarkOptMod(b *testing.B) {
	orig := [][]string{{"a"}, {"b"}, {"c"}, {"d"}, {"e"}, {"f"}}
	for b.Loop() {
		om := newOptMod(orig)
		for _, o := range orig {
			om.add(o)
		}
		om.result()
	}
}

func TestJoin_splitSelect(t *testing.T) {
	joinRev = impossible
	defer func() { joinRev = 0 }()
	q1 := newTestQop([]string{"a", "b", "c"})
	q1.indexes = [][]string{{"a"}}
	q2 := newTestQop([]string{"c", "d", "e"})
	q1.indexes = [][]string{{"c"}}
	q2.fixed = []Fixed{
		{col: "c", values: fixvals("1")},
		{col: "e", values: fixvals("2", "")},
	}
	jn := NewJoin(q1, q2, nil, nil).(*Join)
	assert.This(jn.by).Is([]string{"c"})

	cols := []string{"a", "c"}
	vals := fixvals("9", "1")
	jn.Select(cols, vals)
	assert.This(q1.sel).Is(sel{cols: cols, vals: vals})
}

func fixvals(strs ...string) []string {
	return slc.MapFn(strs, func(s string) string {
		return Pack(SuStr(s))
	})
}

type TestQop struct {
	Nothing
	sel
}

func newTestQop(cols []string) *TestQop {
	q := &TestQop{}
	q.header = SimpleHeader(cols)
	return q
}

type sel struct {
	cols []string
	vals []string
}

func (q *TestQop) Indexes() [][]string { // override Nothing
	return q.indexes
}

func (q *TestQop) Keys() [][]string { // override Nothing
	if q.keys == nil {
		return q.indexes
	}
	return q.keys
}

func (q *TestQop) Select(cols, vals []string) {
	q.sel = sel{cols: cols, vals: vals}
}

func (q *TestQop) fastSingle() bool { // override Nothing
	return false
}
