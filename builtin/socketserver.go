// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

// TODO

func init() {
	Global.Builtin("SocketServer", &SuClass{Lib: "builtin", Name: "SocketServer"})
}
