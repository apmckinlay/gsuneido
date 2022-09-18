// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"runtime"
	"runtime/debug"

	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin(Built, "()")

func Built() Value {
	return SuStr(BuiltStr())
}

func BuiltStr() string {
	return options.BuiltDate +
		" (" + runtime.Version() + " " + runtime.GOARCH + options.BuiltExtra + ")"
}

var _ = builtin(BuildInfo, "()")

func BuildInfo() Value {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return EmptyStr
	}
	return SuStr(bi.String())
}
