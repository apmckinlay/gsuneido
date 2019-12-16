// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin("Libraries()", func(t *Thread, args []Value) Value {
	return t.Dbms().Libraries()
})

var _ = builtin("Use(library)",
	func(t *Thread, args []Value) Value {
		return SuBool(t.Dbms().Use(ToStr(args[0])))
	})

var _ = builtin("Unuse(library)",
	func(t *Thread, args []Value) Value {
		return SuBool(t.Dbms().Unuse(ToStr(args[0])))
	})

var _ = builtin1("Unload(name = false)",
	func(arg Value) Value {
		if arg == False {
			Global.UnloadAll()
		} else {
			Global.Unload(ToStr(arg))
		}
		return nil
	})

var _ = builtin3("LibraryOverride(lib, name, text='')",
	func(lib, nameval, text Value) Value {
		name := ToStr(nameval)
		LibraryOverride(ToStr(lib), name, ToStr(text))
		return nil
	})

var _ = builtin0("LibraryOverrideClear()",
	func() Value {
		LibraryOverrideClear()
		return nil
	})
