// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package str has miscellaneous string functions
// Uses ascii package for lower/upper case
package str

import (
	"cmp"
	"strings"

	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/hacks"
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

// BeforeLast returns s before the last occurrence of sub
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
	return s[i+len(sub):]
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
	sep string
	sb  strings.Builder
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

// CmpLower compares the ascii.ToLower of each character
// returning -1, 0, or +1 similar to strings.Compare
func CmpLower(s1, s2 string) int {
	n1 := len(s1)
	n2 := len(s2)
	for i := range min(n1, n2) {
		c1 := ascii.ToLower(s1[i])
		c2 := ascii.ToLower(s2[i])
		if c1 < c2 {
			return -1
		}
		if c1 > c2 {
			return +1
		}
	}
	return cmp.Compare(n1, n2)
}

// ToLower is an ascii version of strings.ToLower
func ToLower(s string) string {
	for i := range len(s) {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			lower := make([]byte, len(s))
			copy(lower, s[:i])
			lower[i] = c + 32
			for j := i + 1; j < len(s); j++ {
				c := s[j]
				if 'A' <= c && c <= 'Z' {
					c += 32
				}
				lower[j] = c
			}
			return hacks.BStoS(lower)
		}
	}
	return s
}

// ToUpper is an ascii version os strings.ToUpper
func ToUpper(s string) string {
	for i := range len(s) {
		c := s[i]
		if 'a' <= c && c <= 'z' {
			upper := make([]byte, len(s))
			copy(upper, s[:i])
			upper[i] = c - 32
			for j := i + 1; j < len(s); j++ {
				c := s[j]
				if 'a' <= c && c <= 'z' {
					c -= 32
				}
				upper[j] = c
			}
			return hacks.BStoS(upper)
		}
	}
	return s
}

// EqualCI checks if two strings are equal ignoring (ascii) case
func EqualCI(x, y string) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range len(x) {
		if ascii.ToLower(x[i]) != ascii.ToLower(y[i]) {
			return false
		}
	}
	return true
}

// Cut slices s around the first instance of sep,
// returning the text before and after sep.
// If sep does not appear in s, cut returns s, "".
// Similar to the Go 1.18 strings.Cut (except it uses a string separator)
func Cut(s string, sep byte) (before, after string) {
	if i := strings.IndexByte(s, sep); i >= 0 {
		return s[:i], s[i+1:]
	}
	return s, ""
}

// CommonPrefix returns the common prefix of two strings
func CommonPrefix(s, t string) string {
	for i := 0; ; i++ {
		if i >= len(s) || i >= len(t) || s[i] != t[i] {
			return s[:i]
		}
	}
}

// CommonPrefixLen returns the length of the common prefix of two strings
func CommonPrefixLen(s, t string) int {
	for i := 0; ; i++ {
		if i >= len(s) || i >= len(t) || s[i] != t[i] {
			return i
		}
	}
}

// HasPrefix is like the Go strings.HasPrefix
// but it allows mixing string and []byte to avoid conversion allocation
func HasPrefix[T1, T2 ~string|~[]byte](s1 T1, s2 T2) bool {
	if len(s1) < len(s2) {
		return false
	}
	for i := range len(s2) {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
