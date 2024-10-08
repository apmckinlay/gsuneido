// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"
	"io"
	"sync"

	. "github.com/apmckinlay/gsuneido/core"
)

// writer must be implemented by destinations
type writer interface {
	writer() io.Writer
}

var bufs = sync.Pool{New: func() any {
	return new([rwsize]byte)
}}

const rwsize = 32 * 1024 // same as Go io.Copy

// CopyTo copies from src to to, up to nbytes or until src eof.
// Called by CopyTo in file, socket, and runpiped.
func CopyTo(th *Thread, src io.Reader, to, nbytes Value) Value {
	tow, ok := to.(writer)
	if !ok {
		panic("ERROR: CopyTo: can only copy to file, pipe, or socket")
	}
	dst := tow.writer()

	if _, ok = src.(io.WriterTo); !ok {
		fmt.Println("ERROR: CopyTo: src should have WriteTo")
	}
	if _, ok = dst.(io.ReaderFrom); !ok {
		fmt.Println("ERROR: CopyTo: dst should have ReadFrom")
	}

	var n int64
	if nbytes != False {
		n = ToInt64(nbytes)
		if n < 0 {
			panic("ERROR: CopyTo: nbytes cannot be negative")
		}
		src = io.LimitReader(src, int64(n))
	}

	array := bufs.Get().(*[rwsize]byte)
	defer bufs.Put(array)
	nw, err := io.CopyBuffer(dst, src, array[:]) // the actual copy
	if err != nil {
		panic("ERROR: CopyTo: " + err.Error())
	}
	if nbytes != False && nw != n {
		th.ReturnThrow = true
	}
	return Int64Val(nw)
}
