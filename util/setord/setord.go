// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package setord defines functions for sets of lists.
// For sets of sets see the setset package.
package setord

import "github.com/apmckinlay/gsuneido/util/strs"

//go:generate genny -in ../../genny/set/set.go -out setord2.go -pkg setord gen "T=[]string"

func eq(x, y []string) bool {
	// NOTE: list equal, not set equal
	return strs.Equal(x, y)
}
