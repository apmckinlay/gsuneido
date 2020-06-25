// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/verify"
	md "github.com/apmckinlay/gsuneido/database/db19/metadata"
)

type DbState struct {
	store         *stor.Stor
	baseInfo      *md.InfoPacked
	baseSchema    *md.SchemaPacked
	overSchemaOff uint64
	// memMeta includes both schema and info
	memMeta *md.TableInfoHtbl
}

// theState is the central immutable state of the database.
// It must be accessed atomically and only updated via UpdateState.
var theState = unsafe.Pointer(&DbState{})

// stateMutex guards updates to theState
var stateMutex sync.Mutex

func GetState() *DbState {
	return (*DbState)(atomic.LoadPointer(&theState))
}

// UpdateState applies the given update function to a copy of theState
// and sets theState to the result.
// Note: the state passed to the update function is a *shallow* copy,
// it is up to the function to make copies of any nested containers.
func UpdateState(fn func(*DbState)) *DbState {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	newState := *GetState() // shallow copy
	fn(&newState)
	atomic.StorePointer(&theState, unsafe.Pointer(&newState))
	return &newState
}

//-------------------------------------------------------------------

// Merge updates the base fbtree's with the first overlay mbtree
// for the given transaction number.
func Merge(tranNum int) {
	state := GetState()
	updates := state.memMeta.Merge(tranNum) // outside UpdateState
	UpdateState(func(state *DbState) {
		state.memMeta = state.memMeta.WithMerged(updates)
	})
}

//-------------------------------------------------------------------

func Persist() uint64 {
	// NOTE: must not run concurrent instances of this
	state := GetState()
	updates := state.memMeta.SaveIndexes() // outside UpdateState
	state = UpdateState(func(state *DbState) {
		state.memMeta = state.memMeta.WithSaved(updates)
	})
	return state.Write()
}

const magic = "\xfe\xdc\xba\x98\x76\x54\x32\x10"
const stateLen = 2*len(magic) + 4*stor.SmallOffsetLen

func (state *DbState) Write() uint64 {
	// NOTE: indexes should already have been saved
	overInfoOff := state.memMeta.WriteInfo(state.store)
	stateOff, buf := state.store.Alloc(stateLen)
	copy(buf, magic)
	stor.WriteSmallOffset(buf[8:], state.baseSchema.Offset())
	stor.WriteSmallOffset(buf[13:], state.baseInfo.Offset())
	stor.WriteSmallOffset(buf[18:], state.overSchemaOff)
	stor.WriteSmallOffset(buf[23:], overInfoOff)
	copy(buf[28:], magic)
	return stateOff
}

func ReadState(st *stor.Stor, off uint64) *DbState {
	buf := st.Data(off)[:stateLen]
	verify.That(string(buf[:len(magic)]) == magic)
	verify.That(string(buf[28:28+len(magic)]) == magic)
	baseSchemaOff := stor.ReadSmallOffset(buf[8:])
	baseInfoOff := stor.ReadSmallOffset(buf[13:])
	overSchemaOff := stor.ReadSmallOffset(buf[18:])
	overInfoOff := stor.ReadSmallOffset(buf[23:])
	return &DbState{
		store: st,
		overSchemaOff: overSchemaOff,
		baseSchema: md.NewSchemaPacked(st, baseSchemaOff),
		baseInfo: md.NewInfoPacked(st, baseInfoOff),
		memMeta: md.ReadInfo(st, overInfoOff).ReadSchema(st, overSchemaOff),
	}
}
