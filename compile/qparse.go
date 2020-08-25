// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/lexer"
	tok "github.com/apmckinlay/gsuneido/lexer/tokens"
	"github.com/apmckinlay/gsuneido/util/str"
)

type qparser struct {
	parserBase
}

func NewQueryParser(src string) *qparser {
	lxr := NewQueryLexer(src)
	p := &qparser{parserBase{lxr: lxr, Factory: ast.Builder{}}} //TODO query factory
	p.next()
	return p
}

func ParseRequest(src string) interface{} {
	p := NewQueryParser(src)
	result := p.request()
	if p.Token != tok.Eof {
		p.error("did not parse all input")
	}
	return result
}

func (p *qparser) request() interface{} {
	switch p.Token {
	case tok.Create:
		return p.create()
	//TODO: Ensure, Alter, Rename, View, Sview, Drop
	default:
		panic("invalid request")
	}
}

func (p *qparser) create() *Schema {
	p.next()
	table := p.matchIdent()
	columns, rules := p.columns()
	indexes := p.indexes(columns)
	return &Schema{Table: table, Columns: columns, Rules: rules, Indexes: indexes}
}

func (p *qparser) columns() (columns, rules []string) {
	p.match(tok.LParen)
	columns = make([]string, 0, 8)
	for p.Token != tok.RParen {
		if p.matchIf(tok.Sub) {
			columns = append(columns, "-")
		} else {
			col := p.matchIdent()
			if str.Capitalized(col) || strings.HasSuffix(col, "_lower!") {
				rules = append(rules, col)
			} else {
				columns = append(columns, col)
			}

		}
		p.matchIf(tok.Comma)
	}
	p.match(tok.RParen)
	return columns, rules
}

func (p *qparser) indexes(columns []string) []*Index {
	indexes := make([]*Index, 0, 4)
	for ix := p.index(columns); ix != nil; ix = p.index(columns) {
		indexes = append(indexes, ix)
	}
	return indexes
}

func (p *qparser) index(columns []string) *Index {
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
	ixcols := p.indexColumns(columns)
	if mode != 'k' && len(ixcols) == 0 {
		p.error("index columns must not be empty")
	}
	ix := &Index{Fields: ixcols, Mode: mode}
	ix.Fktable, ix.Fkcolumns, ix.Fkmode = p.foreignKey()
	return ix
}

func (p *qparser) indexColumns(columns []string) []int {
	p.match(tok.LParen)
	ixcols := make([]int, 0, 8)
	for p.Token != tok.RParen {
		col := p.matchIdent()
		c := str.List(columns).Index(col)
		if strings.HasSuffix(col, "_lower!") {
			if c = str.List(columns).Index(col[:len(col)-7]); c != -1 {
				c = -c - 2
			}
		}
		if c == -1 {
			p.error("invalid index column: " + col)
		}
		ixcols = append(ixcols, c)
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
	mode = BLOCK
	if p.matchIf(tok.Cascade) {
		mode = CASCADE
		if p.matchIf(tok.Update) {
			mode = CASCADE_UPDATES
		}
	}
	return table, columns, mode
}

// fkmode bits //TODO don't duplicate
const (
	BLOCK           = 0
	CASCADE_UPDATES = 1
	CASCADE_DELETES = 2
	CASCADE         = CASCADE_UPDATES | CASCADE_DELETES
)

type Schema struct {
	Table   string
	Columns []string
	Rules   []string
	Indexes []*Index
}

func (sc *Schema) String() string {
	var cb str.CommaBuilder
	for _, c := range sc.Columns {
		cb.Add(c)
	}
	for _, c := range sc.Rules {
		cb.Add(c)
	}
	s := sc.Table + " (" + cb.String() + ")"
	for _, ix := range sc.Indexes {
		s += " " + ix.string(sc.Columns)
	}
	return s
}

type Index struct {
	Fields []int
	// Mode is 'k' for key, 'i' for index, 'u' for unique index
	Mode      int
	Fktable   string
	Fkmode    int
	Fkcolumns []string
}

func (ix Index) string(columns []string) string {
	var cb str.CommaBuilder
	for _, c := range ix.Fields {
		if c < 0 {
			cb.Add(columns[-c-2] + "_lower!")
		} else {
			cb.Add(columns[c])
		}
	}
	s := map[int]string{'k': "key", 'i': "index", 'u': "index unique"}[ix.Mode] +
		"(" + cb.String() + ")"
	if ix.Fktable != "" {
		s += " in " + ix.Fktable
		if ix.Fkcolumns != nil {
			s += " " + str.Join("(,)", ix.Fkcolumns...)
		}
		if ix.Fkmode&CASCADE != 0 {
			s += " cascade"
			if ix.Fkmode == CASCADE_UPDATES {
				s += " update"
			}
		}
	}
	return s
}

func (ix Index) String() string {
	var cb str.CommaBuilder
	for _, c := range ix.Fields {
		cb.Add(strconv.Itoa(c))
	}
	return map[int]string{'k': "key", 'i': "index", 'u': "index unique"}[ix.Mode] +
		"(" + cb.String() + ")"
}
