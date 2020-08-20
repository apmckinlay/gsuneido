// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package lexer

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func Test_isIdentifier(t *testing.T) {
	assert := assert.T(t).This
	assert(IsIdentifier("")).Is(false)
	assert(IsIdentifier("123")).Is(false)
	assert(IsIdentifier("123bar")).Is(false)
	assert(IsIdentifier("foo123")).Is(true)
	assert(IsIdentifier("foo 123")).Is(false)
	assert(IsIdentifier("_foo")).Is(true)
	assert(IsIdentifier("Bar!")).Is(true)
	assert(IsIdentifier("Bar?")).Is(true)
	assert(IsIdentifier("Bar?x")).Is(false)
}
