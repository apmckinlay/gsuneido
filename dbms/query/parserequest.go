// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/compile/lexer"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/util/str"
)

type reqparser struct {
	compile.ParserBase
}

func NewRequestParser(src string) *reqparser {
	lxr := lexer.NewQueryLexer(src)
	p := &reqparser{compile.ParserBase{Lxr: lxr}}
	p.Next()
	return p
}

type Schema = schema.Schema
type Index = schema.Index

type Request struct {
	Action    string
	SubAction string
	Schema
	Renames []Rename1
}

func ParseRequest(src string) *Request {
	p := NewRequestParser(src)
	result := p.request()
	if p.Token != tok.Eof {
		p.Error("did not parse all input")
	}
	return result
}

func (p *reqparser) request() *Request {
	switch {
	case p.MatchIf(tok.Create):
		return &Request{Action: "create", Schema: p.schema(true)}
	case p.MatchIf(tok.Ensure):
		return &Request{Action: "ensure", Schema: p.schema(false)}
	case p.MatchIf(tok.Drop):
		table := p.MatchIdent()
		return &Request{Action: "drop", Schema: Schema{Table: table}}
	case p.MatchIf(tok.Rename):
		rename := p.rename1()
		return &Request{Action: "rename", Renames: []Rename1{rename}}
	case p.MatchIf(tok.Alter):
		return p.alter()
	//TODO: View, Sview
	default:
		panic("invalid request")
	}
}

func (p *reqparser) rename1() Rename1 {
	from := p.MatchIdent()
	p.Match(tok.To)
	to := p.MatchIdent()
	return Rename1{from: from, to: to}
}

func (p *reqparser) alter() *Request {
	table := p.MatchIdent()
	switch {
	case p.MatchIf(tok.Create):
		return &Request{Action: "alter", SubAction: "create",
			Schema: p.schema2(table, false)}
	case p.MatchIf(tok.Drop):
		return &Request{Action: "alter", SubAction: "drop",
			Schema: p.schema2(table, false)}
	case p.MatchIf(tok.Rename):
		return &Request{Action: "alter", SubAction: "rename",
			Schema: p.schema2(table, false), Renames: p.renames()}
	default:
		panic("invalid request")
	}
}

func (p *reqparser) renames() []Rename1 {
	var renames []Rename1
	for {
		renames = append(renames, p.rename1())
		if !p.MatchIf(tok.Comma) {
			return renames
		}
	}
}

func (p *reqparser) schema(full bool) Schema {
	table := p.MatchIdent()
	return p.schema2(table, full)
}

func (p *reqparser) schema2(table string, full bool) Schema {
	columns, derived := p.columns(full)
	indexes := p.indexes(columns, derived, full)
	return Schema{Table: table, Columns: columns, Derived: derived, Indexes: indexes}
}

func (p *reqparser) columns(full bool) (columns, derived []string) {
	if !full && p.Token != tok.LParen {
		return
	}
	p.Match(tok.LParen)
	columns = make([]string, 0, 8)
	for p.Token != tok.RParen {
		if p.MatchIf(tok.Sub) {
			columns = append(columns, "-")
		} else {
			col := p.MatchIdent()
			if str.Capitalized(col) {
				derived = append(derived, col)
			} else if strings.HasSuffix(col, "_lower!") {
				if full &&
					!str.List(columns).Has(strings.TrimSuffix(col, "_lower!")) {
					panic("_lower! base column not found")
				}
				derived = append(derived, col)
			} else {
				columns = append(columns, col)
			}

		}
		p.MatchIf(tok.Comma)
	}
	p.Match(tok.RParen)
	return columns, derived
}

func (p *reqparser) indexes(columns, derived []string, full bool) []Index {
	hasKey := false
	indexes := make([]Index, 0, 4)
	for ix := p.index(columns, derived, full); ix != nil; ix = p.index(columns, derived, full) {
		indexes = append(indexes, *ix)
		hasKey = hasKey || ix.Mode == 'k'
	}
	if full && !hasKey {
		panic("key required")
	}
	return indexes
}

func (p *reqparser) index(columns, derived []string, full bool) *Index {
	if p.Token != tok.Key && p.Token != tok.Index {
		return nil
	}
	mode := int('i')
	if p.Token == tok.Key {
		mode = 'k'
	}
	p.Next()
	if mode != 'k' && p.MatchIf(tok.Unique) {
		mode = 'u'
	}
	ixcols := p.indexColumns(columns, derived, full)
	if mode != 'k' && len(ixcols) == 0 {
		p.Error("index columns must not be empty")
	}
	ix := &Index{Columns: ixcols, Mode: mode}
	ix.Fktable, ix.Fkcolumns, ix.Fkmode = p.foreignKey()
	return ix
}

func (p *reqparser) indexColumns(columns, derived []string, full bool) []string {
	p.Match(tok.LParen)
	ixcols := make([]string, 0, 8)
	for p.Token != tok.RParen {
		col := p.MatchIdent()
		if full && !str.List(columns).Has(col) &&
			(!strings.HasSuffix(col, "_lower!") || !str.List(derived).Has(col)) {
			p.Error("invalid index column: " + col)
		}
		ixcols = append(ixcols, col)
		p.MatchIf(tok.Comma)
	}
	p.Match(tok.RParen)
	return ixcols
}

func (p *reqparser) foreignKey() (table string, columns []string, mode int) {
	if !p.MatchIf(tok.In) {
		return
	}
	table = p.MatchIdent()
	if p.MatchIf(tok.LParen) {
		for p.Token != tok.RParen {
			columns = append(columns, p.Text)
			p.MatchIdent()
			p.MatchIf(tok.Comma)
		}
		p.Next()
	}
	mode = schema.Block
	if p.MatchIf(tok.Cascade) {
		mode = schema.Cascade
		if p.MatchIf(tok.Update) {
			mode = schema.CascadeUpdates
		}
	}
	return table, columns, mode
}

//--------------------------------------------------------------------

func (rq *Request) String() string {
	s := rq.Action
	switch rq.Action {
	case "drop", "create", "ensure", "alter":
		s += " " + rq.Table
	}
	switch rq.Action {
	case "create", "ensure":
		s += " " + rq.Schema.String()
	}
	switch rq.Action {
	case "rename":
		s += " " + rq.Renames[0].String()
	case "alter":
		s += " " + rq.SubAction
		switch rq.SubAction {
		case "create", "drop":
			s += " " + rq.Schema.String()
		case "rename":
			sep := " "
			for _, rn := range rq.Renames {
				s += sep + rn.String()
				sep = ", "
			}
		}
	}
	return s
}
