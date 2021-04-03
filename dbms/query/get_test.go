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
	test := func(query, expected string) {
		q := ParseQuery(query)
		tran := db.NewReadTran()
		Setup(q, readMode, tran)
		assert.T(t).This(get(q, rt.Next)).Like(expected)
		assert.T(t).This(get(q, rt.Prev)).Like(expected)
	}
	test("customer",
		`id	name	city
		'a'	'axon'	'saskatoon'
		'c'	'calac'	'calgary'
		'e'	'emerald'	'vancouver'
		'i'	'intercon'	'saskatoon'`)
	test("customer sort id",
		`id	name	city
		'a'	'axon'	'saskatoon'
		'c'	'calac'	'calgary'
		'e'	'emerald'	'vancouver'
		'i'	'intercon'	'saskatoon'`)
	test("customer sort reverse id",
		`id	name	city
		'i'	'intercon'	'saskatoon'
		'e'	'emerald'	'vancouver'
		'c'	'calac'	'calgary'
		'a'	'axon'	'saskatoon'`)
	test("hist",
		`date	item	id	cost
		970101	'disk'	'a'	100
		970101	'disk'	'e'	200
		970102	'mouse'	'c'	200
		970103	'pencil'	'e'	300`)
	test("trans",
		`item		id	cost	date
		'disk'		'a'	100	970101
		'eraser'	'c'	150	970201
		'mouse'		'e'	200	960204
		'mouse'		'c'	200	970101`)
}
