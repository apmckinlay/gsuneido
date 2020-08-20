// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestNumberPat(t *testing.T) {
	assert := assert.T(t)
	assert.True(numberPat.Matches("0"))
	assert.True(numberPat.Matches("123"))
	assert.True(numberPat.Matches("+123"))
	assert.True(numberPat.Matches("-123"))
	assert.True(numberPat.Matches(".123"))
	assert.True(numberPat.Matches("123.465"))
	assert.True(numberPat.Matches("-.5"))
	assert.True(numberPat.Matches("-1.5"))
	assert.True(numberPat.Matches("-1.5e2"))
	assert.True(numberPat.Matches("1.5e-23"))

	assert.False(numberPat.Matches(""))
	assert.False(numberPat.Matches("."))
	assert.False(numberPat.Matches("+"))
	assert.False(numberPat.Matches("-"))
	assert.False(numberPat.Matches("-."))
	assert.False(numberPat.Matches("+-."))
	assert.False(numberPat.Matches("1.2.3"))
}
