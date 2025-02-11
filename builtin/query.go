// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/dbms"
)

var _ = builtin(Query1, "(@args)")

func Query1(th *Thread, as *ArgSpec, args []Value) Value {
	return queryOne(th, as, args, Only)
}

var _ = builtin(QueryFirst, "(@args)")

func QueryFirst(th *Thread, as *ArgSpec, args []Value) Value {
	return queryOne(th, as, args, Next)
}

var _ = builtin(QueryLast, "(@args)")

func QueryLast(th *Thread, as *ArgSpec, args []Value) Value {
	return queryOne(th, as, args, Prev)
}

var queryParams = params("(query)")

func queryOne(th *Thread, as *ArgSpec, args []Value, dir Dir) Value {
	query, _ := extractQuery(th, &queryParams, as, args)
	row, hdr, table := th.Dbms().Get(th, query, dir)
	if hdr == nil {
		return False
	}
	return SuRecordFromRow(row, hdr, table, nil) // no transaction
}

// extractQuery does queryWhere then Args and returns the query and the args.
// NOTE: the base query must be the first argument
func extractQuery(
	th *Thread, ps *ParamSpec, as *ArgSpec, args []Value) (string, []Value) {
	where := queryWhere(as, args)
	args = th.Args(ps, as)
	query := AsStr(args[0])
	return query + where, args
}

// queryWhere builds a string of where's for the named arguments
// (except for 'block')
func queryWhere(as *ArgSpec, args []Value) string {
	var sb strings.Builder
	sep := "\nwhere "
	iter := NewArgsIter(as, args)
	for k, v := iter(); v != nil; k, v = iter() {
		if k == nil {
			continue
		}
		field := ToStr(k)
		if field == "query" || (field == "block" && !stringable(v)) {
			continue
		}
		sb.WriteString(sep)
		sep = "\nand "
		sb.WriteString(field)
		sb.WriteString(" is ")
		sb.WriteString(v.String())
	}
	return sb.String()
}

func stringable(v Value) bool {
	_, ok := v.AsStr()
	return ok
}

var _ = exportMethods(&QueryMethods, "query")

var _ = method(query_Close, "()")

func query_Close(this Value) Value {
	this.(ISuQueryCursor).Close()
	return nil
}

var _ = method(query_Columns, "()")

func query_Columns(this Value) Value {
	return this.(ISuQueryCursor).Columns()
}

var _ = method(query_Keys, "()")

func query_Keys(this Value) Value {
	return this.(ISuQueryCursor).Keys()
}

var _ = method(query_Next, "()")

func query_Next(th *Thread, this Value, _ []Value) Value {
	return this.(*SuQuery).GetRec(th, Next)
}

var _ = method(query_Prev, "()")

func query_Prev(th *Thread, this Value, _ []Value) Value {
	return this.(*SuQuery).GetRec(th, Prev)
}

var _ = method(query_Output, "(record)")

func query_Output(th *Thread, this Value, args []Value) Value {
	trace.Dbms.Println("Query Output", this, args[0])
	this.(*SuQuery).Output(th, ToContainer(args[0]))
	return nil
}

var _ = method(query_Order, "()")

func query_Order(this Value) Value {
	return this.(ISuQueryCursor).Order()
}

var _ = method(query_Rewind, "()")

func query_Rewind(this Value) Value {
	this.(ISuQueryCursor).Rewind()
	return nil
}

var _ = method(query_RuleColumns, "()")

func query_RuleColumns(this Value) Value {
	return this.(ISuQueryCursor).RuleColumns()
}

var _ = method(query_Strategy, "(formatted = false)")

func query_Strategy(_ *Thread, this Value, args []Value) Value {
	formatted := ToBool(args[0])
	return this.(ISuQueryCursor).Strategy(formatted)
}

var _ = method(query_Tree, "()")

func query_Tree(this Value) Value {
	return this.(ISuQueryCursor).Tree()
}

var _ = builtin(formatQuery, "(query)")

func formatQuery(th *Thread, args []Value) Value {
	if dbms, ok := th.Dbms().(*dbms.DbmsLocal); ok {
		return SuStr(dbms.FormatQuery(ToStr(args[0])))
	}
	return th.Dbms().Exec(th, SuObjectOf(SuStr("FormatQuery"), args[0]))
}
