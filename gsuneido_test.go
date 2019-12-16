// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

// import (
// 	"runtime"
// 	"testing"

// 	"github.com/apmckinlay/gsuneido/builtin"
// 	_ "github.com/apmckinlay/gsuneido/builtin"
// 	"github.com/apmckinlay/gsuneido/database/dbms"
// 	. "github.com/apmckinlay/gsuneido/runtime"
// )

// func TestGsuneido(*testing.T) {
// 	runtime.LockOSThread()
// 	Global.Builtin("Suneido", new(SuObject))
// 	GetDbms = func() IDbms { return dbms.NewDbmsClient("127.0.0.1:3147") }
// 	Libload = libload // dependency injection
// 	mainThread = NewThread()
// 	defer mainThread.Close()
// 	builtin.UIThread = mainThread
// 	builtin.CmdlineOverride = "0" // prevent loading persistent windows
// 	eval("Init()")
// 	eval("CodeControl(); MessageLoop()")
// }
