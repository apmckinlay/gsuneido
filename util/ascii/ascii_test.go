package ascii

import (
	"testing"
	"unicode"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestIsLower(t *testing.T) {
	for i := 0; i < 256; i++ {
		Assert(t).That(IsLower(byte(i)), Equals('a' <= i && i <= 'z'))
	}
}

func TestIsUpper(t *testing.T) {
	for i := 0; i < 256; i++ {
		Assert(t).That(IsUpper(byte(i)), Equals('A' <= i && i <= 'Z'))
	}
}

func TestToLower(t *testing.T) {
	for i := 0; i < 128; i++ {
		Assert(t).That(ToLower(byte(i)), Equals(byte(unicode.ToLower(rune(i)))))
	}
}

func TestToUpper(t *testing.T) {
	for i := 0; i < 128; i++ {
		Assert(t).That(ToUpper(byte(i)), Equals(byte(unicode.ToUpper(rune(i)))))
	}
}
