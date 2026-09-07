// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"cmp"
	"maps"
	"math"

	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
)

const unknownFrac = .5

// perIndex returns an idxSel for each usable index.
// It is called by optInit (which is called on-demand by several methods)
// Its input (colSpans) is the result of perField (where2.go).
// Its output is used by Nrows, bestIndex, and finally Get.
// It sets w.singleton if the where selects a single row.
func (w *Where) perIndex() []*idxSel {
	colFracs := w.getColFracs(w.colSpans)
	idxSels := make([]*idxSel, 0, 4)
	indexes := w.tbl.schemaIndexes()
	minFrac := 1.0
	for i := range indexes {
		schix := &indexes[i]
		isel := w.buildIdxSel(schix.Fields, schix.Mode, w.colSpans)
		if isel.singleton {
			w.singleton = true
			return []*idxSel{isel}
		}
		if len(isel.prefixRanges) == 0 {
			continue
		}
		// calculate fracs
		isel.prefixFrac = w.prefixFrac(isel)
		frac := isel.calcFracs(colFracs, w.colSpans, w.unspanable, w.dataExprCount)
		if frac < minFrac {
			minFrac = frac
		}

		idxSels = append(idxSels, isel)
	}
	if len(idxSels) == 0 {
		// use a dummy idxSel so we can use calcFracs
		isel := &idxSel{}
		w.wfrac = isel.calcFracs(colFracs, w.colSpans, w.unspanable, w.dataExprCount)
		return nil
	}
	w.wfrac = minFrac
	return idxSels
}

//-------------------------------------------------------------------

type colFrac struct {
	col  string
	frac float64
}

// returns a list of colFracs for colSpans, sorted by frac
func (w *Where) getColFracs(colSpans map[string][]span) []colFrac {
	var colFracs []colFrac
	for col := range colSpans {
		if frac, ok := w.colStatsFrac(w.tbl.Name(), col, colSpans[col]); ok {
			colFracs = append(colFracs, colFrac{col: col, frac: frac})
		}
	}
	return colFracs
}

// colStatsFrac returns the estimated fraction of rows matching any of the
// spans for col, or false if there are no stats for the column.
func (w *Where) colStatsFrac(table, col string, spans []span) (float64, bool) {
	frac := 0.0
	for _, sp := range spans {
		var f float64
		var ok bool
		if sp.isValue() {
			f, ok = w.t.StatsPointFrac(table, col, sp.org.val)
		} else {
			f, ok = w.t.StatsRangeFrac(table, col,
				statsBound(sp.org), statsBound(sp.end))
		}
		if !ok {
			return 0, false
		}
		frac += f
	}
	// cap the sum of the spans (each span fraction is already in [0,1])
	return min(frac, 1.0), true
}

// statsBound adjusts a span bound to the [from, to) semantics
// of StatsRangeFrac. For org, inc means exclusive (>);
// for end, inc means inclusive (<=); in both cases
// appending "\x00" gives the correct boundary (cf. valRaw).
func statsBound(x side) string {
	if x.inc {
		return x.val + "\x00"
	}
	return x.val
}

//-------------------------------------------------------------------

// buildIdxSel constructs an idxSel describing how to use one index:
// prefix point/ranges from leading column equalities, a skip scan range
// on a later column, and any residual index/data filters.
// It does not do anything with selectivity fracs.
func (w *Where) buildIdxSel(index []string, mode byte, colSpans map[string][]span) *idxSel {
	encode := mode != 'k' || len(index) > 1
	isel := idxSel{index: index, encoded: encode, mode: mode}

	// Fast path: all prefix columns have single-value spans
	if prefixLen, org, ok := allSingleValuePrefix(index, encode, colSpans); ok {
		isel.prefixLen = prefixLen
		// A unique 'u' index appends Ixspec.Fields2 (the table's best key)
		// to the physical entry when the index value is entirely empty
		// (to allow multiple empty entries), so a Lookup by value alone
		// only works when the value is not empty.
		lookup := prefixLen == len(index) && (mode != 'u' || org != "")
		if lookup {
			isel.prefixRanges = []pointRange{{Org: org}}
		} else {
			assert.That(encode)
			var enc ixkey.Encoder
			for i := range prefixLen {
				enc.Add(colSpans[index[i]][0].org.val)
			}
			enc.Add(ixkey.Max)
			end := enc.String()
			isel.prefixRanges = []pointRange{{Org: org, End: end}}
		}
		if prefixLen == len(index) {
			if isel.prefixRanges[0].isPoint() {
				isel.singleton = true
				return &isel
			}
			return &isel
		}
	} else if idxSpans := indexSpans(index, colSpans); len(idxSpans) > 0 {
		// prefix range
		exploded := explodeIndexSpans(idxSpans, [][]span{nil})
		comp := makePointRanges(encode, exploded)
		for i := range comp {
			c := &comp[i]
			if c.isPoint() {
				// see comment above about unique 'u' index Lookup and empty values
				lookup := len(exploded[i]) == len(index) && (mode != 'u' || c.Org != "")
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
	if len(index) > 1 {
		isel.skipStart, isel.skipLen, isel.skipRange =
			skipScanSuffix(colSpans, isel.index, max(1, isel.prefixLen))
	}

	if len(isel.prefixRanges) == 0 &&
		(isel.skipLen > 0 || w.hasIndexFilter(index, colSpans, &isel)) {
		// need a range for execution
		isel.prefixRanges = []pointRange{{Org: ixkey.Min, End: ixkey.Max}}
	}

	return &isel
}

// allSingleValuePrefix checks if all prefix columns have exactly one value span.
// If so, it encodes the values directly and returns (prefixLen, org, true).
// Otherwise returns (0, "", false).
func allSingleValuePrefix(index []string, encode bool, colSpans map[string][]span) (int, string, bool) {
	prefixLen := 0
	for i, col := range index {
		cs := colSpans[col]
		if cs == nil {
			break
		}
		if len(cs) != 1 || !cs[0].isValue() {
			return 0, "", false
		}
		prefixLen = i + 1
	}
	if prefixLen == 0 {
		return 0, "", false
	}
	var org string
	if !encode {
		org = colSpans[index[0]][0].org.val
	} else {
		var enc ixkey.Encoder
		for i := 0; i < prefixLen; i++ {
			enc.Add(colSpans[index[i]][0].org.val)
		}
		org = enc.String()
	}
	return prefixLen, org, true
}

// recalcIdxSel rebuilds the idxSel for the current index using merged
// where+select constraints. Returns (isel, conflict).
func (w *Where) recalcIdxSel(index []string, mode byte, sels Sels) (*idxSel, bool) {
	merged, conflict := w.mergeColSpans(index, sels)
	if conflict {
		return &idxSel{}, true
	}
	isel := w.buildIdxSel(index, mode, merged)
	return isel, false
}

// mergeColSpans builds a colSpans map from w.colSpans intersected with equality
// spans for the select cols that appear in the current index.
// Returns (nil, true) if the intersection results in a conflict.
func (w *Where) mergeColSpans(index []string, sels Sels) (map[string][]span, bool) {
	if w.mergeSpans == nil {
		w.mergeSpans = make(map[string][]span, len(w.colSpans)+len(sels))
	} else {
		clear(w.mergeSpans)
	}
	maps.Copy(w.mergeSpans, w.colSpans)
	for _, sel := range sels {
		if !slices.Contains(index, sel.col) {
			continue
		}
		eq := []span{valSpan(sel.val)}
		if existing := w.mergeSpans[sel.col]; existing != nil {
			result := intersectSpans(existing, eq)
			if result == nil {
				return nil, true // conflict
			}
			w.mergeSpans[sel.col] = result
		} else {
			w.mergeSpans[sel.col] = eq
		}
	}
	return w.mergeSpans, false
}

// indexSpans returns the spans for an index
func indexSpans(idx []string, colSpans map[string][]span) [][]span {
	idxSpans := make([][]span, 0, len(idx))
	for i := range idx {
		cs := colSpans[idx[i]]
		if cs == nil {
			break
		}
		idxSpans = append(idxSpans, cs)
		if hasRange(cs) {
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
// NOTE: skip scan only applies to multi-column indexes, so we can assume encoding.
func skipScanSuffix(colSpans map[string][]span, idx []string, prefixLen int) (
	start, size int, sr pointRange) {
	for i := prefixLen; i < len(idx); i++ {
		spans := indexSpans(idx[i:], colSpans)
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

// hasIndexFilter is only called when there is no prefix or skip scan.
// It checks for any other filters on the index columns.
func (w *Where) hasIndexFilter(index []string, colSpans map[string][]span, isel *idxSel) bool {
	assert.That(isel.prefixLen == 0 && isel.skipLen == 0)
	for _, col := range index {
		if _, ok := colSpans[col]; ok {
			return true
		}
		if slices.Contains(w.unspanable, col) {
			return true
		}
	}
	return false
}

//-------------------------------------------------------------------

// prefixFrac estimates the fraction of rows matched by the index prefix
// points/ranges: each point contributes ~0.5/nrows (assuming ~50% existence),
// ranges use the btree's RangeFrac, summed and clamped to 1.
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

// damp is a helper for applying exponential dampening to column fractions
// to account for dependence (overlap) between columns.
func damp(i int, frac float64) float64 {
	for range i {
		frac = math.Sqrt(frac)
	}
	return frac
}

// calcFracs calculates the selectivity of each stage (index range and filter).
// The prefix columns are covered by prefixFrac (from the btree probe), so they
// are excluded here to avoid double counting. This is used for costing in [WhereCost].
func (isel *idxSel) calcFracs(colFracs []colFrac, colSpans map[string][]span,
	unspanable []string, dataExprCount int) float64 {
	const (
		iRange = iota
		iFilter
		dFilter
	)
	type stageFrac struct {
		stage int
		frac  float64
	}
	var fracs []stageFrac
	// prefixFrac is from the btree probe, more reliable than column stats
	if isel.prefixLen > 0 {
		fracs = append(fracs, stageFrac{iRange, isel.prefixFrac})
	}

	for col := range colSpans {
		i := slices.Index(isel.index, col)
		if i >= 0 && i < isel.prefixLen {
			continue // covered by prefixFrac
		}
		frac := getColFrac(colFracs, col)
		stage := dFilter
		if i >= 0 {
			stage = iFilter
			if i >= isel.skipStart && i < isel.skipStart+isel.skipLen {
				stage = iRange
			}
		}
		fracs = append(fracs, stageFrac{stage, frac})
	}

	for _, col := range unspanable {
		stage := dFilter
		if slices.Contains(isel.index, col) {
			stage = iFilter
		}
		fracs = append(fracs, stageFrac{stage, unknownFrac})
	}

	for range dataExprCount {
		fracs = append(fracs, stageFrac{dFilter, unknownFrac})
	}

	// sort by stage and then frac (range first, most selective first)
	slices.SortFunc(fracs, func(x, y stageFrac) int {
		return cmp.Or(cmp.Compare(x.stage, y.stage), cmp.Compare(x.frac, y.frac))
	})
	// group by stage, applying damp
	sfs := [2]float64{1, 1}
	for i, sf := range fracs {
		if sf.stage == dFilter {
			isel.hasDataFilter = true
		} else {
			sfs[sf.stage] *= damp(i, sf.frac)
		}
	}
	// assign to isel
	isel.indexRangeFrac = sfs[iRange]
	isel.indexFilterFrac = sfs[iFilter]

	// overall frac
	// sort by just frac (most selective first)
	slices.SortFunc(fracs, func(x, y stageFrac) int {
		return cmp.Compare(x.frac, y.frac)
	})
	frac := 1.0
	for i, sf := range fracs {
		frac *= damp(i, sf.frac)
	}
	return frac
}

// getColFrac is a helper for calcFracs to get the selectivity of a column
func getColFrac(colFracs []colFrac, col string) float64 {
	for _, cf := range colFracs {
		if cf.col == col {
			return cf.frac
		}
	}
	return unknownFrac
}
