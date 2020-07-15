// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package comp implements encoding composite keys
// so they can be compared as strings of bytes.
// Fields are separated by two zero bytes 0,0
// Zero bytes are encoded as 0,1
package comp

import (
	"bytes"

	"github.com/apmckinlay/gsuneido/util/hacks"
)

func Encode(flds [][]byte) string {
	if len(flds) == 0 {
		return ""
	}
	// estimate result size
	n := 2 * len(flds) // for separators (2 bytes extra)
	for _, b := range flds {
		n += len(b)
	}
	n += n / 16 // allow for some escapes
	buf := make([]byte, 0, n)
	for f := 0; ; {
		b := flds[f]
		for len(b) > 0 {
			i := bytes.IndexByte(b, 0)
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
		if f == len(flds) {
			break
		}
		buf = append(buf, 0, 0) // separator
	}
	return hacks.BStoS(buf)
}
