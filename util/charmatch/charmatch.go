package charmatch

import (
	"strings"
	"unicode"
)

type CharMatch func(rune) bool

// Predefined CharMatch's
var (
	SPACE  CharMatch = AnyOf(" \t\r\n")
	DIGIT  CharMatch = InRange('0', '9')
	LETTER CharMatch = unicode.IsLetter
	LOWER  CharMatch = unicode.IsLower
	UPPER  CharMatch = unicode.IsUpper
)

// Match returns true if the character matches, otherwise false
func (cm CharMatch) Match(c rune) bool {
	return cm(c)
}

// Is returns a CharMatch that matches a specific character
func Is(c1 rune) CharMatch {
	return func(c2 rune) bool { return c1 == c2 }
}

// AnyOf returns a CharMatch that matches any character in a string
func AnyOf(s string) CharMatch {
	return func(c rune) bool { return strings.ContainsRune(s, c) }
}

// InRange returns a CharMatch that matches any character in the range (inclusive)
func InRange(from rune, to rune) CharMatch {
	return func(c rune) bool { return from <= c && c <= to }
}

// Negate returns a CharMatch that matches any character not matched by this one
func (cm CharMatch) Negate() CharMatch {
	return func(c rune) bool { return !cm.Match(c) }
}

// Or returns a CharMatch that matches any character that matches either CharMatch
func (cm1 CharMatch) Or(cm2 CharMatch) CharMatch {
	return func(c rune) bool { return cm1.Match(c) || cm2.Match(c) }
}

// CountIn returns the number of characters in the string that match
func (cm CharMatch) CountIn(s string) int {
	n := 0
	for _, c := range s {
		if cm.Match(c) {
			n++
		}
	}
	return n
}

// IndexIn returns the index of the first character that matches
// or -1 if no match is found
func (cm CharMatch) IndexIn(s string) int {
	return strings.IndexFunc(s, cm)
}

// Trim returns a slice of the string with all the matching characters
// removed from the beginning and end
func (cm CharMatch) Trim(s string) string {
	return strings.TrimFunc(s, cm)
}

// Trim returns a slice of the string with all the matching characters
// removed from the beginning
func (cm CharMatch) TrimLeft(s string) string {
	return strings.TrimLeftFunc(s, cm)
}
