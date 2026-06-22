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
	ReqConflict = -1 // only used by MergeReq
)

type Require struct {
	cols     []string
	frac     float32
	nlookups int32
} // 32 bytes

func NewRequire(cols []string, frac float32, nlookups int32) Require {
	return Require{cols: cols, frac: frac, nlookups: nlookups}
}

func UnorderedReq(frac float32) Require {
	req := Require{frac: frac}
	assert.That(req.Use() == ReqUnordered)
	return req
}

// OrderedReq will return UnorderedReq if cols is empty
func OrderedReq(cols []string, frac float32) Require {
	if len(cols) == 0 {
		return UnorderedReq(frac)
	}
	req := Require{cols: cols, frac: frac}
	assert.That(req.Use() == ReqOrdered)
	return req
}

// GroupedReq will return UnorderedReq if cols is empty
func GroupedReq(cols []string, frac float32, nlookups int32) Require {
	if len(cols) == 0 {
		return UnorderedReq(frac)
	}
	if nlookups <= 0 {
		nlookups = 1
	}
	req := Require{cols: cols, frac: frac, nlookups: nlookups}
	assert.That(req.Use() == ReqGrouped)
	return req
}

func LookupReq(cols []string, nlookups int32) Require {
	req := Require{cols: cols, nlookups: nlookups}
	assert.That(req.Use() == ReqLookup)
	return req
}

func (r Require) Use() Use {
	if len(r.cols) == 0 {
		assert.That(r.nlookups == 0)
		return ReqUnordered
	}
	if r.nlookups > 0 {
		if r.frac > 0 {
			return ReqGrouped
		}
		return ReqLookup
	}
	return ReqOrdered
}

// LookupCount returns how many lookups a per-row-driven downstream source
// (Join/SemiJoin source2, Intersect/Minus source2) should expect from this
// require. ReqLookup/ReqGrouped inherit the parent's count; scanning (ReqOrdered/
// ReqUnordered) estimates from frac × nrows1.
// Never derive the count from frac alone — frac=0 (pure lookup) must not
// collapse the count to zero, since each parent lookup drives one child lookup.
func (r Require) LookupCount(nrows1 int) int32 {
	if r.nlookups > 0 {
		return r.nlookups
	}
	return int32(float64(max(1, nrows1)) * float64(r.frac))
}

func (r Require) String() string {
	s := r.Use().String() + " " + str.Join("(,)", r.cols)
	if r.nlookups > 0 || r.frac > 0 {
		s += fmt.Sprintf(" f%g n%d", r.frac, r.nlookups)
	}
	return s
}

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
