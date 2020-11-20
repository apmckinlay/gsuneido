// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package ixspec defines the type T
// that specifies how to get a particular index key from a record.
// comp.Key and Compare implement how it is used.
package ixspec

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/db19/index/comp"
	"github.com/apmckinlay/gsuneido/runtime"
)

type T = ixspec

type ixspec struct {
	Fields  []int
	Fields2 []int
}

func (is *ixspec) Key(rec runtime.Record) string {
	return comp.Key(rec, is.Fields, is.Fields2)
}

func (is *ixspec) String() string {
	return fmt.Sprint("ixspec ", is.Fields, ",", is.Fields2)
}
