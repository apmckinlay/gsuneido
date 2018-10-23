// Package ascii contains functions for dealing with ASCII characters
package ascii

// IsLower returns whether an ASCII character is lower case (a-z)
func IsLower(c byte) bool {
	return byte('a') <= c && c <= byte('z')
}

// IsUpper returns whether an ASCII character is upper case (A-Z)
func IsUpper(c byte) bool {
	return byte('A') <= c && c <= byte('Z')
}

// IsLower converts an ASCII character to lower case
func ToLower(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		c += 'a' - 'A'
	}
	return c
}

// IsUpper converts an ASCII character to upper case
func ToUpper(c byte) byte {
	if 'a' <= c && c <= 'z' {
		c -= 'a' - 'A'
	}
	return c
}
