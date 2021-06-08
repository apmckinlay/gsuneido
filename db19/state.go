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
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

type DbState struct {
	store *stor.Stor
	Meta  *meta.Meta
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
func (db *Database) GetState() *DbState {
	return db.state.get()
}

func (sh *stateHolder) updateState(fn func(*DbState)) *DbState {
	sh.mutex.Lock()
	defer sh.mutex.Unlock()
	newState := *sh.get() // shallow copy
	fn(&newState)
	sh.set(&newState)
	return &newState
}

// UpdateState applies the given update function to a copy of theState
// and sets theState to the result.
// Guarded by stateMutex so only one thread can execute at a time.
// Note: the state passed to the update function is a *shallow* copy,
// it is up to the function to make copies of any nested containers.
func (db *Database) UpdateState(fn func(*DbState)) *DbState {
	return db.state.updateState(fn)
}

//-------------------------------------------------------------------

// WARNING: Merge and Persist must not run concurrently

type mergefn func(*DbState, *mergeList) []meta.MergeUpdate

// Merge updates the base ixbuf's with the ones from transactions
// It is called by concur.go merger.
func (db *Database) Merge(fn mergefn, merges *mergeList) {
	updates := fn(db.GetState(), merges) // outside UpdateState
	db.UpdateState(func(state *DbState) {
		meta := *state.Meta // copy
		meta.ApplyMerge(updates)
		state.Meta = &meta
	})
}

//-------------------------------------------------------------------

type persistfn func(*DbState) []meta.PersistUpdate

// Persist writes index changes (and a new state) to the database file.
// It is called from concur.go e.g. once per minute.
// flatten applies to the schema and info chains.
func (db *Database) Persist(exec execPersist, flatten bool) uint64 {
	// fmt.Println("Persist", flatten)
	var off uint64
	db.GetState().Meta.Persist(exec.Submit) // outside UpdateState
	updates := exec.Results()
	db.UpdateState(func(state *DbState) {
		meta := *state.Meta // copy
		meta.ApplyPersist(updates)
		state.Meta = &meta
		off = state.Write(flatten)
	})
	return off
}

const magic1 = "\x01\x23\x45\x67\x89\xab\xcd\xef"
const magic2 = "\xfe\xdc\xba\x98\x76\x54\x32\x10"
const dateSize = 8
const stateLen = len(magic1) + dateSize + 2*stor.SmallOffsetLen +
	len(magic2) + cksum.Len
const magic2at = stateLen - len(magic2)

func (state *DbState) Write(flatten bool) uint64 {
	// NOTE: indexes should already have been saved
	offSchema, offInfo := state.Meta.Write(state.store, flatten)
	return writeState(state.store, offSchema, offInfo)
}

func writeState(store *stor.Stor, offSchema, offInfo uint64) uint64 {
	stateOff, buf := store.Alloc(stateLen)
	copy(buf, magic1)
	i := len(magic1)
	t := time.Now().Unix()
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

func ReadState(st *stor.Stor, off uint64) (*DbState, time.Time) {
	offSchema, offInfo, t := readState(st, off)
	return &DbState{store: st, Meta: meta.ReadMeta(st, offSchema, offInfo)}, t
}

func readState(st *stor.Stor, off uint64) (offSchema, offInfo uint64, t time.Time) {
	buf := st.Data(off)[:stateLen]
	i := len(magic1)
	if string(buf[:i]) != magic1 {
		panic("ReadState bad magic1")
	}
	cksum.MustCheck(buf[:magic2at])
	if string(buf[magic2at:magic2at+len(magic2)]) != magic2 {
		panic("ReadState bad magic2")
	}
	t = time.Unix(int64(binary.BigEndian.Uint64(buf[i:])), 0)
	i += dateSize
	offSchema = stor.ReadSmallOffset(buf[i:])
	i += stor.SmallOffsetLen
	offInfo = stor.ReadSmallOffset(buf[i:])
	i += stor.SmallOffsetLen
	return offSchema, offInfo, t
}
