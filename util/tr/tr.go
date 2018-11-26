/*
Package tr is similar to the Unix tr program.
Originally based on the code in Software Tools by Kernighan and Plauger.
Currently it is byte based - it does not handle unicode.
*/
package tr

import (
	"strings"
)

type trset string

var EmptySet = trset("")

// Replace translates, squeezes, or deletes characters from the src string.
// If the first character of from is '^' then the from set is complemented.
// Ranges are specified with '-' between to characters.
// If the to set is empty, then characters in the from set are deleted.
// If the to set is shorter than the from set, then
// the last character in the to set is repeated to make the sets the same length,
// and this repeated character is never put more than once in a row in the output.
func Replace(src string, from trset, to trset) string {
	srclen := len(src)
	if srclen == 0 || len(from) == 0 {
		return src
	}

	allbut := from[0] == '^'
	if allbut {
		from = from[1:]
	}

	si := 0
	for ; si < srclen; si++ {
		c := src[si]
		p := strings.IndexByte(string(from), c)
		if allbut == (p == -1) {
			break
		}
	}
	if si == srclen {
		return src // no changes
	}
	buf := strings.Builder{}
	buf.Grow(srclen)
	buf.WriteString(src[:si])

	lastto := len(to)
	collapse := lastto > 0 && (allbut || lastto < len(from))
	lastto--

scan:
	for ; si < srclen; si++ {
		c := src[si]
		i := xindex(from, c, allbut, lastto)
		if collapse && i >= lastto {
			buf.WriteByte(to[lastto])
			for {
				si++
				if si >= srclen {
					break scan
				}
				c = src[si]
				i = xindex(from, c, allbut, lastto)
				if i < lastto {
					break
				}
			}
		}
		if i < 0 {
			buf.WriteByte(c)
		} else if lastto >= 0 {
			buf.WriteByte(to[i])
		} /* else
		delete */
	}
	return buf.String()
}

func Set(s string) trset {
	if len(s) < 3 {
		return trset(s)
	}
	i := 0
	if s[0] == '^' {
		i++
	}
	dash := strings.IndexByte(s[i+1:len(s)-1], '-')
	if dash == -1 {
		return trset(s) // no ranges to expand
	}
	return expandRanges(s)
}

func expandRanges(s string) trset {
	slen := len(s)
	buf := strings.Builder{}
	buf.Grow(slen)
	if s[0] == '^' {
		buf.WriteByte('^')
		s = s[1:]
		slen--
	}
	for i := 0; i < slen; i++ {
		c := s[i]
		if c == '-' && i > 0 && i+1 < slen {
			for r := s[i-1] + 1; r < s[i+1]; r++ {
				buf.WriteByte(r)
			}
		} else {
			buf.WriteByte(c)
		}
	}
	return trset(buf.String())
}

func xindex(from trset, c byte, allbut bool, lastto int) int {
	i := strings.IndexByte(string(from), c)
	if allbut {
		if i == -1 {
			return lastto + 1
		}
		return -1
	}
	return i
}
