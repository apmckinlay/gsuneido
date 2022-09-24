// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"os/exec"
	"syscall"
)

func cmdSetup(cmd *exec.Cmd, command string, inherit bool) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CmdLine:       command,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
		NoInheritHandles: !inherit,
	}
}
