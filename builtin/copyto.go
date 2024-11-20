// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"io"

	. "github.com/apmckinlay/gsuneido/core"
)

// writer must be implemented by destinations
type writer interface {
	writer() io.Writer
}

// CopyTo copies from src to to, up to nbytes or until src eof.
// Called by CopyTo in file, socket, and runpiped.
func CopyTo(th *Thread, src io.Reader, to, nbytes Value) Value {
	tow, ok := to.(writer)
	if !ok {
		panic("ERROR: CopyTo: can only copy to file, pipe, or socket")
	}
	dst := tow.writer()

	var n int64
	if nbytes != False {
		n = ToInt64(nbytes)
		if n < 0 {
			panic("ERROR: CopyTo: nbytes cannot be negative")
		}
		src = io.LimitReader(src, int64(n))
	}
	nw, err := dst.(io.ReaderFrom).ReadFrom(src)
	if err != nil {
		panic("ERROR: CopyTo: " + err.Error())
	}
	if nbytes != False && nw != n {
		th.ReturnThrow = true
	}
	return Int64Val(nw)
}
