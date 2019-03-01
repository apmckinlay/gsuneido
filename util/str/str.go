// Package str has miscellaneous string functions
// Uses ascii package for lower/upper case
package str

import "github.com/apmckinlay/gsuneido/util/ascii"

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

func IndexFunc(s string, f func (byte) bool) int {
	for i,c := range []byte(s) {
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
