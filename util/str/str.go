// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package str has miscellaneous string functions
// Uses ascii package for lower/upper case
package str

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/ints"
)

// Capitalized returns true is the string starts with A-Z, otherwise false
func Capitalized(s string) bool {
	return len(s) >= 1 && ascii.IsUpper(s[0])
}

// Capitalize returns the string with the first letter converted from a-z to A-Z
func Capitalize(s string) string {
	if len(s) > 0 && ascii.IsLower(s[0]) {
		s = string(ascii.ToUpper(s[0])) + s[1:]
	}
	return s
}

// UnCapitalize returns the string with the first letter converted from A-Z to a-z
func UnCapitalize(s string) string {
	if len(s) > 0 && ascii.IsUpper(s[0]) {
		s = string(ascii.ToLower(s[0])) + s[1:]
	}
	return s
}

// IndexFunc returns the index of the first byte that the func returns true for
// else -1
func IndexFunc(s string, f func(byte) bool) int {
	for i, c := range []byte(s) {
		if f(c) {
			return i
		}
	}
	return -1
}

// Dup is intended to make a copy of a string
// so we don't hold a reference to a large source and prevent garbage collection
func Dup(s string) string {
	s = " " + s
	return s[1:]
}

// IndexNotAny returns the index of the first byte not in chars, else -1
func IndexNotAny(s, chars string) int {
	nc := len(chars)
loop:
	for i := 0; i < len(s); i++ {
		c := s[i]
		for j := 0; j < nc; j++ {
			if c == chars[j] {
				continue loop
			}
		}
		return i
	}
	return -1
}

// LastIndexNotAny returns the last index of the first byte not in chars, else -1
func LastIndexNotAny(s, chars string) int {
	nc := len(chars)
loop:
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		for j := 0; j < nc; j++ {
			if c == chars[j] {
				continue loop
			}
		}
		return i
	}
	return -1
}

// Doesc returns the next byte, interpreting escape sequences
func Doesc(s string, i int) (byte, int) {
	c := s[i]
	if c != '\\' || i+1 >= len(s) {
		return c, i
	}
	c = s[i+1]
	switch c {
	case 'n':
		return '\n', i + 1
	case 't':
		return '\t', i + 1
	case 'r':
		return '\r', i + 1
	case '\\', '"', '\'':
		return c, i + 1
	case 'x':
		if i+2 < len(s) {
			dig1 := ascii.Digit(s[i+1], 16)
			dig2 := ascii.Digit(s[i+2], 16)
			if dig1 != -1 && dig2 != -1 {
				return byte(16*dig1 + dig2), i + 3
			}
		}
	}
	return '\\', i
}

// BeforeFirst returns s up to the first occurrence of sub
// or all of s if sub is not found.
func BeforeFirst(s, sub string) string {
	i := strings.Index(s, sub)
	if i == -1 {
		return s
	}
	return s[:i]
}

// AfterFirst returns s up to the first occurrence of sub
// or all of s if sub is not found.
func AfterFirst(s, sub string) string {
	i := strings.Index(s, sub)
	if i == -1 {
		return s // different from stdlib which returns ""
	}
	return s[i+len(sub):]
}

// BeforerLast returns s before the last occurrence of sub
// or all of s if sub is not found.
func BeforeLast(s, sub string) string {
	i := strings.LastIndex(s, sub)
	if i == -1 {
		return s
	}
	return s[:i]
}

// AfterLast returns s after the last occurrence of sub
// or all of s if sub is not found.
func AfterLast(s, sub string) string {
	i := strings.LastIndex(s, sub)
	if i == -1 {
		return s
	}
	return s[i + len(sub):]
}

// Opt returns "" if any of the strings are ""
// else it returns the concatenation of the strings.
// e.g. Opt("=", s) or Opt(s, ",") or Opt("<", s, ">")
func Opt(strs ...string) string {
	for _, s := range strs {
		if s == "" {
			return ""
		}
	}
	return strings.Join(strs, "")
}

// CommaBuilder builds a comma separated list. Zero value is ready to use.
type CommaBuilder struct {
	sb  strings.Builder
	sep string
}

// Add adds a string to the list
func (cb *CommaBuilder) Add(s string) {
	cb.sb.WriteString(cb.sep)
	cb.sb.WriteString(s)
	cb.sep = ","
}

// String returns the comma separated list
func (cb *CommaBuilder) String() string {
	return cb.sb.String()
}

// Split returns nil if s is "", otherwise it returns strings.Split
func Split(s, sep string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, sep)
}

// Subi returns the substring specified by a starting index and a limit index
// allowing indexes to exceed the string
func Subi(s string, i, j int) string {
	if i >= len(s) {
		return ""
	}
	if j >= len(s) {
		return s[i:]
	}
	return s[i:j]
}

// Subn returns the substring specified by a starting index and a length
// allows index and length to exceed the string
func Subn(s string, i, n int) string {
	if i >= len(s) {
		return ""
	}
	if i+n >= len(s) {
		return s[i:]
	}
	return s[i : i+n]
}

func Min(s1, s2 string) string {
	if s1 <= s2 {
		return s1
	}
	return s2
}

func Max(s1, s2 string) string {
	if s1 >= s2 {
		return s1
	}
	return s2
}

// CmpLower compares the ascii.ToLower of each character
// returning -1, 0, or +1 similar to strings.Compare
func CmpLower(s1, s2 string) int {
	n1 := len(s1)
	n2 := len(s2)
	for i := 0; i < n1 && i < n2; i++ {
		c1 := ascii.ToLower(s1[i])
		c2 := ascii.ToLower(s2[i])
		if c1 < c2 {
			return -1
		}
		if c1 > c2 {
			return +1
		}
	}
	return ints.Compare(n1, n2)
}

// ToLower is an ascii version of strings.ToLower
func ToLower(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))
	for _, c := range []byte(s) {
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}
		sb.WriteByte(c)
	}
	return sb.String()
}

// ToUpper is an ascii version os strings.ToUpper
func ToUpper(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))
	for _, c := range []byte(s) {
		if 'a' <= c && c <= 'z' {
			c -= 'a' - 'A'
		}
		sb.WriteByte(c)
	}
	return sb.String()
}

// Contains returns whether the string contains the given byte.
// Equivalent to strings.ContainsRune but slightly simpler.
func Contains(s string, c byte) bool {
	return strings.IndexByte(s, c) >= 0
}

// EqualCI checks if two strings are equal ignoring (ascii) case
func EqualCI(x, y string) bool {
	if len(x) != len(y) {
		return false
	}
	for i := 0; i < len(x); i++ {
		if ascii.ToLower(x[i]) != ascii.ToLower(y[i]) {
			return false
		}
	}
	return true
}

// RemoveSuffix returns s with suffix removed.
// It returns s unchanged if it does not have suffix.
func RemoveSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}
