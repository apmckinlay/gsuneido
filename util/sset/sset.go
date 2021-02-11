// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package sset

//go:generate genny -in ../../genny/set/set.go -out sset2.go -pkg sset gen "T=string"

func eq(x, y string) bool {
	return x == y
}
