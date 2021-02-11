// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ssset

import "github.com/apmckinlay/gsuneido/util/sset"

//go:generate genny -in ../../genny/set/set.go -out ssset2.go -pkg ssset gen "T=[]string"

func eq(x, y []string) bool {
	return sset.Equal(x, y)
}
