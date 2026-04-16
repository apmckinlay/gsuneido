// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestWhere_perField(t *testing.T) {
	test := func(query string, expected string) {
		t.Helper()
		w := ParseQuery("table where "+query, testTran{}, nil).(*Where)
		actual := "conflict"
		if !w.conflict {
			actual = fmt.Sprint(w.colSels)[3:]
		}
		assert.T(t).Msg(query).This(actual).Is(expected)
	}
	// nothing indexable
	test("Foo()", "[]")
	test("a =~ 'x'", "[]")
	test("a", "[]")
	// binary
	test("a is 123", "[a:[123]]")
	test("a isnt 123", "[a:[<123 >123]]")
	test("a < 123", "[a:[<123]]")
	test("a <= 123", "[a:[<=123]]")
	test("a > 123", "[a:[>123]]")
	test("a >= 123", "[a:[>=123]]")
	// compare to ""
	test("a isnt ''", "[a:[>'']]")
	test("a < ''", "conflict")
	test("a <= ''", "[a:['']]")
	test("a > ''", "[a:[>'']]")
	test("a >= ''", "[a:[<max]]") // everything, always matches
	// in
	test("a in (3,1,2)", "[a:[1 2 3]]")
	// range
	test("a > 3 and a < 6", "[a:[>3_<6]]")
	test("a >= 3 and a <= 6", "[a:[>=3_<=6]]")
	// type
	test("String?(a)", "[a:['' >=PackString_<PackDate]]")
	test("Number?(a)", "[a:[>=PackMinus_<PackString]]")
	test("Date?(a)", "[a:[>=PackDate_<PackDate+1]]")
	// or
	test("a is 0 or b is 0", "[]")
	test("a is 0 or a =~ 'x'", "[]")
	test("a is 0 or a > 6", "[a:[0 >6]]")
	test("a is false and (b is 1 or b is 2)", "[a:[false] b:[1 2]]")
	test("a is 1 or a is 2 or a is 3", "[a:[1 2 3]]")
	test("a is 1 or a is 2 or a is 1", "[a:[1 2]]")
	test("a is 1 or a in (1,2,3) or a is 2", "[a:[1 2 3]]")
	test("a > 3 or a > 6", "[a:[>3]]")
	test("a > 3 or a is 6", "[a:[>3]]")
	test("a > 6 or a > 3", "[a:[>3]]")
	test("a is 6 or a > 3", "[a:[>3]]")
	test("Number?(a) or String?(a)", "[a:['' >=PackMinus_<PackDate]]")
	test("String?(a) or Number?(a)", "[a:['' >=PackMinus_<PackDate]]")
	// multiple
	test("a is 123 and b is 456", "[a:[123] b:[456]]")
	test("a isnt 'm' and a > 'a'", "[a:[>'a'_<'m' >'m']]")
	test("a isnt '' and a > 'm'", "[a:[>'m']]")
	test("a isnt '' and a in ('','m')", "[a:['m']]")
	test("a isnt '' and String?(a)", "[a:[>=PackString_<PackDate]]")
	// intersect
	test("a in (1,2,3) and a in (1,2,3)", "[a:[1 2 3]]")
	test("a in (1,2,3) and a in (2,3,4)", "[a:[2 3]]")
	test("a in (1,2,3) and a is 2", "[a:[2]]")
	test("a in (1,2,3) and a >= 2", "[a:[2 3]]")
	test("a < 5 and a > 2", "[a:[>2_<5]]")
	test("a < 5 and a <= 5", "[a:[<5]]")
	test("a >= 5 and a > 5", "[a:[>5]]")
	// conflict
	test("a in (1,2,3) and a in (4,5,6)", "conflict")
	test("a in (1,3,5) and a in (2,4,6)", "conflict")
	test("a < 5 and a > 6", "conflict")
}

func TestWhere_span_none(t *testing.T) {
	x := span{}
	assert.T(t).That(x.none())
	x.org = side{val: "a"}
	assert.T(t).That(x.none())
}

func TestWhere_indexSpans(t *testing.T) {
	var idx []string
	test := func(query string, expected string) {
		t.Helper()
		w := ParseQuery(query, testTran{}, nil).(*Where)
		pf, _ := perField(w.expr.Exprs, w.source.Header().Physical())
		idxSpans := indexSpans(idx, pf)
		assert.T(t).This(fmt.Sprint(idxSpans)).Is("[" + expected + "]")
	}

	idx = []string{"a", "b", "c"}
	test("comp where a is 1", "[1]")
	test("comp where a is 1 and c is 2", "[1]")
	test("comp where a is 1 and b is 2", "[1] [2]")
	test("comp where a is 1 and b is 2 and c is 3", "[1] [2] [3]")
	test("comp where a >= 4", "[>=4]")
	test("comp where a >= 4 and b is 2", "[>=4]")
	test("comp where a is 2 and b >= 4", "[2] [>=4]")
	test("comp where a in (1,2) and b in (3,4)", "[1 2] [3 4]")
	test("comp where a is '' and b isnt 0", "[''] [<0 >0]")
	test("comp where a > ''", "[>'']")

	idx = []string{"id"}
	test("customer where id is 'e'", "['e']")
}

func TestWhere_explodeIndexSpans(t *testing.T) {
	idx := []string{"a", "b", "c"}
	test := func(query string, expected string) {
		t.Helper()
		w := ParseQuery("comp where "+query, testTran{}, nil).(*Where)
		pf, _ := perField(w.expr.Exprs, w.source.Header().Physical())
		idxSpans := indexSpans(idx, pf)
		exploded := explodeIndexSpans(idxSpans, [][]span{nil})
		assert.T(t).This(fmt.Sprint(exploded)).Is("[" + expected + "]")
	}
	test("a is 1", "[1]")
	test("a is 1 and b is 2", "[1 2]")
	test("a is 1 and b is 2 and c is 3", "[1 2 3]")
	test("a >= 4", "[>=4]")
	test("a is 2 and b >= 4", "[2 >=4]")
	test("a in (1,2) and b in (3,4)", "[1 3] [1 4] [2 3] [2 4]")
	test("a in (1,2) and b >= 4", "[1 >=4] [2 >=4]")
	test("a is '' and b isnt 0", "['' <0] ['' >0]")
}

func TestWhere_perIndex(t *testing.T) {
	table := ""
	test := func(query string, expected string) {
		t.Helper()
		w := ParseQuery(table+" where "+query, testTran{}, nil).(*Where)
		w.optInit() // runs perIndex
		assert.T(t).This(fmt.Sprint(w.idxSels)).Is("[" + expected + "]")
	}

	table = "comp" // key(a,b,c) nrows = 1000
	test("a is 1",
		"(a,b,c) a: <1..1,max> = .1")
	test("a is 1 and b is 2",
		"(a,b,c) a,b: <1,2..1,2,max> = .01")
	test("a is 1 and b is 2 and c is 3",
		"(a,b,c) a,b,c: <1,2,3> = .0005")

	test("a > 4",
		"(a,b,c) a: <4,max..max> = .5")
	test("a <= 4",
		"(a,b,c) a: <..4,max> = .5")
	test("a is 2 and b >= 4",
		"(a,b,c) a,b: <2,4..2,max> = .06")
	test("a in (1,2) and b in (3,4)",
		"(a,b,c) a,b: <1,3..1,3,max | 1,4..1,4,max | "+
			"2,3..2,3,max | 2,4..2,4,max> = .04")
	test("a in (1,2) and b > 4",
		"(a,b,c) a,b: <1,4,max..1,max | 2,4,max..2,max> = .1")
	test("a is 1 or a > 3",
		"(a,b,c) a: <1..1,max | 3,max..max> = .7")
	test("a isnt 5",
		"(a,b,c) a: <..5 | 5,max..max> = .9")
	test("a is '' and b isnt 0",
		"(a,b,c) a,b: <..'',0 | '',0,max..'',max> = .09")

	test("b is 2",
		"(a,b,c) +b: <2..2,max> = .01")
	test("F(b)",
		"(a,b,c) b = 1 .071")
	test("b in (2,3)",
		"(a,b,c) b = 1 .1")
	test("b is 1 and c in (2,3)",
		"(a,b,c) +b: <1..1,max> c = .01 .0032")
	test("b is 2 and c is 3",
		"(a,b,c) +b,c: <2,3..2,3,max> = .01")
	test("a is 1 and c is 3",
		"(a,b,c) a: <1..1,max> +c: <3..3,max> = .01")

	table = "table" // key(a) nrows = 100
	test("a >= ''",
		"(a) a: <''..max> = 1") //TODO skip useless

	table = "comp2" // key(a,b,c) nrows = 0
	test("a is 1 and b is 2 and c is 3",
		"(a,b,c) a,b,c: <1,2,3> = 0")
}

type wtestTran struct {
	testTran
}

func (t wtestTran) RangeFrac(table string, iIndex int, org, end string) float64 {
	return .5
}

func TestWhere_consistent(t *testing.T) {
	assert := assert.T(t)
	strs := []string{"0", "1", "-1", "''", "'foo'", "true", "false", "#20230812"}
	vals := []Value{Zero, One, MinusOne, EmptyStr, SuStr("foo"), True, False,
		DateFromLiteral("20230812")}
	cols := []string{"a"}
	hdr := SimpleHeader(cols)
	test := func(lhs int, op string, rhs int) {
		t.Helper()
		rec := new(RecordBuilder).Add(vals[lhs].(Packable)).Build()
		c := &ast.RowContext{Hdr: hdr, Row: []DbRec{{Record: rec}}}
		query := "table where a " + op + " " + strs[rhs]
		w := ParseQuery(query, wtestTran{}, nil).(*Where)
		w.optInit()
		e := w.expr.Exprs[0]
		assert.That(e.CanEvalRaw(cols))
		eval := e.Eval(c) == True
		// fmt.Println(eval, "\t", strs[lhs], op, strs[rhs])
		var ixrange bool
		if !w.conflict {
			idxsel := w.idxSels[0]
			packed := Pack(vals[lhs].(Packable))
			for _, pr := range idxsel.prefixRanges {
				if pr.isPoint() {
					ixrange = ixrange || packed == pr.Org
					// fmt.Printf("%q == %q\n", packed, pr.org)
				} else { // range
					ixrange = ixrange || pr.Org <= packed && packed < pr.End
					// fmt.Printf("%q <= %q < %q\n", pr.org, packed, pr.end)
				}
			}
		}
		assert.Msg(strs[lhs], op, strs[rhs]).This(ixrange).Is(eval)
	}
	ops := []string{"is", "isnt", "<", "<=", ">", ">="}
	for lhs := range vals {
		for _, op := range ops {
			for rhs := range vals {
				test(lhs, op, rhs)
			}
		}
	}
}

func TestFracPos(t *testing.T) {
	tt := testTran{}
	test := func(expected float64, digits ...int) {
		t.Helper()
		var enc ixkey.Encoder
		for _, d := range digits {
			enc.Add(Pack(SuInt(d)))
		}
		key := enc.String()
		f := tt.fracPos(key, true)
		assert.T(t).That(math.Abs(f-expected) < .0001)
	}
	test(0)
	test(.5, 5)
	test(.234, 2, 3, 4)
}

func TestWhere_Nrows(t *testing.T) {
	test := func(query string, nrows int) {
		t.Helper()
		var tran testTran
		w := ParseQuery(query, tran, nil)
		Setup(w, ReadMode, tran)
		n, p := w.Nrows()
		assert.T(t).This(n).Is(nrows)
		assert.T(t).This(p).Is(100)
	}
	test("table where F()", 50)
	test("inven where item >= 5", 50)
	test("inven where item < 3 and item > 3", 0) // conflict
	test("inven where item is 1", 1)
	test("inven where item in (1,2,3,4)", 2)
	test("inven where item > 2 and item < 4", 20)
	test("inven where item > 2 and item < 4 and qty", 14)
	test("hist where date is 3", 10)
	test("inven extend x where x > 5", 50) // not on table
}

func TestWhere_Select(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create lin(a,b,c) key(a,b)")
	db.act("insert { a: 1, b: 2, c: 3 } into lin")
	db.act("insert { a: 4, b: 5, c: 6 } into lin")
	db.act("insert { a: 7, b: 5, c: 8 } into lin")
	db.act("insert { a: 9, b: 0, c: 3 } into lin")

	query := "lin where b = 5"
	tran := db.NewReadTran()
	q := ParseQuery(query, tran, nil)
	cols := []string{"a", "b"}
	q = SetupIdx(q, CursorMode, tran, cols)
	sels := Sels{{"a", Pack(IntVal(4))}, {"b", Pack(IntVal(5))}}
	q.Select(sels)
	assert.This(queryAll2(q)).Is("a=4 b=5 c=6")
	q.Select(nil)
	assert.This(queryAll2(q)).Is("a=4 b=5 c=6 | a=7 b=5 c=8")

	q = ParseQuery(query, tran, nil)
	q = SetupIdx(q, CursorMode, tran, cols)
	sels = Sels{{"a", Pack(IntVal(1))}, {"b", Pack(IntVal(2))}} // conflict
	q.Select(sels)
	assert.This(queryAll2(q)).Is("")
	q.Select(nil)
	assert.This(queryAll2(q)).Is("a=4 b=5 c=6 | a=7 b=5 c=8")

	// select with col not in index fields (c not in key(a,b))
	q = ParseQuery(query, tran, nil)
	q = SetupIdx(q, CursorMode, tran, cols)
	q.Select(Sels{{"a", Pack(IntVal(4))}, {"c", Pack(IntVal(6))}})
	assert.This(queryAll2(q)).Is("a=4 b=5 c=6")
}

func TestWhere_fixed(t *testing.T) {
	test := func(query, expected string) {
		t.Helper()
		w := ParseQuery("table where "+query, testTran{}, nil).(*Where)
		assert.T(t).This(fixedStr(w.Fixed())).Is(expected)
	}
	test("a", "[]")
	test("a is 1", "[a=(1)]")
	test("a is 1 and b is 2", "[a=(1), b=(2)]")
	test("a in (1,2,3)", "[a=(1,2,3)]")
	test("a is 1 and a is 1", "[a=(1)]")
	test("a in (1,2) and a is 1", "[a=(1)]")
	test("a in (1,2,3) and a in (2,3,4)", "[a=(2,3)]")
	test("a in (1,2) and a in (3,4)", "[]")
}

func TestWhere_indexes(t *testing.T) {
	db := heapDb()
	defer db.Close()
	test := func(schema, where, colSels, idxSels string) {
		const table = "twi"
		t.Helper()
		db.adm("create " + table + " " + schema)
		defer db.adm("drop " + table)
		tran := db.NewReadTran()
		w := ParseQuery(table+" where "+where, tran, nil).(*Where)
		actual := "conflict"
		if !w.conflict {
			actual = fmt.Sprint(w.colSels)[3:]
		}
		assert.T(t).This(actual).Is(colSels)

		w.optInit()
		assert.T(t).This(fmt.Sprint(w.idxSels)).Is(idxSels)
	}
	test("(a,b,c) key(a)", "a = 1", "[a:[1]]", "[(a) a: <1> = 0]")
}

func TestWhere_ixfilter(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create table (a,b,c,d) key(a) index(b,c)")
	db.act("insert { a: 1, b: 2, c: 3, d: 4 } into table")
	db.act("insert { a: 4, b: 5, c: 6, d: 7 } into table")
	db.act("insert { a: 7, b: 5, c: 8, d: 9 } into table")
	db.act("insert { a: 9, b: 0, c: 3, d: 4 } into table")
	tran := db.NewReadTran()
	q := ParseQuery("table where c=8", tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)
	assert.This(Strategy2(q)).Like(`
		table^(b,c)
		where c is 8`)
}

func TestWhere_idxSel_plus_indexFilter(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create table (a,b,c) key(a,b,c)")
	db.act("insert { a: 1, b: 1, c: 2 } into table")
	db.act("insert { a: 2, b: 1, c: 2 } into table")
	db.act("insert { a: 2, b: 2, c: 3 } into table")
	db.act("insert { a: 3, b: 1, c: 2 } into table")
	db.act("insert { a: 4, b: 1, c: 1 } into table")

	tran := db.NewReadTran()
	q := ParseQuery("table where a > 1 and c = 2", tran, nil)
	q, _, _ = Setup(q, CursorMode, tran)
	w, ok := q.(*Where)
	assert.T(t).That(ok)
	assert.T(t).That(w.idxSelBase != nil)
	assert.T(t).That(w.idxSelBase.prefixLen == 1)
	assert.T(t).That(w.ixExpr != nil)
	assert.This(queryAll2(q)).Is("a=2 b=1 c=2 | a=3 b=1 c=2")
}

func TestWhere_skipScan_pure(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create table (a,b) key(a,b)")
	for a := range 3 {
		for b := range 10 {
			db.act(fmt.Sprintf("insert { a: %d, b: %d } into table", a+1, b+1))
		}
	}

	tran := db.NewReadTran()
	q := ParseQuery("table where b = 5", tran, nil)
	q, _, _ = Setup(q, CursorMode, tran)
	w := q.(*Where)
	assert.T(t).This(fmt.Sprint(w.idxSelBase)).Is("(a,b) +b: <5..5,max> = .01")
	assert.This(queryAll2(q)).Is("a=1 b=5 | a=2 b=5 | a=3 b=5")
}

func TestWhere_skipScan_idxSel(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create table (a,b) key(a,b)")
	for a := range 5 {
		for b := range 5 {
			db.act(fmt.Sprintf("insert { a: %d, b: %d } into table", a+1, b+1))
		}
	}

	tran := db.NewReadTran()
	q := ParseQuery("table where a > 2 and b = 3", tran, nil)
	q, _, _ = Setup(q, CursorMode, tran)
	w, ok := q.(*Where)
	assert.T(t).That(ok)
	assert.T(t).That(w.idxSelBase != nil)
	assert.T(t).That(w.idxSelBase.prefixLen == 1)
	// assert.T(t).That(w.skipScan)
	// assert.T(t).That(w.skipPrefixLen == 1)
	assert.This(queryAll2(q)).Is("a=3 b=3 | a=4 b=3 | a=5 b=3")
}

func TestWhere_Select_recalcIdxSel(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create table (a,b,c) key(a,b,c)")
	for a := range 3 {
		for b := range 5 {
			db.act(fmt.Sprintf("insert { a: %d, b: %d, c: 9 } into table", a+1, b+1))
		}
	}

	tran := db.NewReadTran()
	q := ParseQuery("table where b > 2", tran, nil)
	q = SetupIdx(q, CursorMode, tran, []string{"a", "b", "c"})
	w := q.(*Where)
	assert.T(t).That(w.idxSelBase != nil)
	assert.T(t).This(w.idxSelBase.skipLen).Is(1)

	q.Select(Sels{{"a", Pack(IntVal(2))}})
	assert.T(t).This(w.idxSelActive.prefixLen).Is(2)
	assert.T(t).This(w.idxSelActive.skipLen).Is(0)
	assert.This(queryAll2(q)).Is("a=2 b=3 c=9 | a=2 b=4 c=9 | a=2 b=5 c=9")
}

// TestWhere_Select_conflict tests that Select with a value conflicting with
// the where range constraint sets a no-scan conflict marker (not a full scan).
// where a > 1 means a is NOT fixed, so selectFixed doesn't catch a=0;
// the conflict must be detected in mergedPerCol / recalcIdxSel.
func TestWhere_Select_conflict(t *testing.T) {
	setup := func(where string) *Where {
		t.Helper()
		w := ParseQuery("comp where "+where, testTran{}, nil).(*Where)
		w.optInit()
		w.idxSelBase = &w.idxSels[0]
		w.fixed = nil
		w.singleton = false
		return w
	}

	// full recalc path: first select, a=0 conflicts with where a>1
	w := setup("a > 1")
	w.Select(Sels{{"a", Pack(SuInt(0))}})
	assert.T(t).Msg("recalc conflict").That(w.selConflict)

	// fast-path: first select a=2 (non-conflict), then a=0 conflicts
	w = setup("a > 1")
	w.Select(Sels{{"a", Pack(SuInt(2))}})
	assert.T(t).Msg("non-conflict selOrg").That(!w.selConflict)
	w.Select(Sels{{"a", Pack(SuInt(0))}}) // conflict
	assert.T(t).Msg("reRange conflict").This(w.selConflict)
}

func TestWhere_skipScan_gap(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create table (a,b,c) key(a,b,c)")
	for b := range 3 {
		db.act(fmt.Sprintf("insert { a: 1, b: %d, c: 5 } into table", b+1))
		db.act(fmt.Sprintf("insert { a: 1, b: %d, c: 6 } into table", b+1))
	}
	db.act("insert { a: 2, b: 1, c: 5 } into table")

	tran := db.NewReadTran()
	q := ParseQuery("table where a = 1 and c = 5", tran, nil)
	q, _, _ = Setup(q, CursorMode, tran)
	w, ok := q.(*Where)
	assert.T(t).That(ok)
	assert.T(t).That(w.idxSelBase != nil)
	assert.T(t).That(w.idxSelBase.prefixLen == 1)
	// assert.T(t).That(w.skipScan)
	// assert.T(t).That(w.skipPrefixLen == 2)
	assert.This(queryAll2(q)).Is("a=1 b=1 c=5 | a=1 b=2 c=5 | a=1 b=3 c=5")
}

func TestWhere_skipScan_emptyString(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create sktest3 (name, path) key(name, path)")
	const n = 50
	for i := range n {
		name := fmt.Sprintf("n%02d", i)
		db.act(fmt.Sprintf("insert { name: '%s', path: '' } into sktest3", name))
		for j := range 20 {
			db.act(fmt.Sprintf("insert { name: '%s', path: '/x%02d' } into sktest3", name, j))
		}
	}

	tran := db.NewReadTran()
	q := ParseQuery("sktest3 where path = ''", tran, nil)
	q, _, _ = Setup(q, CursorMode, tran)
	// w, ok := q.(*Where)
	// assert.T(t).That(ok)
	// assert.T(t).That(w.skipScan)
	// assert.T(t).That(w.skipPrefixLen == 1)
	th := &Thread{}
	count := 0
	for row := q.Get(th, Next); row != nil; row = q.Get(th, Next) {
		count++
	}
	assert.T(t).This(count).Is(n)
}

// TestWhere_skipScan_rangeQuery simulates a query like start_date <= X and end_date >= X
// on a (start, end) index with an extra key field appended for uniqueness.
// Uses SuDate values to match TestQueryBug2 row types.
func TestWhere_skipScan_rangeQuery(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create events (num, start, end) key(num) index(start, end)")
	// (start, end, num) events sorted by start, end, num:
	db.act("insert { num: 3, start: #20260318, end: #20260320 } into events")
	db.act("insert { num: 4, start: #20260320, end: #20260320 } into events")
	db.act("insert { num: 2, start: #20260320, end: #20260327 } into events")
	db.act("insert { num: 1, start: #20260327, end: #20260327 } into events")

	tran := db.NewReadTran()
	// Query: start <= #20260320 and end >= #20260320 => should match num=3,4,2
	q := ParseQuery("events where start <= #20260320 and end >= #20260320", tran, nil)
	q, _, _ = Setup(q, CursorMode, tran)
	th := &Thread{}
	count := 0
	for row := q.Get(th, Next); row != nil; row = q.Get(th, Next) {
		count++
	}
	assert.T(t).This(count).Is(3)
}

func TestWhere_bug(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create table (a,b,c,d) key(a,b,c) index(b)")
	db.act("insert { a: 1, b: 2, c: 3, d: 4 } into table")
	db.act("insert { a: 4, b: 5, c: 6, d: 7 } into table")
	db.act("insert { a: 7, b: 5, c: 8, d: 9 } into table")
	db.act("insert { a: 9, b: 0, c: 3, d: 4 } into table")
	tran := db.NewReadTran()
	q := ParseQuery("table where a=1 and b=2 and c=3", tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)
	assert.This(Strategy2(q)).Like(`
		table^(a,b,c)
		where*1 a is 1 and b is 2 and c is 3`)
}

func TestWhere_keyfixed(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create table (a,b,c,d) key(a)")
	db.act("insert { a: 1, b: 2, c: 3, d: 4 } into table")
	db.act("insert { a: 4, b: 5, c: 6, d: 7 } into table")
	tran := db.NewReadTran()
	q := ParseQuery("table where a=4", tran, nil)
	q = q.Transform()
	index := []string{"b"}
	Optimize(q, ReadMode, index, 1)
	q = SetApproach(q, index, 1, tran)
	assert.That(q.fastSingle())
	th := &Thread{}
	sels := Sels{{"b", Pack(IntVal(5))}}
	row := q.Lookup(th, sels)
	hdr := q.Header()
	assert.This(row2str(hdr, row)).Is("a=4 b=5 c=6 d=7")

	q.Select(sels)
	row = q.Get(th, Next)
	assert.This(row2str(hdr, row)).Is("a=4 b=5 c=6 d=7")
	row = q.Get(th, Next)
	assert.This(row).Is(nil)
}

func TestSplit(t *testing.T) {
	// Test empty
	sels := Sels{}
	index := []string{"a"}
	isels, osels := Split(false, sels, index)
	assert.T(t).This(isels).Is(nil)
	assert.T(t).This(osels).Is(nil)

	// Test all in index
	sels = Sels{{"a", "1"}, {"b", "2"}}
	index = []string{"a", "b"}
	isels, osels = Split(false, sels, index)
	assert.T(t).This(isels).Is(Sels{{"a", "1"}, {"b", "2"}})
	assert.T(t).This(osels).Is(nil)

	// Test none in index
	sels = Sels{{"c", "3"}, {"d", "4"}}
	index = []string{"a", "b"}
	isels, osels = Split(false, sels, index)
	assert.T(t).This(isels).Is(nil)
	assert.T(t).This(osels).Is(Sels{{"c", "3"}, {"d", "4"}})

	// Test mixed
	sels = Sels{{"a", "1"}, {"c", "3"}, {"b", "2"}, {"d", "4"}}
	index = []string{"a", "b"}
	isels, osels = Split(false, sels, index)
	// iflds should contain "a" and "b", in some order, ivals accordingly
	// oflds "c" and "d"
	assert.T(t).This(len(isels)).Is(2)
	assert.T(t).This(len(osels)).Is(2)
	// Check that iflds are in index
	for _, sel := range isels {
		if !slices.Contains(index, sel.col) {
			t.Errorf("isels contains %s not in index", sel.col)
		}
	}
	for _, sel := range osels {
		if slices.Contains(index, sel.col) {
			t.Errorf("osels contains %s which is in index", sel.col)
		}
	}
	// Check vals match flds order
	expected := map[string]string{"a": "1", "b": "2", "c": "3", "d": "4"}
	for _, sel := range isels {
		assert.T(t).This(sel.val).Is(expected[sel.col])
	}
	for _, sel := range osels {
		assert.T(t).This(sel.val).Is(expected[sel.col])
	}
}
