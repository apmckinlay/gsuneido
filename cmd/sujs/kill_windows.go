// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows

package main

import (
	"fmt"
	"os/exec"
)

func setSysProcAttr(_ *exec.Cmd) {
	// no-op: Windows does not support Setpgid
}

func killProcessGroup(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	// /T kills the process tree; /F forces termination
	_ = exec.Command("taskkill", "/T", "/F", "/PID",
		fmt.Sprint(cmd.Process.Pid)).Run()
}
