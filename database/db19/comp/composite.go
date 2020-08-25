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
		if len(fields2) == 0 || fieldLen(rec, fields[0]) > 0 {
			return getRaw(rec, fields[0]) // don't need to encode single field keys
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
			n += fieldLen(rec, field)
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
		b := getRaw(rec, fields[f])
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

func fieldLen(rec Record, field int) int {
	if field < 0 {
		field = -field - 2 // _lower!
	}
	return len(rec.GetRaw(field))
}

func getRaw(rec Record, field int) string {
	if field >= 0 {
		return rec.GetRaw(field)
	}
	field = -field - 2 // _lower!
	return PackedToLower(rec.GetRaw(field))
}

// Compare compares the specified fields of the two records
// without building keys for them
func Compare(r1, r2 Record, fields, fields2 []int) int {
	empty := true
	for _, f := range fields {
		var x1,x2 string
		var cmp int
		if f < 0 { // _lower!
			f = -f - 2
			x1 = r1.GetRaw(f)
			x2 = r2.GetRaw(f)
			cmp = PackedCmpLower(x1, x2)
		} else {
			x1 = r1.GetRaw(f)
			x2 = r2.GetRaw(f)
			cmp = strings.Compare(x1, x2)
		}
		if cmp != 0 {
			return cmp
		}
		if x1 != "" || x2 != "" {
			empty = false
		}
	}
	if empty {
		for _, f := range fields2 {
			// NOTE: assumes fields2 will not be _lower!
			if cmp := strings.Compare(r1.GetRaw(f), r2.GetRaw(f)); cmp != 0 {
				return cmp
			}
		}
	}
	return 0
}
