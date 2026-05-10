// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"maps"
	"math"
	"strings"

	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/set"
)

// perIndex returns an idxSel for each usable index.
// It is called by optInit (which is called on-demand by several methods)
// Its input (perCol) is the result of perField (where2.go).
// Its output is used by Nrows, bestIndex, and finally Get.
// It sets w.singleton if the where selects a single row.
func (w *Where) perIndex(perCol map[string][]span) []*idxSel {
	idxSels := make([]*idxSel, 0, 4)
	indexes := w.tbl.schemaIndexes()
	for i := range indexes {
		schix := &indexes[i]
		isel, singleton := w.buildIdxSel(schix.Columns, schix.Mode, perCol)
		if singleton {
			w.singleton = true
			isel.indexFrac = w.prefixFrac(isel)
			return []*idxSel{isel}
		}
		if len(isel.prefixRanges) > 0 || isel.indexFilter < 1 {
			isel.indexFrac = w.indexFrac(isel)
			idxSels = append(idxSels, isel)
		}
	}
	return idxSels
}

func (w *Where) buildIdxSel(index []string, mode byte, perCol map[string][]span) (*idxSel, bool) {
	key := mode == 'k'
	encode := !key || len(index) > 1
	isel := idxSel{index: index, encoded: encode, mode: mode,
		indexFilter: 1, dataFilter: 1}

	// Fast path: all prefix columns have single-value spans
	if prefixLen, org, ok := allSingleValuePrefix(index, encode, perCol); ok {
		isel.prefixLen = prefixLen
		uniq := mode == 'u'
		lookup := prefixLen == len(index) && (key || (uniq && org != ""))
		if lookup {
			isel.prefixRanges = []pointRange{{Org: org}}
		} else {
			assert.That(encode)
			end := org + ixkey.Sep + ixkey.Max
			isel.prefixRanges = []pointRange{{Org: org, End: end}}
		}
		if prefixLen == len(index) {
			if isel.prefixRanges[0].isPoint() {
				return &isel, true
			}
			return &isel, false
		}
	} else if idxSpans := indexSpans(index, perCol); len(idxSpans) > 0 {
		// prefix range
		uniq := mode == 'u'
		exploded := explodeIndexSpans(idxSpans, [][]span{nil})
		comp := makePointRanges(encode, exploded)
		for i := range comp {
			c := &comp[i]
			if c.isPoint() {
				lookup := len(exploded[i]) == len(index) &&
					(key || (uniq && c.Org != ""))
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
			skipScanSuffix(perCol, index, max(1, isel.prefixLen))
		if isel.prefixLen == 0 && isel.skipLen > 0 {
			isel.prefixRanges = []pointRange{{Org: ixkey.Min, End: ixkey.Max}}
		}
	}

	// more filters
	isel.indexFilter, isel.dataFilter = w.moreFilters(index, &isel)
	if len(isel.prefixRanges) == 0 && isel.indexFilter < 1 {
		// filter-only index selection needs a range for execution
		isel.prefixRanges = []pointRange{{Org: ixkey.Min, End: ixkey.Max}}
	}

	return &isel, false
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
	isel, _ := w.buildIdxSel(index, mode, merged)
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
	idxFields := w.tbl.IndexCols(index)
	for _, sel := range sels {
		if !slices.Contains(idxFields, sel.col) {
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
func (w *Where) moreFilters(index []string, isel *idxSel) (float64, float64) {
	fields := w.tbl.IndexCols(index)
	unconstrained := fields[isel.prefixLen:]
	if isel.skipStart > 0 {
		unconstrained = fields[isel.prefixLen:isel.skipStart]
		unconstrained = append(slices.Clip(unconstrained),
			fields[isel.skipStart+isel.skipLen:]...)
	}
	indexFilter := 1.0
	dataFilter := 1.0
	for _, e := range w.expr.Exprs {
		exprCols := e.Columns()
		if len(exprCols) == 0 || !set.Subset(fields, exprCols) {
			dataFilter = .5
		} else if !set.Disjoint(exprCols, unconstrained) {
			// e.g. index(a,b) where a>1 and F(a,b)
			// F(a,b) overlaps unconstrained (b)
			// so we know it is in addition to the range/skip
			indexFilter = .5
		} else if _, sp := exprToSpans(e, fields); sp == nil {
			// e.g. index(a) where a=1 and F(a)
			// F(a) is not a span
			// so we know it is in addition to the range/skip
			indexFilter = .5
		}
	}
	//TODO when applicable use more specific values like skipFrac
	return indexFilter, dataFilter
}

//-------------------------------------------------------------------

// indexFrac returns the index and data fractions for an idxSel
func (w *Where) indexFrac(isel *idxSel) float64 {
	p := 1.0
	indexFrac := 1.0
	if isel.prefixLen > 0 {
		indexFrac = w.prefixFrac(isel)
		p = .5
	}
	if isel.skipStart > 0 {
		indexFrac *= math.Pow(w.skipFrac(isel), p)
	}
	return indexFrac
}

func (w *Where) prefixFrac(isel *idxSel) float64 {
	iIndex := w.tbl.indexi(isel.index)
	npoints := 0
	var frac float64
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

const ( // ???
	pointFrac   = .01
	rangeFrac   = .1
	compareFrac = .33
)

func (w *Where) skipFrac(isel *idxSel) float64 {
	// not pointRange.isPoint() because skipScanSuffix converts to range
	if wasPoint(isel.skipRange) {
		return pointFrac
	} else if isel.skipRange.Org != ixkey.Min &&
		isel.skipRange.End != ixkey.Max {
		return rangeFrac
	}
	return compareFrac
}

func wasPoint(pr pointRange) bool {
	return strings.HasSuffix(pr.End, ixkey.Sep+ixkey.Max) &&
		strings.HasPrefix(pr.End, pr.Org)
}
