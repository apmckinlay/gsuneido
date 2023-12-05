// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
)

type Fixed struct {
	col string
	// values are packed
	values []string
}

func NewFixed(col string, val Value) Fixed {
	packed := Pack(val.(Packable))
	return Fixed{col: col, values: []string{packed}}
}

func fixedStr(fixed []Fixed) string {
	s := "["
	sep := ""
	for _, fxd := range fixed {
		s += sep + fxd.String()
		sep = ", "
	}
	return s + "]"
}

func (f *Fixed) String() string {
	s := f.col + "=("
	sep := ""
	for _, v := range f.values {
		s += sep + Unpack(v).String()
		sep = ","
	}
	return s + ")"
}

// combineFixed is used by Where and Join.
// The result is all the ones from fixed1 that are not in fixed2,
// plus the ones from fixed2 that are not in fixed1,
// plus the intersection of values of ones that are in both.
// If an intersection is empty, that is a conflict (none).
// e.g. a: 1, b: 2|3 COMBINE b: 3|4, c: 5 => a: 1, b: 3, c: 5
func combineFixed(fixed1, fixed2 []Fixed) (result []Fixed, none bool) {
	if len(fixed1) == 0 {
		return fixed2, false
	}
	if len(fixed2) == 0 {
		return fixed1, false
	}
	result = make([]Fixed, 0, len(fixed1)+len(fixed2))
	// add fixed1 that are not in fixed2
	for _, sf := range fixed1 {
		if !isFixed(fixed2, sf.col) {
			result = append(result, sf)
		}
	}
	// process fixed2
	for _, f2 := range fixed2 {
		if src1vals := getFixed(fixed1, f2.col); src1vals != nil {
			// field is in both
			vals := set.Intersect(src1vals, f2.values)
			if len(vals) == 0 {
				return nil, true // can't match anything
			}
			result = append(result, Fixed{col: f2.col, values: vals})
		} else {
			// add fixed2 that are not in fixed1
			result = append(result, f2)
		}
	}
	return result, false
}

// FixedIntersect is used by Intersect.
// It returns none = true if there are both fixed1 and fixed2,
// but nothing in common.
func FixedIntersect(fixed1, fixed2 []Fixed) (result []Fixed, none bool) {
	if len(fixed1) == 0 || len(fixed2) == 0 {
		return nil, false
	}
	result = make([]Fixed, len(fixed1))
	for i, f1 := range fixed1 {
		if vals2 := getFixed(fixed2, f1.col); vals2 != nil {
			vals := set.Intersect(f1.values, vals2)
			if len(vals) == 0 {
				return nil, true // can't match anything
			}
			result[i] = Fixed{col: f1.col, values: vals}
		} else {
			result[i] = f1
		}
	}
	return result, false
}

// isSingleFixed returns true if col is fixed with a single value
func isSingleFixed(fixed []Fixed, col string) bool {
	for _, f := range fixed {
		if col == f.col && f.single() {
			return true
		}
	}
	return false
}

func isFixed(fixed []Fixed, col string) bool {
	for _, f := range fixed {
		if col == f.col {
			return true
		}
	}
	return false
}

// getFixed returns the values for a column, or nil if not found
func getFixed(fixed []Fixed, col string) []string {
	for _, f := range fixed {
		if col == f.col {
			return f.values
		}
	}
	return nil
}

func (f *Fixed) single() bool {
	return len(f.values) == 1
}

func withoutFixed(cols []string, fixed []Fixed) []string {
	return slc.WithoutFn(cols,
		func(col string) bool { return isSingleFixed(fixed, col) })
}

func withoutFixed2(cols [][]string, fixed []Fixed) [][]string {
	return slc.MapFn(cols,
		func(cols []string) []string { return withoutFixed(cols, fixed) })
}

func fixedWith(fixed Fixed, val string) Fixed {
	return Fixed{col: fixed.col,
		values: append(slices.Clip(fixed.values), val)}
}

func selectFixed(cols, vals []string, fixed []Fixed) (satisfied, conflict bool) {
	satisfied = true
	for i, col := range cols {
		if fv := getFixed(fixed, col); len(fv) == 1 {
			if fv[0] != vals[i] {
				return false, true // conflict
			}
		} else {
			satisfied = false
		}
	}
	return satisfied, false
}

func conflictFixed(cols, vals []string, fixed []Fixed) bool {
	_, conflict := selectFixed(cols, vals, fixed)
	return conflict
}
