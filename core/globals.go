// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/atomics"
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
//
//	Normal: values[gnum] != nil
//	Unloaded: values[gnum] == nil
//	Error: values[gnum] == nil, errors[gnum] == string
//	Missing: values[gnum] == nil, errors[gnum] == false
type globals struct {
	name2num map[string]Gnum
	builtins map[Gnum]Value
	errors   map[Gnum]any
	noDef    map[string]struct{} // used by FindName
	names    []string
	values   []Value
	lock     sync.RWMutex
}

var g = globals{
	name2num: map[string]Gnum{"Suneido": 1},
	// put ""/nil in first slot so we never use gnum of zero
	names:    []string{"", "Suneido"},
	values:   []Value{nil, nil},
	builtins: make(map[Gnum]Value, 100),
	errors:   make(map[Gnum]any),
	noDef:    make(map[string]struct{}),
}

const GnSuneido = 1

var Suneido *SuneidoObject

var _ = func() int { // needs to be var, init() is run later
	assert.This(Global.Num("Suneido")).Is(GnSuneido)
	Suneido = new(SuneidoObject)
	Suneido.SetConcurrent()
	Global.Set(GnSuneido, Suneido)
	g.builtins[GnSuneido] = Suneido
	return 0
}()

// Builtin is used to set up built-in values
func (typeGlobal) Builtin(name string, value Value) Value {
	// only called by single threaded init so no locking required
	if gn, ok := g.name2num[name]; ok && g.builtins[gn] != nil {
		log.Fatalln("duplicate builtin: " + name)
	}
	gnum := Global.add(name, nil)
	g.builtins[gnum] = value
	return value // return value to allow: var _ = Global.Builtin(...)
}

func GetBuiltinNames() []Value {
	names := make([]Value, len(g.builtins))
	i := 0
	for gn := range g.builtins {
		names[i] = SuStr(Global.Name(gn))
		i++
	}
	return names
}

// Add is used by tests
func (typeGlobal) Add(name string, val Value) Gnum {
	g.lock.Lock()
	defer g.lock.Unlock()
	return Global.add(name, val)
}

// add creates a new name and value to globals.
// This is used for set up of built-in globals
// Callers should write Lock (unless during init)
func (typeGlobal) add(name string, val Value) Gnum {
	if _, ok := g.name2num[name]; ok {
		panic("duplicate global: " + name)
	}
	gnum := len(g.names)
	if gnum > math.MaxUint16 {
		Fatal("too many globals")
	}
	g.name2num[name] = gnum // this is the only place we add to name2num
	g.names = append(g.names, name)
	g.values = append(g.values, val)
	return gnum
}

var _ = AddInfo("core.nGlobal", func() int { return len(g.names) })
var _ = AddInfo("core.lastGlobals", func() string {
	s := ""
	for i := len(g.names) - 1; i >= len(g.names)-10; i-- {
		s += g.names[i] + "\n"
	}
	return s
})

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
func (typeGlobal) GetName(th *Thread, name string) Value {
	return Global.Get(th, Global.Num(name))
}

// Get returns the value for a global number, or panics
func (typeGlobal) Get(th *Thread, gnum Gnum) (result Value) {
	if x := Global.Find(th, gnum); x != nil {
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
func (typeGlobal) FindName(th *Thread, name string) Value {
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
	x, e := Libload(th, name)
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
	} else {
		Global.SetName(name, x)
	}
	return x
}

// Libload requires dependency injection
var Libload = func(*Thread, string) (Value, any) { return nil, nil }

var gnPrint = Global.Num("Print")

// Find returns the value for a global number, or nil if not found.
func (typeGlobal) Find(th *Thread, gnum Gnum) (result Value) {
	if x, ok := g.builtins[gnum]; ok {
		if gnum == GnSuneido {
			if suneido := th.Suneido.Load(); suneido != nil {
				return suneido
			}
		}
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
	var e any
	name := Global.Name(gnum)
	x, e = Libload(th, name)
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

// GetIfPresent is used by LibraryOverride.
// It returns the current value (if there is one) without doing LibLoad.
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
	Global.unload(name)
	LibraryOverrides.Unload(name)
}

func (typeGlobal) unload(name string) {
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
	g.errors = make(map[Gnum]any)
	g.noDef = make(map[string]struct{})
	LibraryOverrides.ClearOriginals()
	LibsList.Store(nil)
}

// LibsList is used by libload
var LibsList atomics.Value[[]string]

func (typeGlobal) SetName(name string, val Value) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.values[Global.num(name)] = val
}

func (typeGlobal) Set(gn Gnum, val Value) {
	g.lock.Lock()
	g.values[gn] = val
	g.lock.Unlock()
}

func (typeGlobal) SetErr(gn Gnum, e any) {
	g.lock.Lock()
	g.errors[gn] = e
	g.lock.Unlock()
}

// Overload is used by compile to handle overload inheritance (_Name).
// It creates a new slot to contain the previous value.
// The original slot will be set to the final visible value.
// name2num points to the original slot.
func (typeGlobal) Overload(name string, prevVal Value) Gnum {
	g.lock.Lock()
	defer g.lock.Unlock()
	Global.num(name[1:]) // ensure original exists
	newgn := len(g.names)
	g.names = append(g.names, name)
	g.values = append(g.values, prevVal)
	return newgn
}
