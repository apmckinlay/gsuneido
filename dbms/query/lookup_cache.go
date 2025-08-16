// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"
	"sync/atomic"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/lrucache"
)

// lookupCache provides an operator-agnostic LRU cache for Query.Lookup.
// - Probes count each cache access attempt (before checking the cache).
// - Misses count only when the provider runs i.e. key not found (cache insert).
// Operators (Join, LeftJoin, Compatible) set their own counters via SetCounters,
// allowing per-operator attribution of probes/misses.
// The cache can auto-disable based on hit rate; Reset clears state but does not re-enable.
type lookupCache struct {
	cache         *lrucache.Cache[lookupKey, Row]
	cacheDisabled bool          // flag to bypass cache when hit rate is too low
	cacheOpCount  int           // operation counter for periodic evaluation
	probes        *atomic.Int64 // optional instrumentation: incremented on each cache probe
	misses        *atomic.Int64 // optional instrumentation: incremented when provider runs (cache miss)
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

// lookupKey is the composite key for caching lookups based on column/value pairs.
type lookupKey struct {
	cols []string
	vals []string
}

func (lk lookupKey) Hash() uint64 {
	h := uint64(0)
	for _, col := range lk.cols {
		h = h*131 + hash.String(col)
	}
	for _, val := range lk.vals {
		h = h*131 + hash.String(val)
	}
	return h
}

func (lk lookupKey) Equal(other any) bool {
	if o, ok := other.(lookupKey); ok {
		return slices.Equal(lk.cols, o.cols) && slices.Equal(lk.vals, o.vals)
	}
	return false
}

// Reset clears the internal cache (if any) and operation counter.
// Note: it does not re-enable a previously disabled cache.
func (lc *lookupCache) Reset() {
	if lc.cache != nil {
		lc.cache.Reset()
	}
	lc.cacheOpCount = 0
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
func (lc *lookupCache) Lookup(th *Thread, q Query, cols, vals []string, st *SuTran) Row {
	if !lc.cacheDisabled && (st == nil || !st.Updatable()) {
		if lc.cache == nil {
			lc.cache = lrucache.New[lookupKey, Row](cacheCapacity)
		}
		lc.cacheOpCount++
		if lc.cacheOpCount%lookupCacheCheckInterval == 0 {
			lc.evaluatePerformance()
			if lc.cacheDisabled {
				goto bypass
			}
		}

		key := lookupKey{cols: cols, vals: vals}
		if lc.probes != nil {
			lc.probes.Add(1)
		}
		return lc.cache.GetPut(key, func(k lookupKey) Row {
			if lc.misses != nil {
				lc.misses.Add(1)
			}
			return q.Lookup(th, k.cols, k.vals)
		})
	}
bypass:
	return q.Lookup(th, cols, vals)
}

// evaluatePerformance uses cache stats since Reset (hits and misses).
// If hit rate < lookupCacheMinHitRate, the cache is disabled and deleted.
// Reset does not re-enable; operators typically re-create per transaction.
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
