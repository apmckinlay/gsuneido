// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package maps2

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

// SortedByKey returns an iterator for the given map that
// yields the key-value pairs in sorted order.
func SortedByKey[Map ~map[K]V, K cmp.Ordered, V any](m Map) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range slices.Sorted(maps.Keys(m)) {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}
