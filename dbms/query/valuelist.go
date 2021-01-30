// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import . "github.com/apmckinlay/gsuneido/runtime"

func vContains(x []Value, v Value) bool {
	for _, xv := range x {
		if xv.Equal(v) {
			return true
		}
	}
	return false
}

func vDisjoint(x, y []Value) bool {
	for _, xv := range x {
		for _, yv := range y {
			if xv.Equal(yv) {
				return false
			}
		}
	}
	return true
}

func vUnion(x, y []Value) []Value {
	z := make([]Value, 0, len(x)+len(y))
	z = append(z, x...)
outer:
	for _, yv := range y {
		for _, xv := range x {
			if yv == xv {
				continue outer
			}
		}
		z = append(z, yv)
	}
	return z
}
