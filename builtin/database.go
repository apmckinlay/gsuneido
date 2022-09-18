// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

type suDatabaseGlobal struct {
	SuBuiltin
}

func init() {
	Global.Builtin("Database", &suDatabaseGlobal{
		SuBuiltin{Fn: Database,
			BuiltinParams: BuiltinParams{ParamSpec: params("(string)")}}})
}

func Database(t *Thread, args []Value) Value {
	t.Dbms().Admin(ToStr(args[0]), nil)
	return nil
}

var databaseMethods = methods()

var _ = method(db_Auth, "(data)")

func db_Auth(t *Thread, this Value, args []Value) Value {
	return SuBool(t.Dbms().Auth(t, ToStr(args[0])))
}

var _ = method(db_Check, "()")

func db_Check(t *Thread, this Value, args []Value) Value {
	return SuStr(t.Dbms().Check())
}

var _ = method(db_Connections, "()")

func db_Connections(t *Thread, this Value, args []Value) Value {
	return t.Dbms().Connections()
}

var _ = method(db_CurrentSize, "()")

func db_CurrentSize(t *Thread, this Value, args []Value) Value {
	return IntVal(int(t.Dbms().Size()))
}

var _ = method(db_Cursors, "()")

func db_Cursors(t *Thread, this Value, args []Value) Value {
	return IntVal(t.Dbms().Cursors())
}

var _ = method(db_Dump, "(table = '')")

func db_Dump(t *Thread, this Value, args []Value) Value {
	return SuStr(t.Dbms().Dump(ToStr(args[0])))
}

var _ = method(db_Final, "()")

func db_Final(t *Thread, this Value, args []Value) Value {
	return IntVal(t.Dbms().Final())
}

var _ = method(db_Info, "()")

func db_Info(t *Thread, this Value, args []Value) Value {
	return t.Dbms().Info()
}

var _ = method(db_Kill, "(sessionId)")

func db_Kill(t *Thread, this Value, args []Value) Value {
	return IntVal(t.Dbms().Kill(ToStr(args[0])))
}

var _ = method(db_Load, "(table)")

func db_Load(t *Thread, this Value, args []Value) Value {
	return IntVal(t.Dbms().Load(ToStr(args[0])))
}

var _ = method(db_Nonce, "()")

func db_Nonce(t *Thread, this Value, args []Value) Value {
	return SuStr(t.Dbms().Nonce(t))
}

var _ = method(db_Schema, "(table)")

func db_Schema(t *Thread, this Value, args []Value) Value {
	return SuStr(t.Dbms().Schema(ToStr(args[0])))
}

var _ = method(db_SessionId, "(id = '')")

func db_SessionId(t *Thread, this Value, args []Value) Value {
	return SuStr(t.SessionId(ToStr(args[0])))
}

var _ = method(db_TempDest, "()")

func db_TempDest(Value) Value {
	return Zero
}

var _ = method(db_Token, "()")

func db_Token(t *Thread, this Value, args []Value) Value {
	return SuStr(t.Dbms().Token())
}

var _ = method(db_Transactions, "()")

func db_Transactions(t *Thread, this Value, args []Value) Value {
	return t.Dbms().Transactions()
}

func (d *suDatabaseGlobal) Get(t *Thread, key Value) Value {
	m := ToStr(key)
	if fn, ok := databaseMethods[m]; ok {
		return fn.(Value)
	}
	if fn, ok := ParamsMethods[m]; ok {
		return NewSuMethod(d, fn.(Value))
	}
	return nil
}

func (d *suDatabaseGlobal) Lookup(t *Thread, method string) Callable {
	if f, ok := databaseMethods[method]; ok {
		return f
	}
	return d.SuBuiltin.Lookup(t, method) // for Params
}

func (d *suDatabaseGlobal) String() string {
	return "Database /* builtin class */"
}

var _ = builtin(DoWithoutTriggers, "(tables, block)")

func DoWithoutTriggers(t *Thread, args []Value) Value {
	dbms := t.Dbms()
	ob := ToContainer(args[0])
	for i := ob.ListSize() - 1; i >= 0; i-- {
		table := ToStr(ob.ListGet(i))
		dbms.DisableTrigger(table)
		defer dbms.EnableTrigger(table)
	}
	return t.Call(args[1])
}
