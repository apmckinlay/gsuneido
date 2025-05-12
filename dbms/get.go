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
	trace.Query.Println(dir, fixcost+varcost, "-", query)
	d := dir
	if dir == Only || dir == Any {
		d = Next
	}
	row := q.Get(th, d)
	if dir == Only || dir == Any {
		if w, ok := q.(*qry.Where); ok && w.Slow() &&
			!(strings.HasPrefix(query, "columns") ||
				strings.HasPrefix(query, "indexes")) {
			Warning(dir, "slow:", query)
		}
	}
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
		return getLookup(th, tran, tbl, flds, packed, dir)
	}
	if dir == Any {
		return getExists(th, tran, tbl, flds, packed, dir)
	}
	panic(assert.ShouldNotReachHere())
}

func getLookup(th *Thread, tran qry.QueryTran, table *qry.Table, flds, vals []string, dir Dir) (Row, *Header) {
	single, getfn := getIndex(th, tran, table, flds, vals, dir)
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

func getExists(th *Thread, tran qry.QueryTran, table *qry.Table, flds, vals []string, dir Dir) (row Row, hdr *Header) {
	_, getfn := getIndex(th, tran, table, flds, vals, dir)
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
	flds, vals []string, dir Dir) (single bool, getfn func() Row) {
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
	if len(flds) == 0 {
		trace.QueryOpt.Println(dir, table.Name(), "no filter")
		return false, func() Row {
			return table.Get(th, Next)
		}
	}
	if table.Single() {
		trace.QueryOpt.Println(dir, table.Name(), "empty key")
		return true, func() Row {
			return filter(table.Get(th, Next))
		}
	}
	if key := findKey(table.Keys(), flds); key != nil {
		trace.QueryOpt.Println(dir, table.Name(), "key", key)
		table.SetIndex(key)
		return true, func() Row {
			return filter(table.Lookup(th, flds, vals))
		}
	}
	if idx := findAll(table.Indexes(), flds); idx != nil {
		trace.QueryOpt.Println(dir, table.Name(), "just", idx)
		table.SetIndex(idx)
		table.Select(flds, vals)
		return false, func() Row {
			return table.Get(th, Next)
		}
	}
	indexes := usableIndexes(table.Indexes(), flds)
	if len(indexes) == 0 {
		return
	}
	if len(indexes) == 1 {
		trace.QueryOpt.Println(dir, table.Name(), "only", indexes[0])
		table.SetIndex(indexes[0])
		table.Select(flds, vals)
		return false, func() Row {
			for n := 0; ; n++ {
				row := table.Get(th, Next)
				if row == nil || nil != filter(row) {
					if n > 100 {
						Warning(dir, "slow", n, table)
					}
					return row
				}
			}
		}
	}
	tables := make([]*qry.Table, len(indexes))
	for i, idx := range indexes {
		tbl := *table // copy
		tables[i] = &tbl
		tables[i].SetIndex(idx)
		tables[i].Select(flds, vals)
	}
	return false, func() Row {
		// iterate the indexes in parallel
		for n := 0; ; n++ {
			for _, tbl := range tables {
				if tbl == nil {
					continue
				}
				row := tbl.Get(th, Next)
				if row == nil || nil != filter(row) {
					// clear the other tables so next get is from this one
					for j := range tables {
						if tables[j] != tbl {
							tables[j] = nil
						}
					}
					trace.QueryOpt.Println(dir, "multi", tbl)
					if n > 100 {
						Warning(dir, "slow", n, tbl)
					}
					return row
				}
			}
		}
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

// findAll returns the first index that contains all the fields
func findAll(indexes [][]string, flds []string) []string {
	for _, idx := range indexes {
		if len(idx) >= len(flds) && set.Equal(idx[:len(flds)], flds) {
			return idx
		}
	}
	return nil
}

func usableIndexes(indexes [][]string, flds []string) [][]string {
	var usable [][]string
	for _, idx := range indexes {
		if hasPrefix(idx, flds) {
			usable = append(usable, idx)
		}
	}
	return usable
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
