// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import "github.com/apmckinlay/gsuneido/runtime"

type Fixed struct {
	col    string
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

func combineFixed(fixed1, fixed2 []Fixed) []Fixed {
	// fixed1 has precedence e.g. combine(f=1, f=2) => f=1
	if len(fixed1) == 0 {
		return fixed2
	}
	if len(fixed2) == 0 {
		return fixed1
	}
	result := make([]Fixed, len(fixed1))
	copy(result, fixed1)
	for _, f2 := range fixed2 {
		if !hasCol(fixed1, f2.col) {
			result = append(result, f2)
		}
	}
	return result
}

func hasCol(fixed []Fixed, col string) bool {
	for _, f := range fixed {
		if col == f.col {
			return true
		}
	}
	return false
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

func getFixed(fixed []Fixed, col string) []string {
	for _, f := range fixed {
		if col == f.col {
			return f.values
		}
	}
	return nil
}
