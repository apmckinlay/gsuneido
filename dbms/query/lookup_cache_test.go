// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestLookupCachePerformanceMonitoring(t *testing.T) {
	var lc lookupCache
	q := &TestQop{} // minimal Query; Lookup returns nil via Nothing

	for i := range lookupCacheCheckInterval + 50 {
		lc.Lookup(nil, q, []string{"id"}, []string{strconv.Itoa(i)}, nil)
	}
	assert.That(lc.cacheDisabled)
}

func TestLookupCacheGoodPerformance(t *testing.T) {
	var lc lookupCache
	q := &TestQop{}

	for i := range lookupCacheCheckInterval + 50 {
		if i%2 == 0 {
			lc.Lookup(nil, q, []string{"id"}, []string{"1"}, nil)
		} else {
			lc.Lookup(nil, q, []string{"id"}, []string{"2"}, nil)
		}
	}
	assert.That(!lc.cacheDisabled)
}

func TestLookupCacheOperationCounter(t *testing.T) {
	var lc lookupCache
	q := &TestQop{}

	lc.Lookup(nil, q, []string{"id"}, []string{"1"}, nil)
	assert.That(lc.cacheOpCount == 1)

	lc.Lookup(nil, q, []string{"id"}, []string{"2"}, nil)
	assert.That(lc.cacheOpCount == 2)

	// Disable cache and verify counter doesn't increment
	lc.cacheDisabled = true
	lc.Lookup(nil, q, []string{"id"}, []string{"3"}, nil)
	assert.That(lc.cacheOpCount == 2)
}

// NOTE: These tests assert minimum deltas to avoid flakiness from other tests
// that may also exercise Join/LeftJoin/Compatible in the same run.

func TestLookupCacheCounters(t *testing.T) {
	db := testDb()
	defer db.Close()

	// Baseline
	baseP := joinCacheProbes.Load()
	baseM := joinCacheMisses.Load()

	// Schema and data: ensure 1:1 join so Join uses cachedLookup path
	doAdmin(db, "create t1 (a) key(a)")
	act(db, "insert { a: '1' } into t1")
	act(db, "insert { a: '2' } into t1")
	act(db, "insert { a: '3' } into t1")

	doAdmin(db, "create t2 (a, b) key(a)")
	act(db, "insert { a: '1', b: 'x' } into t2")
	act(db, "insert { a: '2', b: 'y' } into t2")
	act(db, "insert { a: '3', b: 'z' } into t2")

	// Execute query to drive lookups via cache
	_ = queryAll(db, "t1 join t2")

	afterP := joinCacheProbes.Load()
	afterM := joinCacheMisses.Load()

	// Expect at least 3 probes and 3 misses (new keys)
	if afterP-baseP < 3 || afterM-baseM < 3 {
		t.Fatalf("join counters too small: probes %d->%d misses %d->%d",
			baseP, afterP, baseM, afterM)
	}
}

func TestLeftJoinCacheCounters(t *testing.T) {
	db := testDb()
	defer db.Close()

	baseP := leftJoinCacheProbes.Load()
	baseM := leftJoinCacheMisses.Load()

	// Schema and data: 1:1 (or n:1) so LeftJoin uses cachedLookup path
	doAdmin(db, "create t1 (a) key(a)")
	act(db, "insert { a: '1' } into t1")
	act(db, "insert { a: '2' } into t1")
	act(db, "insert { a: '3' } into t1")

	doAdmin(db, "create t2 (a, b) key(a)")
	act(db, "insert { a: '1', b: 'x' } into t2")
	act(db, "insert { a: '2', b: 'y' } into t2")
	// leave a='3' missing to exercise left side with no match

	_ = queryAll(db, "t1 leftjoin t2")

	afterP := leftJoinCacheProbes.Load()
	afterM := leftJoinCacheMisses.Load()

	// Expect at least 3 probes (one per left row) and >=2 misses (new keys)
	if afterP-baseP < 3 || afterM-baseM < 2 {
		t.Fatalf("leftjoin counters too small: probes %d->%d misses %d->%d",
			baseP, afterP, baseM, afterM)
	}
}

func TestCompatibleIntersectCacheCounters(t *testing.T) {
	db := testDb()
	defer db.Close()

	baseP := compatCacheProbes.Load()
	baseM := compatCacheMisses.Load()

	doAdmin(db, "create t1 (a) key(a)")
	act(db, "insert { a: '1' } into t1")
	act(db, "insert { a: '2' } into t1")
	act(db, "insert { a: '3' } into t1")

	doAdmin(db, "create t2 (a) key(a)")
	act(db, "insert { a: '2' } into t2")
	act(db, "insert { a: '3' } into t2")
	act(db, "insert { a: '4' } into t2")

	// Intersect will iterate one source and Lookup on the other via Compatible
	_ = queryAll(db, "t1 intersect t2")

	afterP := compatCacheProbes.Load()
	afterM := compatCacheMisses.Load()

	// Expect at least some activity; conservatively require >=2
	if afterP-baseP < 2 || afterM-baseM < 2 {
		t.Fatalf("compatible (intersect) counters too small: probes %d->%d misses %d->%d",
			baseP, afterP, baseM, afterM)
	}
}
