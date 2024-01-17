// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"runtime/debug"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/options"
)

var _ = builtin(Built, "()")

func Built() Value {
	return SuStr(options.BuiltStr())
}

var _ = AddInfo("built", func() string {
	return options.BuiltStr()
})

var _ = AddInfo("build_info", func() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return bi.String()
})
