// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/comp"
	"github.com/apmckinlay/gsuneido/database/db19/ixspec"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

type Database struct {
	mode  stor.Mode
	store *stor.Stor

	// state is the central immutable state of the database.
	// It must be accessed atomically and only updated via UpdateState.
	state stateHolder

	ck Checker
}

const magic = "gsndo001"

func CreateDatabase(filename string) *Database {
	store, err := stor.MmapStor(filename, stor.CREATE)
	if err != nil {
		panic("can't create database " + filename)
	}
	var db Database
	db.state.set(&DbState{store: store})

	n := len(magic) + stor.SmallOffsetLen
	_, buf := store.Alloc(n)
	copy(buf, magic)
	stor.WriteSmallOffset(buf[len(magic):], uint64(n))
	db.store = store
	db.mode = stor.CREATE
	return &db
}

func OpenDatabase(filename string) *Database {
	return openDatabase(filename, stor.UPDATE)
}

func OpenDatabaseRead(filename string) *Database {
	return openDatabase(filename, stor.READ)
}

func openDatabase(filename string, mode stor.Mode) *Database {
	var db Database

	store, err := stor.MmapStor(filename, mode)
	if err != nil {
		panic("can't open database " + filename)
	}

	//TODO recovery
	buf := store.Data(0)
	if magic != string(buf[:len(magic)]) {
		panic("not a valid database " + filename)
	}
	size := stor.ReadSmallOffset(buf[len(magic):])
	if size != store.Size() {
		panic("database size mismatch - not shut down properly?")
	}

	db.store = store
	db.mode = mode
	db.state.set(ReadState(db.store, size-uint64(stateLen)))
	//TODO integrity check

	return &db
}

func (db *Database) Close() {
	if db.mode != stor.READ {
		buf := db.store.Data(0)
		stor.WriteSmallOffset(buf[len(magic):], db.store.Size())
	}
	db.store.Close()
}

//-------------------------------------------------------------------

func init() {
	btree.GetLeafKey = getLeafKey
}

func getLeafKey(store *stor.Stor, is *ixspec.T, off uint64) string {
	rec := offToRec(store, off)
	return comp.Key(rt.Record(rec), is.Cols, is.Cols2)
}

func mkcmp(store *stor.Stor, is *ixspec.T) func(x, y uint64) int {
	return func(x, y uint64) int {
		xr := offToRec(store, x)
		yr := offToRec(store, y)
		return comp.Compare(xr, yr, is.Cols, is.Cols2)
	}
}

func offToRec(store *stor.Stor, off uint64) rt.Record {
	buf := store.Data(off)
	return rt.Record(hacks.BStoS(buf))
}
