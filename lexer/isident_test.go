// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package lexer

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func Test_isIdentifier(t *testing.T) {
	Assert(t).That(IsIdentifier(""), Is(false))
	Assert(t).That(IsIdentifier("123"), Is(false))
	Assert(t).That(IsIdentifier("123bar"), Is(false))
	Assert(t).That(IsIdentifier("foo123"), Is(true))
	Assert(t).That(IsIdentifier("foo 123"), Is(false))
	Assert(t).That(IsIdentifier("_foo"), Is(true))
	Assert(t).That(IsIdentifier("Bar!"), Is(true))
	Assert(t).That(IsIdentifier("Bar?"), Is(true))
	Assert(t).That(IsIdentifier("Bar?x"), Is(false))
}
