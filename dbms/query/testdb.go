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
	return db
}

func ck(err error) {
	if err != nil {
		panic(err.Error())
	}
}
