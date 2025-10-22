// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && gui

package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(CreateLexer, "(name)")

func CreateLexer(a Value) Value {
	rtn := goc.CreateLexer(zstrArg(a))
	return intRet(rtn)
}
