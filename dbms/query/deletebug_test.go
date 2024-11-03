// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func init() {
	db19.MakeSuTran = func(ut *db19.UpdateTran) *SuTran {
		return NewSuTran(nil, true)
	}
}

// func TestDeleteBug2(*testing.T) {
// 	// store := stor.HeapStor(8192)
// 	// db, err := db19.CreateDb(store)
// 	db, err := db19.OpenDatabase("../../suneido.db")
// 	ck(err)
// 	defer db.Close()
// 	// db.CheckerSync()
// 	db19.StartConcur(db, 5*time.Millisecond)
// 	act := func(act string) int {
// 		// time.Sleep(1 * time.Microsecond)
// 		ut := db.NewUpdateTran()
// 		defer ut.Commit()
// 		return DoAction(nil, ut, act, nil)
// 	}
// 	// DoAdmin(db, "ensure tmp(a,b,c) key(a,b) key(c)")
// 	// for range 10000 {
// 	// 	act("delete tmp")
// 	// 	act("insert { a: 1, b: 1, c: 1 } into tmp")
// 	// 	act("insert { a: 2, b: 2, c: 2 } into tmp")
// 	// 	n := act("delete tmp")
// 	// 	assert.This(n).Is(2)
// 	// }
// 	for range 10000 {
// 		// fmt.Println(i)
// 		act("delete Test_lib")
// 		act("insert { name: 'One', group: -1, num: 99999 } into Test_lib")
// 		// db.NewReadTran().GetInfo("tmp").Indexes[0].Print()
// 		// time.Sleep(time.Microsecond)
// 		act("insert { name: 'Two', group: -1, num: 99998 } into Test_lib")
// 		n := act("delete Test_lib")
// 		assert.This(n).Is(2)
// 	}
// }

func TestDeleteBug(*testing.T) {
	store := stor.HeapStor(8192)
	db, err := db19.CreateDb(store)
	ck(err)
	db19.StartConcur(db, 50*time.Second)
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		n := DoAction(nil, ut, act)
		assert.This(n).Is(1)
	}
	DoAdmin(db, "create tmp(k) key(k)", nil)
	N := 10000
	if testing.Short() {
		N = 1000
	}
	for range N {
		act("insert { k: 1 } into tmp")
		act("delete tmp")
	}
	db.MustCheck()
}

func TestDeleteSynch(*testing.T) {
	store := stor.HeapStor(8192)
	db, err := db19.CreateDb(store)
	ck(err)
	db.CheckerSync()
	act := func(act string) {
		ut := db.NewUpdateTran()
		DoAction(nil, ut, act)
		db.CommitMerge(ut)
	}
	DoAdmin(db, "create tmp(k) key(k)", nil)
	N := 10000
	if testing.Short() {
		N = 1000
	}
	for range N {
		act("insert { k: 1 } into tmp")
		act("delete tmp")
	}
	db.MustCheck()
}
