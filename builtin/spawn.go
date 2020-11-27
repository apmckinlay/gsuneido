// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("System(command)",
	func(arg Value) Value {
		shell, name, flag := getShell()
		if shell == "" {
			panic("System: can't get " + name)
		}
		cmd := exec.Command(shell, flag, ToStr(arg))
		if options.Mode != "gui" {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		err := cmd.Run()
		if err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				panic("System failed to run: " + err.Error())
			}
		}
		return IntVal(cmd.ProcessState.ExitCode())
	})

func getShell() (string, string, string) {
	if runtime.GOOS == "windows" {
		return os.Getenv("COMSPEC"), "COMSPEC", "/c"
	}
	return os.Getenv("SHELL"), "SHELL", "-c"
}

var _ = builtinRaw("Spawn(@args)",
	func(t *Thread, as *ArgSpec, rawargs []Value) Value {
		const wait = 0
		const nowait = 1
		iter := NewArgsIter(as, rawargs)
		var args []Value
		for k, v := iter(); k == nil && v != nil; k, v = iter() {
			args = append(args, v)
		}
		if len(args) < 2 {
			panic("usage: Spawn(mode, command, @args)")
		}
		mode := IfInt(args[0])
		if mode != wait && mode != nowait {
			panic("Spawn: bad mode")
		}
		exe := ToStr(args[1])
		argstr := make([]string, len(args)-2)
		for i, v := range args[2:] {
			argstr[i] = ToStr(v)
		}
		cmd := exec.Command(exe, argstr...)
		if options.Mode != "gui" {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		err := cmd.Start()
		if err != nil {
			log.Println("Spawn:", err)
			return IntVal(-1)
		}
		if mode == wait {
			cmd.Wait()
			return IntVal(cmd.ProcessState.ExitCode())
		}
		return IntVal(cmd.Process.Pid)
	})
