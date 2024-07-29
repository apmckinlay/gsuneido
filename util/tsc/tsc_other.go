// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !amd64

package tsc

func Read() uint64 {
	return 0 //
}
