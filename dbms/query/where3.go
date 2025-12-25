// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"math"
	"strconv"
	"strings"

	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

// idxSel is the pointRanges for a single index.
type idxSel struct {
	index   []string
	ptrngs  []pointRange
	frac    float64
	nfields int
	encoded bool
}

// pointRange holds either a range or a single key (in org with end = "")
type pointRange struct {
	org string
	end string
}

// perIndex returns an idxSel for each usable index.
// It is called by optInit (which is called on-demand by several methods)
// Its input is the result of perField.
// Its output is used by Nrows, bestIndex, and finally Get
func (w *Where) perIndex(perCol map[string][]span) []idxSel {
	idxSels := make([]idxSel, 0, 4)
	indexes := w.tbl.schemaIndexes()
	for i := range indexes {
		schix := &indexes[i]
		idx := schix.Columns
		key := schix.Mode == 'k'
		uniq := schix.Mode == 'u'
		encode := !key || len(idx) > 1
		if idxSpans := indexSpans(idx, perCol); len(idxSpans) > 0 {
			exploded := explodeIndexSpans(idxSpans, [][]span{nil})
			comp := makePointRanges(encode, exploded)
			for i := range comp {
				c := &comp[i]
				if c.isPoint() {
					lookup := len(exploded[i]) == len(idx) &&
						(key || (uniq && c.org != ""))
					if !lookup {
						// convert point to range
						if !encode {
							c.end = c.org + "\x00"
						} else {
							c.end = c.org + ixkey.Sep + ixkey.Max
						}
					}
				}
			}
			isel := idxSel{index: idx, nfields: len(idxSpans),
				ptrngs: comp, encoded: encode}
			if isel.singleton() {
				w.singleton = true
				idxSels = append(idxSels[:0], isel)
				break
			}
			idxSels = append(idxSels, isel)
		}
	}
	for i := range idxSels {
		is := &idxSels[i]
		is.frac = w.idxFrac(is.index, is.ptrngs)
	}
	return idxSels
}

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
			// Clip so each append will make a new copy (COW)
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

func makePointRanges(encode bool, spans [][]span) []pointRange {
	result := make([]pointRange, len(spans))
outer:
	for i, fs := range spans {
		if !encode {
			assert.That(len(fs) == 1)
			f := fs[0]
			if f.isValue() {
				result[i] = pointRange{org: f.org.val}
			} else { // range
				result[i] = pointRange{org: f.org.valRaw(), end: f.end.valRaw()}
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
					result[i] = pointRange{org: enc.String(), end: enc2.String()}
					continue outer
				}
			}
			result[i] = pointRange{org: enc.String()}
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

func (w *Where) idxFrac(idx []string, ptrngs []pointRange) float64 {
	iIndex := w.tbl.indexi(idx)
	var frac float64
	npoints := 0
	nrows, _ := w.tbl.Nrows()
	for _, pr := range ptrngs {
		if pr.isPoint() {
			npoints++
		} else { // range
			frac += w.t.RangeFrac(w.tbl.Name(), iIndex, pr.org, pr.end)
		}
	}
	if nrows > 0 {
		frac += .5 * float64(npoints) / float64(nrows) // ??? estimate 1/2 exist
	}
	assert.That(!math.IsNaN(frac) && !math.IsInf(frac, 0))
	if frac > 1 {
		frac = 1
	}
	return frac
}

// idxSel -----------------------------------------------------------

func (is idxSel) String() string {
	s := str.Join(",", is.index)
	sep := ": "
	for _, pr := range is.ptrngs {
		s += sep + showKey(is.encoded, pr.org)
		sep = " | "
		if pr.isRange() {
			s += ".." + showKey(is.encoded, pr.end)
		}
	}
	if is.frac != 0 {
		s += " = " + strconv.FormatFloat(is.frac, 'g', 4, 64)
	}
	return s
}

func showKey(encode bool, key string) string {
	if !encode {
		return packToStr(key)
	}
	s := ""
	sep := ""
	for _, t := range ixkey.Decode(key) {
		s += sep + packToStr(t)
		sep = ","
	}
	return s
}

func packToStr(s string) string {
	if s == "" {
		return "''"
	}
	if s[0] == 0xff {
		return "<max>"
	}
	if len(s) == 1 && s[0] == PackString {
		return "PackString"
	}
	if len(s) == 1 && s[0] == PackMinus {
		return "PackMinus"
	}
	if len(s) == 1 && s[0] == PackDate {
		return "PackDate"
	}
	if len(s) == 1 && s[0] == PackDate+1 {
		return "PackDate+1"
	}
	return strings.ReplaceAll(Unpack(s).String(), `"`, `'`)
}

// singleton returns true if we know the result is at most one record
// because there is a single point select on a key
func (is idxSel) singleton() bool {
	return len(is.ptrngs) == 1 && is.ptrngs[0].end == ""
}

// pointRange -------------------------------------------------------

func (pr pointRange) isPoint() bool {
	return pr.end == ""
}

func (pr pointRange) isRange() bool {
	return pr.end != ""
}

func (pr pointRange) String() string {
	if pr.conflict() {
		return "<empty>"
	}
	// WARNING: does NOT decode, intended for explode output
	// use idxSel.String for compositePtrngs output
	s := packToStr(pr.org)
	if pr.isRange() {
		s += ".." + packToStr(pr.end)
	}
	return s
}

// intersect returns a new pointRange restricted to selOrg,selEnd
func (pr pointRange) intersect(selOrg, selEnd string) pointRange {
	if pr.isPoint() {
		if selOrg <= pr.org && pr.org < selEnd {
			return pr
		}
	} else { // range
		pr = pointRange{org: max(pr.org, selOrg), end: min(pr.end, selEnd)}
		if pr.org < pr.end {
			return pr
		}
	}
	return pointRange{org: "z", end: "a"} // conflict
}

func (pr pointRange) conflict() bool {
	return pr.end != "" && pr.end < pr.org
}
