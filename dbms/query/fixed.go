// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/generic/set"
)

type Fixed struct {
	col string
	// values are packed
	values []string
}

func NewFixed(col string, val runtime.Value) Fixed {
	packed := runtime.Pack(val.(runtime.Packable))
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
		s += sep + runtime.Unpack(v).String()
		sep = ","
	}
	return s + ")"
}

// combineFixed is used by Where and Join
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
		if getFixed(fixed2, sf.col) == nil {
			result = append(result, sf)
		}
	}
	for _, f2 := range fixed2 {
		if srcvals := getFixed(fixed1, f2.col); srcvals != nil {
			// field is in both
			vals := set.Intersect(srcvals, f2.values)
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

// FixedIntersect is used by Intersect
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

// isFixed returns true if col is fixed with a single value
func isFixed(fixed []Fixed, col string) bool {
	for _, f := range fixed {
		if col == f.col && len(f.values) == 1 {
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
