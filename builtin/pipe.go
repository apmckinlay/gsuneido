// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"
	"io"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

var _ = builtin(Pipe, "()")

// Pipe is a thin wrapper around the Go io.Pipe
func Pipe(th *Thread, _ []Value) Value {
	rd, wr := io.Pipe()
	// return rd, wr
	th.ReturnMulti = append(th.ReturnMulti[:0],
		&suPipeWriter{wr: wr}, &suPipeReader{rd: rd})
	return nil
}

//-------------------------------------------------------------------

const readMax = 64 * 1024 // no advantage to larger

var suPipeReaderMethods = methods("piper")

var _ = method(piper_Read, "(n)")

func piper_Read(this Value, a Value) Value {
	rd := this.(*suPipeReader).rd
	n := ToInt(a)
	if n > readMax {
		panic("Pipe.Read: too large")
	}
	buf := make([]byte, n)
	nr, err := rd.Read(buf)
	if nr > 0 {
		return SuStr(hacks.BStoS(buf[:nr]))
	}
	if err != nil {
		if err == io.EOF {
			return False
		}
		panic(fmt.Sprint("Pipe.Read: ", err))
	}
	return SuStr("")
}

var _ = method(piper_CopyTo, "(dest, nbytes = false)")

func piper_CopyTo(th *Thread, this Value, args []Value) Value {
	rd := this.(*suPipeReader).rd
	return CopyTo(th, rd, args[0], args[1])
}

var _ = method(piper_Close, "()")

func piper_Close(this Value) Value {
	rd := this.(*suPipeReader).rd
	err := rd.Close()
	if err != nil {
		panic(fmt.Sprint("Pipe.Close: ", err))
	}
	return nil
}

//-------------------------------------------------------------------

var suPipeWriterMethods = methods("pipew")

var _ = method(pipew_Write, "(s)")

func pipew_Write(this Value, a Value) Value {
	wr := this.(*suPipeWriter).wr
	s := ToStr(a)
	_, err := wr.Write(hacks.Stobs(s))
	if err != nil {
		panic(fmt.Sprint("Pipe.Write: ", err))
	}
	return nil
}

var _ = method(pipew_Close, "()")

func pipew_Close(this Value) Value {
	wr := this.(*suPipeWriter).wr
	err := wr.Close()
	if err != nil {
		panic(fmt.Sprint("Pipe.Close: ", err))
	}
	return nil
}

//-------------------------------------------------------------------

// @immutable
type suPipeReader struct {
	ValueBase[suPipeReader]
	rd *io.PipeReader
}

var _ Value = (*suPipeReader)(nil)

func (spr *suPipeReader) Equal(other any) bool {
	return spr == other
}

func (spr *suPipeReader) Lookup(_ *Thread, method string) Value {
	return suPipeReaderMethods[method]
}

func (spr *suPipeReader) SetConcurrent() {
	// safe
}

//-------------------------------------------------------------------

// @immutable
type suPipeWriter struct {
	ValueBase[suPipeWriter]
	wr *io.PipeWriter
}

var _ Value = (*suPipeWriter)(nil)

func (spw *suPipeWriter) Equal(other any) bool {
	return spw == other
}

func (spw *suPipeWriter) Lookup(_ *Thread, method string) Value {
	return suPipeWriterMethods[method]
}

func (spw *suPipeWriter) SetConcurrent() {
	// safe
}

// writer is for CopyTo
func (spw *suPipeWriter) writer() io.Writer {
	return spw.wr
}
