// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

//go:generate stringer -type=Use

// Use specifies the type of requirement placed on a query.
// There is an implicit contract between the Use type and the operations
// that can be performed on the query after SetApproach:
//   - ReqUnordered: iteration only (Get), no Select, no Lookup
//   - ReqOrdered: ordered iteration (Get), no Select, no Lookup
//   - ReqGrouped: enables Select (with the exact index columns, not a prefix)
//   - ReqLookup: enables Lookup (with the exact index columns)
//
// This contract is enforced by the optimizer: Select should only be called
// if the query was optimized with ReqGrouped, and Lookup should only be
// called if optimized with ReqLookup. ReqOrdered is for ordered iteration only.
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
	// Floor nlookups so a degenerate caller (e.g. empty source1 yielding
	// LookupCount==0, or tempindex exploration under a ReqLookup parent)
	// still produces a valid ReqLookup. With nlookups==0 Use() would derive
	// ReqOrdered and the assert below would fail. Matches GroupedReq's floor.
	// The cost distortion (1 seek vs 0) only arises in vacuous cases where
	// the subtree produces no rows anyway.
	if nlookups <= 0 {
		nlookups = 1
	}
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

// SelectFrac returns the fraction of rows/groups the parent will access,
// for scaling a downstream source's iteration varcost or its own child frac.
// Use this — not r.frac directly — whenever a query needs the parent's
// access fraction to build a child Require. r.frac is 0 for ReqLookup
// (frac==0, nlookups>0), so using it naively would compute a child frac of 0,
// collapsing a GroupedReq into ReqLookup. That breaks the implicit contract
// - a node optimized ReqGrouped may only be Select+iterated, never Lookup.
// Deriving from nlookups keeps the child frac > 0, preserving ReqGrouped.
// This is the dual of LookupCount's "never derive from frac alone" warning.
// ReqUnordered/ReqOrdered return the stored frac unchanged.
func (r Require) SelectFrac(nrows int) float32 {
	if r.frac > 0 {
		return r.frac
	}
	if r.nlookups > 0 {
		return min(float32(1), float32(r.nlookups)/float32(max(1, nrows)))
	}
	return 1
}

func (r Require) String() string {
	s := r.Use().String() + str.Join("(,)", r.cols)
	if r.nlookups > 0 || r.frac > 0 {
		s += fmt.Sprintf(" f%g n%d", r.frac, r.nlookups)
	}
	return s
}
