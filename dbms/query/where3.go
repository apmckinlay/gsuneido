// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"log"
	"maps"
	"math"

	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/set"
)

const unknownFrac = .5

// perIndex returns an idxSel for each usable index.
// It is called by optInit (which is called on-demand by several methods)
// Its input (perCol) is the result of perField (where2.go).
// Its output is used by Nrows, bestIndex, and finally Get.
// It sets w.singleton if the where selects a single row.
func (w *Where) perIndex(perCol map[string][]span) []*idxSel {
	var minPreFrac *idxSel
	idxSels := make([]*idxSel, 0, 4)
	indexes := w.tbl.schemaIndexes()
	for i := range indexes {
		schix := &indexes[i]
		isel := w.buildIdxSel(schix.Fields, schix.Mode, perCol)
		if isel.singleton {
			w.singleton = true
			return []*idxSel{isel}
		}
		if len(isel.prefixRanges) > 0 || isel.indexFilter {
			isel.prefixFrac = w.prefixFrac(isel)
			if minPreFrac == nil || betterMinPre(isel, minPreFrac) {
				minPreFrac = isel
			}
			idxSels = append(idxSels, isel)
		}
	}
	if len(idxSels) == 0 {
		w.wfrac = unknownFrac
		return nil
	}

	// now we have minPreFrac we can calculate the other fractions
	// if minPreFrac is prefix only then that is the overall selectivity
	w.wfrac = minPreFrac.prefixFrac
	if !minPreFrac.OnlyPrefix() {
		w.wfrac *= unknownFrac
	}

	for _, isel := range idxSels {
		if isel.OnlyPrefix() || isel.prefixFrac < w.wfrac {
			isel.prefixFrac = w.wfrac
		}
		moreFrac := 1.0
		if w.wfrac < isel.prefixFrac {
			moreFrac = w.wfrac / isel.prefixFrac
		}
		isel.skipFrac, isel.indexFilterFrac, isel.dataFilterFrac =
			splitFrac(moreFrac, isel.HasSkipScan(), isel.indexFilter, isel.dataFilter)

		// all the isels should have the same overall fraction
		// since a Where has the same result regardless of index used
		product := isel.prefixFrac * isel.skipFrac * isel.indexFilterFrac * isel.dataFilterFrac
		if !sameFrac(w.wfrac, product) {
			log.Println("ERROR: Where perIndex frac mismatch")
			for _, is := range idxSels {
				fmt.Println("\tisel:", is, "=", fracStr(is.frac()))
			}
		}
	}
	return idxSels
}

// betterMinPre returns true if a should replace b as minPreFrac.
// OnlyPrefix always wins (btree-derived fraction is more reliable).
// Within the same kind, lower prefixFrac wins.
// Rationale: OnlyPrefix means it covers the entire Where expression.
// This means it should be the most selective (within btree probe accuracy).
// No other idxSel should be significantly better
// unless there is an error somewhere.
func betterMinPre(a, b *idxSel) bool {
	if a.OnlyPrefix() != b.OnlyPrefix() {
		return a.OnlyPrefix()
	}
	return a.prefixFrac < b.prefixFrac
}

func sameFrac(x, y float64) bool {
	const epsilon = 1e-9
	return math.Abs(x-y) < epsilon
}

func (w *Where) buildIdxSel(index []string, mode byte, perCol map[string][]span) *idxSel {
	encode := mode != 'k' || len(index) > 1
	isel := idxSel{index: index, encoded: encode, mode: mode}

	// Fast path: all prefix columns have single-value spans
	if prefixLen, org, ok := allSingleValuePrefix(index, encode, perCol); ok {
		isel.prefixLen = prefixLen
		lookup := prefixLen == len(index)
		if lookup {
			isel.prefixRanges = []pointRange{{Org: org}}
		} else {
			assert.That(encode)
			end := org + ixkey.Sep + ixkey.Max
			isel.prefixRanges = []pointRange{{Org: org, End: end}}
		}
		if prefixLen == len(index) {
			if isel.prefixRanges[0].isPoint() {
				isel.singleton = true
				return &isel
			}
			return &isel
		}
	} else if idxSpans := indexSpans(index, perCol); len(idxSpans) > 0 {
		// prefix range
		exploded := explodeIndexSpans(idxSpans, [][]span{nil})
		comp := makePointRanges(encode, exploded)
		for i := range comp {
			c := &comp[i]
			if c.isPoint() {
				lookup := len(exploded[i]) == len(index)
				if !lookup {
					assert.That(encode)
					c.End = c.Org + ixkey.Sep + ixkey.Max
				}
			}
		}
		isel.prefixLen = len(idxSpans)
		isel.prefixRanges = comp
		assert.That(len(isel.prefixRanges) != 1 || !isel.prefixRanges[0].isPoint())
	}

	// skip scan range
	if encode {
		isel.skipStart, isel.skipLen, isel.skipRange =
			skipScanSuffix(perCol, isel.index, max(1, isel.prefixLen))
		if isel.prefixLen == 0 && isel.skipLen > 0 {
			isel.prefixRanges = []pointRange{{Org: ixkey.Min, End: ixkey.Max}}
		}
	}

	// more filters
	isel.indexFilter, isel.dataFilter = w.moreFilters(index, &isel)
	if len(isel.prefixRanges) == 0 && isel.indexFilter {
		// filter-only index selection needs a range for execution
		isel.prefixRanges = []pointRange{{Org: ixkey.Min, End: ixkey.Max}}
	}

	return &isel
}

// allSingleValuePrefix checks if all prefix columns have exactly one value span.
// If so, it encodes the values directly and returns (prefixLen, org, true).
// Otherwise returns (0, "", false).
func allSingleValuePrefix(index []string, encode bool, perCol map[string][]span) (int, string, bool) {
	prefixLen := 0
	for i, col := range index {
		colSpans := perCol[col]
		if colSpans == nil {
			break
		}
		if len(colSpans) != 1 || !colSpans[0].isValue() {
			return 0, "", false
		}
		prefixLen = i + 1
	}
	if prefixLen == 0 {
		return 0, "", false
	}
	var org string
	if !encode {
		org = perCol[index[0]][0].org.val
	} else {
		var enc ixkey.Encoder
		for i := 0; i < prefixLen; i++ {
			enc.Add(perCol[index[i]][0].org.val)
		}
		org = enc.String()
	}
	return prefixLen, org, true
}

// recalcIdxSel rebuilds the idxSel for the current index using merged
// where+select constraints. Returns (isel, conflict).
func (w *Where) recalcIdxSel(index []string, mode byte, sels Sels) (*idxSel, bool) {
	merged, conflict := w.mergedPerCol(index, sels)
	if conflict {
		return &idxSel{}, true
	}
	isel := w.buildIdxSel(index, mode, merged)
	return isel, false
}

// mergedPerCol builds a perCol map from w.colSels intersected with equality
// spans for the select cols that appear in the current index.
// Returns (nil, true) if the intersection results in a conflict.
func (w *Where) mergedPerCol(index []string, sels Sels) (map[string][]span, bool) {
	if w.mergedBuf == nil {
		w.mergedBuf = make(map[string][]span, len(w.colSels)+len(sels))
	} else {
		clear(w.mergedBuf)
	}
	maps.Copy(w.mergedBuf, w.colSels)
	for _, sel := range sels {
		if !slices.Contains(index, sel.col) {
			continue
		}
		eq := []span{valSpan(sel.val)}
		if existing := w.mergedBuf[sel.col]; existing != nil {
			result := intersectSpans(existing, eq)
			if result == nil {
				return nil, true // conflict
			}
			w.mergedBuf[sel.col] = result
		} else {
			w.mergedBuf[sel.col] = eq
		}
	}
	return w.mergedBuf, false
}

// indexSpans returns the spans for an index
func indexSpans(idx []string, perCol map[string][]span) [][]span {
	idxSpans := make([][]span, 0, len(idx))
	for i := range idx {
		colSpans := perCol[idx[i]]
		if colSpans == nil {
			break
		}
		idxSpans = append(idxSpans, colSpans)
		if hasRange(colSpans) {
			break // can't have anything after a range
		}
	}
	return idxSpans
}

func hasRange(spans []span) bool {
	for _, s := range spans {
		if s.isRange() {
			return true
		}
	}
	return false
}

func (sp span) isRange() bool {
	return !sp.isValue()
}

func (sp span) isValue() bool {
	return sp.org.val == sp.end.val && !sp.org.inc && sp.end.inc
}

const explodeWarn = 10_000

// explodeIndexSpans handles multiple values for an index column.
// For example, a in (1,2) and b in (3,4)
// will be expanded to: a,b in ((1,3) (1,4) (2,3) (2,4))
func explodeIndexSpans(remaining [][]span, prefixes [][]span) [][]span {
	f := remaining[0]
	if len(f) == 1 { // single value or final range
		for i := range prefixes {
			prefixes[i] = append(prefixes[i], f[0])
		}
	} else { // len(f) > 1
		newpre := make([][]span, 0, len(f)*len(prefixes))
		if len(prefixes) < explodeWarn && len(newpre) >= explodeWarn {
			Warning("query where explode large >", explodeWarn)
		}
		for i := range prefixes {
			// Clip so append will make a new copy (COW)
			pre := slices.Clip(prefixes[i])
			for _, v := range f {
				p := append(pre, v)
				newpre = append(newpre, p)
			}
		}
		prefixes = newpre
	}
	if len(remaining) > 1 {
		return explodeIndexSpans(remaining[1:], prefixes) // RECURSE
	}
	return prefixes
}

// makePointRanges converts spans to pointRanges
func makePointRanges(encode bool, spans [][]span) []pointRange {
	result := make([]pointRange, len(spans))
outer:
	for i, fs := range spans {
		if !encode {
			assert.That(len(fs) == 1)
			f := fs[0]
			if f.isValue() {
				result[i] = pointRange{Org: f.org.val}
			} else { // range
				result[i] = pointRange{Org: f.org.valRaw(), End: f.end.valRaw()}
			}
		} else {
			var enc ixkey.Encoder
			for _, f := range fs {
				if f.isValue() {
					enc.Add(f.org.val)
				} else { // final range
					enc2 := enc.Dup()
					enc.Add(f.org.val)
					if f.org.inc {
						enc.Add(ixkey.Max)
					}
					enc2.Add(f.end.val)
					if f.end.inc {
						enc2.Add(ixkey.Max)
					}
					result[i] = pointRange{Org: enc.String(), End: enc2.String()}
					continue outer
				}
			}
			result[i] = pointRange{Org: enc.String()}
		}
	}
	return result
}

// valRaw is for non-encoded (single field keys)
func (x side) valRaw() string {
	if x.inc {
		return x.val + "\x00"
	}
	return x.val
}

// skipScanSuffix looks for a skip scan suffix range for index idx.
// prefixLen is the first column position to consider.
// Skip scan only supports a single contiguous range.
// If the first column has multiple spans (e.g. in-list), skip this position.
// If a later column has multiple spans, truncate the spans there.
func skipScanSuffix(perCol map[string][]span, idx []string, prefixLen int) (
	start, size int, sr pointRange) {
	for i := prefixLen; i < len(idx); i++ {
		spans := indexSpans(idx[i:], perCol)
		if len(spans) == 0 {
			continue
		}
		// if first column is multi-span (e.g. in-list), can't start here
		if len(spans[0]) > 1 {
			continue
		}
		// truncate at the first multi-span column
		for j, s := range spans {
			if len(s) > 1 {
				spans = spans[:j]
				break
			}
		}
		sp := make([]span, len(spans))
		for j, s := range spans {
			sp[j] = s[0]
		}
		pr := makePointRanges(true, [][]span{sp})[0]
		if pr.isPoint() {
			// convert point to range
			pr.End = pr.Org + ixkey.Sep + ixkey.Max
		}
		return i, len(spans), pr
	}
	return
}

// moreFilters returns estimated selectivity fractions for expressions
// not already handled by the prefix points/ranges and skip scan.
// indexFilter covers expressions on index columns, dataFilter on non-index columns.
func (w *Where) moreFilters(index []string, isel *idxSel) (bool, bool) {
	unconstrained := index[isel.prefixLen:]
	if isel.skipStart > 0 {
		unconstrained = index[isel.prefixLen:isel.skipStart]
		unconstrained = append(slices.Clip(unconstrained),
			index[isel.skipStart+isel.skipLen:]...)
	}
	indexFilter := false
	dataFilter := false
	for _, e := range w.expr.Exprs {
		exprCols := e.Columns()
		if len(exprCols) == 0 || !set.Subset(index, exprCols) {
			dataFilter = true
		} else if !set.Disjoint(exprCols, unconstrained) {
			// e.g. index(a,b) where a>1 and F(a,b)
			// F(a,b) overlaps unconstrained (b)
			// so we know it is in addition to the range/skip
			indexFilter = true
		} else if _, sp := exprToSpans(e, index); sp == nil {
			// e.g. index(a) where a=1 and F(a)
			// F(a) is not a span
			// so we know it is in addition to the range/skip
			indexFilter = true
		}
	}
	return indexFilter, dataFilter
}

//-------------------------------------------------------------------

func (w *Where) prefixFrac(isel *idxSel) float64 {
	iIndex := w.tbl.indexi(isel.index)
	npoints := 0
	frac := 0.0
	for _, pr := range isel.prefixRanges {
		if pr.isPoint() {
			npoints++
		} else { // range
			frac += w.t.RangeFrac(w.tbl.Name(), iIndex, pr.Org, pr.End)
		}
	}
	nrows, _ := w.tbl.Nrows()
	if nrows > 0 {
		frac += .5 * float64(npoints) / float64(nrows) // ??? estimate 1/2 exist
	}
	assert.That(!math.IsNaN(frac) && !math.IsInf(frac, 0))
	if frac > 1 {
		frac = 1
	}
	return frac
}

func splitFrac(f float64, b1, b2, b3 bool) (float64, float64, float64) {
	n := 0
	if b1 {
		n++
	}
	if b2 {
		n++
	}
	if b3 {
		n++
	}
	if n == 0 {
		return 1, 1, 1
	}
	v := math.Pow(f, 1.0/float64(n))
	r := [...]float64{1, 1, 1}
	if b1 {
		r[0] = v
	}
	if b2 {
		r[1] = v
	}
	if b3 {
		r[2] = v
	}
	return r[0], r[1], r[2]
}
