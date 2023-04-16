// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync/atomic"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

type iFile interface {
	io.ReadWriteCloser
	io.Seeker
}

type suFile struct {
	ValueBase[*suFile]
	f    iFile
	r    *bufio.Reader // only one of r or w will be used
	w    *bufio.Writer
	name string
	mode string
	// tell is used to track our own position in the file.
	// We can't use f.Tell() because of buffering.
	// Any reads or writes must update this.
	// Not used for "a" (append) mode.
	tell int64
}

var nFile atomic.Int32

var _ = builtin(File, "(filename, mode='r', block=false)")

func File(th *Thread, args []Value) Value {
	name := ToStr(args[0])
	mode := ToStr(args[1])
	sf := newSuFile(name, mode)
	nFile.Add(1)
	if args[2] == False {
		return sf
	}
	// block form
	defer sf.close()
	return th.Call(args[2], sf)
}

func newSuFile(name, mode string) *suFile {
	var flag int
	switch mode {
	case "r":
		flag = os.O_RDONLY
	case "a":
		flag = os.O_WRONLY | os.O_CREATE
	case "w":
		flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	default:
		panic("File: invalid mode")
	}
	var f iFile
	f, err := os.OpenFile(name, flag, 0644)
	if err != nil {
		panic("File: can't " + err.Error())
	}
	if mode == "a" {
		f = appender{f}
	}
	sf := &suFile{name: name, mode: mode, f: f}
	if sf.mode == "r" {
		sf.r = bufio.NewReader(f)
	} else { // "w" or "a"
		sf.w = bufio.NewWriter(f)
	}
	return sf
}

// reset is called by Seek to reset buffering
func (sf *suFile) reset() {
	if sf.mode == "r" {
		sf.r.Reset(sf.f)
	} else {
		err := sf.w.Flush()
		sf.w.Reset(sf.f)
		if err != nil {
			panic("File: " + err.Error())
		}
	}
}

func (sf *suFile) size() int64 {
	info, err := sf.f.(*os.File).Stat()
	if err != nil {
		panic("File: " + err.Error())
	}
	return info.Size()
}

func (sf *suFile) close() {
	nFile.Add(-1)
	if sf.mode != "r" {
		err := sf.w.Flush()
		if err != nil {
			defer panic("File: " + err.Error())
		}
	}
	err := sf.f.Close()
	sf.f = nil
	if err != nil {
		panic("File: " + err.Error())
	}
}

var _ Value = (*suFile)(nil)

func (sf *suFile) String() string {
	return "File(" + sf.name + ", " + sf.mode + ")"
}

func (*suFile) Type() types.Type {
	return types.File
}

func (sf *suFile) Equal(other any) bool {
	return sf == other
}

func (*suFile) Lookup(_ *Thread, method string) Callable {
	return suFileMethods[method]
}

const MaxLine = 4000

var suFileMethods = methods()

var _ = method(file_Close, "()")

func file_Close(this Value) Value {
	sfOpen(this).close()
	return nil
}

var _ = method(file_Flush, "()")

func file_Flush(this Value) Value {
	err := sfOpenWrite(this).w.Flush()
	if err != nil {
		panic("File: " + err.Error())
	}
	return nil
}

var _ = method(file_Read, "(nbytes=false)")

func file_Read(this, arg Value) Value {
	sf := sfOpenRead(this)
	n := int(sf.size() - sf.tell) // remaining
	if n == 0 {                   // at end
		return False
	}
	if arg != False {
		if m := ToInt(arg); m < n {
			n = m
		}
	}
	buf := make([]byte, n)
	_, err := io.ReadFull(sf.r, buf)
	sf.tell += int64(n)
	if err != nil {
		panic("File: Read: " + err.Error())
	}
	return SuStr(string(buf))
}

var _ = method(file_Readline, "()")

func file_Readline(this Value) Value {
	sf := sfOpenRead(this)
	val, nr := readline(sf.r, "File: Readline: ")
	sf.tell += int64(nr)
	return val
}

var _ = method(file_Seek, "(offset, origin='set')")

func file_Seek(this, arg1, arg2 Value) Value {
	sf := sfOpen(this)
	if sf.mode == "a" {
		panic("File: Seek: invalid with mode 'a'")
	}
	sf.reset()
	offset := ToInt64(arg1)
	switch ToStr(arg2) {
	case "set":
		//
	case "cur":
		offset += sf.tell
	case "end":
		offset += sf.size()
	default:
		panic("File: Seek: origin must be 'set', 'end', or 'cur'")
	}
	_, err := sf.f.Seek(offset, io.SeekStart)
	if err != nil {
		panic("File: Seek: " + err.Error())
	}
	sf.tell = offset
	return nil
}

var _ = method(file_Tell, "()")

func file_Tell(this Value) Value {
	sf := sfOpen(this)
	if sf.mode == "a" {
		panic("File: Tell: invalid with mode 'a'")
	}
	return Int64Val(sf.tell)
}

var _ = method(file_Write, "(string)")

func file_Write(this, arg Value) Value {
	s := AsStr(arg)
	sf := sfOpenWrite(this)
	_, err := sf.w.WriteString(s)
	if err != nil {
		panic("File: Write: " + err.Error())
	}
	sf.tell += int64(len(s))
	return arg
}

var _ = method(file_Writeline, "(string)")

func file_Writeline(this, arg Value) Value {
	s := AsStr(arg)
	sf := sfOpenWrite(this)
	sf.w.WriteString(s)
	_, err := sf.w.WriteString("\r\n")
	if err != nil {
		panic("File: Writeline: " + err.Error())
	}
	sf.tell += int64(len(s) + 2)
	return arg
}

func sfOpen(this Value) *suFile {
	sf := this.(*suFile)
	if sf.f == nil {
		panic("can't use a closed file: " + sf.name)
	}
	return sf
}

func sfOpenRead(this Value) *suFile {
	sf := sfOpen(this)
	if sf.mode != "r" {
		panic("File: can't read a file opened for writing")
	}
	return sf
}

func sfOpenWrite(this Value) *suFile {
	sf := sfOpen(this)
	if sf.mode == "r" {
		panic("File: can't write a file opened for reading")
	}
	return sf
}

func Readline(rdr io.Reader, errPrefix string) Value {
	val, _ := readline(rdr, errPrefix)
	return val
}

func readline(rdr io.Reader, errPrefix string) (Value, int) {
	nr := 0
	var buf strings.Builder
	b := make([]byte, 1)
	for {
		_, err := rdr.Read(b)
		if err == io.EOF {
			if buf.Len() == 0 {
				return False, nr
			}
			break
		}
		if err != nil {
			panic(errPrefix + err.Error())
		}
		nr++
		if b[0] == '\n' {
			break
		}
		if buf.Len() < MaxLine {
			buf.WriteByte(b[0])
		}
	}
	s := buf.String()
	s = strings.TrimRight(s, "\r")
	return SuStr(s), nr
}

// appender is a workaround for a Windows bug
// where normal append doesn't work on RDP shares.
// e.g. \\tsclient\C\...
type appender struct {
	f iFile
}

func (a appender) Write(buf []byte) (int, error) {
	_, err := a.f.Seek(0, io.SeekEnd)
	if err != nil {
		panic("File: " + err.Error())
	}
	return a.f.Write(buf)
}

func (a appender) Read([]byte) (int, error) {
	panic("appender Read not implemented")
}

func (a appender) Seek(int64, int) (int64, error) {
	panic("appender Seek not implemented")
}

func (a appender) Close() error {
	return a.f.Close()
}
