// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
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
	get := func(q Query, dir Dir) string {
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
		hdr := q.Header()
		var sb strings.Builder
		for _, col := range hdr.Columns {
			sb.WriteString(col)
			sb.WriteString("\t")
		}
		sb.WriteString("\n")
		for _, row := range rows {
			for _, col := range hdr.Columns {
				if row.GetVal(hdr, col, nil, nil) == nil {
					fmt.Println(col, "is nil in", row)
				}
				sb.WriteString(row.GetVal(hdr, col, nil, nil).String())
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
		tran := sizeTran{db.NewReadTran()}
		q := ParseQuery(query, tran, nil)

		query2 := format(0, q, 0)
		q2 := ParseQuery(query2, tran, nil)
		assert.This(format(0, q2, 0)).Is(query2)

		q, _, _ = Setup(q, ReadMode, tran)
		qs := strings.ReplaceAll(q.String(), `"`, "'")
		assert.T(t).This(qs).Is(strategy)
		assert.T(t).Msg("forward:", query).This(get(q, Next)).Like(expected)
		assert.T(t).Msg("reverse:", query).This(get(q, Prev)).Like(expected)
	}
	test("indexes project table, columns, key",
		"indexes PROJECT-COPY table,columns,key",
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
		`table   column  field
        'alias' 'id'    0
        'alias' 'name2' 1
        'co'    'tnum'  0
        'co'    'signed'        1
        'columns'       'table' 0
        'columns'       'column'        1
        'columns'       'field' 2
        'cus'   'cnum'  0
        'cus'   'abbrev'        1
        'cus'   'name'  2
        'customer'      'id'    0
        'customer'      'name'  1
        'customer'      'city'  2
        'dates' 'date'  0
        'hist'  'date'  0
        'hist'  'item'  1
        'hist'  'id'    2
        'hist'  'cost'  3
        'hist2' 'date'  0
        'hist2' 'item'  1
        'hist2' 'id'    2
        'hist2' 'cost'  3
        'indexes'       'table' 0
        'indexes'       'columns'       1
        'indexes'       'key'   2
        'indexes'		'fktable'	3
        'indexes'		'fkcolumns'	4
        'indexes'		'fkmode'	5
        'inven' 'item'  0
        'inven' 'qty'   1
        'supplier'      'supplier'      0
        'supplier'      'name'  1
        'supplier'      'city'  2
        'tables'        'table' 0
        'tables'        'tablename'     1
        'tables'        'nrows' 2
        'tables'        'totalsize'     3
        'task'  'tnum'  0
        'task'  'cnum'  1
        'trans' 'item'  0
        'trans' 'id'    1
        'trans' 'cost'  2
        'trans' 'date'  3
        'views' 'view_name'     0
        'views' 'view_definition'       1`)
	test("tables",
		"tables",
		`table   tablename       nrows   totalsize
        'alias' 'alias' 2       25
        'co'    'co'    4       55
        'columns'       'columns'       45       0
        'cus'   'cus'   4       64
        'customer'      'customer'      4       98
        'dates' 'dates' 4       52
        'hist'  'hist'  4       91
        'hist2' 'hist2' 3       68
        'indexes'       'indexes'       16      0
        'inven' 'inven' 3       42
        'supplier'      'supplier'      4       128
        'tables'        'tables'        15      0
        'task'  'task'  8       95
        'trans' 'trans' 4       92
		'views' 'views' 0       0`)
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
	test("customer extend up = city[1..] sort up",
		"customer^(id) EXTEND up = city[1..] TEMPINDEX(up)",
		`id	name	city	up
		'c'	'calac'	'calgary'	'algary'
		'e'	'emerald'	'vancouver'	'ancouver'
		'a'	'axon'	'saskatoon'	'askatoon'
		'i'	'intercon'	'saskatoon'	'askatoon'`)
	test("trans minus hist sort id, cost",
		"(trans^(item) MINUS hist^(date,item,id)) TEMPINDEX(id,cost)",
		`item		id	cost	date
		'eraser'	'c'	150		970201
		'mouse'		'c'	200		970101
		'mouse'		'e'	200		960204`)
	test("customer rename id to id_new sort id_new",
		"customer^(id) RENAME id to id_new",
		`id_new	name	city
		'a'		'axon'	'saskatoon'
		'c'		'calac'	'calgary'
		'e'		'emerald'	'vancouver'
		'i'		'intercon'	'saskatoon'`)

	// project
	test("customer project city, id",
		"customer^(id) PROJECT-COPY city,id",
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
		"customer^(id) PROJECT-MAP city",
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
	test("(trans minus trans) where item = 0",
		"trans^(item) WHERE item is 0 MINUS "+
			"(trans^(date,item,id) WHERE item is 0)",
		`item		id	cost	date`)
	test("inven minus (inven where item = 'mouse')",
		"inven^(item) MINUS (inven^(item) WHERE*1 item is 'mouse')",
		`item		qty
		'disk'		5
		'pencil'	7`)

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
	test("(hist intersect hist2) where cost = 100",
		"hist^(date) WHERE cost is 100 INTERSECT (hist2^(date) WHERE cost is 100)",
		`date	item		id	cost
        970101	'disk'		'a'	100`)

	// union
	test("hist2 union hist",
		"hist2^(date) UNION-LOOKUP(date,item,id) hist^(date,item,id)",
		`date	item	 id		cost
		970102	'disk'	 'e'	200
		970101	'disk'	 'a'	100
		970101	'disk'	 'e'	200
		970102	'mouse'	 'c'	200
		970103	'pencil' 'e'	300`)
	test("hist2 union trans",
		"hist2^(date) UNION-LOOKUP(date,item,id) trans^(date,item,id)",
		`date	item	 id		cost
        970102	'disk'	 'e'	200
        970103	'pencil' 'e'	300
        960204	'mouse'	 'e'	200
        970101	'disk'	 'a'	100
        970101	'mouse'	 'c'	200
        970201	'eraser' 'c'	150`)
	test("alias union alias",
		"alias^(id) UNION-MERGE(id) alias^(id)",
		`id name2
        'a'	'abc'
		'c'	'trical'`)
	test("trans union hist",
		"trans^(date,item,id) UNION-MERGE(date,item,id) hist^(date,item,id)",
		`item	id	cost	date
		'mouse'	'e'	200	960204
		'disk'	'a'	100	970101
		'disk'	'e'	200	970101
		'mouse'	'c'	200	970101
		'mouse'	'c'	200	970102
		'pencil'	'e'	300	970103
		'eraser'	'c'	150	970201`)
	test("trans union trans sort id",
		"(trans^(date,item,id) UNION-MERGE(date,item,id) "+
			"trans^(date,item,id)) TEMPINDEX(id)",
		`item		id	cost	date
		'disk'		'a'	100	970101
		'mouse'		'c'	200	970101
		'eraser'	'c'	150	970201
		'mouse'		'e'	200	960204`)
	test("(hist2 rename cost to amt) union (trans rename cost to amt)",
		"hist2^(date) RENAME cost to amt UNION-LOOKUP(date,item,id) "+
			"(trans^(date,item,id) RENAME cost to amt)",
		`date	item	 id		amt
        970102	'disk'	 'e'	200
        970103	'pencil' 'e'	300
        960204	'mouse'	 'e'	200
        970101	'disk'	 'a'	100
        970101	'mouse'	 'c'	200
        970201	'eraser' 'c'	150`)
	test("(hist2 extend x = 1) union (trans extend x = 1)",
		"hist2^(date) EXTEND x = 1 UNION-LOOKUP(date,item,id) "+
			"(trans^(date,item,id) EXTEND x = 1)",
		`date	item	 id		cost  x
        970102	'disk'	 'e'	200   1
        970103	'pencil' 'e'	300   1
        960204	'mouse'	 'e'	200   1
        970101	'disk'	 'a'	100   1
        970101	'mouse'	 'c'	200   1
        970201	'eraser' 'c'	150   1`)
	test("(co where tnum = 100 remove signed) union "+
		"(co where tnum = 100 remove signed)",
		"co^(tnum) WHERE*1 tnum is 100 PROJECT-COPY tnum UNION-MERGE() "+
			"(co^(tnum) WHERE*1 tnum is 100 PROJECT-COPY tnum)",
		`tnum
        100`)
	test("(co where tnum = 100 remove tnum) union "+
		"(co where tnum = 100 remove tnum)",
		"co^(tnum) WHERE*1 tnum is 100 PROJECT-COPY signed UNION-MERGE() "+
			"(co^(tnum) WHERE*1 tnum is 100 PROJECT-COPY signed)",
		`signed
        990101`)
	test("((co where tnum = 100) union (co where tnum = 102)) union "+
		"(co where tnum = 100)",
		"co^(tnum) WHERE*1 tnum is 100 UNION-LOOKUP(tnum) "+
			"(co^(tnum) WHERE*1 tnum is 100 UNION-DISJOINT(tnum)-MERGE(tnum) "+
			"(co^(tnum) WHERE*1 tnum is 102))",
		`tnum	signed
        100		990101
        102		990102`)
	test("(((co where tnum = 100) union (co where tnum = 102)) remove tnum)"+
		" union "+
		"(((co where tnum = 104) union (co where tnum = 106)) remove tnum)",
		"(co^(tnum) WHERE*1 tnum is 100 UNION-DISJOINT(tnum)-MERGE(signed) (co^(tnum) WHERE*1 tnum is 102)) PROJECT-SEQ signed UNION-MERGE(signed) ((co^(tnum) WHERE*1 tnum is 104 UNION-DISJOINT(tnum)-MERGE(signed) (co^(tnum) WHERE*1 tnum is 106)) PROJECT-SEQ signed)",
		`signed
        990101
        990102
        990103
        990104`)
	test(`((co where tnum = 104 remove tnum) union (co where tnum = 106 remove tnum))
		union
		((co where tnum = 104 remove tnum) union (co where tnum = 106 remove tnum))`,
		"(co^(tnum) WHERE*1 tnum is 104 PROJECT-COPY signed UNION-MERGE(signed) (co^(tnum) WHERE*1 tnum is 106 PROJECT-COPY signed)) UNION-MERGE(signed) (co^(tnum) WHERE*1 tnum is 104 PROJECT-COPY signed UNION-MERGE(signed) (co^(tnum) WHERE*1 tnum is 106 PROJECT-COPY signed))",
		`signed
        990103
        990104`)

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
		"hist^(date) JOIN n:1 by(id) customer^(id)",
		`date	item	id	cost	name	city
		970101	'disk'	'a'	100	'axon'	'saskatoon'
		970101	'disk'	'e'	200	'emerald'	'vancouver'
		970102	'mouse'	'c'	200	'calac'	'calgary'
		970103	'pencil'	'e'	300	'emerald'	'vancouver'`)
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
	test("customer leftjoin (trans where date = 970101 and item = 'mouse')",
		"customer^(id) LEFTJOIN 1:1 by(id) "+
			"(trans^(date,item,id) WHERE date is 970101 and item is 'mouse')",
		`id	name	city	item	cost	date
        'a'	'axon'	'saskatoon'	''	''	''
        'c'	'calac'	'calgary'	'mouse'	200	970101
        'e'	'emerald'	'vancouver'	''	''	''
        'i'	'intercon'	'saskatoon'	''	''	''`)

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
	test("inven where item in (1, 'disk', 'mouse', 2, 'disk', 'pencil')", // points
		"inven^(item) WHERE item in (1, 'disk', 'mouse', 2, 'disk', 'pencil')",
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
	test("customer where id not in ('z')",
		"customer^(id) WHERE not id is 'z'",
		`id	name	city
		'a'	'axon'	'saskatoon'
		'c'	'calac'	'calgary'
		'e'	'emerald'	'vancouver'
		'i'	'intercon'	'saskatoon'`)

	// summarize
	test("customer summarize count",
		"customer^(id) SUMMARIZE-TBL count",
		`count
		1000`)
	test("hist summarize max date",
		"hist^(date) SUMMARIZE-IDX max date",
		`max_date
		970103`)
	test("customer summarize max id",
		"customer^(id) SUMMARIZE-IDX* max id",
		`id	name		city		max_id
		'i'	'intercon'	'saskatoon' 'i'`)
	test("hist summarize item, total cost sort item",
		"hist^(date) SUMMARIZE-MAP item, total cost TEMPINDEX(item)",
		`item		total_cost
		'disk'		300
		'mouse'		200
		'pencil'	300`)
	test("hist summarize item, total cost, max id, average cost sort item",
		"hist^(date) SUMMARIZE-MAP item, total cost, "+
			"average cost, max id TEMPINDEX(item)",
		`item		total_cost	average_cost max_id
        'disk'		300		    150			 'e'
        'mouse'		200		    200			 'c'
        'pencil'	300		    300			 'e'`)
	test("hist summarize item, total cost sort total_cost, item",
		"hist^(date) SUMMARIZE-MAP item, total cost "+
			"TEMPINDEX(total_cost,item)",
		`item		total_cost
		'mouse'		200
		'disk'		300
		'pencil'	300`)
	test("customer summarize max name",
		"customer^(id) SUMMARIZE-SEQ max name",
		`max_name
		'intercon'`)
	test("hist summarize min cost, average cost, max cost, sum = total cost",
		"hist^(date) SUMMARIZE-SEQ min cost, "+
			"average cost, max cost, sum = total cost",
		`min_cost	average_cost	max_cost	sum
		100			200				300			800`)
	test("hist summarize item, total cost, count sort item",
		"hist^(date) SUMMARIZE-MAP item, count, total cost TEMPINDEX(item)",
		`item		count	total_cost
		'disk'		2		300
		'mouse'		1		200
		'pencil'	1    	300			`)
	test("inven summarize max item",
		"inven^(item) SUMMARIZE-IDX* max item",
		`item	qty	max_item
		'pencil'	7	'pencil'`)
	test("hist summarize date, list id",
		"hist^(date) SUMMARIZE-SEQ date, list id",
		`date	list_id
		970101	#('a', 'e')
		970102	#('c')
		970103	#('e')`)
	test("hist summarize date, total cost sort total_cost",
		"hist^(date) SUMMARIZE-SEQ date, total cost TEMPINDEX(total_cost)",
		`date	total_cost
        970102	200
        970101	300
        970103	300`)
	test("hist summarize list id",
		"hist^(date) SUMMARIZE-SEQ list id",
		`list_id
		#('a', 'c', 'e')`)
	test("cus summarize max cnum sort name",
		"cus^(cnum) SUMMARIZE-IDX* max cnum",
		`cnum	abbrev	name	max_cnum
		4		'd'		'dick'	4`)
	test("supplier summarize min city",
		"supplier^(city) SUMMARIZE-IDX min city",
		`min_city
		'calgary'`)
	test("supplier summarize max city",
		"supplier^(city) SUMMARIZE-IDX max city",
		`max_city
		'vancouver'`)
	test("supplier summarize min city, max city",
		"supplier^(supplier) SUMMARIZE-SEQ min city, max city",
		`min_city	max_city
		'calgary'	'vancouver'`)
	test("hist summarize max cost",
		"hist^(date) SUMMARIZE-SEQ max cost",
		`max_cost
		300`)

	// tempindex
	test("tables intersect columns",
		"columns INTERSECT (tables TEMPINDEX(table))",
		`table`)
	test("tables minus tables",
		"tables MINUS (tables TEMPINDEX(table))",
		`table	tablename	nrows	totalsize`)

	test("(customer project id) union (customer project id)",
		"customer^(id) PROJECT-COPY id UNION-MERGE(id) (customer^(id) PROJECT-COPY id)",
		`id
		'a'
		'c'
		'e'
		'i'`)
	test("(trans project item,date) union (trans project date,item)",
		"trans^(date,item,id) PROJECT-SEQ item,date UNION-MERGE(item,date) "+
			"(trans^(date,item,id) PROJECT-SEQ date,item)",
		`item		date
		'mouse'		960204
		'disk'		970101
		'mouse'		970101
		'eraser'	970201`)

	test("(customer summarize id, count) union (customer summarize id, count)",
		"customer^(id) SUMMARIZE-SEQ id, count UNION-MERGE(id) "+
			"(customer^(id) SUMMARIZE-SEQ id, count)",
		`id		count
		'a'		1
		'c'		1
		'e'		1
		'i'		1`)
	test("(trans summarize item,date, count) union (trans summarize date,item, count)",
		"trans^(date,item,id) SUMMARIZE-SEQ item, date, count "+
			"UNION-MERGE(item,date) "+
			"(trans^(date,item,id) SUMMARIZE-SEQ date, item, count)",
		`item		date	count
		'mouse'		960204	1
		'disk'		970101	1
		'mouse'		970101	1
		'eraser'	970201	1`)

	test("(customer project id) join (hist project id)",
		"hist^(date) PROJECT-MAP id JOIN 1:1 by(id) (customer^(id) PROJECT-COPY id)",
		`id
		'a'
		'e'
		'c'`)
	test("(customer summarize id,count) join (hist summarize id,count) sort id",
		"customer^(id) SUMMARIZE-SEQ id, count JOIN 1:1 by(id,count) "+
			"(hist^(date) SUMMARIZE-MAP id, count TEMPINDEX(id,count))",
		`id	count
		'a'	1
		'c'	1`)
}
