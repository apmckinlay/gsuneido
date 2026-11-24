// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"regexp"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/regex"
)

func TestIsGlobal(t *testing.T) {
	assert := assert.T(t)
	assert.True(IsGlobal("F"))
	assert.True(IsGlobal("Foo"))
	assert.True(IsGlobal("Foo_123_Bar"))
	assert.True(IsGlobal("Foo!"))
	assert.True(IsGlobal("Foo?"))

	assert.False(IsGlobal(""))
	assert.False(IsGlobal("f"))
	assert.False(IsGlobal("foo"))
	assert.False(IsGlobal("_foo"))
	assert.False(IsGlobal("Foo!bar"))
	assert.False(IsGlobal("Foo?bar"))
	assert.False(IsGlobal("Foo.bar"))
	assert.False(IsGlobal("\nFoobar"))
	assert.False(IsGlobal("Foobar\n"))
}

func BenchmarkNumberPatRegexp(b *testing.B) {
	numPat := regexp.MustCompile(`^[A-Z][_a-zA-Z0-9]*?[!?]?$`)
	big := strings.Repeat("x", 1000)
	for b.Loop() {
		numPat.MatchString("X")
		numPat.MatchString("Foobar")
		numPat.MatchString("foobar")
		numPat.MatchString(big)
	}
}

func BenchmarkNumberPatRegex(b *testing.B) {
	numPat := regex.Compile(`\A[A-Z][_a-zA-Z0-9]*?[!?]?\Z`)
	big := strings.Repeat("x", 1000)
	for b.Loop() {
		numPat.Matches("X")
		numPat.Matches("Foobar")
		numPat.Matches("foobar")
		numPat.Matches(big)
	}
}
