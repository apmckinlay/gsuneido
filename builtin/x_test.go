// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"regexp"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/regex"
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
	assert.False(numberPat.Matches("\n123"))
	assert.False(numberPat.Matches("123\n"))
}

func BenchmarkNumberPatRegexp(b *testing.B) {
	numPat := regexp.MustCompile(`^[+-]?(\d+(\.\d*)?)|(\.\d+)([eE][+-]?\d\d?)?$`)
	big := strings.Repeat("1", 80) + "x"
	for i := 0; i < b.N; i++ {
		numPat.MatchString("0")
		numPat.MatchString("123.45")
		numPat.MatchString("+123.45")
		numPat.MatchString(big)
	}
}

func BenchmarkNumberPatRegex(b *testing.B) {
	numPat := regex.Compile(`\A[+-]?(\d+(\.\d*)?)|(\.\d+)([eE][+-]?\d\d?)?\Z`)
	big := strings.Repeat("1", 80) + "x"
	for i := 0; i < b.N; i++ {
		numPat.Matches("0")
		numPat.Matches("123.45")
		numPat.Matches("+123.45")
		numPat.Matches(big)
	}
}
