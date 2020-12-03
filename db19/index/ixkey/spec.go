// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package ixkey handles specifying and encoding index keys
package ixkey

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/runtime"
)

// Spec specifies the field(s) in an index key
type Spec struct {
	Fields  []int
	Fields2 []int
}

func (spec *Spec) Key(rec runtime.Record) string {
	return Key(rec, spec.Fields, spec.Fields2)
}

func (spec *Spec) String() string {
	return fmt.Sprint("ixspec ", spec.Fields, ",", spec.Fields2)
}
