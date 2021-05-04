// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package setset defines functions for sets of sets.
// For sets of ordered lists see the setord package.
package setset

import "github.com/apmckinlay/gsuneido/util/sset"

//go:generate genny -in ../../genny/set/set.go -out setset2.go -pkg setset gen "T=[]string"

func eq(x, y []string) bool {
	// NOTE: set equal
	return sset.Equal(x, y)
}
