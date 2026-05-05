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
	test("a >= ''", "[a:[<<max>]]") // everything, always matches
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
	idx := []string{"a", "b", "c"}
	test := func(query string, expected string) {
		t.Helper()
		w := ParseQuery(query, testTran{}, nil).(*Where)
		pf := perField(w.expr.Exprs, w.source.Header().Physical())
		idxSpans := indexSpans(idx, pf)
		assert.T(t).This(fmt.Sprint(idxSpans)).Is("[" + expected + "]")
	}
	test("comp where a is 1", "[1]")
	test("comp where a is 1 and c is 2", "[1]")
	test("comp where a is 1 and b is 2", "[1] [2]")
	test("comp where a is 1 and b is 2 and c is 3", "[1] [2] [3]")
	test("comp where a >= 4", "[>=4]")
	test("comp where a >= 4 and b is 2", "[>=4]")
	test("comp where a is 2 and b >= 4", "[2] [>=4]")
	test("comp where a in (1,2) and b in (3,4)", "[1 2] [3 4]")
	test("comp where a is '' and b isnt 0", "[''] [<0 >0]")
	idx = []string{"id"}
	test("customer where id is 'e'", "['e']")
}

func TestWhere_explodeIndexSpans(t *testing.T) {
	idx := []string{"a", "b", "c"}
	test := func(query string, expected string) {
		t.Helper()
		w := ParseQuery("comp where "+query, testTran{}, nil).(*Where)
		pf := perField(w.expr.Exprs, w.source.Header().Physical())
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
		w.optInit()
		pf := perField(w.expr.Exprs, w.source.Header().Physical())
		idxSels := w.perIndex(pf)
		assert.T(t).This(fmt.Sprint(idxSels)).Is("[" + expected + "]")
	}

	table = "comp" // key(a,b,c)
	test("a is 1", "a,b,c: 1..1,<max> = 0.1")
	test("a is 1 and b is 2", "a,b,c: 1,2..1,2,<max> = 0.01")
	test("a is 1 and b is 2 and c is 3", "a,b,c: 1,2,3 = 0.0005")
	test("a > 4", "a,b,c: 4,<max>..<max> = 0.5")
	test("a <= 4", "a,b,c: ..4,<max> = 0.5")
	test("a is 2 and b >= 4", "a,b,c: 2,4..2,<max> = 0.06")
	test("a in (1,2) and b in (3,4)", "a,b,c: 1,3..1,3,<max> | 1,4..1,4,<max> | "+
		"2,3..2,3,<max> | 2,4..2,4,<max> = 0.04")
	test("a in (1,2) and b > 4",
		"a,b,c: 1,4,<max>..1,<max> | 2,4,<max>..2,<max> = 0.1")
	test("a is 1 or a > 3", "a,b,c: 1..1,<max> | 3,<max>..<max> = 0.7")
	test("a isnt 5", "a,b,c: ..5 | 5,<max>..<max> = 0.9")
	test("a is '' and b isnt 0", "a,b,c: ..'',0 | '',0,<max>..'',<max> = 0.09")

	table = "table"
	test("a >= ''", "a: ''..<max> = 1")
	// test("a >= ''", "a: '' | '\\x00'..<max> = 1.005")

	table = "comp2" // key(a,b,c) index(b)
	// singleton, skip b index
	test("a is 1 and b is 2 and c is 3", "a,b,c: 1,2,3")
}

func TestWhere_perIndex_fracs(t *testing.T) {
	test := func(table, expr string, ifFrac float64, dataFilter bool) {
		t.Helper()
		w := ParseQuery(table+" where "+expr, testTran{}, nil).(*Where)
		w.optInit()
		pf := perField(w.expr.Exprs, w.source.Header().Physical())
		idxSels := w.perIndex(pf)
		assert.T(t).That(len(idxSels) == 1)
		is := idxSels[0]
		assert.T(t).Msg(expr + " ifFrac").This(is.ifFrac).Is(ifFrac)
		assert.T(t).Msg(expr + " dataFilter").This(is.dataFilter).Is(dataFilter)
	}
	// no extra expressions: ifFrac=1, no dataFilter
	test("comp", "a is 1", 1.0, false)
	// expression on index field beyond range columns: ifFrac=.5
	// a > 1 becomes range (nfields=1), c is in IndexCols but not in rangeCols
	test("comp", "a > 1 and c is 2", 0.5, false)
	// expression on non-index field: dataFilter=true
	// hist index(item) has IndexCols=[item,date,id]; cost is not in IndexCols
	test("table", "a is 1 and b is 2", 1.0, true)
	
	// F(a) is not covered by a > 1
	test("comp", "a > 1 and F(a)", 0.5, false)
	// zero-column expression: dataFilter=true
	test("comp", "a is 1 and Foo()", 1.0, true)
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
			for _, pr := range idxsel.ptrngs {
				if pr.isPoint() {
					ixrange = ixrange || packed == pr.org
					// fmt.Printf("%q == %q\n", packed, pr.org)
				} else { // range
					ixrange = ixrange || pr.org <= packed && packed < pr.end
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
	test := func(query string, nrows, pop int) {
		t.Helper()
		var tran testTran
		w := ParseQuery(query, tran, nil)
		Setup(w, ReadMode, tran)
		n, p := w.Nrows()
		assert.T(t).This(n).Is(nrows)
		assert.T(t).This(p).Is(pop)
	}
	test("table where F()", 50, 100)
	test("inven where item >= 5", 50, 100)
	test("inven where item < 3 and item > 3", 0, 100) // conflict
	test("inven where item is 1", 1, 100)
	test("inven where item in (1,2,3,4)", 2, 100)
	test("inven where item > 2 and item < 4", 20, 100)
	// dataFilter (non-index column)
	test("inven where item > 2 and item < 4 and qty", 10, 100)
	// ifFrac (index column beyond range)
	test("comp where a > 1 and c is 2", 400, 1000)
	test("comp where a > 1 and F(a)", 400, 1000)
	// zero-column expression => dataFilter
	test("table where Foo()", 50, 100)
	test("comp where a is 1 and Foo()", 50, 1000)
	// combined irFrac + ifFrac + dataFilter
	test("hist where date is 3", 10, 100)
	// not on table
	test("inven extend x where x > 5", 50, 100)
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
	q, _, _ = Setup(q, CursorMode, tran)
	cols := []string{"a", "b"}
	vals := []string{Pack(IntVal(4)), Pack(IntVal(5)), Pack(IntVal(6))}
	q.Select(cols, vals)
	assert.This(queryAll2(q)).Is("a=4 b=5 c=6")
	q.Select(nil, nil)
	assert.This(queryAll2(q)).Is("a=4 b=5 c=6 | a=7 b=5 c=8")

	q = ParseQuery(query, tran, nil)
	q, _, _ = Setup(q, CursorMode, tran)
	vals = []string{Pack(IntVal(1)), Pack(IntVal(2))} // conflict
	q.Select(cols, vals)
	assert.This(queryAll2(q)).Is("")
	q.Select(nil, nil)
	assert.This(queryAll2(q)).Is("a=4 b=5 c=6 | a=7 b=5 c=8")
}

func TestWhere_ptrange(t *testing.T) {
	table := "comp"
	test := func(query string, selCols, selVals []string) {
		t.Helper()
		w := ParseQuery(table+" where "+query, testTran{}, nil).(*Where)
		w.optInit()
		pf := perField(w.expr.Exprs, w.source.Header().Physical())
		w.idxSel = &w.perIndex(pf)[0]
		// fmt.Printf("idxSel ptrange\n\t%q\n\t%q\n",
		// 	w.idxSel.ptrngs[0].org, w.idxSel.ptrngs[0].end)
		w.fixed = nil
		w.singleton = false
		w.Select(selCols, selVals)
		// fmt.Printf("selOrg, selEnd\n\t%q\n\t%q\n", w.selOrg, w.selEnd)
		pr := w.idxSel.ptrngs[0].intersect(w.selOrg, w.selEnd)
		// fmt.Printf("intersect\n\t%q\n\t%q\n", pr.org, pr.end)
		assert.This(pr).Is(w.idxSel.ptrngs[0])
	}
	test("a is '1' and b is '2' and c is '3'",
		[]string{"a"}, []string{Pack(SuStr("1"))})
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
	test("(a,b,c) key(a)", "a = 1", "[a:[1]]", "[a: 1]")
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

func TestSplit(t *testing.T) {
	// Test empty
	flds := []string{}
	vals := []string{}
	index := []string{"a"}
	iflds, ivals, oflds, ovals := Split(false, flds, vals, index)
	assert.T(t).This(iflds).Is(nil)
	assert.T(t).This(ivals).Is(nil)
	assert.T(t).This(oflds).Is(nil)
	assert.T(t).This(ovals).Is(nil)

	// Test all in index
	flds = []string{"a", "b"}
	vals = []string{"1", "2"}
	index = []string{"a", "b"}
	iflds, ivals, oflds, ovals = Split(false, flds, vals, index)
	assert.T(t).This(iflds).Is([]string{"a", "b"})
	assert.T(t).This(ivals).Is([]string{"1", "2"})
	assert.T(t).This(oflds).Is(nil)
	assert.T(t).This(ovals).Is(nil)

	// Test none in index
	flds = []string{"c", "d"}
	vals = []string{"3", "4"}
	index = []string{"a", "b"}
	iflds, ivals, oflds, ovals = Split(false, flds, vals, index)
	assert.T(t).This(iflds).Is(nil)
	assert.T(t).This(ivals).Is(nil)
	assert.T(t).This(oflds).Is([]string{"c", "d"})
	assert.T(t).This(ovals).Is([]string{"3", "4"})

	// Test mixed
	flds = []string{"a", "c", "b", "d"}
	vals = []string{"1", "3", "2", "4"}
	index = []string{"a", "b"}
	iflds, ivals, oflds, ovals = Split(false, flds, vals, index)
	// iflds should contain "a" and "b", in some order, ivals accordingly
	// oflds "c" and "d"
	assert.T(t).This(len(iflds)).Is(2)
	assert.T(t).This(len(ivals)).Is(2)
	assert.T(t).This(len(oflds)).Is(2)
	assert.T(t).This(len(ovals)).Is(2)
	// Check that iflds are in index
	for _, f := range iflds {
		if !slices.Contains(index, f) {
			t.Errorf("iflds contains %s not in index", f)
		}
	}
	for _, f := range oflds {
		if slices.Contains(index, f) {
			t.Errorf("oflds contains %s which is in index", f)
		}
	}
	// Check vals match flds order
	expected := map[string]string{"a": "1", "b": "2", "c": "3", "d": "4"}
	for i, f := range iflds {
		assert.T(t).This(ivals[i]).Is(expected[f])
	}
	for i, f := range oflds {
		assert.T(t).This(ovals[i]).Is(expected[f])
	}
}

func TestWhere_SelOrgNotFull(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create table (a,b,c,d) index(a,c) key(a,b)") // index must be first
	db.act("insert { a: 1, b: 2, c: 3, d: 4 } into table")
	db.act("insert { a: 4, b: 5, c: 6, d: 7 } into table")
	db.act("insert { a: 7, b: 5, c: 8, d: 9 } into table")
	tran := db.NewReadTran()
	q := ParseQuery("table where b=5", tran, nil)
	// b is fixed, so index (a,c) can provide order (a,b)
	// (a,c) and (a,b) have the same cost so Where picks the first (a,c)
	// but (a,c) doesn't support lookups on (a,b) even with fixed
	q = q.Transform()
	key := []string{"a", "b"}
	Optimize(q, ReadMode, key, 1)
	q = SetApproach(q, key, 1, tran)
	q.Lookup(nil, key, []string{Pack(SuInt(4)), Pack(SuInt(5))})
}

func TestWhereCost(t *testing.T) {
	assert := assert.T(t).This
	type Eg struct {
		// irFrac is the selectivity of the index range (or 1 for all)
		irFrac float64
		// ifFrac is the selectivity of the index filter (or 1 for none)
		ifFrac float64
		// dataFilter is whether there is additional filtering of the data
		dataFilter bool
		// inFrac is the amount the caller expects to read
		inFrac float64
	}
	test := func(eg *Eg, expected float64) {
		// defaults
		if eg.irFrac == 0 {
			eg.irFrac = 1
		}
		if eg.ifFrac == 0 {
			eg.ifFrac = 1
		}
		if eg.inFrac == 0 {
			eg.inFrac = 1
		}
		cost :=
			WhereCost(100_000, eg.inFrac, eg.irFrac, eg.ifFrac, eg.dataFilter)
		assert(cost).Is(expected)
	}
	// baseline: all defaults => 100 * srcRows
	test(&Eg{}, 100_000)
	// irFrac only
	test(&Eg{irFrac: 1.0 / 1000}, 100)
	// dataFilter triggers pessimistic guard on inFrac
	test(&Eg{dataFilter: true, inFrac: .01}, 25_750)
	// ifFrac < 1 triggers pessimistic guard
	test(&Eg{ifFrac: 0.5}, 60_000)
	// ifFrac < 1 + inFrac with pessimistic guard
	test(&Eg{ifFrac: 0.1, inFrac: 0.5}, 17_500)
	// both irFrac and ifFrac
	test(&Eg{irFrac: 0.1, ifFrac: 0.5}, 6_000)
	// dataFilter + ifFrac
	test(&Eg{dataFilter: true, ifFrac: 0.5}, 60_000)
	// all parameters non-default
	test(&Eg{irFrac: 0.5, ifFrac: 0.5, dataFilter: true, inFrac: 0.5}, 18_750)
	// inFrac without pessimistic guard (no dataFilter, ifFrac >= 1)
	test(&Eg{inFrac: 0.5}, 50_000)
	// irFrac + inFrac: no guard (pure index range, ifFrac=1, no dataFilter)
	test(&Eg{irFrac: 0.1, inFrac: 0.5}, 5_000)
	// dataFilter=true with inFrac=1: guard is no-op, same cost as baseline
	test(&Eg{dataFilter: true}, 100_000)
	// irFrac + dataFilter + small inFrac: guard applied to narrowed range
	test(&Eg{irFrac: 0.5, dataFilter: true, inFrac: 0.01}, 12_875)
}
