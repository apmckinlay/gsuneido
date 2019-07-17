package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

// TODO

func init() {
	Global.Builtin("SocketServer", &SuClass{Lib: "builtin", Name: "SocketServer"})
}
