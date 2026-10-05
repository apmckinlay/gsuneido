// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

// Permit is used by the application defined Auth function
// to set up permissions, one time only, e.g. the ServerEval whitelist.
type suPermit struct {
	staticClass[suPermit]
}

func init() {
	Global.Builtin("Permit", &suPermit{})
}

func (*suPermit) String() string {
	return "Permit /* builtin class */"
}

func (p *suPermit) Equal(other any) bool {
	return p == other
}

var permMethods = methods("perm")

// Permit.ServerEval adds name to the whitelist of function names
// that clients are allowed to execute on the server with ServerEval.
// It can only be called from the application's Auth function.
var _ = staticMethod(perm_ServerEval, "(name :string) :void")

func perm_ServerEval(th *Thread, args []Value) Value {
	perms := th.NewPerms()
	if perms == nil {
		panic("Permit.ServerEval can only be called from Auth")
	}
	perms.AddServerEval(ToStr(args[0]))
	return nil
}

// Permit.Table sets the permissions for a table. The rights string can be
// "" (no permissions), "read" (read-only), or "write" (read/write).)
var _ = staticMethod(perm_Table, "(table :string, rights :string) :void")

func perm_Table(th *Thread, args []Value) Value {
	perms := th.NewPerms()
	if perms == nil {
		panic("Permit.Table can only be called from Auth")
	}
	perms.AddTable(ToStr(args[0]), ToStr(args[1]))
	return nil
}

// Permit.Schema sets permissions for schema modifications.
// The rights string can be "" (no permissions), "create" (tables or columns),
// or "update" (allows any modifications).
var _ = staticMethod(perm_Schema, "(rights :string) :void")

func perm_Schema(th *Thread, args []Value) Value {
	perms := th.NewPerms()
	if perms == nil {
		panic("Permit.Schema can only be called from Auth")
	}
	perms.SetSchema(ToStr(args[0]))
	return nil
}

var _ = staticMethod(perm_All, "() :void")

func perm_All(th *Thread, _ []Value) Value {
	perms := th.NewPerms()
	if perms == nil {
		panic("Permit.All can only be called from Auth")
	}
	perms.All()
	return nil
}

var _ = staticMethod(perm_None, "() :void")

func perm_None(th *Thread, _ []Value) Value {
	perms := th.NewPerms()
	if perms == nil {
		panic("Permit.None can only be called from Auth")
	}
	perms.None()
	return nil
}

var _ = staticMethod(perm_Members, "() :object")

func perm_Members() Value {
	return methodList(permMethods)
}

func (*suPermit) Lookup(_ *Thread, method string) Value {
	return permMethods[method]
}
