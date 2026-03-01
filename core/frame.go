// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

// Shared holds the shared variable storage for closures.
// It supports concurrent access when the concurrent flag is set.
type Shared struct {
	values []Value
	MayLock
}

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
	for j, sname := range sharedNames {
		for i, lname := range localNames {
			if lname == sname {
				fr.shared.values[j] = fr.locals[i]
				break
			}
		}
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

// setSharedSlot is split off so setSlot is inlined
func (fr *Frame) setSharedSlot(idx int, val Value) {
	if fr.shared.Lock() {
		defer fr.shared.Unlock()
	}
	fr.shared.values[idx-SharedSlotStart] = val
}
