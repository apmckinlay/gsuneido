// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"errors"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"sync/atomic"

	. "github.com/apmckinlay/gsuneido/runtime"
)

type suRunPiped struct {
	ValueBase[*suRunPiped]
	w       io.WriteCloser
	r       io.ReadCloser
	cmd     *exec.Cmd
	command string
}

var nRunPiped atomic.Int32
var _ = AddInfo("builtin.nRunPiped", &nRunPiped)

var _ = builtin(RunPiped, "(command, block=false)")

func RunPiped(th *Thread, args []Value) Value {
	command := ToStr(args[0])
	cmdargs := splitCommand(command)
	cmd := exec.Command(cmdargs[0])
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	if runtime.GOOS == "windows" {
		cmdSetup(cmd, command, true)
	} else {
		cmd.Args = cmdargs
	}
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
	nRunPiped.Add(1)
	if args[1] == False {
		return rp
	}
	// block form
	defer rp.close()
	return th.Call(args[1], rp)
}

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
	nRunPiped.Add(-1)
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

func (rp *suRunPiped) String() string {
	return "RunPiped(" + rp.command + ")"
}

func (rp *suRunPiped) Equal(other any) bool {
	return rp == other
}

func (*suRunPiped) SetConcurrent() {
	//FIXME concurrency
	// panic("RunPiped cannot be set to concurrent")
}

func (*suRunPiped) Lookup(_ *Thread, method string) Callable {
	return suRunPipedMethods[method]
}

var suRunPipedMethods = methods()

var _ = method(runpiped_Close, "()")

func runpiped_Close(this Value) Value {
	rpOpen(this).close()
	return this
}

var _ = method(runpiped_CloseWrite, "()")

func runpiped_CloseWrite(this Value) Value {
	rp := rpOpen(this)
	rp.w.Close()
	rp.w = nil
	return this
}

var _ = method(runpiped_ExitValue, "()")

func runpiped_ExitValue(this Value) Value {
	rp := rpOpen(this)
	defer rp.close()
	cmd := rp.cmd
	err := cmd.Wait()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			panic("RunPiped: ExitValue failed: " + err.Error())
		}
	}
	return IntVal(cmd.ProcessState.ExitCode())
}

var _ = method(runpiped_Flush, "()")

func runpiped_Flush(this Value) Value {
	return nil
}

var _ = method(runpiped_Read, "(nbytes=1024)")

func runpiped_Read(this, arg Value) Value {
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
}

var _ = method(runpiped_Readline, "()")

func runpiped_Readline(this Value) Value {
	return Readline(rpOpen(this).r, "runPiped.Readline: ")
}

var _ = method(runpiped_Write, "(string)")

func runpiped_Write(this, arg Value) Value {
	_, err := io.WriteString(rpWrite(this).w, AsStr(arg))
	if err != nil {
		panic("RunPiped: Write failed " + err.Error())
	}
	return arg
}

var _ = method(runpiped_Writeline, "(string)")

func runpiped_Writeline(this, arg Value) Value {
	w := rpWrite(this).w
	_, err := io.WriteString(w, AsStr(arg)+newline)
	if err != nil {
		panic("RunPiped: Writeline failed " + err.Error())
	}
	return arg
}

var newline = func() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}()

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
