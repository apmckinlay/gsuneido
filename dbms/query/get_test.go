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
		Setup(q, readMode, testTran{})
		row := q.(*Table).Lookup(key)
		assert.T(t).This(fmt.Sprint(row)).Is(expected)
	}
	test("tables", key(123), "[<123>]")
	test("columns", key(12, 34), "[<12, 34>]")
}

func TestTableGet(t *testing.T) {
	db := testDb()
	defer db.Close()
	reverse := func(rows []rt.Row) {
		for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
			rows[i], rows[j] = rows[j], rows[i]
		}
	}
	get := func(q Query, dir rt.Dir) string {
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
		Setup(q, readMode, tran)
		assert.T(t).This(q.String()).Is(strategy)
		assert.T(t).This(get(q, rt.Next)).Like(expected)
		assert.T(t).This(get(q, rt.Prev)).Like(expected)
	}
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

	test("trans rename id to code, date to when",
		"trans^(item) RENAME id to code, date to when",
		`item	code	cost	when
		'disk'	'a'	100	970101
		'eraser'	'c'	150	970201
		'mouse'	'e'	200	960204
		'mouse'	'c'	200	970101`)

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
	test("trans project item", // sequential
		"trans^(item) PROJECT-SEQ item",
		`item
		'disk'
		'eraser'
		'mouse'`)
	// test("customer project city", // lookup
	// "",
	// 	`city
	// 	'saskatoon'
	// 	'calgary'
	// 	'vancouver'`)

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

	// test("trans intersect hist",
	// 	`item	id	cost	date
	// 	'disk'	'a'	100	970101`)

	// test("trans minus hist sort date",
	// 	`item	id	cost	date
	// 	'mouse'	'e'	200	960204
	// 	'mouse'	'c'	200	970101
	// 	'eraser'	'c'	150	970201`)
	// test("(trans minus hist) where id = 9",
	// 	`item	id	cost	date`)

	// test("trans union hist", // merge
	// 	`item	id	cost	date
	// 	'mouse'	'e'	200	960204
	// 	'disk'	'a'	100	970101
	// 	'disk'	'e'	200	970101
	// 	'mouse'	'c'	200	970101
	// 	'mouse'	'c'	200	970102
	// 	'pencil'	'e'	300	970103
	// 	'eraser'	'c'	150	970201`)
	// test("hist2 union hist", // lookup
	// 	`date	item	id	cost
	// 	970102	'disk'	'e'	200
	// 	970101	'disk'	'a'	100
	// 	970101	'disk'	'e'	200
	// 	970102	'mouse'	'c'	200
	// 	970103	'pencil'	'e'	300`)

	// test("hist join customer",
	// 	`date	item	id	cost	name	city
	// 	970101	'disk'	'a'	100	'axon'	'saskatoon'
	// 	970101	'disk'	'e'	200	'emerald'	'vancouver'
	// 	970102	'mouse'	'c'	200	'calac'	'calgary'
	// 	970103	'pencil'	'e'	300	'emerald'	'vancouver'`)
	// test("trans join inven",
	// 	`item	qty	id	cost	date
	// 	'disk'	5	'a'	100	970101
	// 	'mouse'	2	'e'	200	960204
	// 	'mouse'	2	'c'	200	970101`)
	// test("customer join alias",
	// 	`id	name2	name	city
	// 	'a'	'abc'	'axon'	'saskatoon'
	// 	'c'	'trical'	'calac'	'calgary'`)
	// test("customer join supplier",
	// 	`supplier	name	city	id`)
	// test("inven leftjoin trans",
	// 	`item	qty	id	cost	date
	// 	'disk'	5	'a'	100	970101
	// 	'mouse'	2	'e'	200	960204
	// 	'mouse'	2	'c'	200	970101
	// 	'pencil'	7	''	''	''`)
	// test("customer leftjoin hist2",
	// 	`id	name	city	date	item	cost
	// 	'a'	'axon'	'saskatoon'	970101	'disk'	100
	// 	'c'	'calac'	'calgary'	''	''	''
	// 	'e'	'emerald'	'vancouver'	970102	'disk'	200
	// 	'e'	'emerald'	'vancouver'	970103	'pencil'	300
	// 	'i'	'intercon'	'saskatoon'	''	''	''`)

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
	test("inven where item in ('disk', 'mouse', 'pencil')", // points
		`inven^(item) WHERE item in ("disk", "mouse", "pencil")`,
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

	// test("hist summarize count", // by is empty
	// 	`count
	// 	4`)
	// test("hist summarize min cost, average cost, max cost, sum = total cost",
	// 	`min_cost	average_cost	max_cost	sum
	// 	100	200	300	800`)
	// test("hist summarize item, total cost",
	// 	`item	total_cost
	// 	'disk'	300
	// 	'mouse'	200
	// 	'pencil'	300`)
	// test("hist summarize date, list id",
	// 	`date	list_id
	// 	970101	#('a', 'e')
	// 	970102	#('c')
	// 	970103	#('e')`)
	// test("hist summarize list id",
	// 	`list_id
	// 	#('a', 'c', 'e')`)
	// test("cus summarize max cnum sort name", // key so whole record
	// 	`cnum	abbrev	name	max_cnum
	// 	4	'd'	'dick'	4`)
	// test("supplier summarize min city", // indexed
	// 	`min_city
	// 	'calgary'`)
	// test("supplier summarize max city", // indexed
	// 	`max_city
	// 	'vancouver'`)
	// test("supplier summarize min city, max city",
	// 	`min_city	max_city
	// 	'calgary'	'vancouver'`)
	// test("hist summarize max cost", // not indexed
	// 	`max_cost
	// 	300`)
}
