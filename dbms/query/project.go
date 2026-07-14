// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"slices"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/shmap"
	"github.com/apmckinlay/gsuneido/util/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

var (
	projCopyCount atomic.Int64
	projSeqCount  atomic.Int64
	projMapCount  atomic.Int64
)

var _ = AddInfo("query.project.copy", &projCopyCount)
var _ = AddInfo("query.project.seq", &projSeqCount)
var _ = AddInfo("query.project.map", &projMapCount)

type Project struct {
	Query1
	results *mapType
	st      *SuTran
	columns []string
	remove  []string
	prevRow Row
	curRow  Row
	projectApproach
	state
	unique        bool
	indexed       bool
	warned        bool
	derivedWarned bool
	prevDir       Dir
	derived       int
	th            *Thread
}

type mapType = shmap.Map[rowHash, struct{}, shmap.Funcs[rowHash]]

type projectApproach struct {
	strat projectStrategy
	req   Require
}

type projectStrategy int

const (
	// projCopy is when the columns contain a key, so it's just pass through
	projCopy projectStrategy = iota + 1
	// projSeq orders by the columns so duplicates are consecutive
	projSeq
	// projMap uses a map to identify duplicates.
	// It does not care about source order — it produces rows in source order,
	// returning the first occurrence of each group.
	// The map is built incrementally (not up front), so cost is variable.
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
			(strings.HasSuffix(col, "_lower!") &&
				slices.Contains(cols, strings.TrimSuffix(col, "_lower!"))) ||
			(strings.HasSuffix(col, "_deps") &&
				slices.Contains(cols, strings.TrimSuffix(col, "_deps")))
	})
	if len(proj) == 0 {
		panic("remove: can't remove all columns")
	}
	p := newProject2(src, proj, false)
	p.remove = cols
	return p
}

// newProject is used by Transform
func newProject(src Query, cols []string) Query {
	if len(cols) == 0 {
		return &ProjectNone{source: src}
	}
	return newProject2(src, cols, false)
}

func newProject2(src Query, cols []string, includeDeps bool) *Project {
	p := &Project{Query1: Query1{source: src}}
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

// hasKey returns whether cols contains a key
// taking fixed into consideration.
// See also [indexContainsKey]
func hasKey(cols []string, keys [][]string, fixed Fixed) bool {
	for _, key := range keys {
		if indexCovered(key, cols, fixed) {
			return true
		}
	}
	return false
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

func (p *Project) String() string {
	s := "project"
	cols := p.columns
	switch p.strat {
	case 0:
		if p.remove != nil {
			s = "remove"
			cols = p.remove
		}
		if !p.unique {
			s += " /*NOT UNIQUE*/"
		}
	case projSeq:
		s += "-seq"
	case projCopy:
		s += "-copy"
	case projMap:
		s += "-map"
	default:
		assert.ShouldNotReachHere()
	}
	return s + " " + str.Join(", ", cols)
}

func (p *Project) SetTran(t QueryTran) {
	p.st = MakeSuTran(t)
	p.source.SetTran(t)
}

// projectKeys keeps keys that are subsets of cols.
// Also used by Summarize
func projectKeys(keys [][]string, cols []string) [][]string {
	var keys2 [][]string
	for _, k := range keys {
		if set.Subset(cols, k) {
			keys2 = append(keys2, k)
		}
	}
	if len(keys2) == 0 {
		return [][]string{cols} // fallback on all columns
	}
	return keys2
}

// projectIndexes keeps prefixes of indexes that are subsets of cols.
// Also used by Summarize
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

const projGrpDiv = 2 // ???

func (p *Project) getNrows() (int, int) {
	nr, pop := p.source.Nrows()
	if !p.unique {
		nr /= projGrpDiv
	}
	return nr, pop
}

func (p *Project) Transform() Query {
	p.remove = nil
	src := p.source.Transform()
	if _, ok := src.(*Nothing); ok {
		return NewNothing(p)
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
			return q.With(src1, src2).Transform()
		}
	case *SemiJoin:
		if set.Subset(p.columns, q.by) {
			src1 := newProject(q.source1, p.columns)
			return q.With(src1, q.source2).Transform()
		}
	case *LeftJoin:
		if set.Subset(p.columns, q.by) {
			src1, src2 := p.splitOver(&q.Query2)
			return q.With(src1, src2).Transform()
		}
	case *Union:
		if p.splitable(&q.Compatible) {
			return NewUnion(p.splitOver(&q.Query2)).Transform()
		}
		// Intersect and Minus are not eligible
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
	q := newProject(r.source, newProj)
	r = NewRename(q, newFrom, newTo)
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
		return NewExtend(&ProjectNone{source: e.source}, afterCols, afterExprs)
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

func (p *Project) Fixed() Fixed {
	if p.fixed == nil {
		p.fixed = projectFixed(p.source.Fixed(), p.columns)
		assert.That(p.fixed != nil)
	}
	return p.fixed
}

// projectFixed is also used by Summarize
func projectFixed(srcFixed Fixed, cols []string) Fixed {
	fixed := Fixed{}
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

// mapThreshold and mapWarn are used by Project and Summarize
const (
	mapThreshold = 10000 // used by optimize
	mapWarn      = 20000
)

func (p *Project) optimize(mode Mode, req Require) (Cost, Cost, any) {
	if p.unique {
		// no dedup needed — pass req through unchanged
		fixcost, varcost := Optimize(p.source, mode, req)
		return fixcost, varcost, &projectApproach{strat: projCopy, req: req}
	}
	// non-unique: merge incoming req with own ReqGroup(p.columns)
	seqFix, seqVar, seqApp := p.seqCost(mode, req)
	mapFix, mapVar, mapApp := p.mapCost(mode, req)
	if seqFix+seqVar <= mapFix+mapVar {
		return seqFix, seqVar, seqApp
	}
	return mapFix, mapVar, mapApp
}

func (p *Project) seqCost(mode Mode, req Require) (Cost, Cost, any) {
	fixed := p.source.Fixed()
	nColsUnfixed := countUnfixed(p.columns, fixed)
	nrows, _ := p.Nrows()
	srcReq := GroupReq(p.columns, req.SelectFrac(nrows), req.nseeks)
	switch req.use {
	case ReqNone:
		fixcost, varcost := Optimize(p.source, mode, srcReq)
		return fixcost, varcost, &projectApproach{strat: projSeq, req: srcReq}
	case ReqOrder:
		if grouped(req.cols, p.columns, nColsUnfixed, fixed) {
			fixcost, varcost := Optimize(p.source, mode, req)
			return fixcost, varcost, &projectApproach{strat: projSeq, req: req}
		}
	case ReqUnique:
		debug.assert(set.Equal(req.cols, p.columns)) // only key is all columns
		// we can use GroupReq because Lookup is implemented by Select + Get
		srcReq := GroupReq(p.columns, req.SelectFrac(nrows), req.nseeks)
		fixcost, varcost := Optimize(p.source, mode, srcReq)
		return fixcost, varcost, &projectApproach{strat: projSeq, req: srcReq}
	case ReqGroup:
		if len(req.cols) == len(p.columns) {
			debug.assert(set.Equal(req.cols, p.columns)) // only key is all columns
			fixcost, varcost := Optimize(p.source, mode, srcReq)
			return fixcost, varcost, &projectApproach{strat: projSeq, req: srcReq}
		}
		if !eitherSubset(req.cols, p.columns) {
			return impossible, impossible, nil
		}
		// requires are different ReqGroup
		// this can't be handled with a single Require
		// so we need to search here
		nColsUnfixedReq := countUnfixed(req.cols, fixed)
		best := newBest[Require]()
		for _, idx := range p.source.Indexes() {
			if grouped(idx, req.cols, nColsUnfixedReq, fixed) &&
				grouped(idx, p.columns, nColsUnfixed, fixed) {
				// source req must be ordered so it doesn't ignore column order
				// which is necessary to satisfy both groupings
				srcReq := OrderReq(idx, req.SelectFrac(nrows))
				f, v := Optimize(p.source, mode, srcReq)
				v += Cost(req.nseeks) * p.source.lookupCost()
				best.update(f, v, srcReq)
			}
		}
		if best.found() {
			return best.fixcost, best.varcost,
				&projectApproach{strat: projSeq, req: best.data}
		}
	}
	return impossible, impossible, nil
}

// eitherSubset returns true if x is a subset of y or y is a subset of x
func eitherSubset(x, y []string) bool {
	if len(x) > len(y) {
		x, y = y, x
	}
	return set.Subset(x, y)
}

const mapCost = 20 // ???

// mapCost estimates the cost of projMap.
// The map is built incrementally during iteration (not up front),
// so the map build cost is added to varcost, not fixcost.
func (p *Project) mapCost(mode Mode, req Require) (Cost, Cost, any) {
	nrows, _ := p.Nrows()
	if mode != ReadMode || nrows > mapThreshold {
		return impossible, impossible, nil
	}
	if req.use == ReqUnique {
		req = GroupReq(req.cols, req.SelectFrac(nrows), req.nseeks)
	}
	srcFixcost, srcVarcost := Optimize(p.source, mode, req)
	srcNrows, _ := p.source.Nrows()
	mapBuild := Cost(float64(srcNrows) * float64(req.frac) * mapCost)
	return srcFixcost, srcVarcost + mapBuild,
		&projectApproach{strat: projMap, req: req}
}

func (p *Project) setApproach(_ Require, approach any, tran QueryTran) {
	p.projectApproach = *approach.(*projectApproach)
	switch p.strat {
	case projCopy:
		projCopyCount.Add(1)
	case projSeq:
		projSeqCount.Add(1)
	case projMap:
		projMapCount.Add(1)
	}
	p.source = SetApproach(p.source, p.projectApproach.req, tran)
	p.header = p.getHeader()
}

// execution --------------------------------------------------------

func (p *Project) Rewind() {
	p.source.Rewind()
	p.rewind()
}

func (p *Project) rewind() {
	p.state = rewound
	p.curRow = nil
	p.prevRow = nil
}

func (p *Project) Get(th *Thread, dir Dir) Row {
	defer func(t uint64) { p.tget += tsc.Read() - t }(tsc.Read())
	if p.state == eof {
		return nil
	}
	var row Row
	switch p.strat {
	case projCopy:
		row = p.source.Get(th, dir)
	case projSeq:
		row = p.getSeq(th, dir)
	case projMap:
		row = p.getMap(th, dir)
	default:
		panic(assert.ShouldNotReachHere())
	}
	if row != nil {
		p.state = within
		p.ngets++
	} else {
		p.state = eof
	}
	return row
}

func (p *Project) getSeq(th *Thread, dir Dir) Row {
	if dir == Next {
		p.prevDir = dir
		// output the first of each group
		// i.e. skip over rows the same as previous output
		for {
			row := p.source.Get(th, dir)
			if row == nil {
				p.prevRow = p.curRow
				p.curRow = nil
				return nil
			}
			if p.state == rewound || p.curRow == nil ||
				!p.header.EqualRows(row, p.curRow, th, p.st) {
				p.prevRow = p.curRow
				p.curRow = row
				return row
			}
		}
	} else { // Prev
		// output the last of each group
		// i.e. output when next record is different
		// (to get the same records as NEXT)

		if p.state == rewound || (p.prevRow == nil && p.prevDir == Next) {
			p.prevRow = p.source.Get(th, dir)
		}
		p.prevDir = dir
		for {
			if p.prevRow == nil {
				p.curRow = nil
				return nil
			}
			row := p.prevRow
			p.prevRow = p.source.Get(th, dir)
			if p.prevRow == nil {
				// source reached the front; rewind so Next direction works
				p.source.Rewind()
				p.curRow = row
				return row
			}
			if !p.header.EqualRows(row, p.prevRow, th, p.st) {
				// output the last row of a group
				p.curRow = row
				return row
			}
		}
	}
}

type rowHash struct {
	row  Row
	hash uint64
}

// getMap returns rows in source order, skipping duplicates.
// For Next direction, the map is built incrementally — each row is checked
// against the map as it is read, so there is no up-front build cost.
// For Prev direction, the full map must be built first (buildMap).
func (p *Project) getMap(th *Thread, dir Dir) Row {
	p.th = th
	defer func() { p.th = nil }()
	if p.state == rewound {
		if p.results == nil {
			hfn := func(k rowHash) uint64 { return k.hash }
			eqfn := func(x, y rowHash) bool {
				return x.hash == y.hash &&
					equalCols(x.row, y.row, p.source.Header(), p.columns, p.th, p.st)
			}
			p.results = shmap.NewMapFuncs[rowHash, struct{}](hfn, eqfn)
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

func hashCols(row Row, hdr *Header, cols []string, th *Thread, st *SuTran) uint64 {
	h := uint64(31)
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
	k, existed := p.results.GetInit(rh)
	if existed {
		return k.row, true
	} else {
		if !p.warned && p.results.Size() > mapWarn {
			p.warned = true
			Warning("project-map large >", mapWarn)
		}
		p.derived += row.Derived()
		if !p.derivedWarned && p.derived > derivedWarn {
			p.derivedWarned = true
			Warning("project-map derived large >",
				derivedWarn, "average", p.derived/p.results.Size())
		}
		return row, false
	}
}

func (p *Project) Output(th *Thread, rec Record) {
	if p.strat != projCopy {
		panic("can't output to a project that doesn't include a key")
	}
	p.source.Output(th, rec)
}

func (p *Project) Select(sels Sels) {
	p.nsels++
	p.source.Select(sels)
	p.indexed = false
	if p.results != nil {
		p.results.Clear()
	}
	p.rewind()
}

func (p *Project) Lookup(th *Thread, sels Sels) Row {
	p.nlooks++
	if p.strat == projCopy {
		return p.source.Lookup(th, sels)
	}
	return lookupViaSelectGet(p, th, sels)
}

func (p *Project) getLookupCost() Cost {
	srcCost := p.source.lookupCost()
	if p.unique {
		return srcCost
	}
	return projGrpDiv * srcCost // ??? (matches Nrows)
}

func (p *Project) Simple(th *Thread) []Row {
	hdr := p.Header()
	dst := 0
	rows := p.source.Simple(th)
outer:
	for i := range rows {
		for j := range i {
			if hdr.EqualRows(rows[i], rows[j], th, nil) {
				continue outer
			}
		}
		rows[dst] = rows[i]
		dst++
	}
	return rows[:dst]
}
