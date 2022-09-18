// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"errors"
	"log"
	"os"
	"os/exec"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin(Spawn, "(@args)")

func Spawn(t *Thread, as *ArgSpec, rawargs []Value) Value {
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
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	cmdSetup(cmd, "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Println("ERROR: Spawn:", err)
		return IntVal(-1)
	}
	if mode == wait {
		cmd.Wait()
		return IntVal(cmd.ProcessState.ExitCode())
	}
	return IntVal(cmd.Process.Pid)
}
