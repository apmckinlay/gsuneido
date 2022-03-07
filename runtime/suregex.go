// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"github.com/apmckinlay/gsuneido/util/regex"
	"golang.org/x/exp/slices"
)

// SuRegex is a compiled regular expression.
// It is not a general purpose Value and is internal, not exposed.
type SuRegex struct {
	ValueBase[SuRegex]
	Pat regex.Pattern
}

var _ Value = SuRegex{}

func (rx SuRegex) Equal(other interface{}) bool {
	rx2, ok := other.(*SuRegex)
	return ok && slices.Equal(rx.Pat, rx2.Pat)
}

func (SuRegex) SetConcurrent() {
	// immutable so ok
}
