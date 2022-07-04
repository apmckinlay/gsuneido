// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"encoding/binary"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

type DbState struct {
	store *stor.Stor
	Meta  *meta.Meta
	Asof  int64  // unix milli time
	Off   uint64 // offset of this state
}

type stateHolder struct {
	state unsafe.Pointer // *DbState
	mutex sync.Mutex
}

func (sh *stateHolder) get() *DbState {
	return (*DbState)(atomic.LoadPointer(&sh.state))
}

func (sh *stateHolder) set(newState *DbState) {
	atomic.StorePointer(&sh.state, unsafe.Pointer(newState))
}

// GetState returns a snapshot of the state as of a point in time.
// This state must be treated as read-only and must not be modified.
// To modify the state use UpdateState.
//
// GetState is atomic, it is not blocked by UpdateState.
func (db *Database) GetState() *DbState {
	return db.state.get()
}

// Persist forces a persist and returns a persisted state,
// with all ixbuf layer entries merged into the btree.
// This is used by dump, load, and checkdb.
func (db *Database) Persist() *DbState {
	if db.ck == nil { // for tests
		return db.GetState()
	}
	return db.ck.Persist()
}

// UpdateState applies the given update function to a copy of theState
// and sets theState to the result.
// Guarded by stateMutex so only one thread can execute at a time.
// Note: the state passed to the update function is a *shallow* copy,
// it is up to the function to make copies of any nested containers.
//
// UpdateState is guarded by a mutex
func (db *Database) UpdateState(fn func(*DbState)) {
	db.state.updateState(fn)
}

func (sh *stateHolder) updateState(fn func(*DbState)) {
	sh.mutex.Lock()
	defer sh.mutex.Unlock()
	oldState := sh.get()
	if oldState == nil {
		runtime.Fatal("database closed")
	}
	newState := *oldState // shallow copy
	fn(&newState)
	if newState.Meta != oldState.Meta {
		sh.set(&newState)
	}
}

//-------------------------------------------------------------------

// WARNING: Merge and Persist must not run concurrently

type mergefn func(*meta.Meta, *mergeList) []meta.MergeUpdate

// Merge updates the base ixbuf's with the ones from transactions
// It is called by concur.go merger.
func (db *Database) Merge(fn mergefn, merges *mergeList) {
	updates := fn(db.GetState().Meta, merges) // outside UpdateState
	db.UpdateState(func(state *DbState) {
		meta := *state.Meta // copy
		meta.ApplyMerge(updates)
		state.Meta = &meta
	})
}

// CommitMerge is for tests.
// It merges synchronously after each commit.
func (db *Database) CommitMerge(ut *UpdateTran) {
	tables := db.ck.(*Check).commit(ut)
	ut.commit()
	merges := &mergeList{}
	merges.add(tables)
	db.Merge(mergeSingle, merges)
}

//-------------------------------------------------------------------

type persistfn func(*DbState) []meta.PersistUpdate

// persist writes index changes (and a new state) to the database file.
// It is called from concur.go regularly e.g. once per minute.
// flatten applies to the schema and info chains (not indexes).
// NOTE: persist should only be called by the checker.
func (db *Database) persist(exec execPersist) *DbState {
	// fmt.Println("persist")
	var newState *DbState
	db.GetState().Meta.Persist(exec.Submit) // outside UpdateState
	updates := exec.Results()
	db.UpdateState(func(state *DbState) {
		meta := *state.Meta // copy
		meta.ApplyPersist(updates)
		state.Meta = &meta
		// Write modifies schema/info offs,ages,clock
		// so it must be inside UpdateState
		state.Off = state.Write()
		newState = state
	})
	return newState
}

const magic1 = "\x01\x23\x45\x67\x89\xab\xcd\xef"
const magic2 = "\xfe\xdc\xba\x98\x76\x54\x32\x10"
const dateSize = 8
const stateLen = len(magic1) + dateSize + 2*stor.SmallOffsetLen +
	len(magic2) + cksum.Len
const magic2at = stateLen - len(magic2)

func (state *DbState) Write() uint64 {
	// NOTE: indexes should already have been saved
	offSchema, offInfo := state.Meta.Write(state.store)
	return writeState(state.store, offSchema, offInfo)
}

func writeState(store *stor.Stor, offSchema, offInfo uint64) uint64 {
	stateOff, buf := store.Alloc(stateLen)
	copy(buf, magic1)
	i := len(magic1)
	t := time.Now().UnixMilli()
	binary.BigEndian.PutUint64(buf[i:], uint64(t))
	i += dateSize
	stor.WriteSmallOffset(buf[i:], offSchema)
	i += stor.SmallOffsetLen
	stor.WriteSmallOffset(buf[i:], offInfo)
	i += stor.SmallOffsetLen
	i += cksum.Len
	cksum.Update(buf[:i])
	copy(buf[i:], magic2)
	i += len(magic2)
	assert.That(i == stateLen)
	return stateOff
}

func ReadState(st *stor.Stor, off uint64) *DbState {
	offSchema, offInfo, t := readState(st, off)
	if t == 0 {
		panic("bad state")
	}
	return &DbState{store: st, Meta: meta.ReadMeta(st, offSchema, offInfo),
		Asof: t, Off: off}
}

func readState(st *stor.Stor, off uint64) (offSchema, offInfo uint64, t int64) {
	buf := st.Data(off)[:stateLen]
	i := len(magic1)
	if string(buf[:i]) != magic1 {
		return 0, 0, 0
	}
	cksum.MustCheck(buf[:magic2at])
	if string(buf[magic2at:magic2at+len(magic2)]) != magic2 {
		return 0, 0, 0
	}
	t = int64(binary.BigEndian.Uint64(buf[i:]))
	i += dateSize
	offSchema = stor.ReadSmallOffset(buf[i:])
	i += stor.SmallOffsetLen
	offInfo = stor.ReadSmallOffset(buf[i:])
	i += stor.SmallOffsetLen
	if offSchema >= off || offInfo >= off {
		return 0, 0, 0
	}
	return offSchema, offInfo, t
}

//-------------------------------------------------------------------

// StateAsof returns the state <= asof, or the initial state.
func StateAsof(store *stor.Stor, asof int64) *DbState {
	var offSchema, offInfo uint64
	var t int64
	off := store.Size()
	for {
		if off = store.LastOffset(off, magic1); off == 0 {
			break
		}
		if offSchema, offInfo, t = readState(store, off); t == 0 {
			continue // invalid
		}
		if t <= asof {
			break
		}
	}
	if t == 0 {
		panic("no state found")
	}
	return &DbState{store: store, Asof: t, Off: off,
		Meta: meta.ReadMeta(store, offSchema, offInfo)}
}

func NextState(store *stor.Stor, off uint64) *DbState {
	for {
		// +1 so we don't find the same state
		off = store.FirstOffset(off+1, magic1)
		if off == 0 {
			return nil
		}
		if offSchema, offInfo, t := readState(store, off); t != 0 {
			return &DbState{store: store, Asof: t, Off: off,
				Meta: meta.ReadMeta(store, offSchema, offInfo)}
		}
	}
}

func PrevState(store *stor.Stor, off uint64) *DbState {
	if off == 0 {
		off = store.Size()
	}
	for {
		off = store.LastOffset(off, magic1)
		if off == 0 {
			return nil
		}
		if offSchema, offInfo, t := readState(store, off); t != 0 {
			return &DbState{store: store, Asof: t, Off: off,
				Meta: meta.ReadMeta(store, offSchema, offInfo)}
		}
	}
}
