// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"regexp"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util2/regex"
)

func TestIsGlobal(t *testing.T) {
	assert := assert.T(t)
	assert.True(isGlobal("F"))
	assert.True(isGlobal("Foo"))
	assert.True(isGlobal("Foo_123_Bar"))
	assert.True(isGlobal("Foo!"))
	assert.True(isGlobal("Foo?"))

	assert.False(isGlobal(""))
	assert.False(isGlobal("f"))
	assert.False(isGlobal("foo"))
	assert.False(isGlobal("_foo"))
	assert.False(isGlobal("Foo!bar"))
	assert.False(isGlobal("Foo?bar"))
	assert.False(isGlobal("Foo.bar"))
	assert.False(isGlobal("\nFoobar"))
	assert.False(isGlobal("Foobar\n"))
}

func BenchmarkNumberPatRegexp(b *testing.B) {
	numPat := regexp.MustCompile(`^[A-Z][_a-zA-Z0-9]*?[!?]?$`)
	big := strings.Repeat("x", 1000)
	for i := 0; i < b.N; i++ {
		numPat.MatchString("X")
		numPat.MatchString("Foobar")
		numPat.MatchString("foobar")
		numPat.MatchString(big)
	}
}

func BenchmarkNumberPatRegex(b *testing.B) {
	numPat := regex.Compile(`\A[A-Z][_a-zA-Z0-9]*?[!?]?\Z`)
	big := strings.Repeat("x", 1000)
	for i := 0; i < b.N; i++ {
		numPat.Matches("X")
		numPat.Matches("Foobar")
		numPat.Matches("foobar")
		numPat.Matches(big)
	}
}
