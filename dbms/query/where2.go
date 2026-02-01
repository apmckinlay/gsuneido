// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"slices"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/set"
)

/*
Expressions that can be applied to an index:
- <ident> <op> <constant> where op is one of: <, <=, >, >=, is, isnt
	- isnt is treated as (id < constant or id > constant)
- <ident> in (<constants>)
- InRange which is folded from: x >[=] org and x <[=] end
- Date?, Number?, String?
- Nary Or of applicable expressions e.g. (x is 0 or x > 10)
*/

type span struct {
	// org is inclusive (>=)
	org side
	// end is exclusive (<)
	end side
}

type side struct {
	val string
	inc bool
}

// perField returns the spans for each field, or nil if there is a conflict.
// It is called by NewWhere. The result is later used by perIndex.
// It depends on CanEvalRaw already being done.
func perField(exprs []ast.Expr, fields []string) (map[string][]span, bool) {
	exprMore := false
	result := make(map[string][]span)
	for _, expr := range exprs {
		if col, espans := exprToSpans(expr, fields); espans != nil {
			x := intersectSpans(result[col], espans)
			if x == nil {
				return nil, false // conflict
			}
			result[col] = x
		} else {
			exprMore = true // since expr not used for column span
		}
	}
	return result, exprMore
}

// exprToSpans returns the spans for an expression, or nil if not indexable
func exprToSpans(expr ast.Expr, fields []string) (string, []span) {
	switch expr := expr.(type) {
	case *ast.Binary:
		return binarySpan(expr, fields)
	case *ast.InRange:
		return rangeSpan(expr, fields)
	case *ast.In:
		return inSpan(expr, fields)
	case *ast.Nary:
		if expr.Tok == tok.Or {
			return orSpan(expr, fields)
		}
	case *ast.Call:
		return typeSpan(expr, fields)
	case *ast.Unary:
		// TODO unary not
		if expr.Tok == tok.LParen {
			return exprToSpans(expr.E, fields)
		}
	}
	return "", nil
}

var sideMin = side{}
var sideMax = side{val: ixkey.Max}

var conflictSpans = []span{{org: sideMax, end: sideMin}}
var isntqqSpans = []span{{org: side{val: "", inc: true}, end: sideMax}}

func binarySpan(bin *ast.Binary, flds []string) (string, []span) {
	// depends on folder putting field on the left and constant on the right
	col, ok := ast.IsField(bin.Lhs, flds)
	if !ok {
		return "", nil
	}
	c, ok := bin.Rhs.(*ast.Constant)
	if !ok {
		return "", nil
	}
	val := c.Packed
	switch bin.Tok {
	case tok.Lt:
		if val == "" {
			return col, conflictSpans
		}
		return col, []span{{end: side{val: val}}}
	case tok.Lte:
		return col, []span{{end: side{val: val, inc: true}}}
	case tok.Gt:
		return col, []span{{org: side{val: val, inc: true}, end: sideMax}}
	case tok.Gte:
		return col, []span{{org: side{val: val}, end: sideMax}}
	case tok.Is:
		return col, []span{valSpan(val)}
	case tok.Isnt:
		if val == "" {
			// treat isnt "" as > "" for packed usage
			return col, isntqqSpans
		} else {
			return col, []span{
				{end: side{val: val}},                          // < val
				{org: side{val: val, inc: true}, end: sideMax}} // > val
		}
	}
	return "", nil
}

func valSpan(val string) span {
	return span{org: side{val: val}, end: side{val: val, inc: true}}
}

func rangeSpan(r *ast.InRange, flds []string) (string, []span) {
	fld, ok := ast.IsField(r.E, flds)
	if !ok {
		return "", nil
	}
	org, ok := r.Org.(*ast.Constant)
	if !ok {
		return "", nil
	}
	end, ok := r.End.(*ast.Constant)
	if !ok {
		return "", nil
	}
	return fld, []span{{
		org: side{val: org.Packed, inc: r.OrgTok == tok.Gt},
		end: side{val: end.Packed, inc: r.EndTok == tok.Lte}}}
}

func inSpan(in *ast.In, flds []string) (string, []span) {
	fld, ok := ast.IsField(in.E, flds)
	if !ok {
		return "", nil
	}
	spans := make([]span, 0, len(in.Exprs))
	for _, e := range in.Exprs {
		c, ok := e.(*ast.Constant)
		if !ok {
			return "", nil
		}
		spans = set.AddUnique(spans,
			span{org: side{val: c.Packed}, end: side{val: c.Packed, inc: true}})
	}
	sortByOrg(spans)
	return fld, spans
}

var numSpan = []span{{
	org: side{val: string(rune(PackMinus))},
	end: side{val: string(rune(PackPlus + 1))}}}
var dateSpan = []span{{
	org: side{val: string(rune(PackDate))},
	end: side{val: string(rune(PackDate + 1))}}}
var strSpan = []span{
	valSpan(""),
	{org: side{val: string(rune(PackString))},
		end: side{val: string(rune(PackString + 1))}}}

func typeSpan(call *ast.Call, flds []string) (string, []span) {
	if !call.RawEval {
		return "", nil
	}
	fn := call.Fn.(*ast.Ident)
	id, ok := ast.IsField(call.Args[0].E, flds)
	if !ok {
		return "", nil
	}
	switch fn.Name {
	case "Number?":
		return id, numSpan
	case "String?":
		return id, strSpan
	case "Date?":
		return id, dateSpan
	}
	return "", nil
}

func orSpan(nary *ast.Nary, fields []string) (string, []span) {
	var col string
	var spans []span
	for i, expr := range nary.Exprs {
		c, espans := exprToSpans(expr, fields)
		if espans == nil {
			return "", nil
		}
		if i == 0 {
			col = c
		} else if col != c {
			return "", nil
		}
		for _, espan := range espans {
			spans = mergeSpan(spans, espan)
		}
	}
	sortByOrg(spans)
	return col, spans
}

func mergeSpan(spans []span, add span) []span {
	dst := 0
	var ok bool
	for i, sp := range spans {
		if add, ok = mergeSpans(sp, add); !ok {
			spans[dst] = spans[i]
			dst++
		}
	}
	return append(spans[:dst], add)
}

func mergeSpans(x, y span) (span, bool) {
	if spansOverlap(x, y) {
		// union
		return span{org: minSide(x.org, y.org), end: maxSide(x.end, y.end)}, true
	}
	return y, false
}

func spansOverlap(x, y span) bool {
	return sideCmp(maxSide(x.org, y.org), minSide(x.end, y.end)) <= 0
}

//-------------------------------------------------------------------

func intersectSpans(spans1, spans2 []span) []span {
	if spans1 == nil {
		for _, sp := range spans2 {
			if sp.none() {
				return nil // conflict
			}
		}
		return spans2
	}
	sortByOrg(spans1)
	sortByOrg(spans2)
	i1 := 0
	i2 := 0
	var result []span
	for i1 < len(spans1) && i2 < len(spans2) {
		if span, ok := intersectSpan(spans1[i1], spans2[i2]); ok {
			result = append(result, span)
		}
		if sideCmp(spans1[i1].end, spans2[i2].end) < 0 {
			i1++
		} else {
			i2++
		}
	}
	return result
}

func intersectSpan(x, y span) (span, bool) {
	result := span{org: maxSide(x.org, y.org), end: minSide(x.end, y.end)}
	if result.none() {
		return span{}, false
	}
	// fmt.Println(x, "intersect", y, "result", result)
	return result, true
}

func (sp span) none() bool {
	cmp := strings.Compare(sp.org.val, sp.end.val)
	if cmp != 0 {
		return cmp > 0
	}
	return sp.org.inc || !sp.end.inc
}

func sortByOrg(spans []span) {
	slices.SortFunc(spans, func(x, y span) int {
		return sideCmp(x.org, y.org)
	})
}

func sideCmp(x, y side) int {
	cmp := strings.Compare(x.val, y.val)
	if cmp == 0 {
		cmp = cmpBool(x.inc, y.inc)
	}
	return cmp
}

func minSide(x, y side) side {
	if sideCmp(x, y) <= 0 {
		return x
	}
	return y
}

func maxSide(x, y side) side {
	if sideCmp(x, y) >= 0 {
		return x
	}
	return y
}

func cmpBool(x, y bool) int {
	if x == y {
		return 0
	}
	if x {
		return 1
	}
	return -1
}

//-------------------------------------------------------------------

func (sp span) String() string {
	if sp.org.val == sp.end.val && !sp.org.inc && sp.end.inc {
		return packToStr(sp.org.val)
	}
	if sp.org.val == "" && !sp.org.inc {
		if sp.end.inc {
			return "<=" + packToStr(sp.end.val)
		} else {
			return "<" + packToStr(sp.end.val)
		}
	}
	if sp.end == sideMax {
		if sp.org.inc {
			return ">" + packToStr(sp.org.val)
		} else {
			return ">=" + packToStr(sp.org.val)
		}
	}
	s := ">"
	if !sp.org.inc {
		s += "="
	}
	s += packToStr(sp.org.val) + "_<"
	if sp.end.inc {
		s += "="
	}
	return s + packToStr(sp.end.val)
}
