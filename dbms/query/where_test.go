// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stats"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestWhere_perField(t *testing.T) {
	test := func(query string, expectedCols string, expectedUnspanable string, expectedDataExprCount int) {
		t.Helper()
		w := ParseQuery("table where "+query, testTran{}, nil).(*Where)
		actualCols := "conflict"
		actualUnspanable := "conflict"
		actualDataExprCount := 0
		if !w.conflict {
			actualCols = fmt.Sprint(w.colSpans)[3:]
			actualUnspanable = fmt.Sprint(w.unspanable)
			actualDataExprCount = w.dataExprCount
		}
		assert.T(t).Msg(query).This(actualCols).Is(expectedCols)
		assert.T(t).Msg(query).This(actualUnspanable).Is(expectedUnspanable)
		assert.T(t).Msg(query).This(actualDataExprCount).Is(expectedDataExprCount)
	}
	// nothing indexable, no columns
	test("Foo()", "[]", "[]", 1)
	// not indexable, single column
	test("a =~ 'x'", "[]", "[a]", 0)
	test("a", "[]", "[a]", 0)
	// not indexable, multiple columns (treated like no-column)
	test("F(a,b)", "[]", "[]", 1)
	// binary
	test("a is 123", "[a:[123]]", "[]", 0)
	test("a isnt 123", "[a:[<123 >123]]", "[]", 0)
	test("a < 123", "[a:[<123]]", "[]", 0)
	test("a <= 123", "[a:[<=123]]", "[]", 0)
	test("a > 123", "[a:[>123]]", "[]", 0)
	test("a >= 123", "[a:[>=123]]", "[]", 0)
	// compare to ""
	test("a isnt ''", "[a:[>'']]", "[]", 0)
	test("a < ''", "conflict", "conflict", 0)
	test("a <= ''", "[a:['']]", "[]", 0)
	test("a > ''", "[a:[>'']]", "[]", 0)
	test("a >= ''", "[a:[<max]]", "[]", 0) // everything, always matches
	// in
	test("a in (3,1,2)", "[a:[1 2 3]]", "[]", 0)
	// range
	test("a > 3 and a < 6", "[a:[>3_<6]]", "[]", 0)
	test("a >= 3 and a <= 6", "[a:[>=3_<=6]]", "[]", 0)
	// type
	test("String?(a)", "[a:['' >=PackString_<PackDate]]", "[]", 0)
	test("Number?(a)", "[a:[>=PackMinus_<PackString]]", "[]", 0)
	test("Date?(a)", "[a:[>=PackDate_<PackDate+1]]", "[]", 0)
	// or
	test("a is 0 or b is 0", "[]", "[]", 1)
	test("a is 0 or a =~ 'x'", "[]", "[a]", 0)
	test("a is 0 or a > 6", "[a:[0 >6]]", "[]", 0)
	test("a is false and (b is 1 or b is 2)", "[a:[false] b:[1 2]]", "[]", 0)
	test("a is 1 or a is 2 or a is 3", "[a:[1 2 3]]", "[]", 0)
	test("a is 1 or a is 2 or a is 1", "[a:[1 2]]", "[]", 0)
	test("a is 1 or a in (1,2,3) or a is 2", "[a:[1 2 3]]", "[]", 0)
	test("a > 3 or a > 6", "[a:[>3]]", "[]", 0)
	test("a > 3 or a is 6", "[a:[>3]]", "[]", 0)
	test("a > 6 or a > 3", "[a:[>3]]", "[]", 0)
	test("a is 6 or a > 3", "[a:[>3]]", "[]", 0)
	test("Number?(a) or String?(a)", "[a:['' >=PackMinus_<PackDate]]", "[]", 0)
	test("String?(a) or Number?(a)", "[a:['' >=PackMinus_<PackDate]]", "[]", 0)
	// multiple
	test("a is 123 and b is 456", "[a:[123] b:[456]]", "[]", 0)
	test("a isnt 'm' and a > 'a'", "[a:[>'a'_<'m' >'m']]", "[]", 0)
	test("a isnt '' and a > 'm'", "[a:[>'m']]", "[]", 0)
	test("a isnt '' and a in ('','m')", "[a:['m']]", "[]", 0)
	test("a isnt '' and String?(a)", "[a:[>=PackString_<PackDate]]", "[]", 0)
	// intersect
	test("a in (1,2,3) and a in (1,2,3)", "[a:[1 2 3]]", "[]", 0)
	test("a in (1,2,3) and a in (2,3,4)", "[a:[2 3]]", "[]", 0)
	test("a in (1,2,3) and a is 2", "[a:[2]]", "[]", 0)
	test("a in (1,2,3) and a >= 2", "[a:[2 3]]", "[]", 0)
	test("a < 5 and a > 2", "[a:[>2_<5]]", "[]", 0)
	test("a < 5 and a <= 5", "[a:[<5]]", "[]", 0)
	test("a >= 5 and a > 5", "[a:[>5]]", "[]", 0)
	// conflict
	test("a in (1,2,3) and a in (4,5,6)", "conflict", "conflict", 0)
	test("a in (1,3,5) and a in (2,4,6)", "conflict", "conflict", 0)
	test("a < 5 and a > 6", "conflict", "conflict", 0)
	// mixed: indexable + not indexable
	test("a is 123 and F(a)", "[a:[123]]", "[a]", 0)
	test("a is 123 and F(b)", "[a:[123]]", "[b]", 0)
	test("a is 123 and F(a,b)", "[a:[123]]", "[]", 1)
	// multiple non-indexable expressions
	test("F(a) and F(b)", "[]", "[a b]", 0)
	test("F(a,b) and F(c)", "[]", "[c]", 1) // multi-col treated as no-col
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
		pf, _, _ := perField(w.expr.Exprs, w.source.Header().Physical())
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
		pf, _, _ := perField(w.expr.Exprs, w.source.Header().Physical())
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
	d := func(args ...float64) float64 {
		frac := 1.0
		for i, a := range args {
			frac *= damp(i, a)
		}
		return frac
	}
	f := func(x float64) string {
		return fmt.Sprintf("%.2g", x)
	}
	test := func(query string, expected string, wfrac float64) {
		t.Helper()
		w := ParseQuery(query, testTran{}, nil).(*Where)
		w.optInit() // runs perIndex
		assert.T(t).This(fmt.Sprint(w.idxSels)).Is("[" + expected + "]")
		assert.T(t).Msg(w.wfrac).This(f(w.wfrac)).Is(f(wfrac))
	}

	// comp: key(a,b,c) nrows = 1000
	test("comp where a is 1",
		"(a,b,c) a: <1..1,max> = pr: .1", .1)
	test("comp where a is 1 and b is 2",
		"(a,b,c) a,b: <1,2..1,2,max> = pr: .01", .01)
	test("comp where a is 1 and b is 2 and c is 3",
		"(a,b,c) a,b,c: <1,2,3> = singleton", 0)
	test("comp where a is 1 and b >= 2",
		"(a,b,c) a,b: <1,2..1,max> = pr: .08", .08)

	test("comp where a > 4",
		"(a,b,c) a: <4,max..max> = pr: .5", .5)
	test("comp where a <= 4",
		"(a,b,c) a: <..4,max> = pr: .5", .5)
	test("comp where a is 2 and b >= 4",
		"(a,b,c) a,b: <2,4..2,max> = pr: .06", .06)
	test("comp where a in (1,2) and b in (3,4)",
		"(a,b,c) a,b: <1,3..1,3,max | 1,4..1,4,max | "+
			"2,3..2,3,max | 2,4..2,4,max> = pr: .04", .04)
	test("comp where a in (1,2) and b > 4",
		"(a,b,c) a,b: <1,4,max..1,max | 2,4,max..2,max> = pr: .1", .1)
	test("comp where a is 1 or a > 3",
		"(a,b,c) a: <1..1,max | 3,max..max> = pr: .7", .7)
	test("comp where a isnt 5",
		"(a,b,c) a: <..5 | 5,max..max> = pr: .9", .9)
	test("comp where a is '' and b isnt 0",
		"(a,b,c) a,b: <..'',0 | '',0,max..'',max> = pr: .09", .09)

	test("comp where b is 2",
		"(a,b,c) +b: <2..2,max> = ir: .5", .5)
	test("comp where F(b)",
		"(a,b,c) = if: .5", .5)
	test("comp where b in (2,3)",
		"(a,b,c) = if: .5", .5)
	test("comp where b is 1 and c in (2,3)",
		"(a,b,c) +b: <1..1,max> = ir: .5 if: .71", d(.5, .5))
	test("comp where b is 2 and c is 3", // skip scan on (b,c)
		"(a,b,c) +b,c: <2,3..2,3,max> = ir: .35", d(.5, .5))
	test("comp where a is 1 and c is 3", // prefix(a) + skip scan(b)
		"(a,b,c) a: <1..1,max> +c: <3..3,max> = pr: .1 ir: .071", d(.1, .5))

	test("comp where b >= 2 and b <= 4",
		"(a,b,c) +b: <2..4,max> = ir: .5", .5)
	test("comp where b > 2",
		"(a,b,c) +b: <2,max..max> = ir: .5", .5)
	test("comp where b < 5",
		"(a,b,c) +b: <..5> = ir: .5", .5)

	test("comp where a > 1 and F(c)", // prefix(a) + indexFilter(c)
		"(a,b,c) a: <1,max..max> = pr: .8 if: .71", d(.5, .8))
	test("comp where a > 1 and F(a)",
		"(a,b,c) a: <1,max..max> = pr: .8 if: .71", d(.5, .8))
	test("comp where a > 1 and F(a) and G(a)",
		"(a,b,c) a: <1,max..max> = pr: .8 if: .59", d(.5, .5, .8))
	test("comp where a is 1 and Foo()", "(a,b,c) a: <1..1,max> = pr: .1 df", d(.1, .5))

	// table: key(a) nrows = 100
	test("table where a >= ''",
		"(a) a: <''..max> = pr: 1", 1) //TODO skip useless
	test("table where F()",
		"", .5)
	test("table where F(c)",
		"", .5)

	// inven: key(item) nrows = 100
	test("inven where item >= 5",
		"(item) item: <5..max> = pr: .5", .5)
	test("inven where item < 3 and item > 3",
		"", 0) // conflict
	test("inven where item in (1,2,3,4)",
		"(item) item: <1 | 2 | 3 | 4> = pr: .02", .02)
	test("inven where item > 2 and item < 4",
		"(item) item: <2..4> = pr: .2", .2)
	test("inven where item > 2 and item < 4 and qty",
		"(item) item: <2..4> = pr: .2 df", .2*math.Sqrt(.5))
	test("inven where qty > 5",
		"", .5)

	// comp: additional wfrac tests
	test("comp where a > 1",
		"(a,b,c) a: <1,max..max> = pr: .8", .8)
	test("comp where a > 1 and F(a)",
		"(a,b,c) a: <1,max..max> = pr: .8 if: .71", .5*math.Sqrt(.8))
	test("comp where a > 1 and F(c)",
		"(a,b,c) a: <1,max..max> = pr: .8 if: .71", .5*math.Sqrt(.8))
	test("comp where a is 1",
		"(a,b,c) a: <1..1,max> = pr: .1", .1)

	// hist: key(date,item,id) index(item) nrows = 100
	test("hist where date is 3",
		"(date,item,id) date: <3..3,max> = pr: .1 (item,date,id) +date: <3..3,max> = ir: .5", .1)

	// comp2: index(b) key(a,b,c) nrows = 0
	// (b) will be (b,a,c) after adding key columns to make it unique
	// so it should get a skip scan the same as (a,b,c)
	test("comp2 where a is 1 and b is 2 and c is 3",
		"(a,b,c) a,b,c: <1,2,3> = singleton", 0)
	test("comp2 where c is 1",
		"(a,b,c) +c: <1..1,max> = ir: .5 (b,a,c) +c: <1..1,max> = ir: .5", .5)
	test("comp2 where a is 1 and b is 2",
		"(a,b,c) a,b: <1,2..1,2,max> = pr: .01 (b,a,c) b,a: <2,1..2,1,max> = pr: .01", .01)

	test("ixftab where c is 8",
		"(b,c,a) +c: <8..8,max> = ir: .5", .5)

	// prefix bug
	test("comp where a in ('', '1') and b is ''",
		"(a,b,c) a,b: <..'','',max | '1'..'1','',max> = pr: .11", .11)

	// skip scan bug
	test("comp where b = '' and c = ''",
		"(a,b,c) +b,c: <..'','',max> = ir: .35", d(.5, .5))
}

func TestWhere_stats_wfrac(t *testing.T) {
	pack := func(s string) string { return Pack(SuStr(s)) }
	stats := stats.Stats{
		"customer": {
			Count: 100,
			Columns: map[string]stats.ColStats{
				"name": {
					Cardinality: 5,
					Quantiles: []string{
						pack("a"), pack("c"), pack("m"), pack("s")},
					Tops:     []stats.Top{{Value: pack("joe"), Frac: .25}},
					TailFrac: .75,
				},
				"city": {
					Cardinality: 4,
					Quantiles: []string{
						pack("a"), pack("c"), pack("m"), pack("s")},
					Tops:     []stats.Top{{Value: pack("paris"), Frac: .3}},
					TailFrac: .7,
				},
			},
		},
	}
	tran := testTran{stats: stats}
	test := func(query string, expected float64) {
		t.Helper()
		w := ParseQuery(query, tran, nil).(*Where)
		w.optInit()
		assert.T(t).Msg(query).This(w.wfrac).Is(expected)
	}
	// no usable index prefix => use stats
	test("customer where name is 'joe'", .25)
	test("customer where name is 'bob'", .75/4) // tail estimate
	test("customer where name in ('joe','bob')", .25+.75/4)
	test("customer where name >= 'm'", 1.-2./3.)
	test("customer where name < 'c'", 1./3.)
	test("customer where name <= 'c'", 2./3.)
	// city has stats
	test("customer where city is 'paris'", .3)
	test("customer where city is 'london'", .7/3) // tail estimate
	test("customer where city >= 'm'", 1.-2./3.)
	// two columns with stats: fracs sorted, smallest * sqrt(second)
	test("customer where name is 'joe' and city is 'paris'",
		.25*math.Sqrt(.3))
	test("customer where name is 'joe' and city is 'london'",
		.7/3*math.Sqrt(.25))
	test("customer where name is 'bob' and city is 'paris'",
		.75/4*math.Sqrt(.3))
	test("customer where name is 'bob' and city is 'london'",
		.75/4*math.Sqrt(.7/3))
	test("customer where name is 'joe' and city >= 'm'",
		.25*math.Sqrt(1./3.))
	// unspanable
	test("customer where F(name)", .5)
}

type wtestTran struct {
	testTran
}

func (t wtestTran) RangeFrac(table string, iIndex int, org, end string) float64 {
	return .5
}

// Verify index selection ranges match actual expression evaluation
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
	q = setupIndex(q, CursorMode, tran, cols)
	sels := Sels{{"a", Pack(IntVal(4))}, {"b", Pack(IntVal(5))}}
	q.Select(sels)
	assert.This(queryAll2(q)).Is("a=4 b=5 c=6")
	q.Select(nil)
	assert.This(queryAll2(q)).Is("a=4 b=5 c=6 | a=7 b=5 c=8")

	q = ParseQuery(query, tran, nil)
	q = setupIndex(q, CursorMode, tran, cols)
	sels = Sels{{"a", Pack(IntVal(1))}, {"b", Pack(IntVal(2))}} // conflict
	q.Select(sels)
	assert.This(queryAll2(q)).Is("")
	q.Select(nil)
	assert.This(queryAll2(q)).Is("a=4 b=5 c=6 | a=7 b=5 c=8")

	// select with col not in index fields (c not in key(a,b))
	q = ParseQuery(query, tran, nil)
	q = setupIndex(q, CursorMode, tran, cols)
	q.Select(Sels{{"a", Pack(IntVal(4))}, {"c", Pack(IntVal(6))}})
	assert.This(queryAll2(q)).Is("a=4 b=5 c=6")
}

func TestWhere_fixed(t *testing.T) {
	test := func(query, expected string) {
		t.Helper()
		w := ParseQuery("table where "+query, testTran{}, nil).(*Where)
		assert.T(t).This(w.Fixed().String()).Is(expected)
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

func TestWhere_skipScan(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create table (a,b,c) key(a,b,c)")
	for a := range 3 {
		for b := range 5 {
			for c := range 2 {
				db.act(fmt.Sprintf("insert { a: %d, b: %d, c: %d } into table",
					a+1, b+1, c+1))
			}
		}
	}
	test := func(query, idxSel, result string) {
		t.Helper()
		tran := db.NewReadTran()
		q := ParseQuery(query, tran, nil)
		q, _, _ = Setup(q, CursorMode, tran)
		w := q.(*Where)
		assert.T(t).This(w.idxSelBase.String()).Is(idxSel)
		assert.This(queryAll2(q)).Is(result)
	}

	// pure prefix, no skip scan
	test("table where a is 2",
		"(a,b,c) a: <2..2,max> = pr: .33",
		"a=2 b=1 c=1 | a=2 b=1 c=2 | a=2 b=2 c=1 | a=2 b=2 c=2 | "+
			"a=2 b=3 c=1 | a=2 b=3 c=2 | a=2 b=4 c=1 | a=2 b=4 c=2 | "+
			"a=2 b=5 c=1 | a=2 b=5 c=2")
	test("table where a = 2 and b >= 4",
		"(a,b,c) a,b: <2,4..2,max> = pr: .13",
		"a=2 b=4 c=1 | a=2 b=4 c=2 | a=2 b=5 c=1 | a=2 b=5 c=2")

	// pure skip scan, no prefix
	test("table where b = 5",
		"(a,b,c) +b: <5..5,max> = ir: .5",
		"a=1 b=5 c=1 | a=1 b=5 c=2 | a=2 b=5 c=1 | a=2 b=5 c=2 | "+
			"a=3 b=5 c=1 | a=3 b=5 c=2")
	test("table where b = 5 and c = 2",
		"(a,b,c) +b,c: <5,2..5,2,max> = ir: .35",
		"a=1 b=5 c=2 | a=2 b=5 c=2 | a=3 b=5 c=2")

	// prefix + skip on adjacent column
	test("table where a > 1 and b = 3",
		"(a,b,c) a: <1,max..max> +b: <3..3,max> = pr: .67 ir: .41",
		"a=2 b=3 c=1 | a=2 b=3 c=2 | a=3 b=3 c=1 | a=3 b=3 c=2")

	// prefix + range skip on non-adjacent column (gap)
	test("table where a > 1 and c > 1",
		"(a,b,c) a: <1,max..max> +c: <1,max..max> = pr: .67 ir: .41",
		"a=2 b=1 c=2 | a=2 b=2 c=2 | a=2 b=3 c=2 | a=2 b=4 c=2 | "+
			"a=2 b=5 c=2 | a=3 b=1 c=2 | a=3 b=2 c=2 | a=3 b=3 c=2 | "+
			"a=3 b=4 c=2 | a=3 b=5 c=2")

	// empty string skip scan (verifies encoding edge case)
	db.adm("create strtab (name, path) key(name, path)")
	for i := range 5 {
		name := fmt.Sprintf("n%02d", i)
		db.act(fmt.Sprintf("insert { name: '%s', path: '' } into strtab", name))
		for j := range 3 {
			db.act(fmt.Sprintf("insert { name: '%s', path: '/x%02d' } into strtab", name, j))
		}
	}
	test("strtab where path = ''",
		"(name,path) +path: <..'',max> = ir: .5",
		"name=n00 | name=n01 | name=n02 | name=n03 | name=n04")
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
	q = setupIndex(q, CursorMode, tran, []string{"a", "b", "c"})
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
// the conflict must be detected in mergeColSpans / recalcIdxSel.
func TestWhere_Select_conflict(t *testing.T) {
	setup := func(where string) *Where {
		t.Helper()
		w := ParseQuery("comp where "+where, testTran{}, nil).(*Where)
		w.optInit()
		w.idxSelBase = w.idxSels[0]
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
	assert.T(t).Msg("non-conflict").That(!w.selConflict)
	w.Select(Sels{{"a", Pack(SuInt(0))}}) // conflict
	assert.T(t).Msg("conflict").This(w.selConflict)
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
	index := []string{"b"}
	q, _, _ = SetupReq(q, ReadMode, tran, OrderReq(index, 1))
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

func TestWhere_SelOrgNotFull(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create table (a,b,c,d) index(a,c) key(a,b)") // index must be first
	// note: index(a,c) will be (a,c,b) with key(a,b) added for uniqueness
	db.act("insert { a: 1, b: 2, c: 3, d: 4 } into table")
	db.act("insert { a: 4, b: 5, c: 6, d: 7 } into table")
	db.act("insert { a: 7, b: 5, c: 8, d: 9 } into table")
	tran := db.NewReadTran()
	q := ParseQuery("table where b=5", tran, nil)
	// b is fixed, so index (a,c) can provide order (a,b)
	// (a,c) and (a,b) have the same cost so Where picks the first (a,c)
	// but (a,c) doesn't support lookups on (a,b) even with fixed
	key := []string{"a", "b"}
	q, _, _ = SetupReq(q, ReadMode, tran, UniqueReq(key, 1))
	q.Lookup(nil, Sels{{"a", Pack(SuInt(4))}, {"b", Pack(SuInt(5))}})
}

func TestWhere_Cost(t *testing.T) {
	assert := assert.T(t).This
	type Eg struct {
		// irFrac is the selectivity of the index range (or 1 for all)
		irFrac float64
		// ifFrac is the selectivity of the index filter (or 1 for none)
		ifFrac float64
		// df is whether there is additional filtering of the data
		df bool
		// inFrac is the amount the caller expects to read
		inFrac float64
	}
	test := func(eg *Eg, expected Cost) {
		t.Helper()
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
			WhereCost(100_000, eg.inFrac, eg.irFrac, eg.ifFrac, eg.df)
		assert(cost).Is(expected)
	}
	// baseline: all defaults => 100 * srcRows
	test(&Eg{}, 100_000)
	// irFrac only
	test(&Eg{irFrac: 1.0 / 1000}, 100)
	// dataFilter triggers pessimistic guard on inFrac
	test(&Eg{df: true, inFrac: .01}, 25_750)
	// ifFrac < 1 triggers pessimistic guard
	test(&Eg{ifFrac: 0.5}, 60_000)
	// ifFrac < 1 + inFrac with pessimistic guard
	test(&Eg{ifFrac: 0.1, inFrac: 0.5}, 17_500)
	// both irFrac and ifFrac
	test(&Eg{irFrac: 0.1, ifFrac: 0.5}, 6_000)
	// dataFilter + ifFrac
	test(&Eg{df: true, ifFrac: 0.5}, 60_000)
	// all parameters non-default
	test(&Eg{irFrac: 0.5, ifFrac: 0.5, df: true, inFrac: 0.5}, 18_750)
	// inFrac without pessimistic guard (no dataFilter, ifFrac=1)
	test(&Eg{inFrac: 0.5}, 50_000)
	// irFrac + inFrac: no guard (pure index range, ifFrac=1, no dataFilter)
	test(&Eg{irFrac: 0.1, inFrac: 0.5}, 5_000)
	// dataFilter=true with inFrac=1: guard is no-op, same cost as baseline
	test(&Eg{df: true}, 100_000)
	// irFrac + dataFilter + small inFrac: guard applied to narrowed range
	test(&Eg{irFrac: 0.5, df: true, inFrac: 0.01}, 12_875)
}

func TestWhere_Damp(t *testing.T) {
	frac := .25
	for i := range 5 {
		f := damp(i, frac)
		assert.T(t).This(f).Is(math.Pow(frac, 1.0/float64(int(1)<<i)))
	}
}

func TestWhere_CalcFracs(t *testing.T) {
	pc := func(cols ...string) map[string][]span {
		m := make(map[string][]span, len(cols))
		for _, c := range cols {
			m[c] = []span{valSpan(c)}
		}
		return m
	}
	cfs := func(args ...any) []colFrac {
		result := make([]colFrac, 0, len(args)/2)
		for i := 0; i < len(args); i += 2 {
			result = append(result, colFrac{col: args[i].(string), frac: args[i+1].(float64)})
		}
		return result
	}
	d := func(i int, args ...float64) float64 {
		frac := 1.0
		for _, a := range args {
			frac *= damp(i, a)
			i++
		}
		return frac
	}
	test := func(isel *idxSel, cfs []colFrac, pc map[string][]span, irf, iff, of float64) {
		t.Helper()
		overall := isel.calcFracs(cfs, pc, nil, 0)
		// fmt.Println(isel)
		assert.T(t).Msg("indexRangeFrac").This(isel.indexRangeFrac).Is(irf)
		assert.T(t).Msg("indexFilterFrac").This(isel.indexFilterFrac).Is(iff)
		assert.T(t).Msg("overall").This(overall).Is(of)
	}
	var isel *idxSel

	// nothing constrained
	isel = &idxSel{index: []string{"a", "b"}}
	test(isel, cfs(), pc(), 1, 1, 1)

	// prefix only
	isel = &idxSel{index: []string{"a", "b"}, prefixLen: 1, prefixFrac: .1}
	test(isel, cfs(), pc(), .1, 1, .1)

	// prefix only, overlapping
	isel = &idxSel{index: []string{"a", "b"}, prefixLen: 2, prefixFrac: .1}
	test(isel, cfs("a", .1, "b", .2), pc("a", "b"), .1, 1, .1)

	// prefix + filter
	isel = &idxSel{index: []string{"a", "b", "c"}, prefixLen: 2, prefixFrac: .01}
	test(isel, cfs("a", .1, "b", .2, "c", .4), pc("a", "b", "c"),
		.01, d(1, .4), d(0, .01, .4))

	// prefix + skip + filter
	isel = &idxSel{index: []string{"a", "b", "c", "d"},
		prefixLen: 1, prefixFrac: .1, skipStart: 2, skipLen: 1}
	test(isel, cfs("a", .1, "b", .2, "c", .3, "d", .4), pc("a", "b", "c", "d"),
		d(0, .1, .3), d(2, .2, .4), d(0, .1, .2, .3, .4))

	// prefix + filter (no stats for filter col)
	isel = &idxSel{index: []string{"a", "b"}, prefixLen: 1, prefixFrac: .1}
	test(isel, cfs("a", .1), pc("a", "b"),
		.1, d(1, unknownFrac), d(0, .1, unknownFrac))

	// no prefix, filter only
	isel = &idxSel{index: []string{"a", "b"}}
	test(isel, cfs("a", .1, "b", .2), pc("a", "b"),
		1, d(0, .1, .2), d(0, .1, .2))

	// prefix + filter, prefix frac is the range
	isel = &idxSel{index: []string{"a", "b", "c"}, prefixLen: 1, prefixFrac: .3}
	test(isel, cfs("a", .3, "b", .1, "c", .2), pc("a", "b", "c"),
		.3, d(1, .1, .2), d(0, .1, .2, .3))

	// no prefix. index filter with no estimable column (e.g. F(x)) fallback
	isel = &idxSel{index: []string{"a", "b"}}
	test(isel, cfs(), pc("x"), 1, 1, unknownFrac)

	// prefix + index filter with no estimable column (e.g. F(x)) → fallback
	isel = &idxSel{index: []string{"a", "b"}, prefixLen: 1, prefixFrac: .1}
	test(isel, cfs("a", .1), pc("a", "x"), .1, 1, d(0, .1, unknownFrac))
}

func TestWhere_AllSingleValuePrefix(t *testing.T) {
	assert := assert.T(t).This

	// helper to create colSpans map with single value spans
	makeColSpans := func(pairs ...string) map[string][]span {
		m := make(map[string][]span)
		for i := 0; i < len(pairs); i += 2 {
			m[pairs[i]] = []span{valSpan(pairs[i+1])}
		}
		return m
	}

	// basic case: two columns with values
	index := []string{"a", "b"}
	colSpans := makeColSpans("a", "x", "b", "y")
	prefixLen, org, ok := allSingleValuePrefix(index, true, colSpans)
	assert(ok).Is(true)
	assert(prefixLen).Is(2)
	assert(org).Is("x\x00\x00y")

	// trailing empty field: (a="x", b="") — org is trimmed to "x"
	colSpans = makeColSpans("a", "x", "b", "")
	prefixLen, org, ok = allSingleValuePrefix(index, true, colSpans)
	assert(ok).Is(true)
	assert(prefixLen).Is(2)
	assert(org).Is("x") // trailing separator trimmed by Encoder.String()

	// multiple trailing empty fields
	index = []string{"a", "b", "c"}
	colSpans = makeColSpans("a", "x", "b", "", "c", "")
	prefixLen, org, ok = allSingleValuePrefix(index, true, colSpans)
	assert(ok).Is(true)
	assert(prefixLen).Is(3)
	assert(org).Is("x") // both trailing separators trimmed

	// middle empty field (not trailing)
	index = []string{"a", "b", "c"}
	colSpans = makeColSpans("a", "x", "b", "", "c", "z")
	prefixLen, org, ok = allSingleValuePrefix(index, true, colSpans)
	assert(ok).Is(true)
	assert(prefixLen).Is(3)
	assert(org).Is("x\x00\x00\x00\x00z")

	// single column (no encoding)
	index = []string{"a"}
	colSpans = makeColSpans("a", "x")
	prefixLen, org, ok = allSingleValuePrefix(index, false, colSpans)
	assert(ok).Is(true)
	assert(prefixLen).Is(1)
	assert(org).Is("x")

	// partial prefix (b has no span)
	index = []string{"a", "b"}
	colSpans = makeColSpans("a", "x")
	prefixLen, org, ok = allSingleValuePrefix(index, true, colSpans)
	assert(ok).Is(true)
	assert(prefixLen).Is(1)
	assert(org).Is("x")
}

func TestBuildIdxSelTrailingEmpty(t *testing.T) {
	assert := assert.T(t).This

	// Verify the range end construction matches Table's selEnd approach.
	// For prefix values with trailing empty fields, the end must include
	// the empty field separator so that keys with non-empty values are excluded.

	// Case 1: trailing empty field (a="x", b="")
	// Expected end: "x\x00\x00\x00\x00Max" (same as selEnd)
	var enc ixkey.Encoder
	enc.Add("x")
	enc.Add("")
	enc.Add(ixkey.Max)
	end := enc.String()
	assert(end).Is("x\x00\x00\x00\x00\xff\xff\xff\xff\xff\xff\xff\xff")

	// Case 2: no trailing empty (a="x", b="y")
	enc.Add("x")
	enc.Add("y")
	enc.Add(ixkey.Max)
	end = enc.String()
	assert(end).Is("x\x00\x00y\x00\x00\xff\xff\xff\xff\xff\xff\xff\xff")

	// Case 3: multiple trailing empty (a="x", b="", c="")
	enc.Add("x")
	enc.Add("")
	enc.Add("")
	enc.Add(ixkey.Max)
	end = enc.String()
	assert(end).Is("x\x00\x00\x00\x00\x00\x00\xff\xff\xff\xff\xff\xff\xff\xff")

	// Verify that a key with non-empty b is excluded from the trailing-empty range
	// Key for (a="x", b="y"): "x\x00\x00y"
	keyWithNonEmptyB := "x\x00\x00y"
	endTrailingEmpty := "x\x00\x00\x00\x00\xff\xff\xff\xff\xff\xff\xff\xff"
	assert(keyWithNonEmptyB > endTrailingEmpty).Is(true)

	// Key for (a="x", b=""): "x" (trimmed)
	keyWithEmptyB := "x"
	assert(keyWithEmptyB >= "x").Is(true)
	assert(keyWithEmptyB < endTrailingEmpty).Is(true)
}

func TestWhere_nonexistent(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm(`create table (num, name, parent, group) 
		key(num) key(name,group) index(parent,name)`)
	for i := range 100 {
		s := strconv.Itoa(i)
		if i%10 == 0 {
			db.act("insert { num: " + s + ", name: 'a" + s + "', " +
				"parent: " + s + ", group: " + s + " } into table")
		} else {
			p := strconv.Itoa(rand.IntN(10) * 10)
			db.act("insert { num: " + s +
				", name: 'a" + s + "', parent: " + p + ", group: -1 } into table")
		}
	}

	test := func(parent string) {
		t.Helper()
		tran := db.NewReadTran()
		q := ParseQuery("table where group >= -1 and parent = "+parent, tran, nil)
		// Regression: a nonexistent prefix point must not collapse wfrac
		// to zero and cause the optimizer to pick a worse index.
		// Both existing and nonexistent parent values should select the
		// (parent,name,num) index since parent is the leading prefix column.
		q, _, _ = SetupReq(q, ReadMode, tran, NoneReq(1))
		assert.T(t).That(strings.Contains(String(q), "table^(parent,name,num)"))
	}
	test("20")   // exists
	test("9999") // nonexistent
}

func TestWhere_singleton_lookup(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create table (a,b,c,d) key(a,b)")
	db.act("insert { a: 1, b: 2, c: 3, d: 4 } into table")
	db.act("insert { a: 4, b: 5, c: 6, d: 7 } into table")
	tran := db.NewReadTran()
	// where clause matches full key => singleton
	q := ParseQuery("table where a=1 and b=2", tran, nil)
	key := []string{"a", "b"}
	q, _, _ = SetupReq(q, ReadMode, tran, UniqueReq(key, 1))
	w := q.(*Where)
	assert.T(t).That(w.singleton)
	assert.T(t).This(w.idxSelBase.String()).Is("(a,b) a,b: <1,2> = singleton")
	th := &Thread{}
	sels := Sels{{"a", Pack(SuInt(1))}, {"b", Pack(SuInt(2))}}
	row := q.Lookup(th, sels)
	hdr := q.Header()
	assert.This(row2str(hdr, row)).Is("a=1 b=2 c=3 d=4")
}

func TestWhere_indexRanges(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create table (a,b,c,d,e,f) key(a,b,c,d,e,f)")
	vals := []string{"", "1"}
	for _, a := range vals {
		for _, b := range vals {
			for _, c := range vals {
				for _, d := range vals {
					for _, e := range vals {
						for _, f := range vals {
							db.act(fmt.Sprintf(
								`insert 
								{ a: %q, b: %q, c: %q, d: %q, e: %q, f: %q } 
								into table`, a, b, c, d, e, f))
						}
					}
				}
			}
		}
	}
	for _, a := range vals {
		for _, b := range vals {
			query := fmt.Sprintf("table where a is %q and b is %q", a, b)
			db.queryCompare(t, query)
		}
	}
	for _, d := range vals {
		for _, e := range vals {
			query := fmt.Sprintf("table where d is %q and e is %q", d, e)
			db.queryCompare(t, query)
		}
	}
	for _, a := range vals {
		for _, b := range vals {
			for _, d := range vals {
				for _, e := range vals {
					query := fmt.Sprintf(
						"table where a is %q and b is %q and d is %q and e is %q",
						a, b, d, e)
					db.queryCompare(t, query)
				}
			}
		}
	}
}

func (hdb *heapdb) queryCompare(t *testing.T, query string) {
	db := hdb.Database
	tran := db.NewReadTran()
	q := ParseQuery(query, tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)

	th := &Thread{}
	hdr := q.Header()

	simple := q.Simple(th)
	hs := NewQueryHasher(hdr).CheckDups()
	for _, row := range simple {
		hs.Row(row)
	}

	q.Rewind()
	hg := NewQueryHasher(hdr).CheckDups()
	var get []Row
	for row := q.Get(th, Next); row != nil; row = q.Get(th, Next) {
		hg.Row(row)
		get = append(get, row)
	}

	if hg.Result(true) == hs.Result(true) {
		return
	}

	fmt.Println("optimized:", String(q))
	for i, row := range simple {
		fmt.Printf("simple[%d]: ", i)
		for _, fld := range hdr.GetFields() {
			if fld != "-" {
				fmt.Printf("%s=%q ", fld, row.GetRawVal(hdr, fld, nil, nil))
			}
		}
		fmt.Println()
	}
	for i, row := range get {
		fmt.Printf("Get[%d]: ", i)
		for _, fld := range hdr.GetFields() {
			if fld != "-" {
				fmt.Printf("%s=%q ", fld, row.GetRawVal(hdr, fld, nil, nil))
			}
		}
		fmt.Println()
	}
	assert.T(t).This(hg.Result(true)).Is(hs.Result(true))
}
