// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"errors"

	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/sset"
)

type Database struct {
	mode  stor.Mode
	Store *stor.Stor

	// state is the central immutable state of the database.
	// It must be accessed atomically and only updated via UpdateState.
	state stateHolder

	ck Checker
	triggers
}

const magic = "gsndo001"

func CreateDatabase(filename string) (*Database, error) {
	store, err := stor.MmapStor(filename, stor.CREATE)
	if err != nil {
		return nil, err
	}
	return CreateDb(store)
}

func CreateDb(store *stor.Stor) (*Database, error) {
	var db Database
	db.state.set(&DbState{store: store, Meta: &meta.Meta{}})

	n := len(magic) + stor.SmallOffsetLen
	_, buf := store.Alloc(n)
	copy(buf, magic)
	stor.WriteSmallOffset(buf[len(magic):], uint64(n))
	db.Store = store
	db.mode = stor.CREATE
	return &db, nil
}

func OpenDatabase(filename string) (*Database, error) {
	return OpenDb(filename, stor.UPDATE, true)
}

func OpenDatabaseRead(filename string) (*Database, error) {
	return OpenDb(filename, stor.READ, true)
}

func OpenDb(filename string, mode stor.Mode, check bool) (db *Database, err error) {
	store, err := stor.MmapStor(filename, mode)
	if err != nil {
		return nil, err
	}
	return OpenDbStor(store, mode, check)
}

func OpenDbStor(store *stor.Stor, mode stor.Mode, check bool) (db *Database, err error) {
	defer func() {
		if err != nil {
			store.Close()
		}
	}()
	buf := store.Data(0)
	if magic != string(buf[:len(magic)]) {
		return nil, errors.New("bad magic")
	}
	size := stor.ReadSmallOffset(buf[len(magic):])
	if size != store.Size() {
		return nil, errors.New("bad size, not shut down properly?")
	}

	defer func() {
		if e := recover(); e != nil {
			err = newErrCorrupt(e)
			db = nil
		}
	}()
	db = &Database{Store: store, mode: mode}
	state, _ := ReadState(db.Store, size-uint64(stateLen))
	db.state.set(state)
	if check {
		if err := db.QuickCheck(); err != nil {
			return nil, err
		}
	}
	return db, nil
}

// CheckerSync is for tests.
// It assigns a synchronous transaction checker to the database.
func (db *Database) CheckerSync() {
	db.ck = NewCheck(db)
}

// LoadedTable is used to add a loaded table to the state
func (db *Database) LoadedTable(ts *meta.Schema, ti *meta.Info) {
	if err := db.loadedTable(ts, ti); err != nil {
		panic(err)
	}
}

func (db *Database) loadedTable(ts *meta.Schema, ti *meta.Info) error {
	var err error
	db.UpdateState(func(state *DbState) {
		if state.Meta.GetRoSchema(ts.Table) != nil {
			err = errors.New("can't create existing table: " + ts.Table)
		}
		state.Meta = state.Meta.Put(ts, ti)
	})
	return err
}

// schema changes ---------------------------------------------------

// Creating new indexes on an existing table (ensure and alter create)
// must be serialized with check/merge
// to ensure that merge sees a state consistent with the transaction.

func (db *Database) Create(schema *schema.Schema) {
	ts, ti := db.create(schema)
	db.UpdateState(func(state *DbState) {
		if state.Meta.GetRoSchema(ts.Table) != nil {
			panic("can't create existing table: " + ts.Table)
		}
		state.Meta = state.Meta.PutNew(ts, ti, schema)
	})
}

func (db *Database) create(schema *schema.Schema) (*meta.Schema, *meta.Info) {
	ts := &meta.Schema{Schema: *schema}
	ts.Ixspecs(ts.Indexes)
	ov := db.createIndexes(ts.Indexes)
	ti := &meta.Info{Table: schema.Table, Indexes: ov}
	return ts, ti
}

func (db *Database) createIndexes(idxs []schema.Index) []*index.Overlay {
	ov := make([]*index.Overlay, len(idxs))
	for i := range ov {
		bt := btree.CreateBtree(db.Store, &idxs[i].Ixspec)
		ov[i] = index.OverlayFor(bt)
	}
	return ov
}

func (db *Database) Ensure(schema *schema.Schema) {
	handled := false
	db.UpdateState(func(state *DbState) {
		ts := state.Meta.GetRoSchema(schema.Table)
		if ts == nil {
			// TODO check if schema is "full" (see parseadmin.go)
			// else you could get assertion failures
			ts, ti := db.create(schema)
			state.Meta = state.Meta.Put(ts, ti)
			handled = true

		} else if schemaSubset(schema, ts) {
			handled = true // nothing to do, common fast case
		}
	})
	if !handled {
		// TODO only use ck.Run if adding indexes
		err := db.ck.Run(func() error { return db.ensure(schema) })
		if err != nil {
			panic(err)
		}
	}
}

// schemaSubset returns whether the table (ts) already has the ensure schema
func schemaSubset(schema *schema.Schema, ts *meta.Schema) bool {
	if !sset.Subset(ts.Columns, schema.Columns) {
		return false
	}
	for i := range schema.Indexes {
		if nil == ts.FindIndex(schema.Indexes[i].Columns) {
			return false
		}
	}
	return true
}

func (db *Database) ensure(schema *schema.Schema) (err error) {
	db.UpdateState(func(state *DbState) {
		var meta *meta.Meta
		meta, err = state.Meta.Ensure(schema, db.Store)
		if err == nil {
			state.Meta = meta
		}
	})
	return err
}

func (db *Database) RenameTable(from, to string) bool {
	result := false
	db.UpdateState(func(state *DbState) {
		if m := state.Meta.RenameTable(from, to); m != nil {
			state.Meta = m
			result = true
		}
	})
	return result
}

// Drop removes a table or view
func (db *Database) Drop(table string) error {
	var err error
	db.UpdateState(func(state *DbState) {
		if m := state.Meta.Drop(table); m != nil {
			state.Meta = m
		} else {
			err = errors.New("can't drop nonexistent table: " + table)
		}
	})
	return err
}

// AlterRename renames columns
func (db *Database) AlterRename(table string, from, to []string) bool {
	result := false
	db.UpdateState(func(state *DbState) {
		if m := state.Meta.AlterRename(table, from, to); m != nil {
			state.Meta = m
			result = true
		}
	})
	return result
}

// AlterCreate creates columns or indexes
func (db *Database) AlterCreate(schema *schema.Schema) {
	err := db.ck.Run(func() error { return db.alterCreate(schema) })
	if err != nil {
		panic(err)
	}
}

func (db *Database) alterCreate(schema *schema.Schema) (err error) {
	db.UpdateState(func(state *DbState) {
		var meta *meta.Meta
		meta, err = state.Meta.AlterCreate(schema, db.Store)
		if err == nil {
			state.Meta = meta
		}
	})
	return err
}

// AlterCreate removes columns or indexes
func (db *Database) AlterDrop(schema *schema.Schema) bool {
	result := false
	db.UpdateState(func(state *DbState) {
		if m := state.Meta.AlterDrop(schema); m != nil {
			state.Meta = m
			result = true
		}
	})
	return result
}

func (db *Database) AddView(name, def string) bool {
	result := false
	db.UpdateState(func(state *DbState) {
		if m := state.Meta.AddView(name, def); m != nil {
			state.Meta = m
			result = true
		}
	})
	return result
}

func (db *Database) GetView(name string) string {
	return db.GetState().Meta.GetView(name)
}

//-------------------------------------------------------------------

func (db *Database) Schema(table string) string {
	state := db.GetState()
	ts := state.Meta.GetRoSchema(table)
	if ts == nil {
		return ""
	}
	return ts.Schema.String()
}

func (db *Database) Size() uint64 {
	return db.Store.Size()
}

// Close closes the database store, writing the current size to the start.
// NOTE: The state must already be written.
func (db *Database) Close() {
	if db.Store == nil {
		return // already closed
	}
	if db.ck != nil {
		db.ck.Stop()
	} else if db.mode != stor.READ {
		db.Persist(&execPersistSingle{}, true)
	}
	if db.mode != stor.READ {
		db.writeSize()
	}
	db.Store.Close()
	db.Store = nil
}

func (db *Database) writeSize() {
	// need to use Write because all but last chunk are read-only
	buf := make([]byte, stor.SmallOffsetLen)
	stor.WriteSmallOffset(buf, db.Store.Size())
	db.Store.Write(uint64(len(magic)), buf)
}

//-------------------------------------------------------------------

func init() {
	btree.GetLeafKey = getLeafKey
}

func getLeafKey(store *stor.Stor, is *ixkey.Spec, off uint64) string {
	return is.Key(OffToRec(store, off))
}

func OffToRec(store *stor.Stor, off uint64) rt.Record {
	buf := store.Data(off)
	size := rt.RecLen(buf)
	return rt.Record(hacks.BStoS(buf[:size]))
}

// OffToRecCk verifies the checksum following the record
func OffToRecCk(store *stor.Stor, off uint64) rt.Record {
	buf := store.Data(off)
	size := rt.RecLen(buf)
	cksum.MustCheck(buf[:size+cksum.Len])
	return rt.Record(hacks.BStoS(buf[:size]))
}

func (db *Database) MakeLess(is *ixkey.Spec) func(x, y uint64) bool {
	return func(x, y uint64) bool {
		xr := OffToRec(db.Store, x)
		yr := OffToRec(db.Store, y)
		return is.Compare(xr, yr) < 0
	}
}
