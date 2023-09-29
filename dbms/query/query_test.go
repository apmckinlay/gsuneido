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
)

func TestKeys(t *testing.T) {
	test := func(query, expected string) {
		t.Helper()
		q := ParseQuery(query, testTran{}, nil)
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

func TestForeignKeys(*testing.T) {
	store := stor.HeapStor(8192)
	db, err := db19.CreateDb(store)
	ck(err)
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act, nil)
		assert.This(n).Is(1)
	}

	doAdmin(db, "create hdr1 (a,b) key(a)")
	act("insert { a: 1, b: 2 } into hdr1")
	act("insert { a: 3, b: 4 } into hdr1")
	doAdmin(db, "create lin1 (a,c) key(c) index(a) in hdr1")
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

	doAdmin(db, "create hdr2 (m) key(m)")
	act("insert { m: 1 } into hdr2")
	doAdmin(db, "create lin2 (m, d) key(d) index(m) in hdr2 cascade")
	act("insert { m: 1, d: 10 } into lin2")
	act("insert { m: 1, d: 11 } into lin2")
	act("insert { m: 1, d: 12 } into lin2")
	assert.This(db.GetState().Meta.GetRoInfo("lin2").Nrows).Is(3)
	act("delete hdr2") // cascade delete
	assert.This(db.GetState().Meta.GetRoInfo("lin2").Nrows).Is(0)

	doAdmin(db, "create hdr3 (m) key(m)")
	act("insert { m: 1 } into hdr3")
	doAdmin(db, "create lin3 (d, m) key(d) index(m) in hdr3 cascade")
	act("insert { m: 1, d: 10 } into lin3")
	act("insert { m: 1, d: 11 } into lin3")
	act("insert { m: 1, d: 12 } into lin3")
	assert.This(queryAll(db, "lin3")).
		Is("d=10 m=1 | d=11 m=1 | d=12 m=1")
	act("update hdr3 set m = 2") // cascade update
	assert.This(queryAll(db, "lin3")).
		Is("d=10 m=2 | d=11 m=2 | d=12 m=2")

	doAdmin(db, "create hdr4 (a) key(a)")
	doAdmin(db, "create lin4 (b,a) key(b)")
	act("insert { b: 1, a: 1 } into lin4")
	assert.This(func() { doAdmin(db, "alter lin4 create index(a) in hdr4") }).
		Panics("blocked by foreign key")
	act("insert { b: 2, a: 1 } into lin4")
	assert.This(func() { doAdmin(db, "alter lin4 create key(a)") }).
		Panics("duplicate")

	doAdmin(db, "ensure hdr5 (a, b, c) key(a)")
	doAdmin(db, "ensure lin5 (a, d, e) key(e) index(a) in hdr5 cascade")
	act("insert { a: 'a1', b: 'b1', c: 'c1'  } into hdr5")
	act("insert { a: 'a1', d: 'd1', e: 'e1' } into lin5")
	act("delete hdr5 where a is 'a1'")
	assert.This(queryAll(db, "lin5")).Is("")

	// requires encode
	doAdmin(db, "create hdr6 (a, b, c) key(a)")
	doAdmin(db, "create lin6 (a, d, e) key(e) index(a) in hdr6 cascade")
	act("insert { a: #20211110.132155918, b: 'b1', c: 'c1'  } into hdr6")
	act("insert { a: #20211110.132155918, d: 'd1', e: 'e1' } into lin6")
	act("update hdr6 where a is #20211110.132155918 set a = #20211110.132155919")
	act("delete hdr6 where a is #20211110.132155919")
	assert.This(queryAll(db, "lin6")).Is("")

	// requires rangeEnd in fkeyDeleteCascade
	doAdmin(db, "ensure hdr7 (a, b, c) key(a,b)")
	doAdmin(db, "ensure lin7 (a, b, d, e) key(e) index(a,b) in hdr7 cascade")
	act("insert { a: 'a1', c: 'c1'  } into hdr7")
	act("insert { a: 'a1', b: 'b2', c: 'c2'  } into hdr7")
	act("insert { a: 'a1', d: 'd1', e: 'e1' } into lin7")
	act("insert { a: 'a1', b: 'b2', d: 'd2', e: 'e2' } into lin7")
	act("delete hdr7 where a is 'a1' and b is ''")
	assert.This(queryAll(db, "lin7")).Is("a=a1 b=b2 d=d2 e=e2")

	// requires rangeEnd in fkeyUpdateCascade
	doAdmin(db, "ensure hdr8 (a, b, c) key(a,b)")
	doAdmin(db, "ensure lin8 (a, b, d, e) key(e) index(a,b) in hdr8 cascade")
	act("insert { a: 'a1', c: 'c1'  } into hdr8")
	act("insert { a: 'a1', b: 'b2', c: 'c2'  } into hdr8")
	act("insert { a: 'a1', d: 'd1', e: 'e1' } into lin8")
	act("insert { a: 'a1', b: 'b2', d: 'd2', e: 'e2' } into lin8")
	act("update hdr8 where a is 'a1' and b is '' set a = 'a0'")
	assert.This(queryAll(db, "lin8")).Is("a=a0 d=d1 e=e1 | a=a1 b=b2 d=d2 e=e2")

	db.Check()
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
	db, err := db19.CreateDb(stor.HeapStor(8192))
	ck(err)
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act, nil)
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
	db, err := db19.CreateDb(stor.HeapStor(8192))
	ck(err)
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act, nil)
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
	db, err := db19.CreateDb(stor.HeapStor(8192))
	ck(err)
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act, nil)
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
	assert.T(t).This(q.String()).Is("t1^(b) WHERE d is '1' and b < 'z'")
	vals := []string{Pack(SuStr("1"))}
	q.Select(idx, vals)
	assert.T(t).This(queryAll2(q)).
		Is("a=3 b=7 d=1 | a=5 b=7 d=1 | a=2 b=8 d=1 | a=4 b=8 d=1")
}

func TestJoinBug(t *testing.T) {
	db, err := db19.CreateDb(stor.HeapStor(8192))
	ck(err)
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act, nil)
		assert.This(n).Is(1)
	}
	doAdmin(db, "create t1 (a) key(a)")
	doAdmin(db, "create t2 (a, b) key(a,b)")
	act("insert {a: '1'} into t1")
	act("insert {a: '1', b: '2'} into t2")
	assert.T(t).This(queryAll(db, "t1 join t2")).Is("a=1 b=2")
}

func TestSelectOnSingleton(t *testing.T) {
	db, err := db19.CreateDb(stor.HeapStor(8192))
	ck(err)
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act, nil)
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
	db, err := db19.CreateDb(stor.HeapStor(8192))
	ck(err)
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act, nil)
		assert.This(n).Is(1)
	}
	doAdmin(db, "create tmp (a,b) key(a) key(b)")
	act("insert { a: 1, b: 2 } into tmp")
	act("insert { a: 3, b: 4 } into tmp")
	tran := sizeTran{db.NewReadTran()}
	q := ParseQuery("tmp where a = 3", tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)
	assert.This(q.String()).Is("tmp^(a) WHERE*1 a is 3") // singleton
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
	db, err := db19.CreateDb(stor.HeapStor(8192))
	ck(err)
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
	for i := 0; i < b.N; i++ {
		result = make([][]string, len(orig))
		//lint:ignore S1011 testing
		for _, o := range orig {
			result = append(result, o)
		}
	}
}

func BenchmarkOptMod(b *testing.B) {
	orig := [][]string{{"a"}, {"b"}, {"c"}, {"d"}, {"e"}, {"f"}}
	for i := 0; i < b.N; i++ {
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
		{col: "c", values: []string{"1"}},
		{col: "e", values: []string{"2", ""}},
	}
	jn := NewJoin(q1, q2, nil)
	assert.This(jn.by).Is([]string{"c"})
	jn.saIndex = []string{"a"}

	cols := []string{"a", "c"}
	vals := []string{"9", "1"}
	jn.Select(cols, vals)
	assert.This(q1.sel).
		Is(sel{cols: []string{"a"}, vals: []string{"9"}})
	assert.That(!jn.conflict1 && !jn.conflict2)
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
