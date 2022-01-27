// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"fmt"
	"log"
	"math"
	"sync"
)

// globals generally follows the usual style that public methods lock
// and private methods do not.
// builtins does not require locking
// since it is populated by init which is single threaded
// and is immutable after that.

// Gnum is a reference to a global name/value
// 0 is invalid
type Gnum = int

type typeGlobal struct{}

// Global is used to group the global methods
// Would be cleaner to use a package but awkward dependencies.
// This approach doesn't seem to add any overhead.
var Global typeGlobal

// globals stores the value of global names.
//	Normal: values[gnum] != nil
//	Unloaded: values[gnum] == nil
//	Error: values[gnum] == nil, errors[gnum] == string
//	Missing: values[gnum] == nil, errors[gnum] == false
type globals struct {
	lock     sync.RWMutex
	name2num map[string]Gnum
	names    []string
	values   []Value
	builtins map[Gnum]Value
	errors   map[Gnum]interface{}
	noDef    map[string]struct{} // used by FindName
}

var g = globals{
	name2num: make(map[string]Gnum),
	// put ""/nil in first slot so we never use gnum of zero
	names:    []string{""},
	values:   []Value{nil},
	builtins: make(map[Gnum]Value, 100),
	errors:   make(map[Gnum]interface{}),
	noDef:    make(map[string]struct{}),
}

func (typeGlobal) Builtin(name string, value Value) Value {
	// only called by single threaded init so no locking required
	if gn, ok := g.name2num[name]; ok && g.builtins[gn] != nil {
		log.Fatalln("duplicate builtin: " + name)
	}
	gnum := Global.add(name, nil)
	g.builtins[gnum] = value
	return value // return value to allow: var _ = Global.Builtin(...)
}

func BuiltinNames() []Value {
	names := make([]Value, len(g.builtins))
	i := 0
	for gn := range g.builtins {
		names[i] = SuStr(Global.Name(gn))
		i++
	}
	return names
}

// Add adds a new name and value to globals.
// This is used for set up of built-in globals
// The return value is so it can be used like:
// var _ = globals.Add(...)
func (typeGlobal) Add(name string, val Value) Gnum {
	g.lock.Lock()
	defer g.lock.Unlock()
	if _, ok := g.name2num[name]; ok {
		panic("duplicate global: " + name)
	}
	return Global.add(name, val)
}

// add requires caller to write Lock
func (typeGlobal) add(name string, val Value) Gnum {
	gnum := len(g.names)
	if gnum > math.MaxUint16 {
		Fatal("too many globals")
	}
	g.name2num[name] = gnum
	g.names = append(g.names, name)
	g.values = append(g.values, val)
	return gnum
}

// TestDef sets a global for tests.
// WARNING: no locking
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
	defer g.lock.Unlock()
	// have to re-check in case another thread beat us to it
	if gn, ok = g.name2num[name]; ok {
		return gn
	}
	return Global.add(name, nil)
}

// num returns the global number for a name
// adding it if it doesn't exist.
// It is the same as Num but caller must write Lock.
func (typeGlobal) num(name string) Gnum {
	gn, ok := g.name2num[name]
	if ok {
		return gn
	}
	return Global.add(name, nil)
}

// Name returns the name for a global number
func (typeGlobal) Name(gnum Gnum) string {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.names[gnum]
}

// Exists returns whether the name exists - for tests
func (typeGlobal) Exists(name string) bool {
	g.lock.RLock()
	defer g.lock.RUnlock()
	_, ok := g.name2num[name]
	return ok
}

// GetName returns the value for a global name, or panics
func (typeGlobal) GetName(t *Thread, name string) Value {
	return Global.Get(t, Global.Num(name))
}

// Get returns the value for a global number, or panics
func (typeGlobal) Get(t *Thread, gnum Gnum) (result Value) {
	if x := Global.Find(t, gnum); x != nil {
		return x // common fast path
	}
	g.lock.RLock()
	defer g.lock.RUnlock()
	name := g.names[gnum]
	if e, ok := g.errors[gnum]; ok && e != false {
		panic("error loading " + name + " " + fmt.Sprint(e))
	}
	panic("can't find " + name)
}

// FindName returns the value for a global name, or nil if not found.
// Used to check if a trigger or rule exists.
// Avoids creating a global if no definition is found.
// Uses noDef to avoid repeatedly looking up nonexistent names.
func (typeGlobal) FindName(t *Thread, name string) Value {
	g.lock.RLock()
	if gn, ok := g.name2num[name]; ok { // name exists
		x := g.values[gn]
		if x != nil {
			g.lock.RUnlock()
			return x
		}
	}
	if _, ok := g.noDef[name]; ok {
		g.lock.RUnlock()
		return nil
	}
	g.lock.RUnlock()
	// NOTE: can't hold lock during Libload
	// since compile may need to access Global.
	x, e := Libload(t, name)
	if e != nil {
		g.lock.Lock()
		defer g.lock.Unlock()
		gnum := Global.num(name)
		g.errors[gnum] = e
		panic("error loading " + name + " " + fmt.Sprint(e))
	}
	if x == nil {
		g.lock.Lock()
		defer g.lock.Unlock()
		g.noDef[name] = struct{}{}
	}
	return x
}

// Libload requires dependency injection
var Libload = func(*Thread, string) (Value, interface{}) { return nil, nil }

var gnPrint = Global.Num("Print")

// Find returns the value for a global number, or nil if not found.
func (typeGlobal) Find(t *Thread, gnum Gnum) (result Value) {
	if x, ok := g.builtins[gnum]; ok {
		return x // common fast path
	}
	g.lock.RLock()
	x := g.values[gnum]
	if x != nil {
		g.lock.RUnlock()
		return x // common fast path
	}
	if _, ok := g.errors[gnum]; ok {
		g.lock.RUnlock()
		return nil
	}
	g.lock.RUnlock()
	// NOTE: can't hold lock during Libload
	// since compile may need to access Global.
	var e interface{}
	name := Global.Name(gnum)
	x, e = Libload(t, name)
	if e != nil {
		Global.SetErr(gnum, e)
		panic("error loading " + name + " " + fmt.Sprint(e))
	}
	if x == nil {
		if gnum == gnPrint {
			// for development we want Print even if we don't have stdlib
			fmt.Println("using built-in Print")
			return printBuiltin
		}
		Global.SetErr(gnum, false) // avoid further libloads
		return nil
	}
	Global.Set(gnum, x)
	return x
}

// GetIfPresent returns the current value (if there is one)
// without doing LibLoad.
// Used by compiler for _Name references
func (typeGlobal) GetIfPresent(name string) Value {
	g.lock.RLock()
	defer g.lock.RUnlock()
	if gnum, ok := g.name2num[name]; ok {
		if x := g.builtins[gnum]; x != nil {
			return x
		}
		if x := g.values[gnum]; x != nil {
			return x
		}
	}
	return nil
}

func (typeGlobal) Unload(name string) {
	g.lock.Lock()
	defer g.lock.Unlock()
	gnum := Global.num(name)
	g.values[gnum] = nil
	delete(g.errors, gnum)
	delete(g.noDef, name)
}

func (typeGlobal) UnloadAll() {
	g.lock.Lock()
	defer g.lock.Unlock()
	for i := range g.values {
		g.values[i] = nil
	}
	g.errors = make(map[Gnum]interface{})
	g.noDef = make(map[string]struct{})
}

// Set is used by libload
func (typeGlobal) Set(gn Gnum, val Value) {
	g.lock.Lock()
	g.values[gn] = val
	g.lock.Unlock()
}

func (typeGlobal) SetErr(gn Gnum, e interface{}) {
	g.lock.Lock()
	g.errors[gn] = e
	g.lock.Unlock()
}

// Overload is used by compile to handle overload inheritance (_Name)
func (typeGlobal) Overload(name string, val Value) Gnum {
	g.lock.Lock()
	defer g.lock.Unlock()
	Global.num(name[1:]) // ensure original exists
	newgn := len(g.names)
	g.names = append(g.names, name)
	g.values = append(g.values, val)
	return newgn
}
