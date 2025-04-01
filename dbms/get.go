// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/util/generic/set"
)

func get(th *Thread, tran qry.QueryTran, args Value, dir Dir) (Row, *Header, string) {
	defer th.Suneido.Store(th.Suneido.Load())
	th.Suneido.Store(nil) // use main Suneido object

	ob := args.(*SuObject)
	query := getQuery(ob)
	if row, hdr, tbl := fastGet(th, tran, query, ob, dir); row != nil {
		return row, hdr, tbl
	}
	query += getWhere(ob)

	q := qry.ParseQuery(query, tran, th.Sviews())
	qs, sorted := q.(*qry.Sort)
	if dir == Only || dir == Any {
		if sorted {
			q = qs.Source() // remove sort
		}
	} else if !sorted &&
		!strings.Contains(query, "CHECKQUERY SUPPRESS: SORT REQUIRED") {
		panic("QueryFirst and QueryLast require sort")
	}

	q, fixcost, varcost := qry.Setup1(q, qry.ReadMode, tran)
	qry.Warnings(query, q)
	if trace.Query.On() {
		d := map[Dir]string{Only: "one", Next: "first", Prev: "last", Any: "any"}[dir]
		trace.Query.Println(d, fixcost+varcost, "-", query)
	}
	d := dir
	if dir == Only || dir == Any {
		d = Next
	}
	row := q.Get(th, d)
	if row == nil {
		return nil, nil, ""
	} else if dir == Any {
		return exists, nil, ""
	}
	if dir == Only && !single(q) && q.Get(th, Next) != nil {
		panic("Query1 not unique: " + query)
	}
	return row, q.Header(), q.Updateable()
}

var exists Row = []DbRec{{Record: "x"}}

func getQuery(ob *SuObject) string {
	if ob.ListSize() >= 1 {
		return ToStr(ob.ListGet(0))
	} else if q := ob.NamedGet(SuStr("query")); q != nil {
		return ToStr(q)
	}
	return ""
}

func fastGet(th *Thread, tran qry.QueryTran, query string, ob *SuObject, dir Dir) (Row, *Header, string) {
	if dir != Only {
		return nil, nil, ""
	}
	if strings.Contains(query, " ") || tran.GetInfo(query) == nil {
		return nil, nil, ""
	}
	table, ok := qry.NewTable(tran, query).(*qry.Table)
	if !ok {
		return nil, nil, ""
	}
	flds := make([]string, 0, ob.NamedSize())
	vals := make([]Value, 0, ob.NamedSize())
	iter := ob.Iter2(false, true)
	for k, v := iter(); v != nil; k, v = iter() {
		field := ToStr(k)
		if field == "query" {
			continue
		}
		flds = append(flds, field)
		vals = append(vals, v)
	}
	key := findKey(table, flds)
	if key == nil {
		return nil, nil, ""
	}
	if len(key) == 0 {
		row := table.Get(th, Next)
		return row, table.Header(), query
	}
	table.SetIndex(key)
	packed := make([]string, len(vals))
	for i, v := range vals {
		packed[i] = Pack(v.(Packable))
	}
	row := table.Lookup(th, flds, packed)
	return row, table.Header(), query
}

func findKey(table qry.Query, flds []string) []string {
	for _, key := range table.Keys() {
		if set.Equal(flds, key) {
			return key
		}
	}
	return nil
}

func getWhere(ob *SuObject) string {
	var sb strings.Builder
	sep := "\nwhere "
	iter := ob.Iter2(false, true)
	for k, v := iter(); v != nil; k, v = iter() {
		field := ToStr(k)
		if field == "query" {
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

func single(q qry.Query) bool {
	keys := q.Keys()
	return len(keys) == 1 && len(keys[0]) == 0
}
