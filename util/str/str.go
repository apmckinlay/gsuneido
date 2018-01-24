// Package str has miscellaneous string functions
package str

// Capitalized returns true is the string starts with A-Z, otherwise false
func Capitalized(s string) bool {
	return len(s) >= 1 && 'A' <= s[0] && s[0] <= 'Z'
}

// Capitalize returns the string with the first letter converted from a-z to A-Z
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	first := s[0]
	if 'a' <= first && first <= 'z' {
		first -= 'a' - 'A'
		s = string(first) + s[1:]
	}
	return s
}

// UnCapitalize returns the string with the first letter converted from A-Z to a-z
func UnCapitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	first := s[0]
	if 'A' <= first && first <= 'Z' {
		first += 'a' - 'A'
		s = string(first) + s[1:]
	}
	return s
}
