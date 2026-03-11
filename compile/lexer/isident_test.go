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
	assert(IsIdentifier("foo_bar")).Is(true)
	assert(IsIdentifier("Bar!")).Is(true)
	assert(IsIdentifier("Bar?")).Is(true)
	assert(IsIdentifier("Bar?x")).Is(false)
	
	assert(IsIdentifier("_foo")).Is(true)
	assert(IsIdentifier("_Bar")).Is(true)
	assert(IsIdentifier("_")).Is(true)
	assert(IsIdentifier("__")).Is(false)
	assert(IsIdentifier("__foo")).Is(false)
	assert(IsIdentifier("__Bar")).Is(false)
	assert(IsIdentifier("_123")).Is(false)
	assert(IsIdentifier("_!")).Is(false)
	assert(IsIdentifier("_?")).Is(false)
}
