// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package ascii contains functions for dealing with ASCII characters
// Note: the Is... functions are usable with cmatch
package ascii

// IsLower returns whether an ASCII character is lower case (a-z)
func IsLower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// IsUpper returns whether an ASCII character is upper case (A-Z)
func IsUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

// ToLower converts an ASCII character to lower case
func ToLower(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		c += 'a' - 'A'
	}
	return c
}

// ToUpper converts an ASCII character to upper case
func ToUpper(c byte) byte {
	if 'a' <= c && c <= 'z' {
		c -= 'a' - 'A'
	}
	return c
}

// IsLetter returns whether an ASCII character is a letter (a-z or A-Z)
func IsLetter(c byte) bool {
	return IsLower(c) || IsUpper(c)
}

// IsDigit returns whether an ASCII character is a number (0-9)
func IsDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// IsSpace returns whether an ASCII character is white space
func IsSpace(c byte) bool {
	switch c {
	case ' ', '\t', '\r', '\n', '\v':
		return true
	}
	return false
}

// IsHexDigit returns whether an ASCII character is a hexadecimal digit
func IsHexDigit(c byte) bool {
	return IsDigit(c) ||
		('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// Digit returns the integer value of an ascii digit in a specified radix
func Digit(c byte, radix int) int {
	n := 99
	if IsDigit(c) {
		n = int(c - '0')
	} else if IsHexDigit(c) {
		n = int(10 + ToLower(c) - 'a')
	}
	if n < radix {
		return n
	}
	return -1
}
