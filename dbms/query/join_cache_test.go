// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestJoinCachePerformanceMonitoring(t *testing.T) {
	join := &Join{}
	join.source2 = &TestQop{}
	for i := range joinCacheCheckInterval + 50 {
		join.cachedLookup(nil, []string{"id"}, []string{strconv.Itoa(i)})
	}
	assert.That(join.cacheDisabled)
}

func TestJoinCacheGoodPerformance(t *testing.T) {
	join := &Join{}
	join.source2 = &TestQop{}
	for i := range joinCacheCheckInterval + 50 {
		if i%2 == 0 {
			join.cachedLookup(nil, []string{"id"}, []string{"1"})
		} else {
			join.cachedLookup(nil, []string{"id"}, []string{"2"})
		}
	}
	assert.That(!join.cacheDisabled)
}

func TestJoinCacheOperationCounter(t *testing.T) {
	join := &Join{}
	join.source2 = &TestQop{}

	join.cachedLookup(nil, []string{"id"}, []string{"1"})
	assert.That(join.cacheOpCount == 1)

	join.cachedLookup(nil, []string{"id"}, []string{"2"})
	assert.That(join.cacheOpCount == 2)

	// Disable cache and verify counter doesn't increment
	join.cacheDisabled = true
	join.cachedLookup(nil, []string{"id"}, []string{"3"})
	assert.That(join.cacheOpCount == 2)
}
