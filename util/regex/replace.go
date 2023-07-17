// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/ascii"
)

// LiteralRep returns (string,true) if rep is literal,
// otherwise it returns ("",false)
func LiteralRep(rep string) (string, bool) {
	if strings.HasPrefix(rep, "\\=") {
		return rep[2:], true
	}
	if -1 == strings.IndexAny(rep, "&\\") {
		return rep, true
	}
	return "", false
}

// Replacement makes a single replacement, handling case conversion, &, and \#
func Replacement(s, rep string, cap *Captures) string {
	if r, ok := LiteralRep(rep); ok {
		return r
	}
	nr := len(rep)
	tr := byte('E')
	var buf strings.Builder
	buf.Grow(nr)
	trcase := func(c byte) byte {
		switch tr {
		case 'l':
			tr = 'E'
			fallthrough
		case 'L':
			c = ascii.ToLower(c)
		case 'u':
			tr = 'E'
			fallthrough
		case 'U':
			c = ascii.ToUpper(c)
		}
		return c
	}
	add := func(ci int) (tr byte) {
		for i := cap[2*ci]; i < cap[2*ci+1]; i++ {
			buf.WriteByte(trcase(s[i]))
		}
		return
	}
	for i := 0; i < nr; i++ {
		c := rep[i]
		if c == '&' {
			add(0)
		} else if c == '\\' && i+1 < nr {
			i++
			c = rep[i]
			switch rep[i] {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				tr = add(int(c - '0'))
			case 'n':
				buf.WriteByte('\n')
			case 't':
				buf.WriteByte('\t')
			case '\\':
				buf.WriteByte('\\')
			case '&':
				buf.WriteByte('&')
			case 'u', 'l', 'U', 'L', 'E':
				tr = c
			default:
				buf.WriteByte(c) // not affected by tr ???
			}
		} else {
			buf.WriteByte(trcase(c))
		}
	}
	return buf.String()
}
