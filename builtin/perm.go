// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

// Perm is used by the application defined Auth function
// to set up permissions, one time only, e.g. the ServerEval whitelist.
type suPerm struct {
	staticClass[suPerm]
}

func init() {
	Global.Builtin("Permit", &suPerm{})
}

func (*suPerm) String() string {
	return "Perm /* builtin class */"
}

func (p *suPerm) Equal(other any) bool {
	return p == other
}

var permMethods = methods("perm")

// ServerEval adds name to the whitelist of function names
// that clients are allowed to execute on the server with ServerEval.
// It can only be called from the application's Auth function.
var _ = staticMethod(perm_ServerEval, "(name :string) :void")

func perm_ServerEval(th *Thread, args []Value) Value {
	perms := th.Perms()
	if perms == nil {
		panic("Perm.ServerEval can only be called from Auth")
	}
	perms.AddServerEval(ToStr(args[0]))
	return nil
}

var _ = staticMethod(perm_Table, "(table :string, rights :string) :void")

func perm_Table(th *Thread, args []Value) Value {
	perms := th.Perms()
	if perms == nil {
		panic("Perm.Table can only be called from Auth")
	}
	perms.AddTable(ToStr(args[0]), ToStr(args[1]))
	return nil
}

var _ = staticMethod(perm_Members, "() :object")

func perm_Members() Value {
	return methodList(permMethods)
}

func (*suPerm) Lookup(_ *Thread, method string) Value {
	return permMethods[method]
}
