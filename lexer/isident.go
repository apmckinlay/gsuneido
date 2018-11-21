package lexer

import "github.com/apmckinlay/gsuneido/util/ascii"

func IsIdentifier(s string) bool {
	last := len(s) - 1
	if last < 0 {
		return false
	}
	for i, c := range []byte(s) {
		if !(c == '_' || ascii.IsLetter(c) ||
			(i > 0 && ascii.IsDigit(c)) ||
			(i == last && (c == '?' || c == '!'))) {
			return false
		}
	}
	return true
}
