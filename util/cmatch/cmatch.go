// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package cmatch implements a composable way of matching characters.
// Based loosely on Guava CharMatcher for Java
package cmatch

import "strings"

// CharMatch is the type of the matchers
type CharMatch func(byte) bool

// None is a CharMatch that doesn't match anything
var None CharMatch = func(_ byte) bool { return false }

// Match returns true if the character matches, otherwise false
func (cm CharMatch) Match(c byte) bool {
	return cm(c)
}

// Is returns a CharMatch that matches a specific character
func Is(c1 byte) CharMatch {
	return func(c2 byte) bool { return c1 == c2 }
}

// AnyOf returns a CharMatch that matches any character in a string
func AnyOf(s string) CharMatch {
	return func(c byte) bool { return -1 != strings.IndexByte(s, c) }
}

// InRange returns a CharMatch that matches any character in the range (inclusive)
func InRange(from byte, to byte) CharMatch {
	return func(c byte) bool { return from <= c && c <= to }
}

// Negate returns a CharMatch that matches any character not matched by this one
func (cm CharMatch) Negate() CharMatch {
	return func(c byte) bool { return !cm.Match(c) }
}

// Or returns a CharMatch that matches any character that matches either CharMatch
func (cm CharMatch) Or(cm2 CharMatch) CharMatch {
	if cm == nil {
		return cm2
	} else if cm2 == nil {
		return cm
	}
	return func(c byte) bool { return cm.Match(c) || cm2.Match(c) }
}

// CountIn returns the number of characters in the string that match
func (cm CharMatch) CountIn(s string) int {
	n := 0
	for _, c := range []byte(s) {
		if cm.Match(c) {
			n++
		}
	}
	return n
}

// IndexIn returns the index of the first character that matches
// or -1 if no match is found
func (cm CharMatch) IndexIn(s string) int {
	for i, c := range []byte(s) {
		if cm(c) {
			return i
		}
	}
	return -1
}

// Trim returns a slice of the string with all the matching characters
// removed from the beginning and end
func (cm CharMatch) Trim(s string) string {
	i := 0
	for ; i < len(s) && cm(s[i]); i++ {
	}
	j := len(s)
	for ; j > i && cm(s[j-1]); j-- {
	}
	return s[i:j]
}

// TrimLeft returns a slice of the string with all the matching characters
// removed from the beginning
func (cm CharMatch) TrimLeft(s string) string {
	i := 0
	for ; i < len(s) && cm(s[i]); i++ {
	}
	return s[i:]
}
