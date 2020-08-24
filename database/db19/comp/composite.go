// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package comp implements encoding composite keys
// so they can be compared as strings of bytes.
// Fields are separated by two zero bytes 0,0
// Zero bytes are encoded as 0,1
package comp

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

// Key builds a composite key string that is comparable raw.
// fields2 is used for unique indexes.
// fields2 will only be added if all of the fields value are empty.
func Key(rec Record, fields, fields2 []int) string {
	if len(fields) == 0 {
		return ""
	}
	if len(fields) == 1 {
		x := rec.GetRaw(fields[0])
		if x != "" || len(fields2) == 0 {
			return x // don't need to encode single field keys
		}
	}
	n := 0
	lastNonEmpty := -1
	for i, field := range fields {
		fldlen := len(rec.GetRaw(field))
		if fldlen > 0 {
			lastNonEmpty = i
		}
		n += fldlen
	}
	if lastNonEmpty == -1 { // fields all empty
		if len(fields2) == 0 {
			return ""
		}
		for _, field := range fields2 {
			n += len(rec.GetRaw(field))
		}
	} else {
		fields = fields[:lastNonEmpty+1]
	}
	n += 2 * len(fields) // for separators (2 bytes extra)
	n += n / 16          // allow for some escapes
	buf := make([]byte, 0, n)
	if lastNonEmpty == -1 {
		for range fields {
			buf = append(buf, 0, 0) // separator
		}
		fields = fields2
	}
	for f := 0; ; {
		b := rec.GetRaw(fields[f])
		for len(b) > 0 {
			i := strings.IndexByte(b, 0)
			if i == -1 { // no zero bytes
				buf = append(buf, b...)
				break
			}
			// b[i] == 0
			i++
			buf = append(buf, b[:i]...) // copy up to and including zero
			buf = append(buf, 1)
			b = b[i:]
		}
		f++
		if f == len(fields) {
			break
		}
		buf = append(buf, 0, 0) // separator
	}
	return hacks.BStoS(buf)
}

// Compare compares the specified fields of the two records
// without building keys for them
func Compare(r1, r2 Record, fields, fields2 []int) int {
	empty := true
	for _, f := range fields {
		x1 := r1.GetRaw(f)
		x2 := r2.GetRaw(f)
		if cmp := strings.Compare(x1, x2); cmp != 0 {
			return cmp
		}
		if x1 != "" || x2 != "" {
			empty = false
		}
	}
	if empty {
		for _, f := range fields2 {
			if cmp := strings.Compare(r1.GetRaw(f), r2.GetRaw(f)); cmp != 0 {
				return cmp
			}
		}
	}
	return 0
}
