// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
)

func TestTableLookup(t *testing.T) {
	ss := func(args ...string) []string {
		return args
	}
	is := func(args ...int) []string {
		ss := make([]string, len(args))
		for i, n := range args {
			ss[i] = Pack(IntVal(n))
		}
		return ss
	}
	test := func(query string, cols, vals []string, expected string) {
		t.Helper()
		q := ParseQuery(query, testTran{}, nil)
		q, _, _ = Setup(q, ReadMode, testTran{})
		row := q.(*Table).Lookup(nil, cols, vals)
		assert.T(t).This(fmt.Sprint(row)).Is(expected)
	}
	test("table", ss("a"), is(123), "[{123}]")
	test("comp", ss("a", "b", "c"), is(12, 34, 56), "[{12, 34, 56}]")
}

func TestQueryGet(t *testing.T) {
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	th := &Thread{}
	db := testDb()
	defer db.Close()
	str := func(hdr *Header, rows []Row, hasSort bool) string {
		hdr.Columns = slices.Clone(hdr.Columns)
		slices.Sort(hdr.Columns)
		if !hasSort && len(rows) > 1 {
			slices.SortFunc(rows, func(x, y Row) int {
				for _, col := range hdr.Columns {
					xs := x.GetRawVal(hdr, col, nil, nil)
					ys := y.GetRawVal(hdr, col, nil, nil)
					if cmp := strings.Compare(xs, ys); cmp != 0 {
						return cmp
					}
				}
				return 0
			})
		}
		var sb strings.Builder
		for _, col := range hdr.Columns {
			sb.WriteString(col)
			sb.WriteString("\t")
		}
		sb.WriteString("\n")
		for _, row := range rows {
			for _, col := range hdr.Columns {
				sb.WriteString(row.GetVal(hdr, col, nil, nil).String())
				sb.WriteString("\t")
			}
			sb.WriteString("\n")
		}
		s := sb.String()
		s = strings.ReplaceAll(s, `"`, "'")
		return s
	}
	get := func(q Query, dir Dir, hasSort bool) string {
		t.Helper()
		var rows []Row
		q.Rewind()
		for row := q.Get(th, dir); row != nil; row = q.Get(th, dir) {
			assert.Msg("too many").That(len(rows) < 100)
			rows = append(rows, row)
		}
		if dir == Prev {
			slc.Reverse(rows)
		}
		return str(q.Header(), rows, hasSort)
	}
	test := func(query, strategy, expected string) {
		t.Helper()
		tran := sizeTran{db.NewReadTran()}
		q := ParseQuery(query, tran, nil)
		_, hasSort := q.(*Sort)

		query2 := format(0, q, 0)
		q2 := ParseQuery(query2, tran, nil)
		assert.This(format(0, q2, 0)).Is(query2)

		q, _, _ = Setup(q, ReadMode, tran)
		qs := strings.ReplaceAll(String(q), `"`, "'")
		assert.T(t).This(qs).Is(strategy)
		assert.T(t).Msg("forward:", query).This(get(q, Next, hasSort)).Like(expected)
		assert.T(t).Msg("reverse:", query).This(get(q, Prev, hasSort)).Like(expected)

		if !strings.Contains(query, "tables") &&
			!strings.Contains(query, "columns") &&
			!strings.Contains(query, "indexes") &&
			!strings.Contains(query, "summarize") {
			simple := str(q2.Header(), q2.Simple(th), hasSort)
			assert.T(t).Msg("simple:", query).This(simple).Like(expected)
		}
	}

	test("indexes project table, columns, key",
		"indexes project-copy table, columns, key",
		`columns	key	table
		'abbrev'	true	'cus'
		'city'	false	'supplier'
		'cnum'	true	'cus'
		'date'	false	'hist'
		'date'	true	'dates'
		'date'	true	'hist2'
		'date,item,id'	true	'hist'
		'date,item,id'	true	'trans'
		'id'	false	'hist2'
		'id'	true	'alias'
		'id'	true	'customer'
		'item'	false	'trans'
		'item'	true	'inven'
		'supplier'	true	'supplier'
		'tnum'	true	'co'
		'tnum'	true	'task'`)
	test("columns",
		"columns",
		`column	field	table
		'abbrev'	1	'cus'
		'city'	2	'customer'
		'city'	2	'supplier'
		'cnum'	0	'cus'
		'cnum'	1	'task'
		'column'	1	'columns'
		'columns'	1	'indexes'
		'cost'	2	'trans'
		'cost'	3	'hist'
		'cost'	3	'hist2'
		'date'	0	'dates'
		'date'	0	'hist'
		'date'	0	'hist2'
		'date'	3	'trans'
		'field'	2	'columns'
		'fkcolumns'	4	'indexes'
		'fkmode'	5	'indexes'
		'fktable'	3	'indexes'
		'id'	0	'alias'
		'id'	0	'customer'
		'id'	1	'trans'
		'id'	2	'hist'
		'id'	2	'hist2'
		'item'	0	'inven'
		'item'	0	'trans'
		'item'	1	'hist'
		'item'	1	'hist2'
		'key'	2	'indexes'
		'name'	1	'customer'
		'name'	1	'supplier'
		'name'	2	'cus'
		'name2'	1	'alias'
		'nrows'	1	'tables'
		'qty'	1	'inven'
		'signed'	1	'co'
		'supplier'	0	'supplier'
		'table'	0	'columns'
		'table'	0	'indexes'
		'table'	0	'tables'
		'tnum'	0	'co'
		'tnum'	0	'task'
		'totalsize'	2	'tables'
		'view_definition'	1	'views'
		'view_name'	0	'views'`)
	test("tables",
		"tables",
		`nrows	table	totalsize
		0	'views'	0
		2	'alias'	25
		3	'hist2'	68
		3	'inven'	42
		4	'co'	55
		4	'cus'	64
		4	'customer'	98
		4	'dates'	52
		4	'hist'	91
		4	'supplier'	128
		4	'trans'	92
		8	'task'	95
		15	'tables'	0
		16	'indexes'	0
		44	'columns'	0`)
	test("customer",
		"customer^(id)",
		`city	id	name
		'calgary'	'c'	'calac'
		'saskatoon'	'a'	'axon'
		'saskatoon'	'i'	'intercon'
		'vancouver'	'e'	'emerald'`)
	test("hist",
		"hist^(date)",
		`cost	date	id	item
		100	970101	'a'	'disk'
		200	970101	'e'	'disk'
		200	970102	'c'	'mouse'
		300	970103	'e'	'pencil'`)
	test("trans",
		"trans^(item)",
		`cost	date	id	item
		100	970101	'a'	'disk'
		150	970201	'c'	'eraser'
		200	960204	'e'	'mouse'
		200	970101	'c'	'mouse'`)

	// rename
	test("trans rename id to code, date to when",
		"trans^(item) rename id to code, date to when",
		`code	cost	item	when
		'a'	100	'disk'	970101
		'c'	150	'eraser'	970201
		'c'	200	'mouse'	970101
		'e'	200	'mouse'	960204`)

	// sort
	test("customer sort id",
		"customer^(id)",
		`city	id	name
		'saskatoon'	'a'	'axon'
		'calgary'	'c'	'calac'
		'vancouver'	'e'	'emerald'
		'saskatoon'	'i'	'intercon'`)
	test("customer sort reverse id",
		"customer^(id) reverse",
		`city	id	name
		'saskatoon'	'i'	'intercon'
		'vancouver'	'e'	'emerald'
		'calgary'	'c'	'calac'
		'saskatoon'	'a'	'axon'`)
	test("customer sort city",
		"customer^(id) tempindex(city)",
		`city	id	name
		'calgary'	'c'	'calac'
		'saskatoon'	'a'	'axon'
		'saskatoon'	'i'	'intercon'
		'vancouver'	'e'	'emerald'`)
	test("customer sort reverse city, name",
		"customer^(id) tempindex(city,name) reverse",
		`city	id	name
		'vancouver'	'e'	'emerald'
		'saskatoon'	'i'	'intercon'
		'saskatoon'	'a'	'axon'
		'calgary'	'c'	'calac'`)
	test("task sort cnum, tnum",
		"task^(tnum) tempindex(cnum,tnum)",
		`cnum	tnum
		1	100
		1	104
		2	101
		2	105
		3	102
		3	106
		4	103
		4	107`)
	test("customer times inven sort qty, id",
		"(customer^(id) times inven^(item)) tempindex(qty,id)",
		`city	id	item	name	qty
		'saskatoon'	'a'	'mouse'	'axon'	2
		'calgary'	'c'	'mouse'	'calac'	2
		'vancouver'	'e'	'mouse'	'emerald'	2
		'saskatoon'	'i'	'mouse'	'intercon'	2
		'saskatoon'	'a'	'disk'	'axon'	5
		'calgary'	'c'	'disk'	'calac'	5
		'vancouver'	'e'	'disk'	'emerald'	5
		'saskatoon'	'i'	'disk'	'intercon'	5
		'saskatoon'	'a'	'pencil'	'axon'	7
		'calgary'	'c'	'pencil'	'calac'	7
		'vancouver'	'e'	'pencil'	'emerald'	7
		'saskatoon'	'i'	'pencil'	'intercon'	7`)
	test("customer extend up = city[1..] sort up",
		"customer^(id) extend up = city[1..] tempindex(up)",
		`city	id	name	up
		'calgary'	'c'	'calac'	'algary'
		'vancouver'	'e'	'emerald'	'ancouver'
		'saskatoon'	'a'	'axon'	'askatoon'
		'saskatoon'	'i'	'intercon'	'askatoon'`)
	test("trans minus hist sort id, cost",
		"(trans^(item) minus(date,item,id) hist^(date,item,id)) tempindex(id,cost)",
		`cost	date	id	item
		150	970201	'c'	'eraser'
		200	970101	'c'	'mouse'
		200	960204	'e'	'mouse'`)
	test("customer rename id to id_new sort id_new",
		"customer^(id) rename id to id_new",
		`city	id_new	name
		'saskatoon'	'a'	'axon'
		'calgary'	'c'	'calac'
		'vancouver'	'e'	'emerald'
		'saskatoon'	'i'	'intercon'`)

	// project
	test("customer project city, id",
		"customer^(id) project-copy city, id",
		`city	id
		'calgary'	'c'
		'saskatoon'	'a'
		'saskatoon'	'i'
		'vancouver'	'e'`)
	test("supplier project city",
		"supplier^(city) project-seq city",
		`city
		'calgary'
		'saskatoon'
		'vancouver'`)
	test("trans project item",
		"trans^(item) project-seq item",
		`item
		'disk'
		'eraser'
		'mouse'`)
	test("customer project city",
		"customer^(id) project-map city",
		`city
		'calgary'
		'saskatoon'
		'vancouver'`)

	// extend
	test("trans extend newcost = cost * 1.1",
		"trans^(item) extend newcost = cost * 1.1",
		`cost	date	id	item	newcost
		100	970101	'a'	'disk'	110
		150	970201	'c'	'eraser'	165
		200	960204	'e'	'mouse'	220
		200	970101	'c'	'mouse'	220`)
	test("trans extend x = cost * 1.1, y = x $ '*'",
		"trans^(item) extend x = cost * 1.1, y = x $ '*'",
		`cost	date	id	item	x	y
		100	970101	'a'	'disk'	110	'110*'
		150	970201	'c'	'eraser'	165	'165*'
		200	960204	'e'	'mouse'	220	'220*'
		200	970101	'c'	'mouse'	220	'220*'`)

	// times
	test("customer times inven",
		"customer^(id) times inven^(item)",
		`city	id	item	name	qty
		'calgary'	'c'	'disk'	'calac'	5
		'calgary'	'c'	'mouse'	'calac'	2
		'calgary'	'c'	'pencil'	'calac'	7
		'saskatoon'	'a'	'disk'	'axon'	5
		'saskatoon'	'a'	'mouse'	'axon'	2
		'saskatoon'	'a'	'pencil'	'axon'	7
		'saskatoon'	'i'	'disk'	'intercon'	5
		'saskatoon'	'i'	'mouse'	'intercon'	2
		'saskatoon'	'i'	'pencil'	'intercon'	7
		'vancouver'	'e'	'disk'	'emerald'	5
		'vancouver'	'e'	'mouse'	'emerald'	2
		'vancouver'	'e'	'pencil'	'emerald'	7`)

	// minus
	test("trans minus trans",
		"trans^(item) minus(date,item,id) trans^(date,item,id)",
		`cost	date	id	item`)
	test("hist minus hist2",
		"hist^(date) minus(date) hist2^(date)",
		`cost	date	id	item
		200	970101	'e'	'disk'
		200	970102	'c'	'mouse'`)
	test("trans minus hist",
		"trans^(item) minus(date,item,id) hist^(date,item,id)",
		`cost	date	id	item
		150	970201	'c'	'eraser'
		200	960204	'e'	'mouse'
		200	970101	'c'	'mouse'`)
	test("trans minus hist sort date",
		"trans^(date,item,id) minus(date,item,id) hist^(date,item,id)",
		`cost	date	id	item
		200	960204	'e'	'mouse'
		200	970101	'c'	'mouse'
		150	970201	'c'	'eraser'`)
	test("(trans minus trans) where item = 0",
		"trans^(item) where item is 0 minus(date,item,id) "+
			"(trans^(date,item,id) where item is 0)",
		`cost	date	id	item`)
	test("inven minus (inven where item = 'mouse')",
		"inven^(item) minus() (inven^(item) where*1 item is 'mouse')",
		`item		qty
		'disk'		5
		'pencil'	7`)

	// intersect
	test("trans intersect trans",
		"trans^(item) intersect(date,item,id) trans^(date,item,id)",
		`cost	date	id	item
		100	970101	'a'	'disk'
		150	970201	'c'	'eraser'
		200	960204	'e'	'mouse'
		200	970101	'c'	'mouse'`)
	test("trans intersect hist",
		"trans^(item) intersect(date,item,id) hist^(date,item,id)",
		`cost	date	id	item
		100	970101	'a'	'disk'`)
	test("hist intersect hist2",
		"hist^(date) intersect(date) hist2^(date)",
		`cost	date	id	item
		100	970101	'a'	'disk'
		300	970103	'e'	'pencil'`)
	test("(hist intersect hist2) where cost = 100",
		"hist^(date) where cost is 100 intersect(date) (hist2^(date) where cost is 100)",
		`cost	date	id	item
		100	970101	'a'	'disk'`)

	// union
	test("hist2 union hist",
		"hist2^(date) union-lookup(date,item,id) hist^(date,item,id)",
		`cost	date	id	item
		100	970101	'a'	'disk'
		200	970101	'e'	'disk'
		200	970102	'c'	'mouse'
		200	970102	'e'	'disk'
		300	970103	'e'	'pencil'`)
	test("hist2 union trans",
		"hist2^(date) union-lookup(date,item,id) trans^(date,item,id)",
		`cost	date	id	item
		100	970101	'a'	'disk'
		150	970201	'c'	'eraser'
		200	960204	'e'	'mouse'
		200	970101	'c'	'mouse'
		200	970102	'e'	'disk'
		300	970103	'e'	'pencil'`)
	test("alias union alias",
		"alias^(id) union-merge(id) alias^(id)",
		`id name2
        'a'	'abc'
		'c'	'trical'`)
	test("trans union hist",
		"trans^(date,item,id) union-merge(date,item,id) hist^(date,item,id)",
		`cost	date	id	item
		100	970101	'a'	'disk'
		150	970201	'c'	'eraser'
		200	960204	'e'	'mouse'
		200	970101	'c'	'mouse'
		200	970101	'e'	'disk'
		200	970102	'c'	'mouse'
		300	970103	'e'	'pencil'`)
	test("trans union trans sort id,cost",
		"(trans^(date,item,id) union-merge(date,item,id) "+
			"trans^(date,item,id)) tempindex(id,cost)",
		`cost	date	id	item
		100	970101	'a'	'disk'
		150	970201	'c'	'eraser'
		200	970101	'c'	'mouse'
		200	960204	'e'	'mouse'`)
	test("(hist2 rename cost to amt) union (trans rename cost to amt)",
		"hist2^(date) rename cost to amt union-lookup(date,item,id) "+
			"(trans^(date,item,id) rename cost to amt)",
		`amt	date	id	item
		100	970101	'a'	'disk'
		150	970201	'c'	'eraser'
		200	960204	'e'	'mouse'
		200	970101	'c'	'mouse'
		200	970102	'e'	'disk'
		300	970103	'e'	'pencil'`)
	test("(hist2 extend x = 1) union (trans extend x = 1)",
		"hist2^(date) extend x = 1 union-lookup(date,item,id) "+
			"(trans^(date,item,id) extend x = 1)",
		`cost	date	id	item	x
		100	970101	'a'	'disk'	1
		150	970201	'c'	'eraser'	1
		200	960204	'e'	'mouse'	1
		200	970101	'c'	'mouse'	1
		200	970102	'e'	'disk'	1
		300	970103	'e'	'pencil'	1`)
	test("(co where tnum = 100 remove signed) union "+
		"(co where tnum = 100 remove signed)",
		"co^(tnum) where*1 tnum is 100 project-copy tnum union-merge() "+
			"(co^(tnum) where*1 tnum is 100 project-copy tnum)",
		`tnum
        100`)
	test("(co where tnum = 100 remove tnum) union "+
		"(co where tnum = 100 remove tnum)",
		"co^(tnum) where*1 tnum is 100 project-copy signed union-merge() "+
			"(co^(tnum) where*1 tnum is 100 project-copy signed)",
		`signed
        990101`)
	test("((co where tnum = 100) union (co where tnum = 102)) union "+
		"(co where tnum = 100)",
		"co^(tnum) where*1 tnum is 100 union-lookup(tnum) "+
			"(co^(tnum) where*1 tnum is 100 union-disjoint(tnum)-merge(tnum) "+
			"(co^(tnum) where*1 tnum is 102))",
		`signed	tnum
		990101	100
		990102	102`)
	test("(((co where tnum = 100) union (co where tnum = 102)) remove tnum)"+
		" union "+
		"(((co where tnum = 104) union (co where tnum = 106)) remove tnum)",
		"(co^(tnum) where*1 tnum is 100 union-disjoint(tnum)-merge(signed) (co^(tnum) where*1 tnum is 102)) project-seq signed union-merge(signed) ((co^(tnum) where*1 tnum is 104 union-disjoint(tnum)-merge(signed) (co^(tnum) where*1 tnum is 106)) project-seq signed)",
		`signed
        990101
        990102
        990103
        990104`)
	test(`((co where tnum = 104 remove tnum) union (co where tnum = 106 remove tnum))
		union
		((co where tnum = 104 remove tnum) union (co where tnum = 106 remove tnum))`,
		"(co^(tnum) where*1 tnum is 104 project-copy signed union-merge(signed) (co^(tnum) where*1 tnum is 106 project-copy signed)) union-merge(signed) (co^(tnum) where*1 tnum is 104 project-copy signed union-merge(signed) (co^(tnum) where*1 tnum is 106 project-copy signed))",
		`signed
        990103
        990104`)

	// join
	test("customer join alias",
		"customer^(id) join 1:1 by(id) alias^(id)",
		`city	id	name	name2
		'calgary'	'c'	'calac'	'trical'
		'saskatoon'	'a'	'axon'	'abc'`)
	test("trans join inven",
		"inven^(item) join 1:n by(item) trans^(item)",
		`cost	date	id	item	qty
		100	970101	'a'	'disk'	5
		200	960204	'e'	'mouse'	2
		200	970101	'c'	'mouse'	2`)
	test("customer leftjoin alias",
		"customer^(id) leftjoin 1:1 by(id) alias^(id)",
		`city	id	name	name2
		'calgary'	'c'	'calac'	'trical'
		'saskatoon'	'a'	'axon'	'abc'
		'saskatoon'	'i'	'intercon'	''
		'vancouver'	'e'	'emerald'	''`)
	test("inven leftjoin trans",
		"inven^(item) leftjoin 1:n by(item) trans^(item)",
		`cost	date	id	item	qty
		''	''	''	'pencil'	7
		100	970101	'a'	'disk'	5
		200	960204	'e'	'mouse'	2
		200	970101	'c'	'mouse'	2`)
	test("customer leftjoin hist2",
		"customer^(id) leftjoin 1:n by(id) hist2^(id)",
		`city	cost	date	id	item	name
		'calgary'	''	''	'c'	''	'calac'
		'saskatoon'	''	''	'i'	''	'intercon'
		'saskatoon'	100	970101	'a'	'disk'	'axon'
		'vancouver'	200	970102	'e'	'disk'	'emerald'
		'vancouver'	300	970103	'e'	'pencil'	'emerald'`)
	test("hist join customer",
		"hist^(date) join n:1 by(id) customer^(id)",
		`city	cost	date	id	item	name
		'calgary'	200	970102	'c'	'mouse'	'calac'
		'saskatoon'	100	970101	'a'	'disk'	'axon'
		'vancouver'	200	970101	'e'	'disk'	'emerald'
		'vancouver'	300	970103	'e'	'pencil'	'emerald'`)
	test("customer leftjoin (alias where name2 is 'abc')",
		"customer^(id) leftjoin 1:1 by(id) (alias^(id) where name2 is 'abc')",
		`city	id	name	name2
		'calgary'	'c'	'calac'	''
		'saskatoon'	'a'	'axon'	'abc'
		'saskatoon'	'i'	'intercon'	''
		'vancouver'	'e'	'emerald'	''`)
	test("customer leftjoin (alias where id is 'c')",
		"customer^(id) leftjoin 1:1 by(id) (alias^(id) where*1 id is 'c')",
		`city	id	name	name2
		'calgary'	'c'	'calac'	'trical'
		'saskatoon'	'a'	'axon'	''
		'saskatoon'	'i'	'intercon'	''
		'vancouver'	'e'	'emerald'	''`)
	test("customer leftjoin (trans where date = 970101 and item = 'mouse')",
		"customer^(id) leftjoin 1:1 by(id) "+
			"(trans^(date,item,id) where date is 970101 and item is 'mouse')",
		`city	cost	date	id	item	name
		'calgary'	200	970101	'c'	'mouse'	'calac'
		'saskatoon'	''	''	'a'	''	'axon'
		'saskatoon'	''	''	'i'	''	'intercon'
		'vancouver'	''	''	'e'	''	'emerald'`)

	// where
	test("customer where id > 'd'", // range
		"customer^(id) where id > 'd'",
		`city	id	name
		'saskatoon'	'i'	'intercon'
		'vancouver'	'e'	'emerald'`)
	test("customer where id > 'd' and id < 'j'", // range
		"customer^(id) where id > 'd' and id < 'j'",
		`city	id	name
		'saskatoon'	'i'	'intercon'
		'vancouver'	'e'	'emerald'`)
	test("customer where id is 'e'", // point
		"customer^(id) where*1 id is 'e'",
		`city	id	name
		'vancouver'	'e'	'emerald'`)
	test("customer where id is 'd'", // point
		"customer^(id) where*1 id is 'd'",
		`city	id	name`)
	test("inven where qty > 0", // filter
		"inven^(item) where qty > 0",
		`item	qty
		'disk'	5
		'mouse'	2
		'pencil'	7`)
	test("inven where item =~ 'i'", // filter
		"inven^(item) where item =~ 'i'",
		`item	qty
		'disk'	5
		'pencil'	7`)
	test("inven where item in (1, 'disk', 'mouse', 2, 'disk', 'pencil')", // points
		"inven^(item) where item in (1, 'disk', 'mouse', 2, 'disk', 'pencil')",
		`item	qty
		'disk'	5
		'mouse'	2
		'pencil'	7`)
	test("inven where item <= 'e' or item >= 'p'", // filter
		"inven^(item) where item <= 'e' or item >= 'p'",
		`item	qty
		'disk'	5
		'pencil'	7`)
	test("cus where cnum is 2 and abbrev is 'b'", // points
		"cus^(cnum) where*1 cnum is 2 and abbrev is 'b'",
		`abbrev	cnum	name
		'b'	2	'bill'`)
	test("cus where cnum is 2 and abbrev >= 'b' and abbrev < 'c'", // point
		"cus^(cnum) where*1 cnum is 2 and abbrev >= 'b' and abbrev < 'c'",
		`abbrev	cnum	name
		'b'	2	'bill'`)
	test("hist where date in (970101, 970102) and item < 'z'", // ranges
		"hist^(date) where date in (970101, 970102) and item < 'z'",
		`cost	date	id	item
		100	970101	'a'	'disk'
		200	970101	'e'	'disk'
		200	970102	'c'	'mouse'`)
	test("customer where id not in ('z')",
		"customer^(id) where not (id is 'z')",
		`city	id	name
		'calgary'	'c'	'calac'
		'saskatoon'	'a'	'axon'
		'saskatoon'	'i'	'intercon'
		'vancouver'	'e'	'emerald'`)

	// summarize
	test("customer summarize count",
		"customer^(id) summarize-tbl count",
		`count
		1000`)
	test("hist summarize max date",
		"hist^(date) summarize-idx max date",
		`max_date
		970103`)
	test("customer summarize max id",
		"customer^(id) summarize-idx* max id",
		`city	id	max_id	name
		'saskatoon'	'i'	'i'	'intercon'`)
	test("hist summarize item, total cost sort item",
		"hist^(date) summarize-map item, total cost tempindex(item)",
		`item		total_cost
		'disk'		300
		'mouse'		200
		'pencil'	300`)
	test("hist summarize item, total cost, max id, average cost sort item",
		"hist^(date) summarize-map item, total cost, "+
			"average cost, max id tempindex(item)",
		`average_cost	item	max_id	total_cost
		150	'disk'	'e'	300
		200	'mouse'	'c'	200
		300	'pencil'	'e'	300`)
	test("hist summarize item, total cost sort total_cost, item",
		"hist^(date) summarize-map item, total cost "+
			"tempindex(total_cost,item)",
		`item		total_cost
		'mouse'		200
		'disk'		300
		'pencil'	300`)
	test("customer summarize max name",
		"customer^(id) summarize-seq max name",
		`max_name
		'intercon'`)
	test("hist summarize min cost, average cost, max cost, sum = total cost",
		"hist^(date) summarize-seq min cost, "+
			"average cost, max cost, sum = total cost",
		`average_cost	max_cost	min_cost	sum
		200	300	100	800`)
	test("hist summarize item, total cost, count sort item",
		"hist^(date) summarize-map item, count, total cost tempindex(item)",
		`count	item	total_cost
		2	'disk'	300
		1	'mouse'	200
		1	'pencil'	300`)
	test("inven summarize max item",
		"inven^(item) summarize-idx* max item",
		`item	max_item	qty
		'pencil'	'pencil'	7`)
	test("hist summarize date, list id",
		"hist^(date) summarize-seq date, list id",
		`date	list_id
		970101	#('a', 'e')
		970102	#('c')
		970103	#('e')`)
	test("hist summarize date, total cost sort total_cost",
		"hist^(date) summarize-seq date, total cost tempindex(total_cost)",
		`date	total_cost
        970102	200
        970101	300
        970103	300`)
	test("hist summarize list id",
		"hist^(date) summarize-seq list id",
		`list_id
		#('a', 'c', 'e')`)
	test("cus summarize max cnum sort name",
		"cus^(cnum) summarize-idx* max cnum",
		`abbrev	cnum	max_cnum	name
		'd'	4	4	'dick'`)
	test("supplier summarize min city",
		"supplier^(city) summarize-idx min city",
		`min_city
		'calgary'`)
	test("supplier summarize max city",
		"supplier^(city) summarize-idx max city",
		`max_city
		'vancouver'`)
	test("supplier summarize min city, max city",
		"supplier^(supplier) summarize-seq min city, max city",
		`max_city	min_city
		'vancouver'	'calgary'`)
	test("hist summarize max cost",
		"hist^(date) summarize-seq max cost",
		`max_cost
		300`)

	// tempindex
	test("tables intersect columns",
		"columns intersect(table) (tables tempindex(table))",
		`table`)
	test("tables minus tables",
		"tables minus(table) (tables tempindex(table))",
		`nrows	table	totalsize`)

	test("(customer project id) union (customer project id)",
		"customer^(id) project-copy id union-merge(id) (customer^(id) project-copy id)",
		`id
		'a'
		'c'
		'e'
		'i'`)
	test("(trans project item,date) union (trans project date, item)",
		"trans^(date,item,id) project-seq item, date union-merge(item,date) "+
			"(trans^(date,item,id) project-seq date, item)",
		`date	item
		960204	'mouse'
		970101	'disk'
		970101	'mouse'
		970201	'eraser'`)

	test("(customer summarize id, count) union (customer summarize id, count)",
		"customer^(id) summarize-seq id, count union-merge(id) "+
			"(customer^(id) summarize-seq id, count)",
		`count	id
		1	'a'
		1	'c'
		1	'e'
		1	'i'`)
	test("(trans summarize item,date, count) union (trans summarize date,item, count)",
		"trans^(date,item,id) summarize-seq item, date, count "+
			"union-merge(item,date) "+
			"(trans^(date,item,id) summarize-seq date, item, count)",
		`count	date	item
		1	960204	'mouse'
		1	970101	'disk'
		1	970101	'mouse'
		1	970201	'eraser'`)

	test("(customer project id) join (hist project id)",
		"hist^(date) project-map id join 1:1 by(id) (customer^(id) project-copy id)",
		`id
		'a'
		'c'
		'e'`)
	test("(customer summarize id,count) join (hist summarize id,count) sort id",
		"customer^(id) summarize-seq id, count join 1:1 by(id,count) "+
			"(hist^(date) summarize-map id, count tempindex(id,count))",
		`count	id
		1	'a'
		1	'c'`)
}
