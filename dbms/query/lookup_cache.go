// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"sync/atomic"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/dbg"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/lrucache"
)

// lookupCache provides an operator-agnostic LRU cache for Query.Lookup.
// - Probes count each cache access attempt (before checking the cache).
// - Misses count only when the provider runs i.e. key not found (cache insert).
// Operators (Join, LeftJoin, Compatible) set their own counters via SetCounters,
// allowing per-operator attribution of probes/misses.
// The cache auto-disables based on hit rate.
// Reset clears state but does not re-enable.
type lookupCache struct {
	cache         *lrucache.Cache[lookupKey, Row, lookupHelper]
	cacheDisabled bool          // flag to bypass cache when hit rate is too low
	cacheOpCount  int           // operation counter for periodic evaluation
	probes        *atomic.Int64 // optional instrumentation: incremented on each cache probe
	misses        *atomic.Int64 // optional instrumentation: incremented when provider runs (cache miss)
	selsCols      []string      // first-seen column sequence; only used by checkCols in dbg builds
}

// Heuristics for auto-disable:
// - Every lookupCacheCheckInterval operations, compute hit rate from cache stats.
// - If hit rate < lookupCacheMinHitRate, disable and drop the cache.
const (
	lookupCacheMinHitRate    = 0.25 // ???
	lookupCacheCheckInterval = 256  // ???
	cacheCapacity            = 200  // empirical
	// little benefit from 50 to 100, but 10% better 50 to 200
)

type lookupKey Sels

type lookupHelper struct{}

// Hash and Equal use only the values, not the columns,
// because the column sequence is constant per cache.

func (lookupHelper) Hash(lk lookupKey) uint64 {
	h := uint64(0)
	for _, sel := range lk {
		h = h*131 + hash.String(sel.val)
	}
	return h
}

func (lookupHelper) Equal(x, y lookupKey) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i].val != y[i].val {
			return false
		}
	}
	return true
}

// SetCounters configures which global counters to increment for this cache.
// Called by operator constructors to attribute metrics to the correct operator.
func (lc *lookupCache) SetCounters(probes, misses *atomic.Int64) {
	lc.probes = probes
	lc.misses = misses
}

// Lookup:
// - Bypasses cache if disabled or the transaction is updatable.
// - Increments probes for every attempted cache access.
// - Increments misses only when the provider runs (i.e. cache insert).
// - Uses GetPut so cache hits avoid invoking the provider function.
func (lc *lookupCache) Lookup(q Query, sels Sels, th *Thread, st *SuTran) Row {
	dbg.Assert(func() bool { return lc.checkCols(sels) })
	if !lc.cacheDisabled && (st == nil || !st.Updatable()) {
		if lc.cache == nil {
			lc.cache = lrucache.NewWith[lookupKey, Row](cacheCapacity, lookupHelper{})
		}
		lc.cacheOpCount++
		if lc.cacheOpCount%lookupCacheCheckInterval == 0 {
			lc.evaluatePerformance()
			if lc.cacheDisabled {
				goto bypass
			}
		}

		key := lookupKey(sels)
		if lc.probes != nil {
			lc.probes.Add(1)
		}
		return lc.cache.GetPut(key, func(k lookupKey) Row {
			if lc.misses != nil {
				lc.misses.Add(1)
			}
			return lookup(q, Sels(k), th, st)
		})
	}
bypass:
	return lookup(q, sels, th, st)
}

// checkCols asserts that every lookup uses the same column sequence.
// It is only called from dbg.Assert, so it's a no-op in production builds.
func (lc *lookupCache) checkCols(sels Sels) bool {
	if lc.selsCols == nil {
		lc.selsCols = make([]string, len(sels))
		for i, sel := range sels {
			lc.selsCols[i] = sel.col
		}
		return true
	}
	if len(lc.selsCols) != len(sels) {
		return false
	}
	for i, sel := range sels {
		if lc.selsCols[i] != sel.col {
			return false
		}
	}
	return true
}

// evaluatePerformance uses cache stats since Reset (hits and misses).
// If hit rate < lookupCacheMinHitRate, the cache is disabled and deleted.
func (lc *lookupCache) evaluatePerformance() {
	if lc.cache == nil {
		return
	}
	hits, misses := lc.cache.Stats()
	total := hits + misses
	if total > 0 {
		hitRate := float64(hits) / float64(total)
		if hitRate < lookupCacheMinHitRate {
			lc.cacheDisabled = true
			lc.cache = nil
		}
	}
}

// Reset clears the internal cache (if any) and operation counter.
// Note: it does not re-enable a previously disabled cache.
func (lc *lookupCache) Reset() {
	if lc.cache != nil {
		lc.cache.Reset()
	}
	lc.cacheOpCount = 0
}
