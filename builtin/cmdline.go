// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/options"
)

var _ = builtin(Cmdline, "()")

func Cmdline() Value {
	return SuStr(options.CmdLine)
}
