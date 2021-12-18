// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"errors"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/sortlist"
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
	// schemaLock is used to prevent concurrent schema modification
	schemaLock int64
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
	db.lockSchema()
	defer db.unlockSchema()
	db.UpdateState(func(state *DbState) {
		if state.Meta.GetRoSchema(schema.Table) != nil {
			panic("can't create existing table: " + schema.Table)
		}
		db.create(state, schema)
	})
}

func (db *Database) lockSchema() {
	if !atomic.CompareAndSwapInt64(&db.schemaLock, 0, 1) {
		panic("concurrent schema modifications are not allowed")
	}
}

func (db *Database) unlockSchema() {
	atomic.StoreInt64(&db.schemaLock, 0)
}

func (db *Database) create(state *DbState, schema *schema.Schema) {
	schema.Check()
	ts := &meta.Schema{Schema: *schema}
	ts.Ixspecs(ts.Indexes)
	ov := db.createIndexes(ts.Indexes)
	ti := &meta.Info{Table: schema.Table, Indexes: ov}
	state.Meta = state.Meta.PutNew(ts, ti, schema)
}

func (db *Database) createIndexes(idxs []schema.Index) []*index.Overlay {
	ov := make([]*index.Overlay, len(idxs))
	for i := range ov {
		bt := btree.CreateBtree(db.Store, &idxs[i].Ixspec)
		ov[i] = index.OverlayFor(bt)
	}
	return ov
}

func (db *Database) Ensure(sch *schema.Schema) {
	db.lockSchema()
	defer db.unlockSchema()
	handled := false
	var newIdxs []schema.Index
	db.UpdateState(func(state *DbState) {
		ts := state.Meta.GetRoSchema(sch.Table)
		if ts == nil { // table doesn't exist
			db.create(state, sch)
			handled = true

		} else if schemaSubset(sch, ts) {
			handled = true // nothing to do, common fast case
		} else {
			var meta *meta.Meta
			newIdxs, meta = state.Meta.Ensure(sch, db.Store)
			if len(newIdxs) == 0 {
				state.Meta = meta
				handled = true
			}
			// else discard meta and just use newIdxs
		}
	})
	// outside UpdateState
	if !handled {
		db.ensure(sch, newIdxs)
	}
}

// schemaSubset returns whether the table (ts) already has the ensure schema
func schemaSubset(schema *schema.Schema, ts *meta.Schema) bool {
	if !sset.Subset(ts.Columns, schema.Columns) {
		return false
	}
	for i := range schema.Indexes {
		ix := ts.FindIndex(schema.Indexes[i].Columns)
		if ix == nil {
			return false
		}
		if !ix.Equal(&schema.Indexes[i]) {
			panic("ensure: index exists but is different")
		}
	}
	return true
}

func (db *Database) ensure(sch *schema.Schema, newIdxs []schema.Index) {
	db.AddExclusive(sch.Table)
	defer db.EndExclusive(sch.Table)

	sch2 := *sch
	sch2.Indexes = nil
	db.UpdateState(func(state *DbState) { // add columns
		_, meta := state.Meta.Ensure(&sch2, db.Store) // final run
		state.Meta = meta
	})

	ov := db.buildIndexes(sch.Table, newIdxs)

	sch.Columns = nil
	db.UpdateState(func(state *DbState) {
		_, meta := state.Meta.Ensure(sch, db.Store) // final run
		// now meta and table info are copies
		if ov != nil {
			// add newly created indexes
			ti := meta.GetRoInfo(sch.Table) // not actually read-only
			i := len(ti.Indexes) - len(ov)
			copy(ti.Indexes[i:], ov)
		}
		state.Meta = meta
	})
}

func (db *Database) AddExclusive(table string) {
	if !db.ck.AddExclusive(table) {
		panic("index creation: can't get exclusive access to " + table)
	}
}

func (db *Database) EndExclusive(table string) {
	db.ck.EndExclusive(table)
}

// buildIndexes creates the new btrees & overlays
func (db *Database) buildIndexes(table string, newIdxs []schema.Index) []*index.Overlay {
	if len(newIdxs) == 0 {
		return nil
	}
	rt := db.NewReadTran()
	ti := rt.meta.GetRoInfo(table)
	if ti.Nrows == 0 {
		return nil
	}
	rt.meta.GetRoSchema(table).Ixspecs(newIdxs)
	nlayers := ti.Indexes[0].Nlayers()
	list := sortlist.NewUnsorted()
	iter := index.NewOverIter(table, 0)
	for iter.Next(rt); !iter.Eof(); iter.Next(rt) {
		_, off := iter.Cur()
		list.Add(off)
	}
	ov := make([]*index.Overlay, len(newIdxs))
	for i := range newIdxs {
		ix := &newIdxs[i]
		fk := &ix.Fk
		list.Sort(MakeLess(db.Store, &ix.Ixspec))
		bldr := btree.Builder(db.Store)
		iter := list.Iter()
		for off := iter(); off != 0; off = iter() {
			rec := OffToRec(db.Store, off)
			key := ix.Ixspec.Key(rec)
			bldr.Add(key, off)
			// check foreign key
			if fk.Table != "" {
				k := ix.Ixspec.Trunc(len(ix.Columns)).Key(rec)
				if k != "" && !rt.fkeyOutputExists(fk.Table, fk.IIndex, k) {
					panic("output blocked by foreign key: " +
						fk.Table + " " + ix.String())
				}
			}
		}
		bt := bldr.Finish()
		bt.SetIxspec(&ix.Ixspec)
		ov[i] = index.OverlayForN(bt, nlayers)
	}
	return ov
}

// MakeLess handles _lower! but not rules.
// It is used for indexes (which don't support rules).
func MakeLess(store *stor.Stor, is *ixkey.Spec) func(x, y uint64) bool {
	return func(x, y uint64) bool {
		xr := OffToRec(store, x)
		yr := OffToRec(store, y)
		return is.Compare(xr, yr) < 0
	}
}

func (db *Database) RenameTable(from, to string) bool {
	db.lockSchema()
	defer db.unlockSchema()
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
	db.lockSchema()
	defer db.unlockSchema()
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
	db.lockSchema()
	defer db.unlockSchema()
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
func (db *Database) AlterCreate(sch *schema.Schema) {
	db.lockSchema()
	defer db.unlockSchema()
	db.AddExclusive(sch.Table)
	defer db.EndExclusive(sch.Table)

	ov := db.buildIndexes(sch.Table, sch.Indexes)
	db.UpdateState(func(state *DbState) {
		meta := state.Meta.AlterCreate(sch, db.Store)
		// now meta and table info are copies
		if ov != nil {
			// add newly created indexes
			ti := meta.GetRoInfo(sch.Table) // not really read-only
			i := len(ti.Indexes) - len(ov)
			copy(ti.Indexes[i:], ov)
		}
		state.Meta = meta
	})
}

// AlterCreate removes columns or indexes
func (db *Database) AlterDrop(schema *schema.Schema) bool {
	db.lockSchema()
	defer db.unlockSchema()
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
	db.ckOpen()
	return db.Store.Size()
}

func (db *Database) Transactions() []int {
	db.ckOpen()
	return db.ck.Transactions()
}

func (db *Database) ckOpen() {
	if db.Store == nil {
		rt.Fatal("database closed")
	}
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
		db.persist(&execPersistSingle{}, true)
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
