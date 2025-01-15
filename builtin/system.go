// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"errors"
	"os"
	"os/exec"
	"runtime"

	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(System, "(command)")

func System(th *Thread, args []Value) Value {
	shell, flag := getShell()
	command := ToStr(args[0])
	cmd := exec.Command(shell)
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	if runtime.GOOS == "windows" {
		cmdSetup(cmd, shell+" "+flag+" "+command, InheritHandles)
	} else {
		cmd.Args = []string{shell, flag, command}
	}
	if InheritHandles {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			panic("System failed to run: " + err.Error())
		}
	}
	result := cmd.ProcessState.ExitCode()
	if result != 0 {
		th.ReturnThrow = true
	}
	return IntVal(result)
}

func getShell() (string, string) {
	var name, flag string
	if runtime.GOOS == "windows" {
		name, flag = "COMSPEC", "/c"
	} else {
		name, flag = "SHELL", "-c"
	}
	shell := os.Getenv(name)
	if shell == "" {
		panic("System: can't get " + name)
	}
	return shell, flag
}
