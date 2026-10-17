// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import "math"

// This file collects the constants that make up the query optimizer's cost
// model. Values are relative cost units (not time) and are tuned together,
// so they are kept in one place. See ticost, WhereCost, Table.costFor, etc.
// NOTE: these are roughly empirical guesstimates, not exact

type Cost = int // alias, not actual type, to avoid too many conversions

// impossible is larger than any real cost, with room to add a few.
const impossible = math.MaxInt / 64 // allow for adding impossible's

// outOfOrder is added when we reverse sources.
// This discourages reversing the order without a good reason
// which makes tests and debugging easier.
const outOfOrder = 10

//-------------------------------------------------------------------
// table access

const (
	tableColCost = 5         // per index column adjustment to tableFast/Slow
	tableFast    = 150       // cost per row for first index and for < tableLarge
	tableSlow    = 300       // cost per row for tables >= tableLarge
	tableLarge   = 4_000_000 // size (bytes) threshold for fast/slow (cached or not)
)

// tableLookupCost is the cost of a lookup based on the btree levels.
func tableLookupCost(levels int) Cost {
	return levels*800 - 300
}

//-------------------------------------------------------------------
// temp index (see ticost)

const (
	tiBaseCost         = 1000 // fixed overhead to build a TempIndex
	tiFieldCost        = 100  // per field adjustment to prefer less fields
	tiMultiReadCost    = 200  // per source row
	tiSingleReadCost   = 100  // multiple tables cost more to iterate
	tiMultiLookupCost  = 400  // cost per lookup on a multi-table temp index
	tiSingleLookupCost = 200  // cost per lookup on a single-table temp index
)

// the per row cost of building a TempIndex
const (
	tiAllCost  = 105 // source is already sorted
	tiPreCost  = 110 // source is sorted by a prefix
	tiNoneCost = 256 // source is not sorted
)

//-------------------------------------------------------------------
// where index cost (see WhereCost, WhereCost2)

const (
	whereIndexFrac = .2 // the index portion of the table cost
	whereDataFrac  = .8 // the data portion of the table cost
)

// wherePessimisticGuard floors selectivity estimates to avoid overly
// optimistic plans when the filter fraction is very low.
func wherePessimisticGuard(frac float64) float64 {
	return .25 + (.75 * frac)
}

const (
	bloomBaseCost = 1000 // fixed overhead to use a bloom filter
	bloomAddCost  = 5.5
	bloomHasCost  = 4.5
)

//-------------------------------------------------------------------
// map / grouping estimates for Project, Summarize, Join

const (
	projGrpDiv   = 4     // estimated rows per group for Project
	sumGrpDiv    = 10    // estimated rows per group for Summarize
	mapCost      = 20    // per row cost of a map in Project, Join, Summarize
	mapThreshold = 10000 // maximum nrows for map strategies
	mapWarn      = 20000 // warn if map exceeds this
)
