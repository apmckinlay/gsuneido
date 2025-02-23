// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/intern"
)

var _ = builtin(Libraries, "()")

func Libraries(th *Thread, args []Value) Value {
	list := th.Dbms().Libraries()
	var ob SuObject
	for _, s := range list {
		ob.Add(SuStr(s))
	}
	return &ob
}

var _ = builtin(Use, "(library)")

func Use(th *Thread, args []Value) Value {
	if !th.Dbms().Use(ToStr(args[0])) {
		return False
	}
	Global.UnloadAll()
	return True
}

var _ = builtin(Unuse, "(library)")

func Unuse(th *Thread, args []Value) Value {
	Global.UnloadAll()
	return SuBool(th.Dbms().Unuse(ToStr(args[0])))
}

var _ = builtin(Unload, "(name = false)")

func Unload(arg Value) Value {
	if arg == False {
		Global.UnloadAll()
	} else {
		Global.Unload(ToStr(arg))
	}
	return nil
}

var _ = builtin(LibraryOverride, "(lib, name, text='')")

func LibraryOverride(lib, name, text Value) Value {
	LibraryOverrides.Put(ToStr(lib), ToStr(name), ToStr(text))
	return nil
}

var _ = builtin(LibraryOverrideClear, "()")

func LibraryOverrideClear() Value {
	LibraryOverrides.Clear()
	return nil
}

var _ = AddInfo("intern.count", intern.Count)

var _ = AddInfo("intern.bytes", intern.Bytes)

var _ = AddInfo("intern.recent", intern.Recent)
