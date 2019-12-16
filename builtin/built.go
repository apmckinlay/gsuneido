// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"runtime"

	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("Built()", func() Value {
	return SuStr(Built())
})

func Built() string {
	return options.BuiltDate + " (" + runtime.Version() + " " +
		runtime.GOARCH + " " + runtime.GOOS + ")"
}
