// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/util/hll"
	"github.com/apmckinlay/gsuneido/util/ss"
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

var databaseMethods = methods("db")

var _ = staticMethod(db_Auth, "(data)")

func db_Auth(th *Thread, args []Value) Value {
	return SuBool(th.Dbms().Auth(th, ToStr(args[0])))
}

var _ = staticMethod(db_Check, "()")

func db_Check(th *Thread, args []Value) Value {
	return SuStr(th.Dbms().Check(false))
}

var _ = staticMethod(db_FullCheck, "()")

func db_FullCheck(th *Thread, args []Value) Value {
	return SuStr(th.Dbms().Check(true))
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

var _ = staticMethod(db_Dump, "(table = '', to = '', publicKey = '')")

func db_Dump(th *Thread, args []Value) Value {
	if dbms, ok := th.Dbms().(*dbms.DbmsLocal); ok {
		err := dbms.Dump(ToStr(args[0]), ToStr(args[1]), ToStr(args[2]))
		if err != "" {
			th.ReturnThrow = true
			return SuStr(strings.Replace(err, "dump", "Database.Dump", 1))
		}
		return EmptyStr
	}
	return th.Dbms().Exec(th,
		SuObjectOf(SuStr("Database.Dump"), args[0], args[1], args[2]))
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

var _ = staticMethod(db_Load, "(table, from = '', privateKey = '', passphrase = '')")

func db_Load(th *Thread, args []Value) Value {
	if dbms, ok := th.Dbms().(*dbms.DbmsLocal); ok {
		return IntVal(dbms.Load(ToStr(args[0]), ToStr(args[1]), ToStr(args[2]), ToStr(args[3])))
	}
	return th.Dbms().Exec(th,
		SuObjectOf(SuStr("Database.Load"), args[0], args[1], args[2], args[3]))
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

var _ = staticMethod(db_CorruptedQ, "()")

func db_CorruptedQ(th *Thread, args []Value) Value {
	if dbms, ok := th.Dbms().(*dbms.DbmsLocal); ok {
		return SuBool(dbms.Corrupted())
	}
	return th.Dbms().Exec(th, SuObjectOf(SuStr("Database.Corrupted?")))
}

var _ = staticMethod(db_Top10, "(table, column)")

func db_Top10(th *Thread, args []Value) Value {
	table := ToStr(args[0])
	column := ToStr(args[1])

	tran := th.Dbms().Transaction(false)
	defer tran.Complete()

	sk := ss.New[string](128)
	q := tran.Query(table, nil)
	hdr := q.Header()
	for row, _ := q.Get(th, Next); row != nil; row, _ = q.Get(th, Next) {
		sk.Add(row.GetRawVal(hdr, column, nil, nil))
	}

	top := sk.Top()
	if len(top) > 10 {
		top = top[:10]
	}

	result := &SuObject{}
	for _, e := range top {
		result.Set(Unpack(e.Value), IntVal(int(e.Count-e.Error)))
	}
	return result
}

var _ = staticMethod(db_Distinct, "(table)")

func db_Distinct(th *Thread, args []Value) Value {
	table := ToStr(args[0])
	t := th.Dbms().Transaction(false)
	defer t.Complete()
	rt := t.(*dbms.ReadTranLocal).ReadTran
	cols := indexedColumns(rt.GetSchema(table).Indexes)
	hdr := SimpleHeader(cols)
	sketches := make([]*hll.HLL, len(cols))
	for i := range cols {
		sketches[i] = hll.New()
	}
	q := t.Query(table, nil)
	for row, _ := q.Get(th, Next); row != nil; row, _ = q.Get(th, Next) {
		for i, col := range cols {
			sketches[i].Add(row.GetRawVal(hdr, col, nil, nil))
		}
	}
	ob := &SuObject{}
	for i, col := range cols {
		ob.Set(SuStr(col), Int64Val(int64(sketches[i].Count())))
	}
	return ob
}

func indexedColumns(indexes []schema.Index) []string {
	cols := make([]string, 0, len(indexes))
	seen := make(map[string]struct{}, len(indexes))
	for _, ix := range indexes {
		for _, col := range ix.Columns {
			if strings.HasSuffix(col, "_lower!") {
				continue
			}
			if _, ok := seen[col]; ok {
				continue
			}
			seen[col] = struct{}{}
			cols = append(cols, col)
		}
	}
	return cols
}

var _ = staticMethod(db_Members, "()")

func db_Members() Value {
	return db_members
}

var db_members = methodList(databaseMethods)

//-------------------------------------------------------------------

func (d *suDatabaseGlobal) Lookup(th *Thread, method string) Value {
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
