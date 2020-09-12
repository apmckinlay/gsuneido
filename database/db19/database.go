// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"errors"

	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/comp"
	"github.com/apmckinlay/gsuneido/database/db19/ixspec"
	"github.com/apmckinlay/gsuneido/database/db19/meta"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/cksum"
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
	db.state.set(&DbState{store: store, meta: meta.CreateMeta(store)})

	n := len(magic) + stor.SmallOffsetLen
	_, buf := store.Alloc(n)
	copy(buf, magic)
	stor.WriteSmallOffset(buf[len(magic):], uint64(n))
	db.store = store
	db.mode = stor.CREATE
	return &db
}

func OpenDatabase(filename string) *Database {
	return openDatabase(filename, stor.UPDATE, quickCheck)
}

func OpenDatabaseRead(filename string) *Database {
	return openDatabase(filename, stor.READ, quickCheck)
}

type checkType int
const (
	noCheck checkType = iota
	quickCheck
	fullCheck
)

func openDatabase(filename string, mode stor.Mode, ck checkType) *Database {
	var db Database

	store, err := stor.MmapStor(filename, mode)
	if err != nil {
		panic("can't open database " + filename)
	}

	buf := store.Data(0)
	if magic != string(buf[:len(magic)]) {
		panic("not a valid database " + filename)
	}
	size := stor.ReadSmallOffset(buf[len(magic):])
	if size != store.Size() {
		//TODO recovery
		panic("database size mismatch - not shut down properly?")
	}

	db.store = store
	db.mode = mode
	db.state.set(ReadState(db.store, size-uint64(stateLen)))
	ckerr := ""
	switch ck {
	case fullCheck:
		ckerr = db.Check()
	case quickCheck:
		ckerr = db.QuickCheck()
	}
	if ckerr != "" {
		//TODO recovery
		panic("database check failed - corrupt?\n" + ckerr)
	}
	return &db
}

// LoadedTable is used to add a loaded table to the state
func (db *Database) LoadedTable(ts *meta.Schema, ti *meta.Info) error {
	var err error
	db.UpdateState(func(state *DbState) {
		if nil != state.meta.GetRoSchema(ts.Table) {
			err = errors.New("can't create " + ts.Table + " - it already exists")
			return
		}
		state.meta = state.meta.Add(ts, ti)
	})
	return err
}

// Close closes the database store, writing the current size to the start.
// NOTE: The state must already be written.
func (db *Database) Close() {
	if db.mode != stor.READ {
		// need to use Write because all but last chunk are read-only
		buf := make([]byte, stor.SmallOffsetLen)
		stor.WriteSmallOffset(buf, db.store.Size())
		db.store.Write(uint64(len(magic)), buf)
	}
	db.store.Close()
}

//-------------------------------------------------------------------

func init() {
	btree.GetLeafKey = getLeafKey
}

func getLeafKey(store *stor.Stor, is *ixspec.T, off uint64) string {
	rec := offToRec(store, off)
	return comp.Key(rt.Record(rec), is.Fields, is.Fields2)
}

func mkcmp(store *stor.Stor, is *ixspec.T) func(x, y uint64) int {
	return func(x, y uint64) int {
		xr := offToRec(store, x)
		yr := offToRec(store, y)
		return comp.Compare(xr, yr, is.Fields, is.Fields2)
	}
}

func offToRec(store *stor.Stor, off uint64) rt.Record {
	buf := store.Data(off)
	size := rt.RecLen(buf)
	return rt.Record(hacks.BStoS(buf[:size]))
}

// offToRecCk verifies the checksum following the record
func offToRecCk(store *stor.Stor, off uint64) rt.Record {
	buf := store.Data(off)
	size := rt.RecLen(buf)
	cksum.Check(buf[:size+cksum.Len])
	return rt.Record(hacks.BStoS(buf[:size]))
}
