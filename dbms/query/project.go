// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"slices"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/hmap"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Project struct {
	results *mapType
	st      *SuTran
	columns []string
	remove  []string
	prevRow Row
	curRow  Row
	projectApproach
	Query1
	unique        bool
	rewound       bool
	indexed       bool
	warned        bool
	derivedWarned bool
	derived       int
}

type mapType = hmap.Hmap[rowHash, struct{}, hmap.Funcs[rowHash]]

type projectApproach struct {
	index    []string
	strategy projectStrategy
}

type projectStrategy int

const (
	// projCopy is when the columns contain a key, so it's just pass through
	projCopy projectStrategy = iota + 1
	// projSeq orders by the columns so duplicates are consecutive
	projSeq
	// projMap builds a map to identify duplicates
	projMap
)

func NewProject(src Query, cols []string) *Project {
	assert.That(len(cols) > 0)
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
	return newProject2(src, cols, true)
}

func NewRemove(src Query, cols []string) *Project {
	proj := slc.WithoutFn(src.Columns(), func(col string) bool {
		return slices.Contains(cols, col) ||
			strings.HasSuffix(col, "_lower!") ||
			(strings.HasSuffix(col, "_deps") &&
				slices.Contains(cols, strings.TrimSuffix(col, "_deps")))
	})
	if len(proj) == 0 {
		panic("remove: can't remove all columns")
	}
	p := newProject(src, proj)
	p.remove = cols
	return p
}

// newProject is common to NewProject and NewRemove
func newProject(src Query, cols []string) *Project {
	return newProject2(src, cols, false)
}
func newProject2(src Query, cols []string, includeDeps bool) *Project {
	p := &Project{Query1: Query1{source: src}, rewound: true}
	if hasKey(cols, src.Keys(), src.Fixed()) {
		p.unique = true
		if includeDeps {
			cols = p.includeDeps(cols, src.Columns())
		}
	}
	p.columns = cols
	p.header = p.getHeader()
	p.keys = projectKeys(src.Keys(), p.columns)
	p.indexes = projectIndexes(src.Indexes(), p.columns)
	p.setNrows(p.getNrows())
	p.rowSiz.Set(src.rowSize())
	p.fast1.Set(src.fastSingle())
	p.singleTbl.Set(src.SingleTable())
	p.lookCost.Set(p.getLookupCost())
	return p
}

func (*Project) includeDeps(cols, srcCols []string) []string {
	newCols := cols
	for _, f := range cols {
		deps := f + "_deps"
		if slices.Contains(srcCols, deps) {
			newCols = set.AddUnique(newCols, deps)
		}
	}
	return newCols
}

func (p *Project) getHeader() *Header {
	srcFlds := p.source.Header().Fields
	newflds := make([][]string, len(srcFlds))
	for i, fs := range srcFlds {
		newflds[i] = projectFields(fs, p.columns)
	}
	return NewHeader(newflds, p.columns)
}

// hasKey returns whether cols contains a key
// taking fixed into consideration
func hasKey(cols []string, keys [][]string, fixed []Fixed) bool {
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
	case projMap:
		s += "-MAP"
	}
	return s + " " + str.Join(",", p.columns)
}

func (p *Project) format() string {
	warn := ""
	if !p.unique {
		warn = "/*NOT UNIQUE*/ "
	}
	if p.remove != nil {
		return "remove " + warn + str.Join(", ", p.remove)
	}
	return "project " + warn + str.Join(", ", p.columns)
}

func (p *Project) SetTran(t QueryTran) {
	p.st = MakeSuTran(t)
}

// projectKeys is also used by Summarize
func projectKeys(keys [][]string, cols []string) [][]string {
	var keys2 [][]string
	for _, ix := range keys {
		if set.Subset(cols, ix) {
			keys2 = append(keys2, ix)
		}
	}
	if len(keys2) == 0 {
		return [][]string{cols} // fallback on all columns
	}
	return keys2
}

// projectIndexes is also used by Summarize
func projectIndexes(idxs [][]string, cols []string) [][]string {
	var idxs2 [][]string
	for _, ix := range idxs {
		// get the prefix of the index that is in cols
		i := 0
		for ; i < len(ix) && slices.Contains(cols, ix[i]); i++ {
		}
		pre := ix[:i]
		if i > 0 && !slc.ContainsFn(idxs2, pre, slices.Equal) {
			idxs2 = append(idxs2, pre)
		}
	}
	return idxs2
}

func (p *Project) getNrows() (int, int) {
	nr, pop := p.source.Nrows()
	if !p.unique {
		nr /= 2 // ??? (matches lookupCost)
	}
	return nr, pop
}

func (p *Project) Transform() Query {
	src := p.source.Transform()
	if _, ok := src.(*Nothing); ok {
		return NewNothing(p.columns)
	}
	if set.Equal(p.columns, p.source.Columns()) {
		// remove projects of all columns
		return src
	}
	switch q := src.(type) {
	case *Project:
		// combine projects by removing all but the first
		return newProject(q.source, p.columns).Transform()
	case *Summarize:
		cols := make([]string, 0, len(q.cols))
		ops := make([]string, 0, len(q.ops))
		ons := make([]string, 0, len(q.ons))
		for i, col := range q.cols {
			if slices.Contains(p.columns, col) {
				cols = append(cols, col)
				ops = append(ops, q.ops[i])
				ons = append(ons, q.ons[i])
			}
		}
		if len(cols) == 0 { // no summaries left
			return newProject(q.source, p.columns).Transform()
		}
		if set.Subset(p.columns, q.by) {
			return NewSummarize(q.source, q.hint, q.by, cols, ops, ons).Transform()
		}
	case *Rename:
		return p.transformRename(q)
	case *Extend:
		return p.transformExtend(q)
	case *Times:
		return NewTimes(p.splitOver(&q.Query2)).Transform()
	case *Join:
		if set.Subset(p.columns, q.by) {
			src1, src2 := p.splitOver(&q.Query2)
			return NewJoin(src1, src2, q.by).Transform()
		}
	case *LeftJoin:
		if set.Subset(p.columns, q.by) {
			src1, src2 := p.splitOver(&q.Query2)
			return NewLeftJoin(src1, src2, q.by).Transform()
		}
	case *Union:
		if p.splitable(&q.Compatible) {
			return NewUnion(p.splitOver(&q.Query2)).Transform()
		}
	case *Intersect:
		if p.splitable(&q.Compatible) {
			return NewIntersect(p.splitOver(&q.Query2)).Transform()
		}
	}
	return p.transform(src)
}

func (p *Project) splitOver(q2 *Query2) (Query, Query) {
	src1 := newProject(q2.source1,
		set.Intersect(p.columns, q2.source1.Columns()))
	src2 := newProject(q2.source2,
		set.Intersect(p.columns, q2.source2.Columns()))
	return src1, src2
}

func (p *Project) splitable(c *Compatible) bool {
	// don't split if project doesn't include disjoint
	return c.disjoint == "" || slices.Contains(p.columns, c.disjoint)
}

// transformRename moves projects before renames
func (p *Project) transformRename(r *Rename) Query {
	// remove renames not in project
	var newFrom, newTo []string
	from := r.from
	to := r.to
	for i := len(to) - 1; i >= 0; i-- {
		ck := to[i]
		if p.unique {
			ck = strings.TrimSuffix(to[i], "_deps")
		}
		if slices.Contains(p.columns, ck) || slices.Contains(newFrom, ck) {
			newFrom = append(newFrom, from[i])
			newTo = append(newTo, to[i])
		}
	}
	slices.Reverse(newFrom)
	slices.Reverse(newTo)
	newProj := r.renameRev(p.columns)
	p = newProject(r.source, newProj)
	r = NewRename(p, newFrom, newTo)
	return r.Transform()
}

// transformExtend tries to move projects before extends.
func (p *Project) transformExtend(e *Extend) Query {
	if e.hasRules() {
		// rules make it too hard to determine what fields they use
		return p.transform(e)
	}
	extendUses := exprsCols(e.exprs)
	// split the extend into what can go after the project,
	// and what has to stay before the project
	var beforeCols, afterCols []string
	var beforeExprs, afterExprs []ast.Expr
	newProjCols := p.columns
	for i, col := range e.cols {
		// if col is a dependency then we can't move it
		if !slices.Contains(extendUses, col) {
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
	if slices.Equal(beforeCols, e.cols) {
		return p.transform(e)
	}
	var result Query
	if len(beforeCols) > 0 {
		q := NewExtend(e.source, beforeCols, beforeExprs)
		result = newProject(q, newProjCols)
	} else {
		// drop original extend since no columns left
		result = newProject(e.source, newProjCols)
	}
	if len(afterCols) > 0 {
		result = NewExtend(result, afterCols, afterExprs)
	}
	return result.Transform()
}

func (p *Project) transform(src Query) Query {
	if src != p.source {
		return newProject(src, p.columns)
	}
	return p
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
	if p.fixed == nil {
		p.fixed = projectFixed(p.source.Fixed(), p.columns)
		assert.That(p.fixed != nil)
	}
	return p.fixed
}

// projectFixed is also used by Summarize
func projectFixed(srcFixed []Fixed, cols []string) []Fixed {
	fixed := []Fixed{}
	for _, f := range srcFixed {
		if slices.Contains(cols, f.col) {
			fixed = append(fixed, f)
		}
	}
	return fixed
}

func (p *Project) Updateable() string {
	if p.unique {
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
	fixcostMap, varcostMap := p.mapCost(mode, index, frac)
	if fixcostMap+varcostMap < seq.cost() {
		return fixcostMap, varcostMap,
			&projectApproach{strategy: projMap, index: index}
	}
	return seq.fixcost, seq.varcost,
		&projectApproach{strategy: projSeq, index: seq.index}
}

const mapLimit = 16384 // mapLimit is also used by Summarize

func (p *Project) mapCost(mode Mode, index []string, frac float64) (Cost, Cost) {
	nrows, _ := p.Nrows()
	if mode != ReadMode || nrows > mapLimit {
		return impossible, impossible
	}
	// assume we're reading Next (normal)
	fixcost, varcost := Optimize(p.source, mode, index, frac)
	mapCost := nrows * 20 // ???
	return fixcost, varcost + mapCost
}

func (p *Project) setApproach(_ []string, frac float64, approach any, tran QueryTran) {
	p.projectApproach = *approach.(*projectApproach)
	p.source = SetApproach(p.source, p.index, frac, tran)
	p.header = p.getHeader()
}

// execution --------------------------------------------------------

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
	case projMap:
		return p.getMap(th, dir)
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
			if p.rewound || !p.header.EqualRows(row, p.curRow, th, p.st) {
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
				!p.header.EqualRows(row, p.prevRow, th, p.st) {
				// output the last row of a group
				p.curRow = row
				return row
			}
		}
	}
}

type rowHash struct {
	row  Row
	hash uint32
}

func (p *Project) getMap(th *Thread, dir Dir) Row {
	if p.rewound {
		p.rewound = false
		if p.results == nil {
			hfn := func(k rowHash) uint32 { return k.hash }
			eqfn := func(x, y rowHash) bool {
				return x.hash == y.hash &&
					equalCols(x.row, y.row, p.source.Header(), p.columns, th, p.st)
			}
			p.results = hmap.NewHmapFuncs[rowHash, struct{}](hfn, eqfn)
		}
		if dir == Prev && !p.indexed {
			p.buildMap(th)
		}
	}
	for {
		row := p.source.Get(th, dir)
		if row == nil {
			break
		}
		oldRow, existed := p.addResult(th, row)
		if !existed || row.SameAs(oldRow) {
			return row
		}
	}
	if dir == Next {
		p.indexed = true
	}
	return nil
}

func hashCols(row Row, hdr *Header, cols []string, th *Thread, st *SuTran) uint32 {
	h := uint32(31)
	for _, col := range cols {
		x := row.GetRawVal(hdr, col, th, st)
		h = 31*h + hash.String(x)
	}
	return h
}
func equalCols(x, y Row, hdr *Header, cols []string, th *Thread, st *SuTran) bool {
	for _, col := range cols {
		if x.GetRawVal(hdr, col, th, st) != y.GetRawVal(hdr, col, th, st) {
			return false
		}
	}
	return true
}

func (p *Project) buildMap(th *Thread) {
	for {
		row := p.source.Get(th, Next)
		if row == nil {
			break
		}
		p.addResult(th, row)
	}
	p.source.Rewind()
	p.indexed = true
}

// addResult returns the old row and true if it already existed,
// else the new row and false
func (p *Project) addResult(th *Thread, row Row) (Row, bool) {
	rh := rowHash{row: row,
		hash: hashCols(row, p.source.Header(), p.columns, th, p.st)}
	k, _, existed := p.results.GetPut(rh, struct{}{})
	if existed {
		return k.row, true
	} else {
		if !p.warned && p.results.Size() > mapLimit {
			p.warned = true
			Warning("project-map large >", mapLimit)
		}
		if !p.derivedWarned && p.derived > derivedWarn {
			p.derivedWarned = true
			Warning("project-map derived large >",
				derivedWarn, "average", p.derived/p.results.Size())
		}
		return row, false
	}
}

func (p *Project) Output(th *Thread, rec Record) {
	if p.strategy != projCopy {
		panic("can't output to a project that doesn't include a key")
	}
	p.source.Output(th, rec)
}

func (p *Project) Select(cols, vals []string) {
	p.source.Select(cols, vals)
	if p.strategy == projMap {
		p.indexed = false
	}
	p.rewound = true
}

func (p *Project) Lookup(th *Thread, cols, vals []string) Row {
	if p.strategy == projCopy {
		return p.source.Lookup(th, cols, vals)
	}
	p.Select(cols, vals)
	defer p.Select(nil, nil) // clear
	return p.Get(th, Next)
}

func (p *Project) getLookupCost() Cost {
	srcCost := p.source.lookupCost()
	if p.unique {
		return srcCost
	}
	return 2 * srcCost // ??? (matches Nrows)
}
