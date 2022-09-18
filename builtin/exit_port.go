// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !windows || portable || com

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/exit"
)

var _ = builtin(Exit, "(code = 0)")

func Exit(arg Value) Value {
	code := 0
	if arg != True {
		code = IfInt(arg)
	}
	exit.Exit(code)
	return nil
}
