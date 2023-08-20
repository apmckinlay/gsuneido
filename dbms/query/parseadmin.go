// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"slices"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/compile/lexer"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/str"
)

type adminParser struct {
	compile.ParserBase
}

func NewAdminParser(src string) *adminParser {
	lxr := lexer.NewQueryLexer(src)
	p := &adminParser{compile.ParserBase{Lxr: lxr}}
	p.Next()
	return p
}

type Schema = schema.Schema
type Index = schema.Index

type Admin interface {
	execute(*db19.Database, *Sviews)
	String() string
}

type Renames struct {
	From []string
	To   []string
}

func ParseAdmin(src string) Admin {
	p := NewAdminParser(src)
	result := p.admin()
	if p.Token != tok.Eof {
		p.Error("did not parse all input")
	}
	return result
}

func (p *adminParser) admin() Admin {
	switch {
	case p.MatchIf(tok.Create):
		return &createAdmin{p.schema()}
	case p.MatchIf(tok.Ensure):
		return &ensureAdmin{p.schema()}
	case p.MatchIf(tok.Rename):
		from, to := p.rename1()
		return &renameAdmin{from: from, to: to}
	case p.MatchIf(tok.Alter):
		return p.alter()
	case p.MatchIf(tok.View):
		return p.view()
	case p.MatchIf(tok.Sview):
		return p.sview()
	case p.MatchIf(tok.Drop):
		table := p.MatchIdent()
		return &dropAdmin{table}
	default:
		panic("invalid admin")
	}
}

func (p *adminParser) alter() Admin {
	table := p.MatchIdent()
	switch {
	case p.MatchIf(tok.Create):
		return &alterCreateAdmin{p.schema2(table)}
	case p.MatchIf(tok.Drop):
		return &alterDropAdmin{p.schema2(table)}
	case p.MatchIf(tok.Rename):
		return p.alterRename(table)
	default:
		panic("invalid admin")
	}
}

func (p *adminParser) alterRename(table string) Admin {
	var from, to []string
	for {
		f, t := p.rename1()
		from = append(from, f)
		to = append(to, t)
		if !p.MatchIf(tok.Comma) {
			break
		}
	}
	return &alterRenameAdmin{table: table, from: from, to: to}
}

func (p *adminParser) rename1() (string, string) {
	from := p.MatchIdent()
	p.Match(tok.To)
	to := p.MatchIdent()
	return from, to
}

func (p *adminParser) Schema() Schema {
	result := p.schema()
	if p.Token != tok.Eof {
		p.Error("did not parse all input")
	}
	return result
}

func (p *adminParser) schema() Schema {
	table := p.MatchIdent()
	return p.schema2(table)
}

func (p *adminParser) schema2(table string) Schema {
	columns, derived := p.columns()
	indexes := p.indexes()
	return Schema{Table: table, Columns: columns, Derived: derived, Indexes: indexes}
}

func (p *adminParser) columns() (columns, derived []string) {
	if !p.MatchIf(tok.LParen) {
		return
	}
	columns = make([]string, 0, 8)
	for p.Token != tok.RParen {
		if p.MatchIf(tok.Sub) {
			columns = append(columns, "-")
		} else {
			col := p.MatchIdent()
			if str.Capitalized(col) {
				derived = p.appendUnique(derived, col)
			} else if strings.HasSuffix(col, "_lower!") {
				derived = p.appendUnique(derived, col)
			} else {
				columns = p.appendUnique(columns, col)
			}

		}
		p.MatchIf(tok.Comma)
	}
	p.Match(tok.RParen)
	return columns, derived
}

func (p *adminParser) appendUnique(columns []string, col string) []string {
	if slices.Contains(columns, col) {
		p.Error("duplicate column:", col)
	}
	return append(columns, col)
}

func (p *adminParser) indexes() []Index {
	indexes := make([]Index, 0, 4)
	for ix := p.index(); ix != nil; ix = p.index() {
		indexes = append(indexes, *ix)
	}
	return indexes
}

func (p *adminParser) index() *Index {
	mode := p.indexMode()
	if mode == 0 {
		return nil
	}
	ixcols := p.indexColumns()
	ix := &Index{Columns: ixcols, Mode: mode}
	ix.Fk.Table, ix.Fk.Columns, ix.Fk.Mode = p.foreignKey()
	if ix.Fk.Columns == nil {
		ix.Fk.Columns = ixcols
	}
	return ix
}

func (p *adminParser) indexMode() byte {
	switch {
	case p.MatchIf(tok.Key):
		return 'k'
	case p.MatchIf(tok.Index):
		if p.MatchIf(tok.Unique) {
			return 'u'
		}
		return 'i'
	case p.MatchIf(tok.Unique):
		p.Match(tok.Index)
		return 'u'
	}
	return 0
}

func (p *adminParser) indexColumns() []string {
	p.Match(tok.LParen)
	ixcols := make([]string, 0, 8)
	for p.Token != tok.RParen {
		col := p.MatchIdent()
		ixcols = append(ixcols, col)
		p.MatchIf(tok.Comma)
	}
	p.Match(tok.RParen)
	return ixcols
}

func (p *adminParser) foreignKey() (table string, columns []string, mode byte) {
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

func (p *adminParser) view() Admin {
	name := p.viewName()
	def := strings.TrimSpace(p.Lxr.Remainder())
	return &viewAdmin{name: name, def: def}
}

func (p *adminParser) sview() Admin {
	name := p.viewName()
	def := strings.TrimSpace(p.Lxr.Remainder())
	return &sviewAdmin{name: name, def: def}
}

func (p *adminParser) viewName() string {
	name := p.MatchIdent()
	p.MustMatch(tok.Eq)
	p.Token = tok.Eof
	return name
}
