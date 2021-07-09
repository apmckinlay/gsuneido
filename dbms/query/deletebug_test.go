// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestDeleteBug(*testing.T) {
	store := stor.HeapStor(8192)
	db, err := db19.CreateDb(store)
	ck(err)
	db19.StartConcur(db, 50*time.Second)
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(ut, act)
		assert.This(n).Is(1)
	}
	DoAdmin(db, "create tmp(k) key(k)")
	for i := 0; i < 10000; i++ {
		act("insert { k: 1 } into tmp")
		act("delete tmp")
	}
}

func TestDeleteSynch(*testing.T) {
	store := stor.HeapStor(8192)
	db, err := db19.CreateDb(store)
	ck(err)
	db.CheckerSync()
	act := func(act string) {
		ut := db.NewUpdateTran()
		DoAction(ut, act)
		db.CommitMerge(ut)
	}
	DoAdmin(db, "create tmp(k) key(k)")
	for i := 0; i < 100000; i++ {
		act("insert { k: 1 } into tmp")
		act("delete tmp")
	}
}
