// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package strs has miscellaneous functions for slices of strings
package strs

import "strings"

// Join joins strings with the specified format.
// The format may include delimiters e.g. "(,)"
func Join(fmt string, list []string) string {
	prefix := ""
	suffix := ""
	nf := len(fmt)
	if nf > 0 && (fmt[0] == '(' || fmt[0] == '{' || fmt[0] == '[') {
		prefix = fmt[0:1]
		suffix = fmt[nf-1:]
		fmt = fmt[1 : nf-1]
	}
	n := len(prefix) + len(suffix)
	sn := 0
	for _, s := range list {
		n += sn + len(s)
		sn = len(fmt)
	}
	sep := ""
	var sb strings.Builder
	sb.Grow(n)
	sb.WriteString(prefix)
	for _, s := range list {
		sb.WriteString(sep)
		sb.WriteString(s)
		sep = fmt
	}
	sb.WriteString(suffix)
	return sb.String()
}
