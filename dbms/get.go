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
)

// var total, fast int

// var _ = exit.Add("get", func() {
// 	fmt.Println("total", total, "fast", fast)
// })

func get(th *Thread, tran qry.QueryTran, args Value, dir Dir) (Row, *Header, string) {
	defer th.Suneido.Store(th.Suneido.Load())
	th.Suneido.Store(nil) // use main Suneido object

	ob := args.(*SuObject)
	query := getQuery(ob)
	if dir == Only || dir == Any {
		query = qry.StripSort(query)
	}
	// total++
	if row, hdr := fastGet(th, tran, query, ob, dir); hdr != nil {
		// fast++
		return row, hdr, query
	}
	if where := getWhere(ob, dir); where != "" {
		// need the newline in case the query ends with //comment
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
	if dir == Next || dir == Prev {
		return nil, nil
	}
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
	if dir == Only {
		return getLookup(th, tran, tbl, flds, vals)
	}
	if dir == Any {
		return getExists(th, tran, tbl, flds, vals)
	}
	panic(assert.ShouldNotReachHere())
}

func getLookup(th *Thread, tran qry.QueryTran, table *qry.Table, flds []string, vals []Value) (Row, *Header) {
	key := findKey(table.Keys(), flds)
	if key != nil {
		if len(key) == 0 {
			return table.Get(th, Next), table.Header()
		}
		table.SetIndex(key)
		return table.Lookup(th, flds, packVals(vals)), table.Header()
	}
	getfn := getIndex(th, tran, table, flds, vals)
	row, hdr := getfn()
	if row != nil {
		if r2, _ := getfn(); r2 != nil {
			panic("Query1 not unique")
		}
	}
	return row, hdr
}

// findKey finds the first key in table that is the same set as flds
func findKey(keys [][]string, flds []string) []string {
	for _, key := range keys {
		if set.Equal(flds, key) {
			return key
		}
	}
	return nil
}

func packVals(vals []Value) []string {
	packed := make([]string, len(vals))
	for i, v := range vals {
		packed[i] = Pack(v.(Packable))
	}
	return packed
}

func getExists(th *Thread, tran qry.QueryTran, table *qry.Table, flds []string, vals []Value) (Row, *Header) {
	getfn := getIndex(th, tran, table, flds, vals)
	row, hdr := getfn()
	if row != nil {
		return existsRow, existsHdr
	}
	return row, hdr
}

func getIndex(th *Thread, tran qry.QueryTran, table *qry.Table, flds []string, vals []Value) func() (Row, *Header) {
	idx, idxlen := findIndex(table.Indexes(), flds)
	if idx == nil {
		return func() (Row, *Header) { return nil, nil }
	}
	table.SetIndex(idx)

	// Split fields and values into those that match the index and those that don't
	idxFlds := make([]string, 0, idxlen)
	idxVals := make([]string, 0, idxlen)
	remainFlds := make([]string, 0, len(flds)-idxlen)
	remainVals := make([]string, 0, len(vals)-idxlen)

	// Collect the fields and values that match the index prefix
	for i, fld := range flds {
		val := Pack(vals[i].(Packable))
		if slices.Contains(idx[:idxlen], fld) {
			idxFlds = append(idxFlds, fld)
			idxVals = append(idxVals, val)
		} else {
			remainFlds = append(remainFlds, fld)
			remainVals = append(remainVals, val)
		}
	}

	if len(idxFlds) > 0 {
		table.Select(idxFlds, idxVals)
	}
	st := qry.MakeSuTran(tran)
	hdr := table.Header()
	return func() (Row, *Header) {
	outer:
		for {
			row := table.Get(th, Next)
			if row == nil {
				break
			}
			for i, fld := range remainFlds {
				if row.GetRawVal(hdr, fld, th, st) != remainVals[i] {
					continue outer // row does not match the additional filter
				}
			}
			return row, hdr
		}
		return nil, hdr
	}
}

// findIndex finds the index that has the most flds
// (in any order) as a prefix
func findIndex(indexes [][]string, flds []string) ([]string, int) {
	var best []string
	bestLen := 0
	for _, idx := range indexes {
		prefixLen := 0
		for _, field := range idx {
			if !slices.Contains(flds, field) {
				break
			}
			prefixLen++
		}
		if prefixLen > 0 && prefixLen > bestLen {
			best = idx
			bestLen = prefixLen
		}
	}
	return best, bestLen
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
