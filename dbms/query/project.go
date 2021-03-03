// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Project struct {
	Query1
	columns  []string
	strategy projectStrategy
}

type projectStrategy int

const (
	projCopy projectStrategy = iota + 1
	projSeq
	projLookup
)

func (p *Project) Init() {
	p.Query1.Init()
	srcCols := p.source.Columns()
	if !sset.Subset(srcCols, p.columns) {
		panic("project: nonexistent column(s): " +
			str.Join(", ", sset.Difference(p.columns, srcCols)))
	}
	for _, col := range p.columns {
		if strings.HasSuffix(col, "_lower!") {
			panic("can't project _lower! fields")
		}
	}
	if hasKey(p.source.Keys(), p.columns) {
		p.strategy = projCopy
		p.includeDeps(srcCols)
	}
}

// hasKey returns true if there is a key containing cols
func hasKey(keys [][]string, cols []string) bool {
	for _, k := range keys {
		if sset.Subset(k, cols) {
			return true
		}
	}
	return false
}

func (p *Project) includeDeps(cols []string) {
	for _, f := range p.columns {
		deps := f + "_deps"
		if sset.Contains(cols, deps) {
			p.columns = sset.AddUnique(p.columns, deps)
		}
	}
}

func (p *Project) String() string {
	s := paren(p.source) + " PROJECT"
	switch p.strategy {
	case projSeq:
		s += "-SEQ"
	case projCopy:
		s += "-COPY"
	case projLookup:
		s += "-LOOKUP"
	}
	return s + " " + str.Join(", ", p.columns)
}

func (p *Project) Columns() []string {
	return p.columns
}

func (p *Project) Keys() [][]string {
	return projectKeys(p.source.Keys(), p.columns)
}

func projectKeys(keys [][]string, cols []string) [][]string {
	keys2 := projectIndexes(keys, cols)
	if len(keys2) == 0 {
		return [][]string{cols}
	}
	return keys2
}

func (p *Project) Indexes() [][]string {
	return projectIndexes(p.source.Indexes(), p.columns)
}

func projectIndexes(idxs [][]string, cols []string) [][]string {
	var idxs2 [][]string
	for _, k := range idxs {
		if sset.Subset(cols, k) {
			idxs2 = append(idxs2, k)
		}
	}
	return idxs2
}

func (p *Project) Transform() Query {
	moved := false
	for {
		if sset.Equal(p.columns, p.source.Columns()) {
			// remove projects of all columns
			return p.source.Transform()
		}
		switch q := p.source.(type) {
		case *Project:
			// combine projects
			p.columns = sset.Intersect(p.columns, q.columns)
			p.source = q.source
			continue
		case *Rename:
			return p.transformRename(q)
		case *Extend:
			if e := p.transformExtend(q); e != nil {
				return e
			}
		case *Times:
			p.splitOver(&q.Query2)
			moved = true
		case *Join:
			if sset.Subset(p.columns, q.by) {
				p.splitOver(&q.Query2)
				moved = true
			}
		case *LeftJoin:
			if sset.Subset(p.columns, q.by) {
				p.splitOver(&q.Query2)
				moved = true
			}
		case *Union:
			if p.splitOver2(&q.Compatible) {
				return p.source.Transform()
			}
		case *Intersect:
			if p.splitOver2(&q.Compatible) {
				return p.source.Transform()
			}
		}
		p.source = p.source.Transform()
		if moved {
			return p.source
		}
		return p
	}
}

func (p *Project) splitOver(q2 *Query2) {
	q2.source = &Project{Query1: Query1{source: q2.source},
		columns: sset.Intersect(p.columns, q2.source.Columns())}
	q2.source2 = &Project{Query1: Query1{source: q2.source2},
		columns: sset.Intersect(p.columns, q2.source2.Columns())}
}

func (p *Project) splitOver2(c *Compatible) bool {
	if c.disjoint != "" && !sset.Contains(p.columns, c.disjoint) {
		cols := append(sset.Copy(p.columns), c.disjoint)
		c.source = &Project{Query1: Query1{source: c.source},
			columns: sset.Intersect(cols, c.source.Columns())}
		c.source2 = &Project{Query1: Query1{source: c.source2},
			columns: sset.Intersect(cols, c.source2.Columns())}
		return false
	}
	c.source = &Project{Query1: Query1{source: c.source},
		columns: sset.Intersect(p.columns, c.source.Columns())}
	c.source2 = &Project{Query1: Query1{source: c.source2},
		columns: sset.Intersect(p.columns, c.source2.Columns())}
	return true
}

// transformRename moves projects before renames
func (p *Project) transformRename(r *Rename) Query {
	// remove renames not in project
	var newFrom, newTo []string
	from := r.from
	to := r.to
	for i := range to {
		if sset.Contains(p.columns, to[i]) {
			newFrom = append(newFrom, from[i])
			newTo = append(newTo, to[i])
		}
	}
	r.from = newFrom
	r.to = newTo

	// rename fields
	var newCols []string
	for _, col := range p.columns {
		if i := str.List(to).Index(col); i != -1 {
			newCols = append(newCols, from[i])
		} else {
			newCols = append(newCols, col)
		}
	}
	p.columns = newCols

	p.source = r.source
	r.source = p
	return r.Transform()
}

// transformExtend moves projects before extends
func (p *Project) transformExtend(e *Extend) Query {
	// remove portions of extend not included in project
	var newCols []string
	var newExprs []ast.Expr
	for i, col := range e.cols {
		if sset.Contains(p.columns, col) {
			newCols = append(newCols, col)
			newExprs = append(newExprs, e.exprs[i])
		}
	}
	origCols := e.cols
	e.cols = newCols
	origExprs := e.exprs
	e.exprs = newExprs

	// project must include all fields required by extend
	// there must be no rules left
	// since we don't know what fields are required by rules
	if !e.hasRules() {
		var exprCols []string
		for _, x := range e.exprs {
			exprCols = sset.Union(exprCols, x.Columns())
		}
		if sset.Subset(p.columns, exprCols) {
			// remove extend fields from project
			var newCols []string
			for _, col := range p.columns {
				if !sset.Contains(e.cols, col) {
					newCols = append(newCols, col)
				}
			}
			p.columns = newCols

			p.source = e.source
			e.source = p
			e.init()
			return e.Transform()
		}
	}
	e.cols = origCols
	e.exprs = origExprs
	return nil
}

func (p *Project) Fixed() []Fixed {
	//TODO cache like extend and union ???
	var fixed []Fixed
	for _, f := range p.source.Fixed() {
		if sset.Contains(p.columns, f.col) {
			fixed = append(fixed, f)
		}
	}
	return fixed
}

func (p *Project) Updateable() bool {
		return p.strategy == projCopy && p.source.Updateable()
}
