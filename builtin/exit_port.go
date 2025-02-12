// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !windows || portable

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(exit, "(code = 0)")

func exit(arg Value) Value {
	code := 0
	if arg != True {
		code = IfInt(arg)
	}
	Exit(code)
	return nil
}
