// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
	"github.com/apmckinlay/gsuneido/util/str"
)

// Use specifies the type of requirement placed on index selection.
// There is an implicit contract between the Use type and the operations
// that can be performed:
//   - ReqNone: unordered Get, no Lookup, no Select
//   - ReqOrder: ordered Get, Select, no Lookup
//   - ReqGroup: grouped Get, Select, no Lookup
//   - ReqUnique: Lookup (cols must have an index as a subset), no Select
//
// ReqOrder and ReqGroup do not support Lookup
// because they allow a longer index to be chosen
// which will not work with Lookup.
type Use byte

const (
	UsePrefix = 0b100 // cols can be a prefix of index
	UseFull   = 0b101 // cols must be the entire index
	UseSet    = 0b100 // cols are a set (order is not significant)
	UseOrder  = 0b110 // cols are ordered
)

const (
	ReqNone   Use = 0
	ReqGroup      = UsePrefix | UseSet   // used with Select
	ReqUnique     = UseFull | UseSet     // used with Lookup
	ReqOrder      = UsePrefix | UseOrder // used to Get in a particular order
)

func (use Use) String() string {
	switch use {
	case ReqNone:
		return "ReqNone"
	case ReqGroup:
		return "ReqGroup"
	case ReqUnique:
		return "ReqUnique"
	case ReqOrder:
		return "ReqOrder"
	}
	panic("invalid Use")
}

type Require struct {
	cols   []string
	frac   float32
	nseeks int32
	use    Use
}

func NoneReq(frac float32) Require {
	return Require{use: ReqNone, frac: frac}
}

// OrderReq will return NoneReq if cols is empty
func OrderReq(cols []string, frac float32) Require {
	if len(cols) == 0 {
		return NoneReq(frac)
	}
	return Require{use: ReqOrder, cols: cols, frac: frac}
}

// GroupReq will return NoneReq if cols is empty
func GroupReq(cols []string, frac float32, nseeks int32) Require {
	if len(cols) == 0 {
		return NoneReq(frac)
	}
	return Require{use: ReqGroup, cols: cols, frac: frac, nseeks: nseeks}
}

// UniqueReq will return NoneReq if cols is empty
func UniqueReq(cols []string, nseeks int32) Require {
	if len(cols) == 0 {
		return NoneReq(0)
	}
	return Require{use: ReqUnique, cols: cols, nseeks: nseeks}
}

func (r Require) SeekCount(nrows1 int) int32 {
	if r.nseeks > 0 {
		return r.nseeks
	}
	return int32(float64(max(1, nrows1)) * float64(r.frac))
}

// SelectFrac returns the fraction of rows/groups the parent will access,
// for scaling a downstream source's iteration varcost or its own child frac.
func (r Require) SelectFrac(nrows int) float32 {
	if r.frac > 0 {
		return r.frac
	}
	if r.nseeks > 0 {
		return min(float32(1), float32(r.nseeks)/float32(max(1, nrows)))
	}
	return 1
}

func (r Require) String() string {
	s := r.use.String() + str.Join("(,)", r.cols)
	if r.frac > 0 {
		s += fmt.Sprintf(" f%g", r.frac)
	}
	if r.nseeks > 0 {
		s += fmt.Sprintf(" s%d", r.nseeks)
	}
	return s
}

// SatisfiedBy does not check that ReqUnique cols include a key
func (r Require) SatisfiedBy(index []string) bool {
	switch r.use {
	case ReqNone:
		return true
	case ReqGroup:
		return set.StartsWithSet(index, r.cols)
	case ReqOrder:
		return slc.HasPrefix(index, r.cols)
	case ReqUnique:
		return set.Subset(r.cols, index)
	}
	panic("invalid Require use")
}

// SatisfiedByWithFixed does not check that ReqUnique cols include a key
func (r Require) SatisfiedByWithFixed(index []string, fixed Fixed) bool {
	switch r.use {
	case ReqNone:
		return true
	case ReqGroup:
		nColsUnfixed := countUnfixed(r.cols, fixed)
		return grouped(index, r.cols, nColsUnfixed, fixed)
	case ReqOrder:
		return ordered(index, r.cols, fixed)
	case ReqUnique:
		return indexCovered(index, r.cols, fixed)
	}
	panic("invalid Require use")
}
