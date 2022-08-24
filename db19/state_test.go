// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestStateReadWrite(*testing.T) {
	store := stor.HeapStor(1024)
	store.Alloc(500)
	off := writeState(store, 123, 456)
	offSchema, offInfo, _ := readState(store, off)
	assert.This(offSchema).Is(123)
	assert.This(offInfo).Is(456)
}

// func TestStateAsof(*testing.T) {
// 	db, err := OpenDatabase("../suneido.db")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	asof := time.Date(2022, 8, 22, 17, 50, 31, 0, time.Local)
// 	state := StateAsof(db.Store, asof.UnixMilli())
// 	fmt.Println(asof, "=>", time.UnixMilli(state.Asof))
// }
