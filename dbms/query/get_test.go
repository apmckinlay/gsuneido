// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTableLookup(t *testing.T) {
	pack := func(n int) string {
		return rt.Pack(rt.IntVal(n).(rt.Packable))
	}
	key := func(vals ...int) string {
		if len(vals) == 1 {
			return pack(vals[0])
		}
		var enc ixkey.Encoder
		for _, v := range vals {
			enc.Add(pack(v))
		}
		return enc.String()
	}
	test := func(query, key, expected string) {
		t.Helper()
		q := ParseQuery(query)
		q = Setup(q, ReadMode, testTran{})
		row := q.(*Table).Lookup(key)
		assert.T(t).This(fmt.Sprint(row)).Is(expected)
	}
	test("table", key(123), "[<123>]")
	test("customer", key(12, 34), "[<12, 34>]")
}

func TestQueryGet(t *testing.T) {
	db := testDb()
	defer db.Close()
	reverse := func(rows []rt.Row) {
		for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
			rows[i], rows[j] = rows[j], rows[i]
		}
	}
	get := func(q Query, dir rt.Dir) string {
		t.Helper()
		var rows []rt.Row
		q.Rewind()
		for row := q.Get(dir); row != nil; row = q.Get(dir) {
			rows = append(rows, row)
		}
		if dir == rt.Prev {
			reverse(rows)
		}
		hdr := q.Header()
		var sb strings.Builder
		for _, col := range hdr.Columns {
			sb.WriteString(col)
			sb.WriteString("\t")
		}
		sb.WriteString("\n")
		for _, row := range rows {
			for _, col := range hdr.Columns {
				if row.Get(hdr, col) == nil {
					fmt.Println(col, "is nil in", row)
				}
				sb.WriteString(row.Get(hdr, col).String())
				sb.WriteString("\t")
			}
			sb.WriteString("\n")
		}
		s := sb.String()
		s = strings.ReplaceAll(s, `"`, "'")
		return s
	}
	test := func(query, strategy, expected string) {
		t.Helper()
		q := ParseQuery(query)
		tran := sizeTran{db.NewReadTran()}
		q = Setup(q, ReadMode, tran)
		qs := strings.ReplaceAll(q.String(), `"`, "'")
		assert.T(t).This(qs).Is(strategy)
		assert.T(t).Msg("forward").This(get(q, rt.Next)).Like(expected)
		assert.T(t).Msg("reverse").This(get(q, rt.Prev)).Like(expected)
	}
	test("indexes",
		"indexes",
		`table		columns			key
        'alias'		'id'			true
        'co'		'tnum'			true
        'cus'		'cnum'			true
        'cus'		'abbrev'		true
        'customer'	'id'			true
        'dates'		'date'			true
        'hist'		'date'			false
        'hist'		'date,item,id'	true
        'hist2'		'date'			true
        'hist2'		'id'			false
        'inven'		'item'			true
        'supplier'	'supplier'		true
        'supplier'	'city'			false
        'task'		'tnum'			true
        'trans'		'item'			false
        'trans'		'date,item,id'	true`)
	test("columns",
		"columns",
		`table		column
        'alias'		'id'
        'alias'		'name2'
        'co'		'tnum'
        'co'		'signed'
        'columns'	'table'
        'columns'	'column'
        'cus'		'cnum'
        'cus'		'abbrev'
        'cus'		'name'
        'customer'	'id'
        'customer'	'name'
        'customer'	'city'
        'dates'		'date'
        'hist'		'date'
        'hist'		'item'
        'hist'		'id'
        'hist'		'cost'
        'hist2'		'date'
        'hist2'		'item'
        'hist2'		'id'
        'hist2'		'cost'
        'indexes'	'table'
        'indexes'	'columns'
        'indexes'	'key'
        'inven'		'item'
        'inven'		'qty'
        'supplier'	'supplier'
        'supplier'	'name'
        'supplier'	'city'
        'tables'	'table'
        'tables'	'tablename'
        'tables'	'nrows'
        'tables'	'totalsize'
        'task'		'tnum'
        'task'		'cnum'
        'trans'		'item'
        'trans'		'id'
        'trans'		'cost'
        'trans'		'date'`)
	test("tables",
		"tables",
		`table   tablename       nrows   totalsize
        'alias' 'alias' 2       25
        'co'    'co'    4       55
        'columns'       'columns'       0       0
        'cus'   'cus'   4       64
        'customer'      'customer'      4       98
        'dates' 'dates' 4       52
        'hist'  'hist'  4       91
        'hist2' 'hist2' 3       68
        'indexes'       'indexes'       16      0
        'inven' 'inven' 3       42
        'supplier'      'supplier'      4       128
        'tables'        'tables'        14      0
        'task'  'task'  8       95
        'trans' 'trans' 4       92`)
	test("customer",
		"customer^(id)",
		`id	name	city
		'a'	'axon'	'saskatoon'
		'c'	'calac'	'calgary'
		'e'	'emerald'	'vancouver'
		'i'	'intercon'	'saskatoon'`)
	test("hist",
		"hist^(date)",
		`date	item	id	cost
		970101	'disk'	'a'	100
		970101	'disk'	'e'	200
		970102	'mouse'	'c'	200
		970103	'pencil'	'e'	300`)
	test("trans",
		"trans^(item)",
		`item		id	cost	date
		'disk'		'a'	100	970101
		'eraser'	'c'	150	970201
		'mouse'		'e'	200	960204
		'mouse'		'c'	200	970101`)

	// rename
	test("trans rename id to code, date to when",
		"trans^(item) RENAME id to code, date to when",
		`item	code	cost	when
		'disk'	'a'	100	970101
		'eraser'	'c'	150	970201
		'mouse'	'e'	200	960204
		'mouse'	'c'	200	970101`)

	// sort
	test("customer sort id",
		"customer^(id)",
		`id	name	city
		'a'	'axon'	'saskatoon'
		'c'	'calac'	'calgary'
		'e'	'emerald'	'vancouver'
		'i'	'intercon'	'saskatoon'`)
	test("customer sort reverse id",
		"customer^(id) reverse",
		`id	name	city
		'i'	'intercon'	'saskatoon'
		'e'	'emerald'	'vancouver'
		'c'	'calac'	'calgary'
		'a'	'axon'	'saskatoon'`)
	test("customer sort city",
		"customer^(id) TEMPINDEX(city)",
		`id	name	city
		'c'	'calac'	'calgary'
		'a'	'axon'	'saskatoon'
		'i'	'intercon'	'saskatoon'
		'e'	'emerald'	'vancouver'`)
	test("customer sort reverse city",
		"customer^(id) TEMPINDEX(city) reverse",
		`id	name	city
		'e'	'emerald'	'vancouver'
		'i'	'intercon'	'saskatoon'
		'a'	'axon'	'saskatoon'
		'c'	'calac'	'calgary'`)
	test("task sort cnum, tnum",
		"task^(tnum) TEMPINDEX(cnum,tnum)",
		`tnum	cnum
		100	1
		104	1
		101	2
		105	2
		102	3
		106	3
		103	4
		107	4`)
	test("customer times inven sort qty, id",
		"(customer^(id) TIMES inven^(item)) TEMPINDEX(qty,id)",
		`id	name	city	item	qty
		'a'	'axon'	'saskatoon'	'mouse'	2
		'c'	'calac'	'calgary'	'mouse'	2
		'e'	'emerald'	'vancouver'	'mouse'	2
		'i'	'intercon'	'saskatoon'	'mouse'	2
		'a'	'axon'	'saskatoon'	'disk'	5
		'c'	'calac'	'calgary'	'disk'	5
		'e'	'emerald'	'vancouver'	'disk'	5
		'i'	'intercon'	'saskatoon'	'disk'	5
		'a'	'axon'	'saskatoon'	'pencil'	7
		'c'	'calac'	'calgary'	'pencil'	7
		'e'	'emerald'	'vancouver'	'pencil'	7
		'i'	'intercon'	'saskatoon'	'pencil'	7`)
	test("customer extend up = city.Upper() sort up",
		"customer^(id) EXTEND up = city.Upper() TEMPINDEX(up)",
		`id	name	city	up
		'c'	'calac'	'calgary'	'CALGARY'
		'a'	'axon'	'saskatoon'	'SASKATOON'
		'i'	'intercon'	'saskatoon'	'SASKATOON'
		'e'	'emerald'	'vancouver'	'VANCOUVER'`)
	test("trans minus hist sort id, cost",
		"(trans^(item) MINUS hist^(date,item,id)) TEMPINDEX(id,cost)",
		`item		id	cost	date
		'eraser'	'c'	150		970201
		'mouse'		'c'	200		970101
		'mouse'		'e'	200		960204`)

	// project
	test("customer project city, id",
		"customer^(id) PROJECT-COPY city, id",
		`city	id
		'saskatoon'	'a'
		'calgary'	'c'
		'vancouver'	'e'
		'saskatoon'	'i'`)
	test("supplier project city",
		"supplier^(city) PROJECT-SEQ city",
		`city
		'calgary'
		'saskatoon'
		'vancouver'`)
	test("trans project item",
		"trans^(item) PROJECT-SEQ item",
		`item
		'disk'
		'eraser'
		'mouse'`)
	test("customer project city",
		"customer^(id) PROJECT-HASH city",
		`city
		'saskatoon'
		'calgary'
		'vancouver'`)

	// extend
	test("trans extend newcost = cost * 1.1",
		"trans^(item) EXTEND newcost = cost * 1.1",
		`item	id	cost	date	newcost
		'disk'	'a'	100	970101	110
		'eraser'	'c'	150	970201	165
		'mouse'	'e'	200	960204	220
		'mouse'	'c'	200	970101	220`)
	test("trans extend x = cost * 1.1, y = x $ '*'",
		"trans^(item) EXTEND x = cost * 1.1, y = x $ '*'",
		`item	id	cost	date	x	y
		'disk'	'a'	100	970101	110	'110*'
		'eraser'	'c'	150	970201	165	'165*'
		'mouse'	'e'	200	960204	220	'220*'
		'mouse'	'c'	200	970101	220	'220*'`)

	// times
	test("customer times inven",
		"customer^(id) TIMES inven^(item)",
		`id	name	city	item	qty
		'a'	'axon'	'saskatoon'	'disk'	5
		'a'	'axon'	'saskatoon'	'mouse'	2
		'a'	'axon'	'saskatoon'	'pencil'	7
		'c'	'calac'	'calgary'	'disk'	5
		'c'	'calac'	'calgary'	'mouse'	2
		'c'	'calac'	'calgary'	'pencil'	7
		'e'	'emerald'	'vancouver'	'disk'	5
		'e'	'emerald'	'vancouver'	'mouse'	2
		'e'	'emerald'	'vancouver'	'pencil'	7
		'i'	'intercon'	'saskatoon'	'disk'	5
		'i'	'intercon'	'saskatoon'	'mouse'	2
		'i'	'intercon'	'saskatoon'	'pencil'	7`)

	// minus
	test("trans minus trans",
		"trans^(item) MINUS trans^(date,item,id)",
		`item		id	cost	date`)
	test("hist minus hist2",
		"hist^(date) MINUS hist2^(date)",
		`date	item	id	cost
        970101	'disk'	'e'	200
        970102	'mouse'	'c'	200`)
	test("trans minus hist",
		"trans^(item) MINUS hist^(date,item,id)",
		`item		id	cost	date
		'eraser'	'c'	150		970201
		'mouse'		'e'	200		960204
		'mouse'		'c'	200		970101`)
	test("trans minus hist sort date",
		"trans^(date,item,id) MINUS hist^(date,item,id)",
		`item		id	cost	date
		'mouse'		'e'	200		960204
		'mouse'		'c'	200		970101
		'eraser'	'c'	150		970201`)

	// intersect
	test("trans intersect trans",
		"trans^(item) INTERSECT trans^(date,item,id)",
		`item		id	cost	date
		'disk'		'a'	100	970101
		'eraser'	'c'	150	970201
		'mouse'		'e'	200	960204
		'mouse'		'c'	200	970101`)
	test("trans intersect hist",
		"trans^(item) INTERSECT hist^(date,item,id)",
		`item	id	cost	date
		'disk'	'a'	100		970101`)
	test("hist intersect hist2",
		"hist^(date) INTERSECT hist2^(date)",
		`date	item		id	cost
        970101	'disk'		'a'	100
        970103	'pencil'	'e'	300`)

	// union
	test("hist2 union hist",
		"hist2^(date) UNION-LOOKUP hist^(date,item,id)",
		`date	item	 id		cost
		970102	'disk'	 'e'	200
		970101	'disk'	 'a'	100
		970101	'disk'	 'e'	200
		970102	'mouse'	 'c'	200
		970103	'pencil' 'e'	300`)
	test("hist2 union trans",
		"hist2^(date) UNION-LOOKUP trans^(date,item,id)",
		`date	item	 id		cost
        970102	'disk'	 'e'	200
        970103	'pencil' 'e'	300
        960204	'mouse'	 'e'	200
        970101	'disk'	 'a'	100
        970101	'mouse'	 'c'	200
        970201	'eraser' 'c'	150`)
	test("alias union alias",
		"alias^(id) UNION-MERGE alias^(id)",
		`id name2
        'a'	'abc'
		'c'	'trical'`)
	test("trans union hist",
		"trans^(date,item,id) UNION-MERGE hist^(date,item,id)",
		`item	id	cost	date
		'mouse'	'e'	200	960204
		'disk'	'a'	100	970101
		'disk'	'e'	200	970101
		'mouse'	'c'	200	970101
		'mouse'	'c'	200	970102
		'pencil'	'e'	300	970103
		'eraser'	'c'	150	970201`)
	test("trans union trans sort id",
		"(trans^(date,item,id) UNION-MERGE trans^(date,item,id)) TEMPINDEX(id)",
		`item		id	cost	date
		'disk'		'a'	100	970101
		'mouse'		'c'	200	970101
		'eraser'	'c'	150	970201
		'mouse'		'e'	200	960204`)

	// join
	test("customer join alias",
		"customer^(id) JOIN 1:1 by(id) alias^(id)",
		`id	name	city	name2
        'a'	'axon'	'saskatoon'	'abc'
        'c'	'calac'	'calgary'	'trical'`)
	test("trans join inven",
		"inven^(item) JOIN 1:n by(item) trans^(item)",
		`item	qty	id	cost	date
		'disk'	5	'a'	100	970101
		'mouse'	2	'e'	200	960204
		'mouse'	2	'c'	200	970101`)
	test("customer leftjoin alias",
		"customer^(id) LEFTJOIN 1:1 by(id) alias^(id)",
		`id	name	city	name2
        'a'	'axon'	'saskatoon'	'abc'
        'c'	'calac'	'calgary'	'trical'
        'e'	'emerald'	'vancouver'	''
        'i'	'intercon'	'saskatoon'	''`)
	test("inven leftjoin trans",
		"inven^(item) LEFTJOIN 1:n by(item) trans^(item)",
		`item	qty	id	cost date
		'disk'	 5	'a'	100	 970101
		'mouse'	 2	'e'	200	 960204
		'mouse'	 2	'c'	200	 970101
		'pencil' 7	''	''	 ''`)
	test("customer leftjoin hist2",
		"customer^(id) LEFTJOIN 1:n by(id) hist2^(id)",
		`id	name	city	date	item	cost
		'a'	'axon'	'saskatoon'	970101	'disk'	100
		'c'	'calac'	'calgary'	''	''	''
		'e'	'emerald'	'vancouver'	970102	'disk'	200
		'e'	'emerald'	'vancouver'	970103	'pencil'	300
		'i'	'intercon'	'saskatoon'	''	''	''`)
	test("hist join customer",
		"customer^(id) JOIN 1:n by(id) (hist^(date) TEMPINDEX(id))",
		`id name	  city			date	item	 cost
		'a'	'axon'	  'saskatoon'	970101	'disk'	 100
		'c'	'calac'	  'calgary'		970102	'mouse'	 200
		'e'	'emerald' 'vancouver'	970101	'disk'	 200
		'e'	'emerald' 'vancouver'	970103	'pencil' 300`)
	test("customer leftjoin (alias where name2 is 'abc')",
		"customer^(id) LEFTJOIN 1:1 by(id) (alias^(id) WHERE name2 is 'abc')",
		`id	name	city	name2
        'a'	'axon'	'saskatoon'	'abc'
        'c'	'calac'	'calgary'	''
        'e'	'emerald'	'vancouver'	''
        'i'	'intercon'	'saskatoon'	''`)
	test("customer leftjoin (alias where id is 'c')",
		"customer^(id) LEFTJOIN 1:1 by(id) (alias^(id) WHERE*1 id is 'c')",
		`id	name	city	name2
        'a'	'axon'	'saskatoon'	''
        'c'	'calac'	'calgary'	'trical'
        'e'	'emerald'	'vancouver'	''
        'i'	'intercon'	'saskatoon'	''`)

	// where
	test("customer where id > 'd'", // range
		"customer^(id) WHERE id > 'd'",
		`id	name	city
		'e'	'emerald'	'vancouver'
		'i'	'intercon'	'saskatoon'`)
	test("customer where id > 'd' and id < 'j'", // range
		"customer^(id) WHERE id > 'd' and id < 'j'",
		`id	name	city
		'e'	'emerald'	'vancouver'
		'i'	'intercon'	'saskatoon'`)
	test("customer where id is 'e'", // point
		"customer^(id) WHERE*1 id is 'e'",
		`id	name	city
		'e'	'emerald'	'vancouver'`)
	test("customer where id is 'd'", // point
		"customer^(id) WHERE*1 id is 'd'",
		`id	name	city`)
	test("inven where qty > 0", // filter
		"inven^(item) WHERE qty > 0",
		`item	qty
		'disk'	5
		'mouse'	2
		'pencil'	7`)
	test("inven where item =~ 'i'", // filter
		"inven^(item) WHERE item =~ 'i'",
		`item	qty
		'disk'	5
		'pencil'	7`)
	test("inven where item in ('disk', 'mouse', 'disk', 'pencil')", // points
		"inven^(item) WHERE item in ('disk', 'mouse', 'disk', 'pencil')",
		`item	qty
		'disk'	5
		'mouse'	2
		'pencil'	7`)
	test("inven where item <= 'e' or item >= 'p'", // filter
		"inven^(item) WHERE item <= 'e' or item >= 'p'",
		`item	qty
		'disk'	5
		'pencil'	7`)
	test("cus where cnum is 2 and abbrev is 'b'", // points
		"cus^(cnum) WHERE*1 cnum is 2 and abbrev is 'b'",
		`cnum	abbrev	name
		2	'b'	'bill'`)
	test("cus where cnum is 2 and abbrev >= 'b' and abbrev < 'c'", // point
		"cus^(cnum) WHERE*1 cnum is 2 and abbrev >= 'b' and abbrev < 'c'",
		`cnum	abbrev	name
		2	'b'	'bill'`)
	test("hist where date in (970101, 970102) and item < 'z'", // ranges
		"hist^(date) WHERE date in (970101, 970102) and item < 'z'",
		`date	item	id	cost
        970101	'disk'	'a'	100
        970101	'disk'	'e'	200
        970102	'mouse'	'c'	200`)
	test("customer where id not in ()",
		"customer^(id) WHERE not id in ()",
		`id	name	city
		'a'	'axon'	'saskatoon'
		'c'	'calac'	'calgary'
		'e'	'emerald'	'vancouver'
		'i'	'intercon'	'saskatoon'`)

	// summarize
	test("customer summarize count",
		"customer^(id) SUMMARIZE-TBL count = count",
		`count
		1000`)
	test("hist summarize max date",
		"hist^(date) SUMMARIZE-IDX max_date = max date",
		`max_date
		970103`)
	test("customer summarize max id",
		"customer^(id) SUMMARIZE-IDX* max_id = max id",
		`id	name		city		max_id
		'i'	'intercon'	'saskatoon' 'i'`)
	test("hist summarize item, total cost",
		"hist^(date) SUMMARIZE-MAP item, total_cost = total cost",
		`item		total_cost
		'disk'		300
		'mouse'		200
		'pencil'	300`)
	test("hist summarize item, total cost, max id, average cost",
		"hist^(date) SUMMARIZE-MAP item, total_cost = total cost, "+
			"max_id = max id, average_cost = average cost",
		`item		total_cost	max_id	average_cost
        'disk'		300			'e'		150
        'mouse'		200			'c'		200
        'pencil'	300			'e'		300`)
	test("hist summarize item, total cost sort total_cost, item",
		"hist^(date) SUMMARIZE-MAP item, total_cost = total cost "+
			"TEMPINDEX(total_cost,item)",
		`item		total_cost
		'mouse'		200
		'disk'		300
		'pencil'	300`)
	test("customer summarize max name",
		"customer^(id) SUMMARIZE-SEQ max_name = max name",
		`max_name
		'intercon'`)
	test("hist summarize min cost, average cost, max cost, sum = total cost",
		"hist^(date) SUMMARIZE-SEQ min_cost = min cost, "+
			"average_cost = average cost, max_cost = max cost, sum = total cost",
		`min_cost	average_cost	max_cost	sum
		100			200				300			800`)
	test("hist summarize item, total cost, count",
		"hist^(date) SUMMARIZE-MAP item, total_cost = total cost, count = count",
		`item		total_cost	count
		'disk'		300			2
		'mouse'		200			1
		'pencil'	300			1`)
	test("hist summarize date, item, max id",
		"hist^(date,item,id) SUMMARIZE-SEQ* date, item, max_id = max id",
		`date	item		id	cost	max_id
        970101	'disk'		'e'	200		'e'
        970102	'mouse'		'c'	200		'c'
        970103	'pencil'	'e'	300		'e'`)
	test("hist summarize date, list id",
		"hist^(date) SUMMARIZE-SEQ date, list_id = list id",
		`date	list_id
		970101	#('a', 'e')
		970102	#('c')
		970103	#('e')`)
	test("hist summarize list id",
		"hist^(date) SUMMARIZE-SEQ list_id = list id",
		`list_id
		#('a', 'c', 'e')`)
	test("cus summarize max cnum sort name",
		"cus^(cnum) SUMMARIZE-IDX* max_cnum = max cnum",
		`cnum	abbrev	name	max_cnum
		4		'd'		'dick'	4`)
	test("supplier summarize min city",
		"supplier^(city) SUMMARIZE-IDX min_city = min city",
		`min_city
		'calgary'`)
	test("supplier summarize max city",
		"supplier^(city) SUMMARIZE-IDX max_city = max city",
		`max_city
		'vancouver'`)
	test("supplier summarize min city, max city",
		"supplier^(supplier) SUMMARIZE-SEQ min_city = min city, max_city = max city",
		`min_city	max_city
		'calgary'	'vancouver'`)
	test("hist summarize max cost",
		"hist^(date) SUMMARIZE-SEQ max_cost = max cost",
		`max_cost
		300`)

	// tempindex
	test("tables intersect columns",
		"tables INTERSECT (columns TEMPINDEX(table,column))",
		`table`)
	test("tables minus tables",
		"tables MINUS (tables TEMPINDEX(tablename))",
		`table	tablename	nrows	totalsize`)
}
