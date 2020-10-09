// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/lexer"
	tok "github.com/apmckinlay/gsuneido/lexer/tokens"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Schema = schema.Schema
type Index = schema.Index

type qparser struct {
	parserBase
}

func NewQueryParser(src string) *qparser {
	lxr := lexer.NewQueryLexer(src)
	p := &qparser{parserBase{lxr: lxr, Factory: ast.Builder{}}} //TODO query factory
	p.next()
	return p
}

type Request struct {
	Action    string
	SubAction string
	Schema
	Renames []Rename
}

type Rename struct {
	From string
	To   string
}

func ParseRequest(src string) *Request {
	p := NewQueryParser(src)
	result := p.request()
	if p.Token != tok.Eof {
		p.error("did not parse all input")
	}
	return result
}

func (p *qparser) request() *Request {
	switch {
	case p.matchIf(tok.Create):
		return &Request{Action: "create", Schema: p.schema(true)}
	case p.matchIf(tok.Ensure):
		return &Request{Action: "ensure", Schema: p.schema(false)}
	case p.matchIf(tok.Drop):
		table := p.matchIdent()
		return &Request{Action: "drop", Schema: Schema{Table: table}}
	case p.matchIf(tok.Rename):
		rename := p.rename()
		return &Request{Action: "rename", Renames: []Rename{rename}}
	case p.matchIf(tok.Alter):
		return p.alter()
	//TODO: View, Sview
	default:
		panic("invalid request")
	}
}

func (p *qparser) rename() Rename {
	from := p.matchIdent()
	p.match(tok.To)
	to := p.matchIdent()
	return Rename{From: from, To: to}
}

func (p *qparser) alter() *Request {
	table := p.matchIdent()
	switch {
	case p.matchIf(tok.Create):
		return &Request{Action: "alter", SubAction: "create",
			Schema: p.schema2(table, false)}
	case p.matchIf(tok.Drop):
		return &Request{Action: "alter", SubAction: "drop",
			Schema: p.schema2(table, false)}
	case p.matchIf(tok.Rename):
		return &Request{Action: "alter", SubAction: "rename",
			Schema: p.schema2(table, false), Renames: p.renames()}
	default:
		panic("invalid request")
	}
}

func (p *qparser) renames() []Rename {
	var renames []Rename
	for {
		renames = append(renames, p.rename())
		if !p.matchIf(tok.Comma) {
			return renames
		}
	}
}

func (p *qparser) schema(full bool) Schema {
	table := p.matchIdent()
	return p.schema2(table, full)
}

func (p *qparser) schema2(table string, full bool) Schema {
	columns, derived := p.columns(full)
	indexes := p.indexes(columns, derived, full)
	return Schema{Table: table, Columns: columns, Derived: derived, Indexes: indexes}
}

func (p *qparser) columns(full bool) (columns, derived []string) {
	if !full && p.Token != tok.LParen {
		return
	}
	p.match(tok.LParen)
	columns = make([]string, 0, 8)
	for p.Token != tok.RParen {
		if p.matchIf(tok.Sub) {
			columns = append(columns, "-")
		} else {
			col := p.matchIdent()
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
		p.matchIf(tok.Comma)
	}
	p.match(tok.RParen)
	return columns, derived
}

func (p *qparser) indexes(columns, derived []string, full bool) []Index {
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

func (p *qparser) index(columns, derived []string, full bool) *Index {
	if p.Token != tok.Key && p.Token != tok.Index {
		return nil
	}
	mode := int('i')
	if p.Token == tok.Key {
		mode = 'k'
	}
	p.next()
	if mode != 'k' && p.matchIf(tok.Unique) {
		mode = 'u'
	}
	ixcols := p.indexColumns(columns, derived, full)
	if mode != 'k' && len(ixcols) == 0 {
		p.error("index columns must not be empty")
	}
	ix := &Index{Columns: ixcols, Mode: mode}
	ix.Fktable, ix.Fkcolumns, ix.Fkmode = p.foreignKey()
	return ix
}

func (p *qparser) indexColumns(columns, derived []string, full bool) []string {
	p.match(tok.LParen)
	ixcols := make([]string, 0, 8)
	for p.Token != tok.RParen {
		col := p.matchIdent()
		if full && !str.List(columns).Has(col) &&
			(!strings.HasSuffix(col, "_lower!") || !str.List(derived).Has(col)) {
			p.error("invalid index column: " + col)
		}
		ixcols = append(ixcols, col)
		p.matchIf(tok.Comma)
	}
	p.match(tok.RParen)
	return ixcols
}

func (p *qparser) foreignKey() (table string, columns []string, mode int) {
	if !p.matchIf(tok.In) {
		return
	}
	table = p.matchIdent()
	if p.matchIf(tok.LParen) {
		for p.Token != tok.RParen {
			columns = append(columns, p.Text)
			p.matchIdent()
			p.matchIf(tok.Comma)
		}
		p.next()
	}
	mode = schema.Block
	if p.matchIf(tok.Cascade) {
		mode = schema.Cascade
		if p.matchIf(tok.Update) {
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

func (rn Rename) String() string {
	return rn.From + " to " + rn.To
}
