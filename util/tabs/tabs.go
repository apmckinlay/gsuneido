// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tabs

import (
	"strings"
)

const TabWidth = 4

// Detab converts tabs to spaces
func Detab(s string) string {
	if -1 == strings.IndexByte(s, '\t') {
		return s
	}
	n := len(s)
	var sb strings.Builder
	sb.Grow(n + n/4)
	col := 0
	for i := range n {
		switch s[i] {
		case '\t':
			sb.WriteByte(' ')
			for col++; (col % TabWidth) != 0; col++ {
				sb.WriteByte(' ')
			}
		case '\r', '\n':
			col = -1
			fallthrough
		default:
			sb.WriteByte(s[i])
			col++
		}
	}
	return sb.String()
}

// Entab converts leading spaces & tabs to tabs
// and strips trailing spaces & tabs
func Entab(s string) string {
	i := 0
	n := len(s)
	if n == 0 {
		return s
	}
	var sb strings.Builder
	sb.Grow(n)
	for i < n { // for each line
		// find the indent
		col := 0
		for ; i < n; i++ {
			if s[i] == ' ' { //nolint
				col++
			} else if s[i] == '\t' {
				for col++; (col % TabWidth) != 0; col++ {
				}
			} else {
				break
			}
		}
		if i < n && s[i] != '\r' && s[i] != '\n' {
			// output the indent using tabs possibly followed by space
			dstcol := 0
			for ; dstcol+TabWidth <= col; dstcol += TabWidth {
				sb.WriteByte('\t')
			}
			for ; dstcol < col; dstcol++ {
				sb.WriteByte(' ')
			}
			// find the end of the current line
			end := strings.IndexAny(s[i:], "\r\n")
			if end == -1 {
				end = n - i
			}
			// back up over trailing spaces and tabs
			j := i + end
			for ; j-1 > i && (s[j-1] == ' ' || s[j-1] == '\t'); j-- {
			}
			// copy the line
			sb.WriteString(s[i:j])
			i += end
		}
		// copy newlines
		for ; i < n && (s[i] == '\r' || s[i] == '\n'); i++ {
			sb.WriteByte(s[i])
		}
	}
	return sb.String()
}
