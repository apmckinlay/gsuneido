// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
)

// In the query system "fixed" means a column has a constant value
// e.g. in `where x = 1` or `extend x = 1` x is fixed
// Fixed can track multiple values e.g. `where x in (1,2)`
// but some uses require a single value.

type Fixed []Fix

type Fix struct {
	col string
	// values are packed
	values []string
}

// NewFix returns a new Fix with the given value (packed)
func NewFix(col string, val Value) Fix {
	packed := Pack(val.(Packable))
	return Fix{col: col, values: []string{packed}}
}

// With returns a new Fix with the given value appended
func (f Fix) With(val string) Fix {
	return Fix{col: f.col,
		values: append(slices.Clip(f.values), val)}
}

func (fixed Fixed) String() string {
	fixed = slices.Clone(fixed) // so sort doesn't modify original
	slices.SortFunc(fixed,
		func(x, y Fix) int { return strings.Compare(x.col, y.col) })
	var s strings.Builder
	s.WriteString("[")
	sep := ""
	for _, f := range fixed {
		s.WriteString(sep + f.String())
		sep = ", "
	}
	return s.String() + "]"
}

func (f Fix) String() string {
	var s strings.Builder
	s.WriteString(f.col + "=(")
	sep := ""
	for _, v := range f.values {
		s.WriteString(sep + Unpack(v).String())
		sep = ","
	}
	return s.String() + ")"
}

// Combine is used by Where, Join, and Intersect.
// The result is all the ones from fixed1 that are not in fixed2,
// plus the ones from fixed2 that are not in fixed1,
// plus the intersection of values of ones that are in both.
// If an intersection is empty, that is a conflict (none).
// e.g. a: 1, b: 2|3 COMBINE b: 3|4, c: 5 => a: 1, b: 3, c: 5
func (fixed Fixed) Combine(fixed2 Fixed) (result Fixed, none bool) {
	if len(fixed) == 0 {
		return fixed2, false
	}
	if len(fixed2) == 0 {
		return fixed, false
	}
	result = make(Fixed, 0, len(fixed)+len(fixed2))
	// add fixed1 that are not in fixed2
	for _, sf := range fixed {
		if !fixed2.Has(sf.col) {
			result = append(result, sf)
		}
	}
	// process fixed2
	for _, f2 := range fixed2 {
		if src1vals := fixed.Get(f2.col); src1vals != nil {
			// field is in both
			vals := set.Intersect(src1vals, f2.values)
			if len(vals) == 0 {
				return nil, true // can't match anything
			}
			result = append(result, Fix{col: f2.col, values: vals})
		} else {
			// add fixed2 that are not in fixed1
			result = append(result, f2)
		}
	}
	return result, false
}

// Single returns true if col is fixed with a single value
func (fixed Fixed) Single(col string) bool {
	for _, f := range fixed {
		if col == f.col && f.Single() {
			return true
		}
	}
	return false
}

func (f Fix) Single() bool {
	return len(f.values) == 1
}

// Has returns true if col is fixed, possibly with multiple values
func (fixed Fixed) Has(col string) bool {
	for _, f := range fixed {
		if col == f.col {
			return true
		}
	}
	return false
}

// Get returns the values for a column, or nil if not found
func (fixed Fixed) Get(col string) []string {
	for _, f := range fixed {
		if col == f.col {
			return f.values
		}
	}
	return nil
}

// RemoveFrom returns cols without the single fixed columns
func (fixed Fixed) RemoveFrom(cols []string) []string {
	return slc.WithoutFn(cols,
		func(col string) bool { return fixed.Single(col) })
}

// RemoveFrom2 returns cols without the single fixed columns
func (fixed Fixed) RemoveFrom2(cols [][]string) [][]string {
	return slc.MapFn(cols,
		func(cols []string) []string { return fixed.RemoveFrom(cols) })
}

// Match returns satisfied=true if sels are satisfied by fixed
// and conflict=true if sels conflict with fixed
func (fixed Fixed) Match(sels Sels) (satisfied, conflict bool) {
	// fixed 1,2,3 val 5 => conflict
	// fixed 2 val 2 => satisfied
	// fixed 1,2,3 val 2 => not conflict, not satisfied
	if len(fixed) == 0 {
		return false, false
	}
	satisfied = true
	for _, sel := range sels {
		fv := fixed.Get(sel.col)
		if fv != nil && !slices.Contains(fv, sel.val) {
			return false, true // conflict
		}
		if fv == nil || len(fv) > 1 {
			satisfied = false
		}
	}
	return satisfied, false
}

// Conflicts returns true if sels conflict with fixed
func (fixed Fixed) Conflicts(sels Sels) bool {
	_, conflict := fixed.Match(sels)
	return conflict
}

// All returns true if all cols are single fixed
func (fixed Fixed) All(cols []string) bool {
	for _, col := range cols {
		if !fixed.Single(col) {
			return false
		}
	}
	return true
}

// Equal returns true if fixed and fixed2 are equal
func (fixed Fixed) Equal(fixed2 Fixed) bool {
	return slices.EqualFunc(fixed, fixed2, func(x, y Fix) bool {
		return x.col == y.col && slices.Equal(x.values, y.values)
	})
}
