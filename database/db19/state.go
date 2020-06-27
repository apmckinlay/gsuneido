// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/apmckinlay/gsuneido/database/db19/meta"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/verify"
)

type DbState struct {
	store *stor.Stor
	meta  *meta.Overlay
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
	updates := state.meta.Merge(tranNum) // outside UpdateState
	UpdateState(func(state *DbState) {
		state.meta.ApplyMerge(updates)
	})
}

//-------------------------------------------------------------------

func Persist() uint64 {
	// NOTE: must not run concurrent instances of this
	state := GetState()
	updates := state.meta.SaveIndexes() // outside UpdateState
	state = UpdateState(func(state *DbState) {
		state.meta.ApplySave(updates)
	})
	return state.Write()
}

const magic = "\xfe\xdc\xba\x98\x76\x54\x32\x10"
const stateLen = 2*len(magic) + meta.Noffsets*stor.SmallOffsetLen

func (state *DbState) Write() uint64 {
	// NOTE: indexes should already have been saved
	stateOff, buf := state.store.Alloc(stateLen)
	copy(buf, magic)
	i := len(magic)
	offsets := state.meta.Write(state.store)
	fmt.Println(offsets)
	for _, o := range offsets {
		stor.WriteSmallOffset(buf[i:], o)
		i += stor.SmallOffsetLen
	}
	copy(buf[i:], magic)
	return stateOff
}

func ReadState(st *stor.Stor, off uint64) *DbState {
	buf := st.Data(off)[:stateLen]
	verify.That(string(buf[:len(magic)]) == magic)
	verify.That(string(buf[stateLen-len(magic):]) == magic)
	i := len(magic)
	var offsets [meta.Noffsets]uint64
	for j := range offsets {
		offsets[j] = stor.ReadSmallOffset(buf[i:])
		i += stor.SmallOffsetLen
	}
	fmt.Println(offsets)
	return &DbState{store: st, meta: meta.FromOffsets(st, offsets)}
}
