// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"time"

	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
)

func testDb() *db19.Database {
	db, err := db19.CreateDb(stor.HeapStor(8192))
	ck(err)
	db19.StartConcur(db, 50*time.Millisecond)
	req := func(req string) {
		DoRequest(db, req)
	}
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		DoAction(ut, act)
	}
	req("create customer (id, name, city) key(id)")
	act("insert {id: 'a', name: 'axon', city: 'saskatoon'} into customer")
	act("insert {id: 'c', name: 'calac', city: 'calgary'} into customer")
	act("insert {id: 'e', name: 'emerald', city: 'vancouver'} into customer")
	act("insert {id: 'i', name: 'intercon', city: 'saskatoon'} into customer")

	req("create hist (date, item, id, cost) index(date) key(date,item,id)")
	act("insert{date: 970101, item: 'disk', id: 'a', cost: 100} into hist")
	act("insert{date: 970101, item: 'disk', id: 'e', cost: 200} into hist")
	act("insert{date: 970102, item: 'mouse', id: 'c', cost: 200} into hist")
	act("insert{date: 970103, item: 'pencil', id: 'e', cost: 300} into hist")

	req("create hist2 (date, item, id, cost) key(date) index(id)")
	act("insert{date: 970101, item: 'disk', id: 'a', cost: 100} into hist2")
	act("insert{date: 970102, item: 'disk', id: 'e', cost: 200} into hist2")
	act("insert{date: 970103, item: 'pencil', id: 'e', cost: 300} into hist2")

	req("create trans (item, id, cost, date) index(item) key(date,item,id)")
	act("insert{item: 'mouse', id: 'e', cost: 200, date: 960204} into trans")
	act("insert{item: 'disk', id: 'a', cost: 100, date: 970101} into trans")
	act("insert{item: 'mouse', id: 'c', cost: 200, date: 970101} into trans")
	act("insert{item: 'eraser', id: 'c', cost: 150, date: 970201} into trans")

	req("create supplier (supplier, name, city) key(supplier) index(city)")
	act("insert{supplier: 'mec', name: 'mtnequipcoop', city: 'calgary'} into supplier")
	act("insert{supplier: 'hobo', name: 'hoboshop', city: 'saskatoon'} into supplier")
	act("insert{supplier: 'ebs', name: 'ebssail&sport', city: 'saskatoon'} into supplier")
	act("insert{supplier: 'taiga', name: 'taigaworks', city: 'vancouver'} into supplier")

	req("create inven (item, qty) key(item)")
	act("insert{item: 'disk', qty: 5} into inven")
	act("insert{item: 'mouse', qty:2} into inven")
	act("insert{item: 'pencil', qty: 7} into inven")

	req("create alias(id, name2) key(id)")
	act("insert{id: 'a', name2: 'abc'} into alias")
	act("insert{id: 'c', name2: 'trical'} into alias")

	req("create cus(cnum, abbrev, name) key(cnum) key(abbrev)")
	act("insert { cnum: 1, abbrev: 'a', name: 'axon' } into cus")
	act("insert { cnum: 2, abbrev: 'b', name: 'bill' } into cus")
	act("insert { cnum: 3, abbrev: 'c', name: 'cron' } into cus")
	act("insert { cnum: 4, abbrev: 'd', name: 'dick' } into cus")

	req("create task(tnum, cnum) key(tnum)")
	act("insert { tnum: 100, cnum: 1 } into task")
	act("insert { tnum: 101, cnum: 2 } into task")
	act("insert { tnum: 102, cnum: 3 } into task")
	act("insert { tnum: 103, cnum: 4 } into task")
	act("insert { tnum: 104, cnum: 1 } into task")
	act("insert { tnum: 105, cnum: 2 } into task")
	act("insert { tnum: 106, cnum: 3 } into task")
	act("insert { tnum: 107, cnum: 4 } into task")

	req("create co(tnum, signed) key(tnum)")
	act("insert { tnum: 100, signed: 990101 } into co")
	act("insert { tnum: 102, signed: 990102 } into co")
	act("insert { tnum: 104, signed: 990103 } into co")
	act("insert { tnum: 106, signed: 990104 } into co")

	req("create dates(date) key(date)")
	act("insert { date: #20010101 } into dates")
	act("insert { date: #20010102 } into dates")
	act("insert { date: #20010301 } into dates")
	act("insert { date: #20010401 } into dates")

	return db
}

func ck(err error) {
	if err != nil {
		panic(err.Error())
	}
}
