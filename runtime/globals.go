// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"fmt"
	"sync"
)

// Gnum is a reference to a global name/value
// 0 is invalid
type Gnum = int

type typeGlobal struct{}

// Global is used to group the global methods
// Would be cleaner to use a package but awkward dependencies.
// This approach doesn't seem to add any overhead.
var Global typeGlobal

type globals struct {
	lock     sync.RWMutex
	name2num map[string]Gnum
	names    []string
	values   []Value
	missing  Value
	builtins map[Gnum]Value
	errors   map[Gnum]interface{}
}

var g = globals{
	name2num: make(map[string]Gnum),
	// put ""/nil in first slot so we never use gnum of zero
	names:    []string{""},
	values:   []Value{nil},
	missing:  &SuExcept{}, // type doesn't matter, just has to be unique
	builtins: make(map[Gnum]Value, 100),
	errors:   make(map[Gnum]interface{}),
}

func (typeGlobal) Builtin(name string, value Value) Value {
	// only called by single threaded init so no locking required
	if gn, ok := g.name2num[name]; ok && g.builtins[gn] != nil {
		panic("duplicate builtin: " + name)
	}
	gnum := Global.add(name, nil)
	g.builtins[gnum] = value
	return value // return value to allow: var _ = Global.Builtin(...)
}

func BuiltinNames() []Value {
	names := make([]Value, 0, len(g.builtins))
	for gn := range g.builtins {
		names = append(names, SuStr(Global.Name(gn)))
	}
	return names
}

// Add adds a new name and value to globals.
// This is used for set up of built-in globals
// The return value is so it can be used like:
// var _ = globals.Add(...)
func (typeGlobal) Add(name string, val Value) Gnum {
	g.lock.Lock()
	if _, ok := g.name2num[name]; ok {
		g.lock.Unlock()
		panic("duplicate global: " + name)
	}
	gn := Global.add(name, val)
	g.lock.Unlock()
	return gn
}

func (typeGlobal) add(name string, val Value) Gnum {
	gnum := len(g.names)
	g.name2num[name] = gnum
	g.names = append(g.names, name)
	g.values = append(g.values, val)
	return gnum
}

// TestDef sets a global for tests
func (typeGlobal) TestDef(name string, val Value) {
	g.values[Global.Num(name)] = val
}

// Num returns the global number for a name
// adding it if it doesn't exist.
func (typeGlobal) Num(name string) Gnum {
	// common case, already exists, just need read lock
	g.lock.RLock()
	gn, ok := g.name2num[name]
	g.lock.RUnlock()
	if ok {
		return gn
	}
	// less common case, doesn't exist, need write lock to add
	g.lock.Lock()
	// have to re-check in case another thread beat us to it
	gn, ok = g.name2num[name]
	if !ok {
		gn = Global.add(name, nil)
	}
	g.lock.Unlock()
	return gn
}

// Name returns the name for a global number
func (typeGlobal) Name(gnum Gnum) string {
	g.lock.RLock()
	name := g.names[gnum]
	g.lock.RUnlock()
	return name
}

// Exists returns whether the name exists - for tests
func (typeGlobal) Exists(name string) bool {
	_, ok := g.name2num[name]
	return ok
}

// GetName returns the value for a global name, or panics
func (typeGlobal) GetName(t *Thread, name string) Value {
	return Global.Get(t, Global.Num(name))
}

// Get returns the value for a global number, or panics
func (typeGlobal) Get(t *Thread, gnum Gnum) (result Value) {
	x := Global.Find(t, gnum)
	if x == nil {
		if err, ok := g.errors[gnum]; ok {
			panic(err)
		}
		panic("can't find " + Global.Name(gnum))
	}
	return x
}

// FindName returns the value for a global name, or nil if not found.
// Used by SuRecord to check if a rule exists.
func (typeGlobal) FindName(t *Thread, name string) (result Value) {
	return Global.Find(t, Global.Num(name))
}

// Libload requires dependency injection
var Libload = func(*Thread, Gnum, string) Value { return nil }

var gnPrint = Global.Num("Print")

// Find returns the value for a global number, or nil if not found.
// The two branches on the happy path (x==nil/missing) should be predicted.
func (typeGlobal) Find(t *Thread, gnum Gnum) (result Value) {
	if x, ok := g.builtins[gnum]; ok {
		return x
	}
	g.lock.RLock()
	x := g.values[gnum]
	g.lock.RUnlock()
	if x == nil {
		if _, ok := g.errors[gnum]; ok {
			return nil
		}
		// NOTE: can't hold lock during Libload
		// since compile may need to access Global.
		// That means two threads could both load
		// but they should both get the same value.
		defer func() {
			if err := recover(); err != nil {
				g.errors[gnum] = err // remember the error
				result = nil
			}
		}()
		x = Libload(t, gnum, Global.Name(gnum))
		// for development we want Print even if we don't have stdlib
		if x == nil && gnum == gnPrint {
			fmt.Println("using built-in Print")
			return printBuiltin
		}
		if x == nil {
			x = g.missing // avoid further libloads
		}
		g.lock.Lock()
		g.values[gnum] = x
		g.lock.Unlock()
	}
	if x == g.missing {
		return nil
	}
	return x
}

// GetIfPresent returns the current value (if there is one)
// without doing LibLoad.
// Used by compiler for _Name references
func (typeGlobal) GetIfPresent(name string) (x Value) {
	g.lock.RLock()
	if gnum, ok := g.name2num[name]; ok {
		if x = g.builtins[gnum]; x == nil {
			x = g.values[gnum]
		}
	}
	g.lock.RUnlock()
	return
}

func (typeGlobal) Unload(name string) {
	gnum := Global.Num(name)
	g.lock.Lock()
	g.values[gnum] = nil
	g.lock.Unlock()
}

func (typeGlobal) UnloadAll() {
	g.lock.Lock()
	for i := range g.values {
		g.values[i] = nil
	}
	g.lock.Unlock()
}

// Set is used by LibLoad
func (typeGlobal) Set(gn Gnum, val Value) {
	g.lock.Lock()
	g.values[gn] = val
	g.lock.Unlock()
}

// Copy is used by compile to handle overload inheritance (_Name)
// It copies the value of a slot to a new slot (without a name)
func (typeGlobal) Copy(name string) Gnum {
	g.lock.Lock()
	gn, ok := g.name2num[name]
	if !ok || g.values[gn] == nil {
		g.lock.Unlock()
		panic("can't find " + name)
	}
	newgn := len(g.names)
	g.names = append(g.names, name)
	g.values = append(g.values, g.values[gn])
	g.lock.Unlock()
	return newgn
}
