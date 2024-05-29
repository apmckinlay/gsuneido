// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ints

func Abs[T int | int8 | int16 | int32 | int64](x T) T {
	if x < 0 {
        return -x
    }
    return x
}
