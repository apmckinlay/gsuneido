// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable && !com

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/exit"
)

var _ = builtin(Exit, "(code = 0)")

func Exit(arg Value) Value {
	if arg == True {
		exit.Exit(0) // immediate exit
	}
	postQuit(uintptr(IfInt(arg)))
	return nil
}
