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
	"github.com/apmckinlay/gsuneido/util/verify"
)

// Key builds a composite key string that is comparable raw.
// If ui (unique index) is true
// then the final field will only be added if the other fields are all empty.
func Key(rec Record, fields []int, ui bool) string {
	verify.That(ui == false || len(fields) >= 2)
	if len(fields) == 0 {
		return ""
	}
	if len(fields) == 1 {
		// don't need to encode single field keys
		return rec.GetRaw(fields[0])
	}
	n := 0
	lastNonEmpty := 0
	for i, field := range fields {
		if ui && i == len(fields)-1 && n > 0 {
			break
		}
		fldlen := len(rec.GetRaw(field))
		if fldlen > 0 {
			lastNonEmpty = i
		}
		n += fldlen
	}
	fields = fields[:lastNonEmpty+1]
	n += 2 * len(fields) // for separators (2 bytes extra)
	n += n / 16          // allow for some escapes
	buf := make([]byte, 0, n)
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
func Compare(r1, r2 Record, fields []int) int {
	for _, f := range fields {
		if cmp := strings.Compare(r1.GetRaw(f), r2.GetRaw(f)); cmp != 0 {
			return cmp
		}
	}
	return 0
}
