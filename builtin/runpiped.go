package builtin

import (
	"io"
	"os"
	"os/exec"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

var _ = builtin("RunPiped(command, block=false)",
	func(t *Thread, args []Value) Value {
		command := ToStr(args[0])
		cmdargs := splitCommand(command)
		cmd := exec.Command(cmdargs[0], cmdargs[1:]...)
		in, err := cmd.StdinPipe()
		pr, pw, err := os.Pipe()
		if err != nil {
			panic("RunPiped create pipe failed: " + err.Error())
		}
		cmd.Stdout = pw
		cmd.Stderr = pw

		err = cmd.Start()
		if err != nil {
			panic("Runpiped failed to start: " + err.Error())
		}
		pw.Close()

		rp := &suRunPiped{command: command, cmd: cmd, w: in, r: pr}
		if args[2] == False {
			return rp
		}
		// block form
		defer rp.close()
		return t.CallWithArgs(args[1], rp)
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

type suRunPiped struct {
	CantConvert
	command string
	cmd     *exec.Cmd
	w       io.WriteCloser
	r       io.ReadCloser
}

func (rp *suRunPiped) close() {
	rp.r.Close()
}

// Value ------------------------------------------------------------

var _ Value = (*suRunPiped)(nil)

func (*suRunPiped) Get(*Thread, Value) Value {
	panic("RunPiped does not support get")
}

func (*suRunPiped) Put(*Thread, Value, Value) {
	panic("RunPiped does not support put")
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
	return 0 //TODO types.RunPiped
}

func (rp *suRunPiped) Equal(other interface{}) bool {
	if rp2, ok := other.(*suRunPiped); ok {
		return rp == rp2
	}
	return false
}

func (*suRunPiped) Lookup(_ *Thread, method string) Callable {
	return suRunPipedMethods[method]
}

var suRunPipedMethods = Methods{
	"Close": method0(func(this Value) Value {
		this.(*suRunPiped).close()
		return this
	}),
	"CloseWrite": method0(func(this Value) Value {
		this.(*suRunPiped).w.Close()
		return this
	}),
	"ExitValue": method0(func(this Value) Value {
		cmd := this.(*suRunPiped).cmd
		err := cmd.Wait()
		if err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				panic("System failed to run: " + err.Error())
			}
		}
		return IntVal(cmd.ProcessState.ExitCode())
	}),
	"Flush": method0(func(this Value) Value {
		return nil
	}),
	"Read": method1("(nbytes=1024)", func(this, arg Value) Value {
		f := this.(*suRunPiped).r
		n := IfInt(arg)
		buf := make([]byte, n)
		m, err := f.Read(buf)
		if err != nil {
			panic("runpiped.Read " + err.Error())
		}
		return SuStr(string(buf[:m]))
	}),
	"Readline": method0(func(this Value) Value {
		f := this.(*suRunPiped).r
		var buf strings.Builder
		b := make([]byte, 1)
		for {
			n, err := f.Read(b)
			if n == 0 {
				if buf.Len() == 0 {
					return False
				}
				break
			}
			if err != nil {
				panic("runpiped.Readline failed " + err.Error())
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
	}),
	"Write": method1("(string)", func(this, arg Value) Value {
		_, err := io.WriteString(this.(*suRunPiped).w, AsStr(arg))
		if err != nil {
			panic("runpiped.Write failed " + err.Error())
		}
		return arg
	}),
	"Writeline": method1("(string)", func(this, arg Value) Value {
		_, err := io.WriteString(this.(*suRunPiped).w, AsStr(arg)+"\n")
		if err != nil {
			panic("runpiped.Writeline failed " + err.Error())
		}
		return arg
	}),
}
