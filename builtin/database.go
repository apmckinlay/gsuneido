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

func Database(th *Thread, args []Value) Value {
	th.Dbms().Admin(ToStr(args[0]), nil)
	return nil
}

var databaseMethods = methods()

var _ = method(db_Auth, "(data)")

func db_Auth(th *Thread, this Value, args []Value) Value {
	return SuBool(th.Dbms().Auth(th, ToStr(args[0])))
}

var _ = method(db_Check, "()")

func db_Check(th *Thread, this Value, args []Value) Value {
	return SuStr(th.Dbms().Check())
}

var _ = method(db_Connections, "()")

func db_Connections(th *Thread, this Value, args []Value) Value {
	return th.Dbms().Connections()
}

var _ = method(db_CurrentSize, "()")

func db_CurrentSize(th *Thread, this Value, args []Value) Value {
	return IntVal(int(th.Dbms().Size()))
}

var _ = method(db_Cursors, "()")

func db_Cursors(th *Thread, this Value, args []Value) Value {
	return IntVal(th.Dbms().Cursors())
}

var _ = method(db_Dump, "(table = '')")

func db_Dump(th *Thread, this Value, args []Value) Value {
	return SuStr(th.Dbms().Dump(ToStr(args[0])))
}

var _ = method(db_Final, "()")

func db_Final(th *Thread, this Value, args []Value) Value {
	return IntVal(th.Dbms().Final())
}

var _ = method(db_Info, "()")

func db_Info(th *Thread, this Value, args []Value) Value {
	return th.Dbms().Info()
}

var _ = method(db_Kill, "(sessionId)")

func db_Kill(th *Thread, this Value, args []Value) Value {
	return IntVal(th.Dbms().Kill(ToStr(args[0])))
}

var _ = method(db_Load, "(table)")

func db_Load(th *Thread, this Value, args []Value) Value {
	return IntVal(th.Dbms().Load(ToStr(args[0])))
}

var _ = method(db_Nonce, "()")

func db_Nonce(th *Thread, this Value, args []Value) Value {
	return SuStr(th.Dbms().Nonce(th))
}

var _ = method(db_Schema, "(table)")

func db_Schema(th *Thread, this Value, args []Value) Value {
	return SuStr(th.Dbms().Schema(ToStr(args[0])))
}

var _ = method(db_SessionId, "(id = '')")

func db_SessionId(th *Thread, this Value, args []Value) Value {
	return SuStr(th.SessionId(ToStr(args[0])))
}

var _ = method(db_TempDest, "()")

func db_TempDest(Value) Value {
	return Zero
}

var _ = method(db_Token, "()")

func db_Token(th *Thread, this Value, args []Value) Value {
	return SuStr(th.Dbms().Token())
}

var _ = method(db_Transactions, "()")

func db_Transactions(th *Thread, this Value, args []Value) Value {
	return th.Dbms().Transactions()
}

func (d *suDatabaseGlobal) Get(th *Thread, key Value) Value {
	m := ToStr(key)
	if fn, ok := databaseMethods[m]; ok {
		return fn.(Value)
	}
	if fn, ok := ParamsMethods[m]; ok {
		return NewSuMethod(d, fn.(Value))
	}
	return nil
}

func (d *suDatabaseGlobal) Lookup(th *Thread, method string) Callable {
	if f, ok := databaseMethods[method]; ok {
		return f
	}
	return d.SuBuiltin.Lookup(th, method) // for Params
}

func (d *suDatabaseGlobal) String() string {
	return "Database /* builtin class */"
}

var _ = builtin(DoWithoutTriggers, "(tables, block)")

func DoWithoutTriggers(th *Thread, args []Value) Value {
	dbms := th.Dbms()
	ob := ToContainer(args[0])
	for i := ob.ListSize() - 1; i >= 0; i-- {
		table := ToStr(ob.ListGet(i))
		dbms.DisableTrigger(table)
		defer dbms.EnableTrigger(table)
	}
	return th.Call(args[1])
}
