// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ascii

import (
	"testing"
	"unicode"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestIsLower(t *testing.T) {
	for i := 0; i < 256; i++ {
		Assert(t).That(IsLower(byte(i)), Is('a' <= i && i <= 'z'))
	}
}

func TestIsUpper(t *testing.T) {
	for i := 0; i < 256; i++ {
		Assert(t).That(IsUpper(byte(i)), Is('A' <= i && i <= 'Z'))
	}
}

func TestToLower(t *testing.T) {
	for i := 0; i < 128; i++ {
		Assert(t).That(ToLower(byte(i)), Is(byte(unicode.ToLower(rune(i)))))
	}
}

func TestToUpper(t *testing.T) {
	for i := 0; i < 128; i++ {
		Assert(t).That(ToUpper(byte(i)), Is(byte(unicode.ToUpper(rune(i)))))
	}
}
