/*
Package tr is similar to the Unix tr program.
Originally based on the code in Software Tools by Kernighan and Plauger.
Currently it is byte based - it does not handle unicode.
*/
package tr

import (
	"bytes"
	"strings"

	"github.com/apmckinlay/gsuneido/util/ptest"
)

// Replace translates, squeezes, or deletes characters from the src string.
// If the first character of from is '^' then the from set is complemented.
// Ranges are specified with '-' between to characters.
// If the to set is shorter than the from set, then
// the last character in the to set is repeated to make the sets the same length,
// and this repeated character is never put more than once in a row in the output.
func Replace(src string, from string, to string) string {
	srclen := len(src)
	if srclen == 0 || len(from) == 0 {
		return src
	}

	allbut := from[0] == '^'
	if allbut {
		from = from[1:]
	}
	fromset := makset(from)

	si := 0
	for ; si < srclen; si++ {
		c := src[si]
		p := strings.IndexByte(fromset, c)
		if allbut == (p == -1) {
			break
		}
	}
	if si == srclen {
		return src // no changes
	}
	buf := new(bytes.Buffer)
	buf.Grow(srclen)
	buf.WriteString(src[:si])

	toset := makset(to)
	lastto := len(toset)
	collapse := lastto > 0 && (allbut || lastto < len(fromset))
	lastto--

scan:
	for ; si < srclen; si++ {
		c := src[si]
		i := xindex(fromset, c, allbut, lastto)
		if collapse && i >= lastto {
			buf.WriteByte(toset[lastto])
			for {
				si++
				if si >= srclen {
					break scan
				}
				c = src[si]
				i = xindex(fromset, c, allbut, lastto)
				if i < lastto {
					break
				}
			}
		}
		if i < 0 {
			buf.WriteByte(c)
		} else if lastto >= 0 {
			buf.WriteByte(toset[i])
		} /* else
		delete */
	}
	return buf.String()
}

func makset(s string) string {
	if len(s) < 3 {
		return s
	}
	dash := strings.IndexByte(s[1:len(s)-1], '-')
	if dash == -1 {
		return s // no ranges to expand
	}
	return expandRanges(s) // TODO cache
}

func expandRanges(s string) string {
	slen := len(s)
	buf := new(bytes.Buffer)
	buf.Grow(slen)
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
	return buf.String()
}

func xindex(fromset string, c byte, allbut bool, lastto int) int {
	i := strings.IndexByte(fromset, c)
	if allbut {
		if i == -1 {
			return lastto + 1
		} else {
			return -1
		}
	} else {
		return i
	}
}

// ptest support ---------------------------------------------------------------

// pt_tr is a ptest for matching
// usage: "string", "from", "to", "result"
func pt_replace(args []string) bool {
	return Replace(args[0], args[1], args[2]) == args[3]
}

var _ = ptest.Add("tr_replace", pt_replace)
