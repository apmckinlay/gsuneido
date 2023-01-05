// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
)

type Project struct {
	Query1
	columns []string
	unique  bool
	projectApproach
	rewound bool
	prevRow Row
	curRow  Row
	srcHdr  *Header
	projHdr *Header
	indexed bool
	results map[string]Row
	st      *SuTran
}

type projectApproach struct {
	strategy projectStrategy
	index    []string
}

type projectStrategy int

const (
	// projCopy is when the columns contain a key, so it's just pass through
	projCopy projectStrategy = iota + 1
	// projSeq orders by the columns so duplicates are consecutive
	projSeq
	// projHash builds a temp hash index to identify duplicates
	projHash
)

func NewProject(src Query, cols []string) *Project {
	cols = set.Unique(cols)
	srcCols := src.Columns()
	if !set.Subset(srcCols, cols) {
		panic("project: nonexistent column(s): " +
			str.Join(", ", set.Difference(cols, srcCols)))
	}
	for _, col := range cols {
		if strings.HasSuffix(col, "_lower!") {
			panic("can't project _lower! fields")
		}
	}
	p := &Project{Query1: Query1{source: src}, columns: cols, rewound: true}
	if hasKey(p.source.Keys(), cols, p.source.Fixed()) {
		p.unique = true
		p.includeDeps(srcCols)
	}
	return p
}

func NewRemove(src Query, cols []string) *Project {
	cols = set.Difference(src.Columns(), cols)
	if len(cols) == 0 {
		panic("remove: can't remove all columns")
	}
	cols = slc.WithoutFn(cols,
		func(s string) bool { return strings.HasSuffix(s, "_lower!") })
	return NewProject(src, cols)
}

// hasKey returns whether cols contains a key
// taking fixed into consideration
func hasKey(keys [][]string, cols []string, fixed []Fixed) bool {
outer:
	for _, key := range keys {
		for _, k := range key {
			if !isFixed(fixed, k) && !slices.Contains(cols, k) {
				continue outer
			}
		}
		return true
	}
	return false
}

func (p *Project) includeDeps(cols []string) {
	for _, f := range p.columns {
		deps := f + "_deps"
		if slices.Contains(cols, deps) {
			p.columns = set.AddUnique(p.columns, deps)
		}
	}
}

func (p *Project) String() string {
	return parenQ2(p.source) + " " + p.stringOp()
}

func (p *Project) stringOp() string {
	s := "PROJECT"
	switch p.strategy {
	case projSeq:
		s += "-SEQ"
	case projCopy:
		s += "-COPY"
	case projHash:
		s += "-HASH"
	}
	return s + " " + str.Join(",", p.columns)
}

func (p *Project) SetTran(t QueryTran) {
	p.st = MakeSuTran(t)
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
		if set.Subset(cols, k) {
			idxs2 = append(idxs2, k)
		}
	}
	return idxs2
}

func (p *Project) Nrows() (int, int) {
	nr, pop := p.source.Nrows()
	if !p.unique {
		nr /= 2 // ???
	}
	return nr, pop
}

func (p *Project) Transform() Query {
	p.source = p.source.Transform()
	if set.Equal(p.columns, p.source.Columns()) {
		// remove projects of all columns
		return p.source
	}
	switch q := p.source.(type) {
	case *Project:
		// combine projects
		p.columns = set.Intersect(p.columns, q.columns)
		p.source = q.source
	case *Rename:
		return p.transformRename(q)
	case *Extend:
		return p.transformExtend(q)
	case *Times:
		return p.splitOver(&q.Query2)
	case *Join:
		if set.Subset(p.columns, q.by) {
			return p.splitOver(&q.Query2)
		}
	case *LeftJoin:
		if set.Subset(p.columns, q.by) {
			return p.splitOver(&q.Query2)
		}
	case *Union:
		return p.splitOver2(&q.Compatible)
	case *Intersect:
		return p.splitOver2(&q.Compatible)
	case *Nothing:
		// propagate Nothing
		return NewNothing(p.Columns())
	}
	return p
}

func (p *Project) splitOver(q2 *Query2) Query {
	q2.source1 = NewProject(q2.source1,
		set.Intersect(p.columns, q2.source1.Columns()))
	q2.source2 = NewProject(q2.source2,
		set.Intersect(p.columns, q2.source2.Columns()))
	return p.source
}

func (p *Project) splitOver2(c *Compatible) Query {
	if c.disjoint != "" && !slices.Contains(p.columns, c.disjoint) {
		// don't split if project doesn't include disjoint
		return p
	}
	c.source1 = NewProject(c.source1,
		set.Intersect(p.columns, c.source1.Columns()))
	c.source2 = NewProject(c.source2,
		set.Intersect(p.columns, c.source2.Columns()))
	return p.source.Transform()
}

// transformRename moves projects before renames
func (p *Project) transformRename(r *Rename) Query {
	// remove renames not in project
	var newFrom, newTo []string
	from := r.from
	to := r.to
	for i := range to {
		ck := to[i]
		if p.unique {
			ck = strings.TrimSuffix(to[i], "_deps")
		}
		if slices.Contains(p.columns, ck) {
			newFrom = append(newFrom, from[i])
			newTo = append(newTo, to[i])
		}
	}
	r.from = newFrom
	r.to = newTo

	// rename fields
	var newCols []string
	for _, col := range p.columns {
		if i := slices.Index(to, col); i != -1 {
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

// transformExtend tries to move projects before extends.
// If it doesn't return nil, it must call transform on the source.
func (p *Project) transformExtend(e *Extend) Query {
	if e.hasRules() {
		// rules make it too hard to determine what fields they use
		return p
	}
	// orig := Format(p)
	extendDeps := exprsCols(e.exprs)
	// split the extend into what can go after the project,
	// and what has to stay before the project
	var beforeCols, afterCols []string
	var beforeExprs, afterExprs []ast.Expr
	newProjCols := p.columns
	for i, col := range e.cols {
		// if col is a dependency then we can't move it
		if !slices.Contains(extendDeps, col) {
			if slices.Contains(p.columns, col) {
				if isConstant(e.exprs[i]) {
					newProjCols = slc.Without(newProjCols, col)
					afterCols = append(afterCols, col)
					afterExprs = append(afterExprs, e.exprs[i])
					continue
				}
			} else {
				// not in the project so remove it
				continue
			}
		}
		// else keep in before
		beforeCols = append(beforeCols, col)
		beforeExprs = append(beforeExprs, e.exprs[i])
	}
	if len(newProjCols) == 0 {
		// the before extend is irrelevant with ProjectNone
		return NewExtend(&ProjectNone{}, afterCols, afterExprs)
	}
	p.columns = newProjCols
	var result Query
	if len(beforeCols) > 0 {
		p.source = NewExtend(e.source, beforeCols, beforeExprs)
		result = p
	} else {
		p.source = e.source // drop original extend since no columns left
		result = p.Transform()
	}
	if len(afterCols) > 0 {
		result = NewExtend(result, afterCols, afterExprs)
	}
	return result
}

func exprsCols(exprs []ast.Expr) []string {
	var cols []string
	for _, x := range exprs {
		cols = set.Union(cols, x.Columns())
	}
	return cols
}

func isConstant(e ast.Expr) bool {
	_, ok := e.(*ast.Constant)
	return ok
}

func (p *Project) Fixed() []Fixed {
	//TODO cache like extend and union ???
	var fixed []Fixed
	for _, f := range p.source.Fixed() {
		if slices.Contains(p.columns, f.col) {
			fixed = append(fixed, f)
		}
	}
	return fixed
}

func (p *Project) Updateable() string {
	if p.strategy == projCopy {
		return p.source.Updateable()
	}
	return ""
}

// optimize ---------------------------------------------------------

func (p *Project) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	if p.unique {
		approach := &projectApproach{strategy: projCopy, index: index}
		fixcost, varcost := Optimize(p.source, mode, index, frac)
		return fixcost, varcost, approach
	}
	seq := bestGrouped(p.source, mode, index, frac, p.columns)
	fixcostHash, varcostHash := p.hashCost(mode, index, frac)
	if fixcostHash+varcostHash < seq.cost() {
		return fixcostHash, varcostHash,
			&projectApproach{strategy: projHash, index: index}
	}
	return seq.fixcost, seq.varcost,
		&projectApproach{strategy: projSeq, index: seq.index}
}

const hashLimit = 10000

func (p *Project) hashCost(mode Mode, index []string, frac float64) (Cost, Cost) {
	nrows, _ := p.Nrows()
	if mode != ReadMode || nrows > hashLimit {
		return impossible, impossible
	}
	// assume we're reading Next (normal)
	fixcost, varcost := Optimize(p.source, mode, index, frac)
	hashCost := nrows * 20 // ???
	return fixcost, varcost + hashCost
}

func (p *Project) setApproach(_ []string, frac float64, approach any, tran QueryTran) {
	p.projectApproach = *approach.(*projectApproach)
	p.source = SetApproach(p.source, p.index, frac, tran)
	p.projHdr = p.Header() // cache for Get
	p.srcHdr = p.source.Header()
}

// execution --------------------------------------------------------

func (p *Project) Header() *Header {
	if p.projHdr != nil {
		return p.projHdr
	}
	srcFlds := p.source.Header().Fields
	newflds := make([][]string, len(srcFlds))
	for i, fs := range srcFlds {
		newflds[i] = projectFields(fs, p.columns)
	}
	return NewHeader(newflds, p.columns)
}

func projectFields(fs []string, pcols []string) []string {
	flds := make([]string, len(fs))
	for i, f := range fs {
		if slices.Contains(pcols, f) {
			flds[i] = f
		} else {
			flds[i] = "-"
		}
	}
	return flds
}

func (p *Project) Rewind() {
	p.rewound = true
	p.source.Rewind()
}

func (p *Project) Get(th *Thread, dir Dir) Row {
	switch p.strategy {
	case projCopy:
		return p.source.Get(th, dir)
	case projSeq:
		return p.getSeq(th, dir)
	case projHash:
		return p.getHash(th, dir)
	}
	panic("should not reach here")
}

func (p *Project) getSeq(th *Thread, dir Dir) Row {
	if dir == Next {
		// output the first of each group
		// i.e. skip over rows the same as previous output
		for {
			row := p.source.Get(th, dir)
			if row == nil {
				return nil
			}
			if p.rewound || !p.projHdr.EqualRows(row, p.curRow, th, p.st) {
				p.rewound = false
				p.prevRow = p.curRow
				p.curRow = row
				return row
			}
		}
	} else { // Prev
		// output the last of each group
		// i.e. output when next record is different
		// (to get the same records as NEXT)
		if p.rewound {
			p.prevRow = p.source.Get(th, dir)
		}
		p.rewound = false
		for {
			if p.prevRow == nil {
				return nil
			}
			row := p.prevRow
			p.prevRow = p.source.Get(th, dir)
			if p.prevRow == nil ||
				!p.projHdr.EqualRows(row, p.prevRow, th, p.st) {
				// output the last row of a group
				p.curRow = row
				return row
			}
		}
	}
}

func (p *Project) getHash(th *Thread, dir Dir) Row {
	//TODO limit size of results
	if p.rewound {
		p.rewound = false
		if p.results == nil {
			p.results = make(map[string]Row)
		}
		if dir == Prev && !p.indexed {
			p.buildHash(th)
		}
	}
	for {
		row := p.source.Get(th, dir)
		if row == nil {
			break
		}
		key := ixkey.Make(row, p.srcHdr, p.columns, th, p.st)
		result, ok := p.results[key]
		if !ok {
			p.results[key] = row
			return row
		}
		if row.SameAs(result) {
			return row
		}
	}
	if dir == Next {
		p.indexed = true
	}
	return nil
}

func (p *Project) buildHash(th *Thread) {
	for {
		row := p.source.Get(th, Next)
		if row == nil {
			break
		}
		key := ixkey.Make(row, p.srcHdr, p.columns, th, p.st)
		if _, ok := p.results[key]; !ok {
			p.results[key] = row
		}
	}
	p.source.Rewind()
	p.indexed = true
}

func (p *Project) Output(th *Thread, rec Record) {
	if p.strategy != projCopy {
		panic("can't output to a project that doesn't include a key")
	}
	p.source.Output(th, rec)
}

func (p *Project) Select(cols, vals []string) {
	p.source.Select(cols, vals)
	if p.strategy == projHash {
		p.indexed = false
	}
	p.rewound = true
}

func (p *Project) Lookup(th *Thread, cols, vals []string) Row {
	if p.strategy == projCopy {
		return p.source.Lookup(th, cols, vals)
	}
	p.Select(cols, vals)
	row := p.Get(th, Next)
	p.Select(nil, nil) // clear select
	return row
}
