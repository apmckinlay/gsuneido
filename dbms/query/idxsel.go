// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strconv"
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/str"
)

// idxSel specifies usage of a single index
type idxSel struct {
	index   []string
	encoded bool
	mode    byte

	singleton bool

	// prefix points/ranges
	prefixLen    int
	prefixRanges []pointRange
	prefixFrac   float64 // from index probe

	// skip scan range
	skipStart int // 0 means no skip scan; indexes into fields
	skipLen   int
	skipRange pointRange

	// hasDataFilter is where expressions that are not covered by the index
	hasDataFilter bool

	// indexRangeFrac is the fraction of the index
	// that is selected by the prefix and skip pointRange
	// It reduces the amount of the index that must be iterated.
	// indexRangeFrac includes prefixFrac
	indexRangeFrac float64

	// indexFilterFrac is the fraction of the index
	// that is selected by the indexFilter
	// It reduces the number of data rows that need to be filtered.
	indexFilterFrac float64
}

func (is *idxSel) HasSkipScan() bool {
	return is.skipStart > 0
}

func (is idxSel) String() string {
	sb := &strings.Builder{}
	sb.WriteString(str.Join("(,)", is.index))

	if is.prefixLen > 0 {
		prefixCols := is.index
		sb.WriteString(" ")
		sb.WriteString(str.Join(",", prefixCols[:is.prefixLen]))
		sb.WriteString(":")
		sep := " <"
		for _, pr := range is.prefixRanges {
			sb.WriteString(sep)
			showKey(sb, is.encoded, pr.Org)
			if pr.isRange() {
				sb.WriteString("..")
				showKey(sb, is.encoded, pr.End)
			}
			sep = " | "
		}
		if sep != " <" {
			sb.WriteString(">")
		}
	}
	if is.singleton {
		sb.WriteString(" = singleton")
	} else {
		if is.skipStart > 0 {
			skipCols := is.index
			sb.WriteString(" +")
			sb.WriteString(str.Join(",", skipCols[is.skipStart:is.skipStart+is.skipLen]))
			sb.WriteString(": <")
			showKey(sb, is.encoded, is.skipRange.Org)
			if is.skipRange.isRange() {
				sb.WriteString("..")
				showKey(sb, is.encoded, is.skipRange.End)
			}
			sb.WriteString(">")
		}
		sb.WriteString(" =")
		if is.prefixLen > 0 {
			sb.WriteString(" pr: ")
			sb.WriteString(fracStr(is.prefixFrac))
		}
		if is.indexRangeFrac != 0 && is.indexRangeFrac != 1 &&
			is.indexRangeFrac != is.prefixFrac {
			sb.WriteString(" ir: ")
			sb.WriteString(fracStr(is.indexRangeFrac))
		}
		if is.indexFilterFrac != 0 && is.indexFilterFrac != 1 {
			sb.WriteString(" if: ")
			sb.WriteString(fracStr(is.indexFilterFrac))
		}
		if is.hasDataFilter {
			sb.WriteString(" df")
		}
	}
	return sb.String()
}

// fracStr formats a fraction to 2 significant digits, without a leading zero
func fracStr(frac float64) string {
	s := strconv.FormatFloat(frac, 'g', 2, 64)
	if strings.HasPrefix(s, "0.") {
		s = s[1:]
	}
	return s
}

func showKey(sb *strings.Builder, encode bool, key string) {
	if !encode {
		sb.WriteString(packToStr(key))
	} else {
		sep := ""
		for _, s := range ixkey.Decode(key) {
			sb.WriteString(sep)
			sb.WriteString(packToStr(s))
			sep = ","
		}
	}
}

func packToStr(s string) string {
	if s == "" {
		return "''"
	}
	if s[0] == 0xff {
		return "max"
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

//-------------------------------------------------------------------

// pointRange holds either a range or a single key (in org with end = "")
type pointRange iface.Range

func (pr pointRange) isPoint() bool {
	return pr.End == ""
}

func (pr pointRange) isRange() bool {
	return pr.End != ""
}

func (pr pointRange) String() string {
	if pr.conflict() {
		return "<empty>"
	}
	// WARNING: does NOT decode, intended for explode output
	// use idxSel.String for compositePtrngs output
	s := packToStr(pr.Org)
	if pr.isRange() {
		s += ".." + packToStr(pr.End)
	}
	return s
}

// conflict returns true if the pointRange cannot match anything
func (pr pointRange) conflict() bool {
	return pr.End != "" && pr.End < pr.Org
}
