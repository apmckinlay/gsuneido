package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtinRaw("Query1(@args)",
	func(t *Thread, as *ArgSpec, args ...Value) Value {
		return queryOne("Query1", t, false, true, as, args...)
	})

var _ = builtinRaw("QueryFirst(@args)",
	func(t *Thread, as *ArgSpec, args ...Value) Value {
		return queryOne("QueryFirst", t, false, false, as, args...)
	})

var _ = builtinRaw("QueryLast(@args)",
	func(t *Thread, as *ArgSpec, args ...Value) Value {
		return queryOne("QueryLast", t, true, false, as, args...)
	})

const noTran = 0

func queryOne(which string, t *Thread, prev bool, single bool,
	as *ArgSpec, args ...Value) Value {
	query := buildQuery(which, as, args)
	row, hdr := t.Dbms().Get(noTran, query, prev, single)
	if hdr == nil {
		return False
	}
	return SuRecordFromRow(row, hdr, nil)
}

func buildQuery(which string, as *ArgSpec, args []Value) string {
	//TODO insert before "sort" or "into"
	iter := NewArgsIter(as, args)
	k, v := iter()
	if k != nil || v == nil {
		panic("usage: " + which + "(query, [field: value, ...])")
	}
	var sb strings.Builder
	sb.WriteString(IfStr(v))
	for {
		k, v := iter()
		if v == nil {
			break
		}
		if k == nil {
			if which == "Query" {
				continue // Query can have additional unnamed block argument
			}
			panic("usage: " + which + "(query, [field: value, ...])")
		}
		field := IfStr(k)
		if which == "Query" && field == "block" {
			continue
		}
		sb.WriteString("\nwhere ")
		sb.WriteString(field)
		sb.WriteString(" = ")
		sb.WriteString(v.String())
	}
	return sb.String()
}

func init() {
	QueryMethods = Methods{
		"Close": method0(func(this Value) Value {
			this.(*SuQuery).Close()
			return nil
		}),
		"Columns": method0(func(this Value) Value {
			return this.(*SuQuery).Columns()
		}),
		"Explain": method0(func(this Value) Value { // deprecated
			return this.(*SuQuery).Strategy()
		}),
		"Keys": method0(func(this Value) Value {
			return this.(*SuQuery).Keys()
		}),
		"Next": method0(func(this Value) Value {
			return this.(*SuQuery).GetRec(Next)
		}),
		"Prev": method0(func(this Value) Value {
			return this.(*SuQuery).GetRec(Prev)
		}),
		"Order": method0(func(this Value) Value {
			return this.(*SuQuery).Order()
		}),
		"Rewind": method0(func(this Value) Value {
			this.(*SuQuery).Rewind()
			return nil
		}),
		"RuleColumns": method0(func(this Value) Value {
			return this.(*SuQuery).RuleColumns()
		}),
		"Strategy": method0(func(this Value) Value {
			return this.(*SuQuery).Strategy()
		}),
	}
}
