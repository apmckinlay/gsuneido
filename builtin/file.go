// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"bufio"
	"io"
	"os"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

type iFile interface {
	io.ReadWriteCloser
	io.Seeker
}

type suFile struct {
	CantConvert
	name string
	mode string
	f    iFile
	r    *bufio.Reader // only one of r or w will be used
	w    *bufio.Writer
	// tell is used to track our own position in the file.
	// We can't use f.Tell() because of buffering.
	// Any reads or writes must update this.
	// Not used for "a" (append) mode.
	tell int64
}

var nFile = 0

var _ = builtin("File(filename, mode='r', block=false)",
	func(t *Thread, args []Value) Value {
		name := ToStr(args[0])
		mode := ToStr(args[1])
		sf := newSuFile(name, mode)
		nFile++
		if args[2] == False {
			return sf
		}
		// block form
		defer sf.close()
		return t.Call(args[2], sf)
	})

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
	nFile--
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

func (*suFile) Get(*Thread, Value) Value {
	panic("File does not support get")
}

func (*suFile) Put(*Thread, Value, Value) {
	panic("File does not support put")
}

func (*suFile) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("File does not support update")
}

func (*suFile) RangeTo(int, int) Value {
	panic("File does not support range")
}

func (*suFile) RangeLen(int, int) Value {
	panic("File does not support range")
}

func (*suFile) Hash() uint32 {
	panic("File hash not implemented")
}

func (*suFile) Hash2() uint32 {
	panic("File hash not implemented")
}

func (*suFile) Compare(Value) int {
	panic("File compare not implemented")
}

func (*suFile) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call File")
}

func (sf *suFile) String() string {
	return "File(" + sf.name + ", " + sf.mode + ")"
}

func (*suFile) Type() types.Type {
	return types.File
}

func (sf *suFile) Equal(other interface{}) bool {
	sf2, ok := other.(*suFile)
	return ok && sf == sf2
}

func (*suFile) Lookup(_ *Thread, method string) Callable {
	return suFileMethods[method]
}

const MaxLine = 4000

var suFileMethods = Methods{
	"Close": method0(func(this Value) Value {
		sfOpen(this).close()
		return nil
	}),
	"Flush": method0(func(this Value) Value {
		err := sfOpenWrite(this).w.Flush()
		if err != nil {
			panic("File: " + err.Error())
		}
		return nil
	}),
	"Read": method1("(nbytes=false)", func(this, arg Value) Value {
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
			panic("file.Read " + err.Error())
		}
		return SuStr(string(buf))
	}),
	"Readline": method0(func(this Value) Value {
		sf := sfOpenRead(this)
		val, nr := readline(sf.r, "file.Readline: ")
		sf.tell += int64(nr)
		return val
	}),
	"Seek": method2("(offset, origin='set')", func(this, arg1, arg2 Value) Value {
		sf := sfOpen(this)
		if sf.mode == "a" {
			panic("file.Seek: invalid with mode 'a'")
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
			panic("file.Seek: origin must be 'set', 'end', or 'cur'")
		}
		_, err := sf.f.Seek(offset, io.SeekStart)
		if err != nil {
			panic("file.Seek: " + err.Error())
		}
		sf.tell = offset
		return nil
	}),
	"Tell": method0(func(this Value) Value {
		sf := sfOpen(this)
		if sf.mode == "a" {
			panic("file.Tell: invalid with mode 'a'")
		}
		return Int64Val(sf.tell)
	}),
	"Write": method1("(string)", func(this, arg Value) Value {
		s := AsStr(arg)
		sf := sfOpenWrite(this)
		_, err := sf.w.WriteString(s)
		if err != nil {
			panic("File: Write: " + err.Error())
		}
		sf.tell += int64(len(s))
		return arg
	}),
	"Writeline": method1("(string)", func(this, arg Value) Value {
		s := AsStr(arg)
		sf := sfOpen(this)
		sf.w.WriteString(s)
		_, err := sf.w.WriteString("\r\n")
		if err != nil {
			panic("File: Writeline: " + err.Error())
		}
		sf.tell += int64(len(s) + 2)
		return arg
	}),
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
