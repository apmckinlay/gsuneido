// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"errors"
	"os"
	"os/exec"
	"runtime"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin(System, "(command)")

func System(arg Value) Value {
	shell, flag := getShell()
	command := ToStr(arg)
	cmd := exec.Command(shell)
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	if runtime.GOOS == "windows" {
		cmdSetup(cmd, shell+" "+flag+" "+command)
	} else {
		cmd.Args = []string{shell, flag, command}
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			panic("System failed to run: " + err.Error())
		}
	}
	return IntVal(cmd.ProcessState.ExitCode())
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
