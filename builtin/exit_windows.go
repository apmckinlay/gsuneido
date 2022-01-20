// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !portable && !com
// +build !portable,!com

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/exit"
)

var _ = builtin1("Exit(code = 0)",
	func(arg Value) Value {
		if arg == True {
			exit.Exit(0) // immediate exit
		}
		PostQuitMessage(uintptr(IfInt(arg)))
		return nil
	})
