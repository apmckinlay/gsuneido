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

func Key(rec Record, fields []int) string {
	if len(fields) == 0 {
		return ""
	}
	n := 2 * len(fields) // for separators (2 bytes extra)
	for _, field := range fields {
		n += len(rec.GetRaw(field))
	}
	n += n / 16 // allow for some escapes
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
