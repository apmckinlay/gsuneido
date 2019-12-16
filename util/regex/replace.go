// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/ascii"
)

func Replace(s, rep string, res *Result) string {
	if strings.HasPrefix(rep, "\\=") {
		return rep[2:]
	}
	if -1 == strings.IndexAny(rep, "&\\") {
		return rep
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
	add := func(p part) (tr byte) {
		if p.pos1 == 0 {
			return
		}
		t := p.Part(s)
		for i := 0; i < len(t); i++ {
			buf.WriteByte(trcase(t[i]))
		}
		return
	}
	for i := 0; i < nr; i++ {
		c := rep[i]
		if c == '&' {
			add(res[0])
		} else if c == '\\' && i+1 < nr {
			i++
			c = rep[i]
			switch rep[i] {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				tr = add(res[c-'0'])
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
