// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
)

var _ = builtinRaw("Query1(@args)",
	func(t *Thread, as *ArgSpec, args []Value) Value {
		return queryOne(t, as, args, Only)
	})

var _ = builtinRaw("QueryFirst(@args)",
	func(t *Thread, as *ArgSpec, args []Value) Value {
		return queryOne(t, as, args, Next)
	})

var _ = builtinRaw("QueryLast(@args)",
	func(t *Thread, as *ArgSpec, args []Value) Value {
		return queryOne(t, as, args, Prev)
	})

var queryParams = params("(query)")

func queryOne(th *Thread, as *ArgSpec, args []Value, dir Dir) Value {
	query, _ := extractQuery(th, queryParams, as, args)
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

func init() {
	QueryMethods = Methods{
		"Close": method0(func(this Value) Value {
			this.(ISuQueryCursor).Close()
			return nil
		}),
		"Columns": method0(func(this Value) Value {
			return this.(ISuQueryCursor).Columns()
		}),
		"Explain": method0(func(this Value) Value { // deprecated
			return this.(ISuQueryCursor).Strategy()
		}),
		"Keys": method0(func(this Value) Value {
			return this.(ISuQueryCursor).Keys()
		}),
		"Next": method("()", func(th *Thread, this Value, _ []Value) Value {
			return this.(*SuQuery).GetRec(th, Next)
		}),
		"NewRecord": method1("(@args)", // deprecated
			func(_ Value, arg Value) Value {
				return SuRecordFromObject(arg.(*SuObject))
			}),
		"Prev": method("()", func(th *Thread, this Value, _ []Value) Value {
			return this.(*SuQuery).GetRec(th, Prev)
		}),
		"Output": method("(record)",
			func(th *Thread, this Value, args []Value) Value {
				trace.Dbms.Println("Query Output", this, args[0])
				this.(*SuQuery).Output(th, ToContainer(args[0]))
				return nil
			}),
		"Order": method0(func(this Value) Value {
			return this.(ISuQueryCursor).Order()
		}),
		"Rewind": method0(func(this Value) Value {
			this.(ISuQueryCursor).Rewind()
			return nil
		}),
		"RuleColumns": method0(func(this Value) Value {
			return this.(ISuQueryCursor).RuleColumns()
		}),
		"Strategy": method0(func(this Value) Value {
			return this.(ISuQueryCursor).Strategy()
		}),
	}
}
