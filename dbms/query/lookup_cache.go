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

 // lookupCache encapsulates reusable lookup caching for Query.Lookup calls.
 // It tracks hits/misses and can auto-disable itself if the hit rate is too low.
type lookupCache struct {
	cache         *lrucache.Cache[lookupKey, Row]
	cacheDisabled bool // flag to bypass cache when hit rate is too low
	cacheOpCount  int  // operation counter for periodic evaluation
}

var (
	joinCacheProbes atomic.Int64
	joinCacheMisses atomic.Int64
)

var _ = AddInfo("query.join.cacheProbes", &joinCacheProbes)
var _ = AddInfo("query.join.cacheMisses", &joinCacheMisses)

const (
	joinCacheMinHitRate    = 0.25 // ???
	joinCacheCheckInterval = 256  // ???
	cacheCapacity          = 100
)

// Reset clears the internal cache (if any) and operation counter.
// Note: it does not re-enable a previously disabled cache.
func (lc *lookupCache) Reset() {
	if lc.cache != nil {
		lc.cache.Reset()
	}
	lc.cacheOpCount = 0
}

// Lookup performs a possibly-cached lookup on q.Lookup(cols, vals).
// Caching is bypassed when lc.cacheDisabled is true or when st is updatable.
func (lc *lookupCache) Lookup(th *Thread, q Query, cols, vals []string, st *SuTran) Row {
	if !lc.cacheDisabled && (st == nil || !st.Updatable()) {
		if lc.cache == nil {
			lc.cache = lrucache.New[lookupKey, Row](cacheCapacity)
		}
		lc.cacheOpCount++
		if lc.cacheOpCount%joinCacheCheckInterval == 0 {
			lc.evaluatePerformance()
			if lc.cacheDisabled {
				goto bypass
			}
		}

		key := lookupKey{cols: cols, vals: vals}
		joinCacheProbes.Add(1)
		return lc.cache.GetPut(key, func(k lookupKey) Row {
			joinCacheMisses.Add(1)
			return q.Lookup(th, k.cols, k.vals)
		})
	}
bypass:
	return q.Lookup(th, cols, vals)
}

// evaluatePerformance checks cache hit rate and disables cache if performance is poor.
func (lc *lookupCache) evaluatePerformance() {
	if lc.cache == nil {
		return
	}
	hits, misses := lc.cache.Stats()
	total := hits + misses
	if total > 0 {
		hitRate := float64(hits) / float64(total)
		if hitRate < joinCacheMinHitRate {
			lc.cacheDisabled = true
			lc.cache = nil
		}
	}
}
