// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/util/str"
)

type queryParser struct {
	compile.Parser
}

func NewQueryParser(src string) *queryParser {
	return &queryParser{*compile.QueryParser(src)}
}

func ParseQuery(src string) Query {
	p := NewQueryParser(src)
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
		q = &Sort{Query1: Query1{source: q},
			reverse: reverse, columns: p.commaList()}
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
	return NewTable(p.MatchIdent())
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
	return &Extend{Query1: Query1{source: q}, cols: cols, exprs: exprs}
}

func (p *queryParser) intersect(q Query) Query {
	return &Intersect{Compatible: Compatible{
		Query2: Query2{Query1: Query1{source: q}, source2: p.source()}}}
}

func (p *queryParser) join(q Query) Query {
	by := p.joinBy()
	return &Join{Query2: Query2{Query1: Query1{source: q},
		source2: p.source()}, by: by}
}

func (p *queryParser) leftjoin(q Query) Query {
	by := p.joinBy()
	return &LeftJoin{Join: Join{Query2: Query2{Query1: Query1{source: q},
		source2: p.source()}, by: by}}
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
	return &Minus{Compatible: Compatible{
		Query2: Query2{Query1: Query1{source: q}, source2: p.source()}}}
}

func (p *queryParser) project(q Query) Query {
	return &Project{Query1: Query1{source: q}, columns: p.commaList()}
}

func (p *queryParser) remove(q Query) Query {
	return &Remove{Query1: Query1{source: q}, columns: p.commaList()}
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
	return &Rename{Query1: Query1{source: q}, from: from, to: to}
}

func (p *queryParser) summarize(q Query) Query {
	su := &Summarize{Query1: Query1{source: q}}
	su.by = p.sumBy()
	p.sumOps(su)
	return su
}

func (p *queryParser) sumBy() []string {
	var by []string
	for p.Token == tok.Identifier &&
		!isSumOp(p.Token) &&
		p.Lxr.Ahead(1).Token != tok.Eq {
		by = append(by, p.MatchIdent())
		p.Match(tok.Comma)
	}
	return by
}

func (p *queryParser) sumOps(su *Summarize) {
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
		su.cols = append(su.cols, col)
		su.ops = append(su.ops, op)
		su.ons = append(su.ons, on)
		if !p.MatchIf(tok.Comma) {
			break
		}
	}
}

func isSumOp(t tok.Token) bool {
	return tok.SummarizeStart < t && t < tok.SummarizeEnd
}

func (p *queryParser) times(q Query) Query {
	return &Times{Query2: Query2{Query1: Query1{source: q},
		source2: p.source()}}
}

func (p *queryParser) union(q Query) Query {
	return &Union{Compatible: Compatible{
		Query2: Query2{Query1: Query1{source: q}, source2: p.source()}}}
}

func (p *queryParser) where(q Query) Query {
	expr := p.Expression()
	if nary, ok := expr.(*ast.Nary); !ok || nary.Tok != tok.And {
		expr = &ast.Nary{Tok: tok.And, Exprs: []ast.Expr{expr}}
	}
	return &Where{Query1: Query1{source: q}, expr: expr.(*ast.Nary)}
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
