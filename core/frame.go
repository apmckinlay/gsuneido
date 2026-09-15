// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"sync"
	"sync/atomic"
)

// Frame is the context for a function/method/block invocation.
type Frame struct {

	// this is the instance if we're running a method
	this Value

	// fn is the Function being executed
	fn *SuFunc

	// blockParent is used for block returns
	blockParent *Frame

	// locals are the local variables (including arguments)
	// They are on the thread stack.
	// Note: shared arguments are moved to shared.
	locals []Value

	// shared is the shared variable storage for closures.
	// Slot indexes >= SharedSlotStart map to shared[idx-SharedSlotStart].
	// It is set when entering a function with shared variables,
	// or captured from the parent frame for closure blocks.
	shared *Shared

	// ip is the current index into the Function's code
	ip int

	catchJump int
	catchSp   int
}

// moveLocalsToShared copies shared parameters
// (which will be in both localNames and sharedNames)
// from local to shared
// It is used by Thread.invoke (interp.go) and SuClosure.Call
func (fr *Frame) moveLocalsToShared() {
	if fr.fn == nil || fr.shared == nil || len(fr.shared.values) == 0 {
		return
	}
	localNames := fr.fn.Names[:fr.fn.Nstack]
	sharedNames := fr.fn.Names[fr.fn.Nstack:]
	wrote := false
	for j, sname := range sharedNames {
		for i, lname := range localNames {
			if lname == sname {
				fr.shared.values[j] = fr.locals[i]
				wrote = true
				break
			}
		}
	}
	if wrote && fr.shared.concurrent {
		fr.shared.modified.Store(true)
	}
}

// lookupName finds a variable by name in the frame.
// Shared slots are checked before locals so shared parameters
// are preferred (after moveLocalsToShared).
// It is used by Thread.dyload (interp.go) and Thread.dyn (args.go).
func (fr *Frame) lookupName(name string) (Value, bool) {
	if fr.fn == nil {
		return nil, false
	}
	names := fr.fn.Names
	sharedStart := int(fr.fn.Nstack)
	if fr.shared == nil {
		sharedStart = len(names)
	}
	if sharedStart > len(names) {
		sharedStart = len(names)
	}
	if fr.shared != nil {
		for j := sharedStart; j < len(names); j++ {
			if names[j] != name {
				continue
			}
			k := j - sharedStart
			if k < 0 || k >= len(fr.shared.values) {
				continue
			}
			return fr.getSlot(int(SharedSlotStart) + k), true
		}
	}
	localEnd := min(sharedStart, len(fr.locals))
	for j := range localEnd {
		if names[j] != name {
			continue
		}
		if x := fr.locals[j]; x != nil {
			return x, true
		}
		return nil, false
	}
	return nil, false
}

// getSlot returns the value at the given slot index.
// For indexes < SharedSlotStart, it reads from locals.
// For indexes >= SharedSlotStart, it reads from shared.
func (fr *Frame) getSlot(idx int) Value {
	if idx < SharedSlotStart {
		return fr.locals[idx]
	}
	return fr.getSharedSlot(idx)
}

// getSharedSlot is split off so getSlot is inlined
func (fr *Frame) getSharedSlot(idx int) Value {
	if fr.shared.Lock() {
		defer fr.shared.Unlock()
	}
	return fr.shared.values[idx-SharedSlotStart]
}

// setSlot sets the value at the given slot index.
// For indexes < SharedSlotStart, it writes to locals.
// For indexes >= SharedSlotStart, it writes to shared.
func (fr *Frame) setSlot(idx int, val Value) {
	if idx < SharedSlotStart {
		fr.locals[idx] = val
		return
	}
	fr.setSharedSlot(idx, val)
}

// getSetSlot updates a slot using op(orig, val), and returns either the
// original value (retOrig) or the updated value.
// For shared slots this is done under one lock for atomic read-modify-write.
func (fr *Frame) getSetSlot(idx int, val Value,
	op func(x, y Value) Value, retOrig bool) Value {
	if idx < SharedSlotStart {
		orig := fr.locals[idx]
		if orig == nil {
			panic("uninitialized variable: " + fr.fn.VarName(idx))
		}
		val = op(orig, val)
		fr.locals[idx] = val
		if retOrig {
			return orig
		}
		return val
	}
	return fr.getSetSharedSlot(idx, val, op, retOrig)
}

// getSetSharedSlot is split off so getSetSlot is inlined
func (fr *Frame) getSetSharedSlot(idx int, val Value,
	op func(x, y Value) Value, retOrig bool) Value {
	if fr.shared.Lock() {
		defer fr.shared.Unlock()
	}
	if fr.shared.concurrent {
		fr.shared.modified.Store(true)
	}
	i := idx - SharedSlotStart
	orig := fr.shared.values[i]
	if orig == nil {
		panic("uninitialized variable: " + fr.fn.VarName(idx))
	}
	val = op(orig, val)
	fr.shared.values[i] = val
	if retOrig {
		return orig
	}
	return val
}

// setSharedSlot is split off so setSlot is inlined
func (fr *Frame) setSharedSlot(idx int, val Value) {
	if fr.shared.Lock() {
		defer fr.shared.Unlock()
	}
	if fr.shared.concurrent {
		fr.shared.modified.Store(true)
	}
	fr.shared.values[idx-SharedSlotStart] = val
}

// Shared holds the shared variable storage for closures.
// It supports concurrent access when the concurrent flag is set.
type Shared struct {
	values []Value
	MayLock
	// seenMu guards seen and warned
	seenMu sync.Mutex
	// seen maps a block's SuFunc to the first SuClosure instance
	// dispatched to a thread with this Shared. A second, different
	// instance of the same SuFunc means the block was created in a loop
	// and its instances share the variable. Different SuFunc's sharing
	// the variables (e.g. separate blocks) are not flagged.
	seen map[*SuFunc]*SuClosure
	// warned is set after emitting a warning for this Shared so the
	// warning is only given once.
	warned bool
	// modified is set if a shared variable is assigned while concurrent.
	// A closure running on another thread may then change or observe the
	// value between uses, e.g. a loop variable captured by a deferred block.
	modified atomic.Bool
}

// noteThread records a closure being dispatched to another thread.
// It returns true if this is a repeat instance of a block (created in a
// loop) whose shared variables have been modified.
func (sh *Shared) noteThread(c *SuClosure) bool {
	sh.seenMu.Lock()
	defer sh.seenMu.Unlock()
	if sh.seen == nil {
		sh.seen = make(map[*SuFunc]*SuClosure)
	}
	prev, ok := sh.seen[c.SuFunc]
	if !ok {
		sh.seen[c.SuFunc] = c
		return false
	}
	if prev == c || sh.warned || !sh.modified.Load() {
		return false
	}
	sh.warned = true
	return true
}
