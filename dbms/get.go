// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"slices"
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
)

func get(th *Thread, tran qry.QueryTran, args Value, dir Dir) (Row, *Header, string) {
	defer th.Suneido.Store(th.Suneido.Load())
	th.Suneido.Store(nil) // use main Suneido object

	ob := args.(*SuObject)
	query := getQuery(ob)
	if dir == Only || dir == Any {
		query = qry.StripSort(query)
		if row, hdr := fastGet(th, tran, query, ob, dir); hdr != nil {
			return row, hdr, query
		}
	}
	if where := getWhere(ob, dir); where != "" {
		// need a newline in case the query ends with //comment
		query += "\n" + where
	}
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
		return existsRow, existsHdr, ""
	}
	if dir == Only && !single(q) && q.Get(th, Next) != nil {
		panic("Query1 not unique: " + query)
	}
	return row, q.Header(), q.Updateable()
}

var existsRow Row = []DbRec{{Record: "x"}}
var existsHdr = SimpleHeader([]string{"x"})

func getQuery(ob *SuObject) string {
	if ob.ListSize() >= 1 {
		return ToStr(ob.ListGet(0))
	} else if q := ob.NamedGet(SuStr("query")); q != nil {
		return ToStr(q)
	}
	return ""
}

// fastGet returns a nil Header to indicate it was not applicable
func fastGet(th *Thread, tran qry.QueryTran, query string, ob *SuObject, dir Dir) (Row, *Header) {
	table := qry.JustTable(query)
	if table == "" || tran.GetInfo(table) == nil { // could be a view
		return nil, nil
	}
	tbl, ok := qry.NewTable(tran, table).(*qry.Table)
	if !ok {
		return nil, nil
	}
	flds := make([]string, 0, ob.NamedSize())
	vals := make([]Value, 0, ob.NamedSize())
	iter := ob.Iter2(false, true)
	for k, v := iter(); v != nil; k, v = iter() {
		field := ToStr(k)
		if field == "query" || field == "sort" {
			continue
		}
		flds = append(flds, field)
		vals = append(vals, v)
	}
	packed := slc.MapFn(vals,
		func(v Value) string { return Pack(v.(Packable)) })
	if dir == Only {
		return getLookup(th, tran, tbl, flds, packed)
	}
	if dir == Any {
		return getExists(th, tran, tbl, flds, packed)
	}
	panic(assert.ShouldNotReachHere())
}

func getLookup(th *Thread, tran qry.QueryTran, table *qry.Table, flds, vals []string) (Row, *Header) {
	single, getfn := getIndex(th, tran, table, flds, vals)
	if getfn == nil {
		return nil, nil
	}
	row := getfn()
	if row != nil && !single {
		if r2 := getfn(); r2 != nil {
			panic("Query1 not unique")
		}
	}
	return row, table.Header()
}

func getExists(th *Thread, tran qry.QueryTran, table *qry.Table, flds, vals []string) (row Row, hdr *Header) {
	_, getfn := getIndex(th, tran, table, flds, vals)
	if getfn == nil {
		return nil, nil
	}
	row = getfn()
	if row != nil {
		return existsRow, existsHdr
	}
	return nil, existsHdr
}

func getIndex(th *Thread, tran qry.QueryTran, table *qry.Table,
	flds, vals []string) (single bool, getfn func() Row) {
	st := qry.MakeSuTran(tran)
	hdr := table.Header()
	filter := func(row Row) Row {
		if row != nil {
			for i, fld := range flds {
				if row.GetRawVal(hdr, fld, th, st) != vals[i] {
					return nil
				}
			}
		}
		return row
	}
	if table.Single() {
		trace.QueryOpt.Println("Query1/Exists?", table.Name(), "empty key")
		return true, func() Row {
			return filter(table.Get(th, Next))
		}
	}
	if len(flds) == 0 {
		trace.QueryOpt.Println("Query1/Exists?", table.Name(), "no filter")
		return false, func() Row {
			return table.Get(th, Next)
		}
	}
	if key := findKey(table.Keys(), flds); key != nil {
		trace.QueryOpt.Println("Query1/Exists?", table.Name(), "key", key)
		table.SetIndex(key)
		return true, func() Row {
			return filter(table.Lookup(th, flds, vals))
		}
	}
	idx := chooseIndex(table, flds, vals)
	if idx == nil {
		return
	}
	table.SetIndex(idx)
	table.Select(flds, vals)
	return false, func() Row {
		n := 0
		for {
			if n++; n == 100 {
				Warning("Query1/Exists? slow query on", table)
			}
			row := table.Get(th, Next)
			if row == nil {
				break
			}
			if nil == filter(row) {
				continue
			}
			return row
		}
		return nil
	}
}

func findKey(keys [][]string, flds []string) []string {
	for _, key := range keys {
		if set.Subset(flds, key) {
			return key
		}
	}
	return nil
}

func chooseIndex(table *qry.Table, flds, vals []string) []string {
	if idx := onlyUsableIndex(table.Indexes(), flds); idx != nil {
		trace.QueryOpt.Println("Query1/Exists?", table.Name(), "only", idx)
		return idx
	}
	return bestIndexFrac(table, flds, vals)
}

func onlyUsableIndex(indexes [][]string, flds []string) []string {
	var onlyIdx []string
	for _, idx := range indexes {
		if hasPrefix(idx, flds) {
			if onlyIdx != nil {
				return nil // multiple usable indexes
			}
			onlyIdx = idx
		}
	}
	return onlyIdx
}

func bestIndexFrac(table *qry.Table, flds, vals []string) []string {
	var bestIdx []string
	bestFrac := 1.0
	for _, idx := range table.Indexes() {
		if !hasPrefix(idx, flds) {
			continue
		}
		frac := table.RangeFrac(idx, flds, vals)
		if frac < bestFrac {
			bestIdx = idx
			bestFrac = frac
		}
	}
	if bestIdx != nil {
		trace.QueryOpt.Println("Query1/Exists?", table.Name(), "best", bestIdx)
	}
	return bestIdx
}

func hasPrefix(idx []string, flds []string) bool {
	return slices.Contains(flds, idx[0])
}

// getWhere builds a where and sort for the named arguments.
// It should be eqivalent to builtin queryWhere
func getWhere(ob *SuObject, dir Dir) string {
	var sb strings.Builder
	sort := ""
	sep := "where "
	iter := ob.Iter2(false, true)
	for k, v := iter(); v != nil; k, v = iter() {
		field := ToStr(k)
		if field == "query" {
			continue
		} else if field == "sort" {
			if dir == Next || dir == Prev {
				sort = " sort " + ToStr(v)
			}
			continue
		}
		sb.WriteString(sep)
		sep = "\nand "
		sb.WriteString(field)
		sb.WriteString(" is ")
		sb.WriteString(v.String())
	}
	sb.WriteString(sort)
	return sb.String()
}

func single(q qry.Query) bool {
	keys := q.Keys()
	return len(keys) == 1 && len(keys[0]) == 0
}
