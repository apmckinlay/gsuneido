// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ascii

import (
	"testing"
	"unicode"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestIsLower(*testing.T) {
	for i := range 256 {
		assert.This(IsLower(byte(i))).Is('a' <= i && i <= 'z')
	}
}

func TestIsUpper(*testing.T) {
	for i := range 256 {
		assert.This(IsUpper(byte(i))).Is('A' <= i && i <= 'Z')
	}
}

func TestToLower(*testing.T) {
	for i := range 128 {
		assert.This(ToLower(byte(i))).Is(byte(unicode.ToLower(rune(i))))
	}
}

func TestToUpper(*testing.T) {
	for i := range 128 {
		assert.This(ToUpper(byte(i))).Is(byte(unicode.ToUpper(rune(i))))
	}
}
