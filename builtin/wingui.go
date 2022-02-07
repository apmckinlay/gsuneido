// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable && !com

package builtin

import (
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

// Startup for Windows gui
// relaunches the process with stdout and stderr redirected to Errlog.
// This is the only way I found to capture built-in output like crashes.
// Note: This breaks start/w
func Startup() {
	if options.NoRelaunch || options.Redirected() {
		return // to avoid infinite loop
	}
	f, err := os.OpenFile(options.Errlog,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		Fatal(err.Error())
	}
	path, _ := os.Executable()
	var cmd exec.Cmd
	cmd.Path = path
	// avoid splitting and joining command line arguments
	// with potential quoting issues
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: GetCommandLine()}
	cmd.Stdout = f
	cmd.Stderr = f
	err = cmd.Start()
	if err != nil {
		Fatal(err.Error())
	}
	os.Exit(0)
}

var getCommandLine = kernel32.MustFindProc("GetCommandLineA").Addr()

func GetCommandLine() string {
	rtn, _, _ := syscall.Syscall(getCommandLine, 0, 0, 0, 0)
	return string(bufStrZ(unsafe.Pointer(rtn), 4096).(SuStr))
}
