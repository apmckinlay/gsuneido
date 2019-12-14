package builtin

import (
	"os/exec"
	"syscall"
)

func cmdSetup(cmd *exec.Cmd, command string) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CmdLine:       command,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
}
