// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ord

import . "golang.org/x/exp/constraints"

func Min[T Ordered](x, y T) T {
	if x <= y {
		return x
	}
	return y
}

func Max[T Ordered](x, y T) T {
	if x >= y {
		return x
	}
	return y
}

func Compare[T Ordered](x, y T) int {
	if x < y {
		return -1
	} else if x > y {
		return +1
	} else {
		return 0
	}
}
