// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package interleave

import "iter"

// All returns an iterator over all possible interleavings of two slices.
// An interleaving preserves the relative order of elements within each input slice.
// For slices a=[1,2] and b=[3], it generates: [1,2,3], [1,3,2], [3,1,2].
// The number of interleavings is C(len(a)+len(b), len(b)).
func All[T any](a []T, b []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		// when one slice is empty, there's only one interleaving
		if len(a) == 0 {
			yield(b)
			return
		}
		if len(b) == 0 {
			yield(a)
			return
		}
		n := len(b)
		m := n + len(a)
		patt := (1 << n) - 1
		// iterate through patt values with n 1 bits
		for 0 == (patt >> m) {
			if !yield(shuffle(patt, a, b)) {
				return
			}
			lowb := patt & -patt               // extract the lowest bit of the pattern
			incr := patt + lowb                // increment the lowest bit
			diff := patt ^ incr                // extract the bits flipped by the increment
			patt = incr + ((diff / lowb) >> 2) // restore bit count after increment
		}
	}
}

func shuffle[T any](pattern int, a []T, b []T) []T {
	var seq []T
	i := 0
	j := 0
	for range len(a) + len(b) {
		bit := pattern & 1
		pattern >>= 1
		if bit == 0 {
			seq = append(seq, a[i])
		} else {
			seq = append(seq, b[j])
		}
		i += 1 - bit
		j += bit
	}
	return seq
}
