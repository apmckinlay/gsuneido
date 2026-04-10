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
	index     []string
	encoded   bool
	mode      byte
	indexFrac float64 // not including filters
	dataFrac  float64 // including filters

	// prefix points/ranges
	prefixLen    int
	prefixRanges []pointRange

	// skip scan
	skipStart int // 0 means no skip scan
	skipLen   int
	skipRange pointRange

	moreFilters []string
}

// Used returns a list of columns "used" by the idxSel
func (is idxSel) Used() []string {
	var used []string
	if is.prefixLen > 0 {
		used = append(used, is.index[:is.prefixLen]...)
	}
	if is.skipStart > 0 {
		used = append(used, is.index[is.skipStart:is.skipStart+is.skipLen]...)
	}
	used = append(used, is.moreFilters...)
	return used
}

func (is idxSel) String() string {
	sb := &strings.Builder{}
	sb.WriteString(str.Join("(,)", is.index))

	if is.prefixLen > 0 {
		sb.WriteString(" ")
		sb.WriteString(str.Join(",", is.index[:is.prefixLen]))
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
		sb.WriteString(">")
	}
	if is.skipStart > 0 {
		sb.WriteString(" +")
		sb.WriteString(str.Join(",", is.index[is.skipStart:is.skipStart+is.skipLen]))
		sb.WriteString(": <")
		showKey(sb, is.encoded, is.skipRange.Org)
		if is.skipRange.isRange() {
			sb.WriteString("..")
			showKey(sb, is.encoded, is.skipRange.End)
		}
		sb.WriteString(">")
	}
	if len(is.moreFilters) > 0 {
		sb.WriteString(" ")
		sb.WriteString(str.Join(",", is.moreFilters))
	}
	sb.WriteString(" = ")
	sb.WriteString(fracStr(is.indexFrac))
	if is.dataFrac != is.indexFrac {
		sb.WriteString(" ")
		sb.WriteString(fracStr(is.dataFrac))
	}
	return sb.String()
}

// fracStr formats a fraction to 2 significant digits, without a leading zero
func fracStr(frac float64) string {
	s := strconv.FormatFloat(frac, 'g', 2, 64)
	return strings.Replace(s, "0.", ".", 1)
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

// intersect returns a new pointRange restricted to selOrg,selEnd
func (pr pointRange) intersect(selOrg, selEnd string) pointRange {
	if pr.isPoint() {
		if selOrg <= pr.Org && pr.Org < selEnd {
			return pr
		}
	} else { // range
		pr = pointRange{Org: max(pr.Org, selOrg), End: min(pr.End, selEnd)}
		if pr.Org < pr.End {
			return pr
		}
	}
	return pointRange{Org: "z", End: "a"} // conflict
}

// conflict returns true if the pointRange cannot match anything
func (pr pointRange) conflict() bool {
	return pr.End != "" && pr.End < pr.Org
}
