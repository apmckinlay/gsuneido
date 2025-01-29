// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
)

type heapdb struct {
	stor *stor.Stor
	*db19.Database
}

func heapDb() heapdb {
	stor := stor.HeapStor(8192)
	db, err := db19.CreateDb(stor)
	ck(err)
	db19.StartConcur(db, 50*time.Millisecond)
	db19.MakeSuTran = func(ut *db19.UpdateTran) *SuTran {
		return NewSuTran(nil, true)
	}
	MakeSuTran = func(qt QueryTran) *SuTran {
		return nil
	}
	return heapdb{stor: stor, Database: db}
}

func (hdb heapdb) adm(admin string) {
	DoAdmin(hdb.Database, admin, nil)
}

func (hdb heapdb) act(act string) {
	ut := hdb.NewUpdateTran()
	defer ut.Commit()
	DoAction(nil, ut, act)
}

func (hdb heapdb) reopen() heapdb {
	hdb.Close()
	db, err := db19.OpenDbStor(hdb.stor, stor.Read, true)
	ck(err)
	db19.StartConcur(db, 50*time.Millisecond)
	return heapdb{stor: hdb.stor, Database: db}
}

func testDb() *db19.Database {
	db := heapDb()
	db.adm("create customer (id, name, city) key(id)")
	db.act("insert {id: 'a', name: 'axon', city: 'saskatoon'} into customer")
	db.act("insert {id: 'c', name: 'calac', city: 'calgary'} into customer")
	db.act("insert {id: 'e', name: 'emerald', city: 'vancouver'} into customer")
	db.act("insert {id: 'i', name: 'intercon', city: 'saskatoon'} into customer")

	db.adm("create hist (date, item, id, cost) index(date) key(date,item,id)")
	db.act("insert{date: 970101, item: 'disk', id: 'a', cost: 100} into hist")
	db.act("insert{date: 970101, item: 'disk', id: 'e', cost: 200} into hist")
	db.act("insert{date: 970102, item: 'mouse', id: 'c', cost: 200} into hist")
	db.act("insert{date: 970103, item: 'pencil', id: 'e', cost: 300} into hist")

	db.adm("create hist2 (date, item, id, cost) key(date) index(id)")
	db.act("insert{date: 970101, item: 'disk', id: 'a', cost: 100} into hist2")
	db.act("insert{date: 970102, item: 'disk', id: 'e', cost: 200} into hist2")
	db.act("insert{date: 970103, item: 'pencil', id: 'e', cost: 300} into hist2")

	db.adm("create trans (item, id, cost, date) index(item) key(date,item,id)")
	db.act("insert{item: 'mouse', id: 'e', cost: 200, date: 960204} into trans")
	db.act("insert{item: 'disk', id: 'a', cost: 100, date: 970101} into trans")
	db.act("insert{item: 'mouse', id: 'c', cost: 200, date: 970101} into trans")
	db.act("insert{item: 'eraser', id: 'c', cost: 150, date: 970201} into trans")

	db.adm("create supplier (supplier, name, city) key(supplier) index(city)")
	db.act("insert{supplier: 'mec', name: 'mtnequipcoop', city: 'calgary'} into supplier")
	db.act("insert{supplier: 'hobo', name: 'hoboshop', city: 'saskatoon'} into supplier")
	db.act("insert{supplier: 'ebs', name: 'ebssail&sport', city: 'saskatoon'} into supplier")
	db.act("insert{supplier: 'taiga', name: 'taigaworks', city: 'vancouver'} into supplier")

	db.adm("create inven (item, qty) key(item)")
	db.act("insert{item: 'disk', qty: 5} into inven")
	db.act("insert{item: 'mouse', qty:2} into inven")
	db.act("insert{item: 'pencil', qty: 7} into inven")

	db.adm("create alias(id, name2) key(id)")
	db.act("insert{id: 'a', name2: 'abc'} into alias")
	db.act("insert{id: 'c', name2: 'trical'} into alias")

	db.adm("create cus(cnum, abbrev, name) key(cnum) key(abbrev)")
	db.act("insert { cnum: 1, abbrev: 'a', name: 'axon' } into cus")
	db.act("insert { cnum: 2, abbrev: 'b', name: 'bill' } into cus")
	db.act("insert { cnum: 3, abbrev: 'c', name: 'cron' } into cus")
	db.act("insert { cnum: 4, abbrev: 'd', name: 'dick' } into cus")

	db.adm("create task(tnum, cnum) key(tnum)")
	db.act("insert { tnum: 100, cnum: 1 } into task")
	db.act("insert { tnum: 101, cnum: 2 } into task")
	db.act("insert { tnum: 102, cnum: 3 } into task")
	db.act("insert { tnum: 103, cnum: 4 } into task")
	db.act("insert { tnum: 104, cnum: 1 } into task")
	db.act("insert { tnum: 105, cnum: 2 } into task")
	db.act("insert { tnum: 106, cnum: 3 } into task")
	db.act("insert { tnum: 107, cnum: 4 } into task")

	db.adm("create co(tnum, signed) key(tnum)")
	db.act("insert { tnum: 100, signed: 990101 } into co")
	db.act("insert { tnum: 102, signed: 990102 } into co")
	db.act("insert { tnum: 104, signed: 990103 } into co")
	db.act("insert { tnum: 106, signed: 990104 } into co")

	db.adm("create dates(date) key(date)")
	db.act("insert { date: #20010101 } into dates")
	db.act("insert { date: #20010102 } into dates")
	db.act("insert { date: #20010301 } into dates")
	db.act("insert { date: #20010401 } into dates")

	// close and reopen to force persist
	db = db.reopen()
	return db.Database
}

func ck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

//-------------------------------------------------------------------

// sizeTran wraps an actual transaction and overrides Nrows and Size
// since testDb does not have enough data to test query optimization
type sizeTran struct {
	QueryTran
}

func (t sizeTran) GetInfo(table string) *meta.Info {
	info := t.QueryTran.GetInfo(table)
	if info == nil {
		return nil
	}
	ti := *info // copy
	ti.Nrows = 1000
	if table == "trans" || table == "hist" || table == "hist2" {
		ti.Nrows = 10_000
	}
	ti.Size = int64(ti.Nrows) * 100
	return &ti
}
