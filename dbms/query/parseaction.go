// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
)

type Action interface {
	execute(th *Thread, ut *db19.UpdateTran) int
	String() string
}

type actionParser struct {
	queryParser
}

// ParseAction parses insert, update, and delete actions
func ParseAction(src string, t QueryTran, sv *Sviews) Action {
	p := actionParser{*NewQueryParser(src, t, sv)}
	result := p.action()
	if p.Token != tok.Eof {
		p.Error("did not parse all input")
	}
	return result
}

func (p *actionParser) action() Action {
	switch {
	case p.MatchIf(tok.Insert):
		return p.insert()
	case p.MatchIf(tok.Update):
		return p.update()
	case p.MatchIf(tok.Delete):
		return p.delete()
	default:
		panic(p.Error("action must be insert, update, or delete"))
	}
}

func (p *actionParser) insert() Action {
	if p.Token == tok.LCurly || p.Token == tok.LBracket {
		return p.insertRecord()
	}
	return p.insertQuery()
}

func (p *actionParser) insertRecord() Action {
	record := p.record()
	p.Match(tok.Into)
	query := p.baseQuery()
	return &insertRecordAction{record: record, query: query}
}

func (p *actionParser) record() *SuRecord {
	if p.Token != tok.LCurly && p.Token != tok.LBracket {
		p.Error("record expected e.g. { a: 1, b: 2 }")
	}
	return p.Const().(*SuRecord)
}

func (p *actionParser) insertQuery() Action {
	query := p.baseQuery()
	p.Match(tok.Into)
	table := p.Text
	p.Match(tok.Identifier)
	return &insertQueryAction{query: query, table: table}
}

func (p *actionParser) update() Action {
	query := p.baseQuery()
	p.Match(tok.Set)
	var cols []string
	var exprs []ast.Expr
	for p.Token.IsIdent() {
		cols = append(cols, p.MatchIdent())
		p.Match(tok.Eq)
		exprs = append(exprs, p.Expression())
		p.MatchIf(tok.Comma)
	}
	return &updateAction{query: query, cols: cols, exprs: exprs}
}

func (p *actionParser) delete() Action {
	query := p.baseQuery()
	return &deleteAction{query: query}
}
