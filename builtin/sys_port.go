// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !windows || portable

package builtin

import (
	"os"
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
)

func Run() int {
	return 0
}

var _ = builtin(GetComputerName, "()")

func GetComputerName() Value {
	name, err := os.Hostname()
	if err != nil {
		panic("GetComputerName: " + err.Error())
	}
	return SuStr(name)
}

func CallbacksCount() int {
	return 0
}

func WndProcCount() int {
	return 0
}

func GetGuiResources() (int, int) {
	return 0, 0
}

func ErrlogDir() string {
	dir := strings.ReplaceAll(os.TempDir(), `\`, "/")
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	return dir
}

func OnUIThread() bool {
	return false
}
