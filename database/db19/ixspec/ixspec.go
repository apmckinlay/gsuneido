// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package ixspec defines the type T
// that specifies how to get a particular index key from a record.
// comp.Key and Compare implement how it is used.
package ixspec

import "fmt"

type T struct {
	Cols  []int
	Cols2 []int
}

func (is *T) String() string {
	return fmt.Sprint("ixspec ", is.Cols, ",", is.Cols2)
}
