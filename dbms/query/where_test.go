// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"math"
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
		pf, _ := perField(w.expr.Exprs, w.source.Header().Physical())
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
		w.optInit()
		pf, _ := perField(w.expr.Exprs, w.source.Header().Physical())
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
	
	table = "comp2"
	test("a is 1 and b is 2 and c is 3", "a,b,c: 1,2,3")
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
		c := &ast.Context{Hdr: hdr, Row: []DbRec{{Record: rec}}}
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
	test("inven where item > 2 and item < 4 and qty", 10)
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
		pf, _ := perField(w.expr.Exprs, w.source.Header().Physical())
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
		w := ParseQuery(table + " where " + where, tran, nil).(*Where)
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
