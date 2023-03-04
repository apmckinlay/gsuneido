// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import "github.com/apmckinlay/gsuneido/util/generic/slc"

// SparseSet is a sparse array based on:
// https://research.swtch.com/sparse
// A zero value is ready to go.
// Add, Has, and Clear are all constant time operations.
// WARNING: The size is base on the largest *value*
// so it should only be used for small values.
type SparseSet struct {
	// dense is a list of the values
	dense []int16
	// sparse is a sparse array pointing to the value in dense
	sparse []int16
}

func (ss *SparseSet) Add(n int16) {
	ss.dense = append(ss.dense, n)
	if int(n) >= len(ss.sparse) {
		ss.sparse = slc.Grow(ss.sparse, int(n)+1-len(ss.sparse))
	}
	ss.sparse[n] = int16(len(ss.dense) - 1)
}

func (ss *SparseSet) Has(n int16) bool {
	if int(n) >= len(ss.sparse) {
		return false
	}
	j := ss.sparse[n]
	return int(j) < len(ss.dense) && ss.dense[j] == n
}

func (ss *SparseSet) AddNew(n int16) bool {
	if ss.Has(n) {
		return false
	}
	ss.Add(n)
	return true
}

func (ss *SparseSet) Clear() {
	ss.dense = ss.dense[:0]
}
