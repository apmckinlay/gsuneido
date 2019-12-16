// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package lexer

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func Test_isIdentifier(t *testing.T) {
	Assert(t).That(IsIdentifier(""), Equals(false))
	Assert(t).That(IsIdentifier("123"), Equals(false))
	Assert(t).That(IsIdentifier("123bar"), Equals(false))
	Assert(t).That(IsIdentifier("foo123"), Equals(true))
	Assert(t).That(IsIdentifier("foo 123"), Equals(false))
	Assert(t).That(IsIdentifier("_foo"), Equals(true))
	Assert(t).That(IsIdentifier("Bar!"), Equals(true))
	Assert(t).That(IsIdentifier("Bar?"), Equals(true))
	Assert(t).That(IsIdentifier("Bar?x"), Equals(false))
}
