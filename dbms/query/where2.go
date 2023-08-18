// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/exp/slices"
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
	if !expr.CanEvalRaw(fields) {
		return "", nil
	}
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
	// TODO unary e.g. not Number?(x)
	}
	return "", nil
}

var sideMin = side{}
var sideMax = side{val: ixkey.Max}

var conflictSpans = []span{{org: sideMax, end: sideMin}}
var isntqqSpans = []span{{org: side{val: "", inc: true}, end: sideMax}}

func binarySpan(bin *ast.Binary, fields []string) (string, []span) {
	// TOTO make LT, GT, GTE match the language
	// currently they are consistent with previous behavior
	// Changing them will require application code changes.
	col := bin.Lhs.(*ast.Ident).Name
	val := bin.Rhs.(*ast.Constant).Packed
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

func rangeSpan(r *ast.InRange, fields []string) (string, []span) {
	return r.E.(*ast.Ident).Name, []span{{
		org: side{val: r.Org.(*ast.Constant).Packed, inc: r.OrgTok == tok.Gt},
		end: side{val: r.End.(*ast.Constant).Packed, inc: r.EndTok == tok.Lte}}}
}

func inSpan(in *ast.In, fields []string) (string, []span) {
	spans := make([]span, 0, len(in.Packed))
	for _, v := range in.Packed {
		spans = append(spans,
			span{org: side{val: v}, end: side{val: v, inc: true}})
	}
	sortByOrg(spans)
	return in.E.(*ast.Ident).Name, spans
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

func typeSpan(call *ast.Call, fields []string) (string, []span) {
	fn := call.Fn.(*ast.Ident)
	id := call.Args[0].E.(*ast.Ident).Name
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

func isColumn(e ast.Expr, cols []string) string {
	// see also: ast.IsColumn
	if id, ok := e.(*ast.Ident); ok && (slices.Contains(cols, id.Name) ||
		(strings.HasSuffix(id.Name, "_lower!") &&
			slices.Contains(cols, strings.TrimSuffix(id.Name, "_lower!")))) {
		return id.Name
	}
	return ""
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
	return sideLte(maxSide(x.org, y.org), minSide(x.end, y.end))
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
		if sideLt(spans1[i1].end, spans2[i2].end) {
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
	slices.SortFunc(spans, func(x, y span) bool {
		return sideLt(x.org, y.org)
	})
}

func sideLte(x, y side) bool {
	return !sideLt(y, x)
}

func sideLt(x, y side) bool {
	cmp := strings.Compare(x.val, y.val)
	if cmp == 0 {
		cmp = cmpBool(x.inc, y.inc)
	}
	return cmp < 0
}

func minSide(x, y side) side {
	if sideLt(x, y) {
		return x
	}
	return y
}

func maxSide(x, y side) side {
	if sideLt(x, y) {
		return y
	}
	return x
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
