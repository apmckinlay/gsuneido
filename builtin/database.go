// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
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
	th.Dbms().Admin(ToStr(args[0]), th.Sviews())
	return nil
}

var databaseMethods = methods()

var _ = staticMethod(db_Auth, "(data)")

func db_Auth(th *Thread, args []Value) Value {
	return SuBool(th.Dbms().Auth(th, ToStr(args[0])))
}

var _ = staticMethod(db_Check, "()")

func db_Check(th *Thread, args []Value) Value {
	return SuStr(th.Dbms().Check())
}

var _ = staticMethod(db_Connections, "()")

func db_Connections(th *Thread, args []Value) Value {
	return th.Dbms().Connections()
}

var _ = staticMethod(db_CurrentSize, "()")

func db_CurrentSize(th *Thread, args []Value) Value {
	return IntVal(int(th.Dbms().Size()))
}

var _ = staticMethod(db_Cursors, "()")

func db_Cursors(th *Thread, args []Value) Value {
	return IntVal(th.Dbms().Cursors())
}

var _ = staticMethod(db_Dump, "(table = '')")

func db_Dump(th *Thread, args []Value) Value {
	if err := th.Dbms().Dump(ToStr(args[0])); err != "" {
		panic(strings.Replace(err, "dump", "Database.Dump", 1))
	}
	return nil
}

var _ = staticMethod(db_Final, "()")

func db_Final(th *Thread, args []Value) Value {
	return IntVal(th.Dbms().Final())
}

var _ = staticMethod(db_Info, "()")

func db_Info(th *Thread, args []Value) Value {
	return th.Dbms().Info()
}

var _ = staticMethod(db_Kill, "(sessionId)")

func db_Kill(th *Thread, args []Value) Value {
	return IntVal(th.Dbms().Kill(ToStr(args[0])))
}

var _ = staticMethod(db_Load, "(table)")

func db_Load(th *Thread, args []Value) Value {
	return IntVal(th.Dbms().Load(ToStr(args[0])))
}

var _ = staticMethod(db_Nonce, "()")

func db_Nonce(th *Thread, args []Value) Value {
	return SuStr(th.Dbms().Nonce(th))
}

var _ = staticMethod(db_Schema, "(table)")

func db_Schema(th *Thread, args []Value) Value {
	return SuStr(th.Dbms().Schema(ToStr(args[0])))
}

var _ = staticMethod(db_SessionId, "(id = '')")

func db_SessionId(th *Thread, args []Value) Value {
	return SuStr(th.SessionId(ToStr(args[0])))
}

var _ = staticMethod(db_TempDest, "()")

func db_TempDest() Value {
	return Zero
}

var _ = staticMethod(db_Token, "()")

func db_Token(th *Thread, args []Value) Value {
	return SuStr(th.Dbms().Token())
}

var _ = staticMethod(db_Transactions, "()")

func db_Transactions(th *Thread, args []Value) Value {
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
