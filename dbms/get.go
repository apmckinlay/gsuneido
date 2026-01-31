// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"fmt"
	"slices"
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
)

var slow = map[Dir]int{Only: 100, Any: 2000}

func get(th *Thread, tran qry.QueryTran, args Value, dir Dir) (Row, *Header, string) {
	defer UseMainSuneido(th)()

	// for dir == Strat
	// if the query has a sort, assume QueryFirst or QueryLast
	// if the query does not have a sort, assume Query1 or QueryEmpty?

	ob := args.(*SuObject)
	query := getQuery(ob)
	if dir == Only || dir == Any ||
		(dir == Strat && qry.GetSort(query) == "") {
		query = qry.StripSort(query)
		if row, hdr, s := fastGet(th, tran, query, ob, dir); hdr != nil {
			return row, hdr, s
		}
	}
	if where := getWhere(ob); where != "" {
		// need a newline in case the query ends with //comment
		query += "\n" + where
	}
	q := qry.ParseQuery(query, tran, th.Sviews())
	qs, sorted := q.(*qry.Sort)
	if dir == Only || dir == Any {
		if sorted {
			q = qs.Source() // remove sort
		}
	} else if !sorted && dir != Strat &&
		!strings.Contains(query, "CHECKQUERY SUPPRESS: SORT REQUIRED") {
		panic("QueryFirst and QueryLast require sort")
	}

	q, fixcost, varcost := qry.Setup1(q, qry.ReadMode, tran)
	qry.Warnings(query, q)
	if dir == Strat {
		n, _ := q.Nrows()
		return existsRow, existsHdr, fmt.Sprint(qry.Strategy(q), "\n",
			"[nrecs~ ", trace.Number(n),
			" cost~ ", trace.Number(fixcost+varcost), "]")
	}
	trace.Query.Println(dir, fixcost+varcost, "-", query)
	d := dir
	if dir == Only || dir == Any {
		d = Next
	}
	row := q.Get(th, d)
	if dir == Only || dir == Any {
		if w, ok := q.(*qry.Where); ok && w.InCount() > slow[dir] &&
			!(strings.HasPrefix(query, "columns") ||
				strings.HasPrefix(query, "indexes") ||
				strings.HasPrefix(query, "views")) {
			Warning(dir, "slow:", w.InCount(), query)
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
func fastGet(th *Thread, tran qry.QueryTran, query string, ob *SuObject, dir Dir) (row Row, hdr *Header, strarg string) {
	strarg = query
	table := qry.JustTable(query)
	if table == "" || tran.GetInfo(table) == nil { // could be a view
		return
	}
	tbl, ok := qry.NewTable(tran, table).(*qry.Table)
	if !ok {
		return
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
	packed := slc.MapFn(vals,
		func(v Value) string { return Pack(v.(Packable)) })
	switch dir {
	case Strat:
		_, strarg, _ = getIndex(th, tran, tbl, flds, packed, dir)
		if strarg == "" {
			return nil, nil, ""
		}
		return existsRow, existsHdr, strarg
	case Only:
		row, hdr = getLookup(th, tran, tbl, flds, packed, dir)
		return
	case Any:
		row, hdr = getExists(th, tran, tbl, flds, packed, dir)
		return
	}
	panic(assert.ShouldNotReachHere())
}

func getLookup(th *Thread, tran qry.QueryTran, table *qry.Table, flds, vals []string, dir Dir) (Row, *Header) {
	single, _, getfn := getIndex(th, tran, table, flds, vals, dir)
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
	_, _, getfn := getIndex(th, tran, table, flds, vals, dir)
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
	flds, vals []string, dir Dir) (single bool, strat string,
	getfn func() Row) {
	st := qry.MakeSuTran(tran)
	hdr := table.Header()
	filter := func(oflds, ovals []string, row Row) Row {
		if row != nil {
			for i, fld := range oflds {
				if row.GetRawVal(hdr, fld, th, st) != ovals[i] {
					return nil
				}
			}
		}
		return row
	}
	if len(flds) == 0 {
		table.SetIndex(table.Indexes()[0])
		strat = "no select: " + table.String()
		trace.QueryOpt.Println(dir, strat)
		return false, strat, func() Row {
			return table.Get(th, Next)
		}
	}
	if key := findKey(table.Keys(), flds); key != nil {
		// selecting on a key so only one record in the result
		table.SetIndex(key)
		strat = "key: " + table.String()
		trace.QueryOpt.Println(dir, strat)
		iflds, ivals, oflds, ovals := qry.Split(flds, vals, key)
		return true, strat, func() Row {
			return filter(oflds, ovals, table.Lookup(th, iflds, ivals))
		}
	}
	if idx := findAll(table.Indexes(), flds); idx != nil {
		table.SetIndex(idx)
		strat = "just index: " + table.String()
		trace.QueryOpt.Println(dir, strat)
		table.Select(flds, vals)
		return false, strat, func() Row {
			return table.Get(th, Next)
		}
	}
	indexes := usableIndexes(table.Indexes(), flds)
	if len(indexes) == 0 {
		return
	}
	if len(indexes) == 1 {
		table.SetIndex(indexes[0])
		strat = "only " + table.String()
		trace.QueryOpt.Println(dir, strat)
		table.Select(flds, vals)
		return false, strat, func() Row {
			for n := 0; ; n++ {
				row := table.Get(th, Next)
				if row == nil || nil != filter(flds, vals, row) {
					if n > slow[dir] {
						Warning(dir, "slow:", n, table, formatFieldsVals(flds, vals))
					}
					return row
				}
			}
		}
	}
	strat = "multiple indexes: " + table.Name()
	tables := make([]*qry.Table, len(indexes))
	for i, idx := range indexes {
		tbl := *table // copy
		tables[i] = &tbl
		tables[i].SetIndex(idx)
		tables[i].Select(flds, vals)
		strat += " " + str.Join("(,)", idx)
	}
	var prevRow Row
	return false, strat, func() Row {
		// iterate the indexes in parallel
		for n := 0; ; n++ {
			for _, tbl := range tables {
				if tbl == nil {
					continue
				}
				row := tbl.Get(th, Next)
				if row == nil ||
					(!row.SameAs(prevRow) && nil != filter(flds, vals, row)) {
					trace.QueryOpt.Println(dir, "multi", tbl)
					if n > slow[dir] {
						Warning(dir, "slow:", n, tbl, formatFieldsVals(flds, vals))
					}
					prevRow = row
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

// getWhere builds a where for the named arguments.
// It should be eqivalent to builtin queryWhere
func getWhere(ob *SuObject) string {
	var sb strings.Builder
	sep := "where "
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

func formatFieldsVals(flds, vals []string) string {
	var sb strings.Builder
	for i, fld := range flds {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fld)
		sb.WriteString(": ")
		sb.WriteString(Unpack(vals[i]).String())
	}
	return sb.String()
}

func single(q qry.Query) bool {
	keys := q.Keys()
	return len(keys) == 1 && len(keys[0]) == 0
}
