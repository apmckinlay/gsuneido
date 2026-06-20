// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"slices"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
	"github.com/apmckinlay/gsuneido/util/str"
)

//go:generate stringer -type=Use

type Use int

const (
	ReqUnordered Use = iota
	ReqOrdered
	ReqGrouped
	ReqLookup
	ReqConflict = -1
)

type Require struct {
	use      Use
	cols     []string
	frac     float32
	nlookups int32
} // 32 bytes

func (r Require) Use() Use {
    if len(r.cols) == 0 {
        return ReqUnordered
    }
    if r.frac > 0 {
        if r.nlookups > 0 {
            return ReqGrouped
        }
        return ReqOrdered
    }
    return ReqLookup // frac == 0, nlookups > 0
}

func (r Require) String() string {
	s := r.use.String() + " " + str.Join("(,)", r.cols)
	if r.nlookups > 0 || r.frac > 0 {
		s += fmt.Sprintf(" f%g n%d", r.frac, r.nlookups)
	}
	return s
}

var reqUnordered = Require{}

// MergeReq combines index usage requirements symmetrically.
// Only really applicable to Query1, Project and Summarize
func MergeReq(req1 Use, cols1 []string, req2 Use, cols2 []string) (Use, []string) {
	if req1 > req2 { // order to simplify the cases
		req1, req2 = req2, req1
		cols1, cols2 = cols2, cols1
	}
	switch {
	case req1 == ReqUnordered:
		return req2, cols2
	case req1 == ReqOrdered && req2 == ReqOrdered:
		if slc.HasPrefix(cols1, cols2) {
			return ReqOrdered, cols1
		}
		if slc.HasPrefix(cols2, cols1) {
			return ReqOrdered, cols2
		}
	case req1 == ReqGrouped && req2 == ReqGrouped:
		if set.Equal(cols1, cols2) {
			return ReqGrouped, cols1
		}
		// COULD merge to Ordered
	case req1 == ReqOrdered:
		if req2 == ReqGrouped {
			return orderedPlusGrouped(cols1, cols2)
		}
		if req2 == ReqLookup {
			if set.Equal(cols1, cols2) {
				return ReqOrdered, cols1
			}
		}
	case req1 == ReqGrouped: // && req2 == ReqLookup
		if set.StartsWithSet(cols2, cols1) {
			return ReqLookup, cols2
		}
	case req1 == ReqLookup: // && req2 == ReqLookup
		if set.Equal(cols1, cols2) {
			return ReqLookup, cols1
		}
	default:
		assert.ShouldNotReachHere()
	}
	return ReqConflict, nil
}

// orderedPlusGrouped - cols1 is ordered, cols2 is grouped
func orderedPlusGrouped(cols1 []string, cols2 []string) (Use, []string) {
	// result must start with (be ordered by) cols1
	cols := slices.Clip(cols1)
	if len(cols1) >= len(cols2) {
		if !set.StartsWithSet(cols1, cols2) {
			return ReqConflict, nil
		}
	} else {
		if !set.StartsWithSet(cols2, cols1) {
			return ReqConflict, nil
		}
		// add any cols2 not already in cols1
		for _, col := range cols2 {
			if !slices.Contains(cols, col) {
				cols = append(cols, col)
			}
		}
	}
	assert.That(slc.HasPrefix(cols, cols1))
	assert.That(set.StartsWithSet(cols, cols2))
	return ReqOrdered, cols
}
