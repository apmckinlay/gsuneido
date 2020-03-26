// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"io"
	"os"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

type suFile struct {
	CantConvert
	name string
	mode string
	f    *os.File
}

var nFile = 0

var _ = builtin("File(filename, mode='r', block=false)",
	func(t *Thread, args []Value) Value {
		name := ToStr(args[0])
		mode := ToStr(args[1])

		f, err := os.OpenFile(name, modeToFlags(mode), 0644)
		if err != nil {
			panic("File: can't " + err.Error())
		}
		sf := &suFile{name: name, mode: mode, f: f}
		nFile++
		if args[2] == False {
			return sf
		}
		// block form
		defer sf.close()
		return t.Call(args[2], sf)
	})

func modeToFlags(mode string) int {
	switch mode {
	case "r":
		return os.O_RDONLY
	case "a":
		return os.O_RDWR | os.O_CREATE | os.O_APPEND
	case "w":
		return os.O_RDWR | os.O_CREATE | os.O_TRUNC
	default:
		panic("File: invalid mode")
	}
}

func (sf *suFile) close() {
	nFile--
	sf.f.Close()
	sf.f = nil
}

var _ Value = (*suFile)(nil)

func (*suFile) Get(*Thread, Value) Value {
	panic("File does not support get")
}

func (*suFile) Put(*Thread, Value, Value) {
	panic("File does not support put")
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
		// no buffering so nothing to do
		sfOpen(this).f.Sync()
		return nil
	}),
	"Read": method1("(nbytes=false)", func(this, arg Value) Value {
		f := sfOpen(this).f
		pos, _ := f.Seek(0, io.SeekCurrent)
		info, _ := f.Stat()
		n := int(info.Size() - pos)
		if n == 0 {
			return False
		}
		if arg != False {
			if m := ToInt(arg); m < n {
				n = m
			}
		}
		buf := make([]byte, n)
		_, err := io.ReadFull(f, buf)
		if err != nil {
			panic("file.Read " + err.Error())
		}
		return SuStr(string(buf))
	}),
	"Readline": method0(func(this Value) Value {
		return Readline(sfOpen(this).f, "file.Readline: ")
	}),
	"Seek": method2("(offset, origin='set')", func(this, arg1, arg2 Value) Value {
		offset := ToInt64(arg1)
		if offset < 0 {
			offset = 0
		}
		var whence int
		switch ToStr(arg2) {
		case "cur":
			whence = io.SeekCurrent
		case "set":
			whence = io.SeekStart
		case "end":
			whence = io.SeekEnd
		default:
			panic("file.Seek: origin must be 'set', 'end', or 'cur'")
		}
		_, err := sfOpen(this).f.Seek(offset, whence)
		if err != nil {
			panic("file.Seek " + err.Error())
		}
		return nil
	}),
	"Tell": method0(func(this Value) Value {
		pos, _ := sfOpen(this).f.Seek(0, io.SeekCurrent)
		return Int64Val(pos)
	}),
	"Write": method1("(string)", func(this, arg Value) Value {
		sfOpen(this).f.WriteString(AsStr(arg))
		return arg
	}),
	"Writeline": method1("(string)", func(this, arg Value) Value {
		f := sfOpen(this).f
		f.WriteString(AsStr(arg))
		f.WriteString("\n")
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

func Readline(rdr io.Reader, errPrefix string) Value {
	var buf strings.Builder
	b := make([]byte, 1)
	for {
		n, err := rdr.Read(b)
		if n == 0 {
			if buf.Len() == 0 {
				return False
			}
			break
		}
		if err != nil {
			panic(errPrefix + err.Error())
		}
		if b[0] == '\n' {
			break
		}
		if buf.Len() < MaxLine {
			buf.WriteByte(b[0])
		}
	}
	s := buf.String()
	s = strings.TrimRight(s, "\r")
	return SuStr(s)
}
