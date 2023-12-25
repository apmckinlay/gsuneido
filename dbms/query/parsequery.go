// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/str"
)

type queryParser struct {
	t         QueryTran
	sviews    *Sviews
	viewNest  []string
	wrapViews bool
	compile.Parser
}

func NewQueryParser(src string, t QueryTran, sv *Sviews) *queryParser {
	return &queryParser{Parser: *compile.QueryParser(src), t: t, sviews: sv}
}

func ParseQuery(src string, t QueryTran, sv *Sviews) Query {
	return parseQuery(src, t, sv, nil, false)
}

func parseQuery(src string, t QueryTran, sv *Sviews, viewNest []string, wrapViews bool) Query {
	p := NewQueryParser(src, t, sv)
	p.viewNest = viewNest
	p.wrapViews = wrapViews
	result := p.sort()
	if p.Token != tok.Eof {
		p.Error("did not parse all input")
	}
	return result
}

func (p *queryParser) sort() Query {
	q := p.baseQuery()
	if p.MatchIf(tok.Sort) {
		reverse := p.MatchIf(tok.Reverse)
		cols := p.commaList()
		q = NewSort(q, reverse, cols)
	}
	return q
}

func (p *queryParser) baseQuery() Query {
	q := p.source()
	for p.operation(&q) {
	}
	return q
}

func (p *queryParser) source() Query {
	if p.MatchIf(tok.LParen) {
		q := p.baseQuery()
		p.Match(tok.RParen)
		return q
	}
	return p.table()
}

func (p *queryParser) table() Query {
	table := p.MatchIdent()
	if !slices.Contains(p.viewNest, table) {
		if def := p.getView(table); def != "" {
			q := parseQuery(def, p.t, p.sviews, append(p.viewNest, table),
				p.wrapViews)
			if p.wrapViews {
				q = NewView(table, q)
			}
			return q
		}
	}
	return NewTable(p.t, table)
}

func (p *queryParser) getView(name string) string {
	def := ""
	if p.sviews != nil { // should only be nil from tests
		def = p.sviews.GetSview(name)
	}
	if def == "" {
		def = p.t.GetView(name)
	}
	return def
}

func (p *queryParser) operation(pq *Query) bool {
	switch {
	case p.MatchIf(tok.Extend):
		*pq = p.extend(*pq)
	case p.MatchIf(tok.Intersect):
		*pq = p.intersect(*pq)
	case p.MatchIf(tok.Join):
		*pq = p.join(*pq)
	case p.MatchIf(tok.Leftjoin):
		*pq = p.leftjoin(*pq)
	case p.MatchIf(tok.Minus):
		*pq = p.minus(*pq)
	case p.MatchIf(tok.Project):
		*pq = p.project(*pq)
	case p.MatchIf(tok.Remove):
		*pq = p.remove(*pq)
	case p.MatchIf(tok.Rename):
		*pq = p.rename(*pq)
	case p.MatchIf(tok.Summarize):
		*pq = p.summarize(*pq)
	case p.MatchIf(tok.Times):
		*pq = p.times(*pq)
	case p.MatchIf(tok.Union):
		*pq = p.union(*pq)
	case p.MatchIf(tok.Where):
		*pq = p.where(*pq)
	default:
		return false
	}
	return true
}

func (p *queryParser) extend(q Query) Query {
	cols := make([]string, 0, 4)
	exprs := make([]ast.Expr, 0, 4)
	for {
		cols = append(cols, p.MatchIdent())
		var expr ast.Expr
		if p.MatchIf(tok.Eq) {
			expr = p.Expression()
		}
		exprs = append(exprs, expr)
		if !p.MatchIf(tok.Comma) {
			break
		}
	}
	return NewExtend(q, cols, exprs)
}

func (p *queryParser) intersect(q Query) Query {
	q2 := p.source()
	return NewIntersect(q, q2)
}

func (p *queryParser) join(q Query) Query {
	by := p.joinBy()
	q2 := p.source()
	return NewJoin(q, q2, by, p.t)
}

func (p *queryParser) leftjoin(q Query) Query {
	by := p.joinBy()
	q2 := p.source()
	return NewLeftJoin(q, q2, by, p.t)
}

func (p *queryParser) joinBy() []string {
	if p.MatchIf(tok.By) {
		by := p.parenList()
		if len(by) == 0 {
			p.Error("invalid empty join by")
		}
		return by
	}
	return nil
}

func (p *queryParser) minus(q Query) Query {
	q2 := p.source()
	return NewMinus(q, q2)
}

func (p *queryParser) project(q Query) Query {
	cols := p.commaList()
	return NewProject(q, cols)
}

func (p *queryParser) remove(q Query) Query {
	cols := p.commaList()
	return NewRemove(q, cols)
}

func (p *queryParser) rename(q Query) Query {
	var from, to []string
	for {
		from = append(from, p.MatchIdent())
		p.Match(tok.To)
		to = append(to, p.MatchIdent())
		if !p.MatchIf(tok.Comma) {
			break
		}
	}
	return NewRename(q, from, to)
}

func (p *queryParser) summarize(q Query) Query {
	var hint sumHint
	remainder := strings.TrimSpace(p.Lxr.Source()[p.EndPos:])
	if strings.HasPrefix(remainder, "/*small*/") {
		hint = "small"
	} else if strings.HasPrefix(remainder, "/*large*/") {
		hint = "large"
	}
	by := p.sumBy()
	cols, ops, ons := p.sumOps()
	return NewSummarize(q, hint, by, cols, ops, ons)
}

func (p *queryParser) sumBy() []string {
	var by []string
	for p.Token.IsIdent() &&
		!isSumOp(p.Token) &&
		p.Lxr.Ahead(1).Token != tok.Eq {
		by = append(by, p.MatchIdent())
		p.Match(tok.Comma)
	}
	return by
}

func (p *queryParser) sumOps() (cols, ops, ons []string) {
	for {
		var col, op, on string
		if p.Lxr.Ahead(1).Token == tok.Eq {
			col = p.MatchIdent()
			p.Match(tok.Eq)
		}
		if !isSumOp(p.Token) {
			p.Error("expected count, total, min, max, or list")
		}
		op = str.ToLower(p.MatchIdent())
		if op != "count" {
			on = p.MatchIdent()
		}
		cols = append(cols, col)
		ops = append(ops, op)
		ons = append(ons, on)
		if !p.MatchIf(tok.Comma) {
			break
		}
	}
	return
}

func isSumOp(t tok.Token) bool {
	return tok.SummarizeStart < t && t < tok.SummarizeEnd
}

func (p *queryParser) times(q Query) Query {
	q2 := p.source()
	return NewTimes(q, q2)
}

func (p *queryParser) union(q Query) Query {
	q2 := p.source()
	return NewUnion(q, q2)
}

func (p *queryParser) where(q Query) Query {
	p.EqToIs = true
	defer func() { p.EqToIs = false }()
	expr := p.Expression()
	return NewWhere(q, expr, p.t)
}

func (p *queryParser) parenList() []string {
	p.Match(tok.LParen)
	if p.MatchIf(tok.RParen) {
		return nil
	}
	list := p.commaList()
	p.Match(tok.RParen)
	return list
}

func (p *queryParser) commaList() []string {
	list := make([]string, 0, 4)
	for {
		list = append(list, p.MatchIdent())
		if !p.MatchIf(tok.Comma) {
			break
		}
	}
	return list
}

func (p *queryParser) Expression() ast.Expr {
	p.InitFuncInfo()
	return p.Parser.Expression()
}
