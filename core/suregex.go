// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"github.com/apmckinlay/gsuneido/util/regex"
)

// SuRegex is a compiled regular expression.
// It is not a general purpose Value and is internal, not exposed.
type SuRegex struct {
	ValueBase[SuRegex]
	Pat regex.Pattern
}

var _ Value = SuRegex{}

func (rx SuRegex) Equal(other any) bool {
	rx2, ok := other.(*SuRegex)
	return ok && rx.Pat == rx2.Pat
}

func (SuRegex) SetConcurrent() {
	// immutable so ok
}

// RegexMethods is initialized by the builtin package
var RegexMethods Methods

func (SuRegex) Lookup(_ *Thread, method string) Callable {
	return RegexMethods[method]
}
