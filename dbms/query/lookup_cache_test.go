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

	for i := range joinCacheCheckInterval + 50 {
		lc.Lookup(nil, q, []string{"id"}, []string{strconv.Itoa(i)}, nil)
	}
	assert.That(lc.cacheDisabled)
}

func TestLookupCacheGoodPerformance(t *testing.T) {
	var lc lookupCache
	q := &TestQop{}

	for i := range joinCacheCheckInterval + 50 {
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
