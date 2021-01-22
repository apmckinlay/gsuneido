// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"io"
	"os/exec"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

type suRunPiped struct {
	CantConvert
	command string
	cmd     *exec.Cmd
	w       io.WriteCloser
	r       io.ReadCloser
}

var nRunPiped = 0

var _ = builtin("RunPiped(command, block=false)",
	func(t *Thread, args []Value) Value {
		command := ToStr(args[0])
		cmdargs := splitCommand(command)
		cmd := exec.Command(cmdargs[0], cmdargs[1:]...)
		cmdSetup(cmd, command)
		w, err := cmd.StdinPipe()
		if err != nil {
			panic("RunPiped: create pipe failed: " + err.Error())
		}
		r, err := cmd.StdoutPipe()
		if err != nil {
			panic("RunPiped: create pipe failed: " + err.Error())
		}
		cmd.Stderr = cmd.Stdout

		err = cmd.Start()
		if err != nil {
			panic("Runpiped: failed to start: " + err.Error())
		}
		rp := &suRunPiped{command: command, cmd: cmd, w: w, r: r}
		nRunPiped++
		if args[1] == False {
			return rp
		}
		// block form
		defer rp.close()
		return t.Call(args[1], rp)
	})

func splitCommand(s string) []string {
	args := []string{}
	for {
		s = strings.TrimLeft(s, " \t")
		if s == "" {
			return args
		}
		delim := byte(' ')
		if s[0] == '"' {
			delim = '"'
			s = s[1:]
		}
		i := strings.IndexByte(s, delim)
		if i == -1 {
			return append(args, s)
		}
		args = append(args, s[:i])
		s = s[i+1:]
	}
}

func (rp *suRunPiped) close() {
	if rp.r == nil {
		return
	}
	nRunPiped--
	rp.r.Close()
	rp.r = nil
	if rp.w != nil {
		rp.w.Close()
		rp.w = nil
	}
	rp.cmd.Process.Release()
}

// Value ------------------------------------------------------------

var _ Value = (*suRunPiped)(nil)

func (*suRunPiped) Get(*Thread, Value) Value {
	panic("RunPiped does not support get")
}

func (*suRunPiped) Put(*Thread, Value, Value) {
	panic("RunPiped does not support put")
}

func (*suRunPiped) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("RunPiped does not support update")
}

func (*suRunPiped) RangeTo(int, int) Value {
	panic("RunPiped does not support range")
}

func (*suRunPiped) RangeLen(int, int) Value {
	panic("RunPiped does not support range")
}

func (*suRunPiped) Hash() uint32 {
	panic("RunPiped hash not implemented")
}

func (*suRunPiped) Hash2() uint32 {
	panic("RunPiped hash not implemented")
}

func (*suRunPiped) Compare(Value) int {
	panic("RunPiped compare not implemented")
}

func (*suRunPiped) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call RunPiped")
}

func (rp *suRunPiped) String() string {
	return "RunPiped(" + rp.command + ")"
}

func (*suRunPiped) Type() types.Type {
	return types.BuiltinClass
}

func (rp *suRunPiped) Equal(other interface{}) bool {
	rp2, ok := other.(*suRunPiped)
	return ok && rp == rp2
}

func (*suRunPiped) Lookup(_ *Thread, method string) Callable {
	return suRunPipedMethods[method]
}

var suRunPipedMethods = Methods{
	"Close": method0(func(this Value) Value {
		rpOpen(this).close()
		return this
	}),
	"CloseWrite": method0(func(this Value) Value {
		rp := rpOpen(this)
		rp.w.Close()
		rp.w = nil
		return this
	}),
	"ExitValue": method0(func(this Value) Value {
		rp := rpOpen(this)
		cmd := rp.cmd
		rp.r = nil
		rp.w = nil
		err := cmd.Wait()
		if err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				panic("RunPiped: ExitValue failed: " + err.Error())
			}
		}
		return IntVal(cmd.ProcessState.ExitCode())
	}),
	"Flush": method0(func(this Value) Value {
		return nil
	}),
	"Read": method1("(nbytes=1024)", func(this, arg Value) Value {
		f := rpOpen(this).r
		n := IfInt(arg)
		buf := make([]byte, n)
		m, err := f.Read(buf)
		if m == 0 || err == io.EOF {
			return False
		}
		if err != nil {
			panic("RunPiped: Read " + err.Error())
		}
		return SuStr(string(buf[:m]))
	}),
	"Readline": method0(func(this Value) Value {
		return Readline(rpOpen(this).r, "runPiped.Readline: ")
	}),
	"Write": method1("(string)", func(this, arg Value) Value {
		_, err := io.WriteString(rpWrite(this).w, AsStr(arg))
		if err != nil {
			panic("RunPiped: Write failed " + err.Error())
		}
		return arg
	}),
	"Writeline": method1("(string)", func(this, arg Value) Value {
		_, err := io.WriteString(rpWrite(this).w, AsStr(arg)+"\n")
		if err != nil {
			panic("RunPiped: Writeline failed " + err.Error())
		}
		return arg
	}),
}

func rpOpen(this Value) *suRunPiped {
	rp := this.(*suRunPiped)
	if rp.r == nil {
		panic("RunPiped: can't use after Close")
	}
	return rp
}

func rpWrite(this Value) *suRunPiped {
	rp := this.(*suRunPiped)
	if rp.w == nil {
		panic("RunPiped: can't use after Close")
	}
	return rp
}
