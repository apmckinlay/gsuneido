// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"github.com/apmckinlay/gsuneido/db19/btree"
	"github.com/apmckinlay/gsuneido/db19/comp"
	"github.com/apmckinlay/gsuneido/db19/ixspec"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
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

func CreateDatabase(filename string) (*Database, error) {
	store, err := stor.MmapStor(filename, stor.CREATE)
	if err != nil {
		return nil, err
	}
	return createDatabase(store)
}

func createDatabase(store *stor.Stor) (*Database, error) {
	var db Database
	db.state.set(&DbState{store: store, meta: meta.CreateMeta(store)})

	n := len(magic) + stor.SmallOffsetLen
	_, buf := store.Alloc(n)
	copy(buf, magic)
	stor.WriteSmallOffset(buf[len(magic):], uint64(n))
	db.store = store
	db.mode = stor.CREATE
	return &db, nil
}

func OpenDatabase(filename string) (*Database, error) {
	return openDatabase(filename, stor.UPDATE, true)
}

func OpenDatabaseRead(filename string) (*Database, error) {
	return openDatabase(filename, stor.READ, true)
}

func openDatabase(filename string, mode stor.Mode, check bool) (db *Database, err error) {
	store, err := stor.MmapStor(filename, mode)
	if err != nil {
		return nil, err
	}
	buf := store.Data(0)
	if magic != string(buf[:len(magic)]) {
		return nil, &ErrCorrupt{}
	}
	size := stor.ReadSmallOffset(buf[len(magic):])
	if size != store.Size() {
		return nil, &ErrCorrupt{}
	}

	defer func() {
		if e := recover(); e != nil {
			err = newErrCorrupt(e)
			db = nil
		}
	}()
	db = &Database{store: store, mode: mode}
	state, _ := ReadState(db.store, size-uint64(stateLen))
	db.state.set(state)
	if check {
		if err := db.QuickCheck(); err != nil {
			return nil, err
		}
	}
	return db, nil
}

// LoadedTable is used to add a loaded table to the state
func (db *Database) LoadedTable(ts *meta.Schema, ti *meta.Info) {
	db.UpdateState(func(state *DbState) {
		state.meta = state.meta.Put(ts, ti)
	})
}

func (db *Database) DropTable(table string) bool {
	result := false
	db.UpdateState(func(state *DbState) {
		if m := state.meta.DropTable(table); m != nil {
			state.meta = m
			result = true
		}
	})
	return result
}

// Close closes the database store, writing the current size to the start.
// NOTE: The state must already be written.
func (db *Database) Close() {
	if db.store == nil {
		return // already closed
	}
	if db.ck != nil {
		db.ck.Stop()
	}
	if db.mode != stor.READ {
		// need to use Write because all but last chunk are read-only
		buf := make([]byte, stor.SmallOffsetLen)
		stor.WriteSmallOffset(buf, db.store.Size())
		db.store.Write(uint64(len(magic)), buf)
	}
	db.store.Close()
	db.store = nil
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
	cksum.MustCheck(buf[:size+cksum.Len])
	return rt.Record(hacks.BStoS(buf[:size]))
}
