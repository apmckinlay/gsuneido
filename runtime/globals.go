package runtime

import (
	"sync"

	"github.com/apmckinlay/gsuneido/util/verify"
)

// Gnum is a reference to a global name/value
// 0 is invalid
type Gnum = int

type globals struct {
	lock     sync.RWMutex
	name2num map[string]Gnum
	names    []string
	values   []Value
	missing  Value
}

var Global = globals{
	name2num: make(map[string]Gnum),
	// put ""/nil in first slot so we never use gnum of zero
	names:   []string{""},
	values:  []Value{nil},
	missing: &SuExcept{}, // type doesn't matter, just has to be unique
}

// Add adds a new name and value to globals.
//
// This is used for set up of built-in globals
// The return value is so it can be used like:
// var _ = globals.Add(...)
func (g *globals) Add(name string, val Value) Gnum {
	g.lock.Lock()
	defer g.lock.Unlock()
	if _, ok := g.name2num[name]; ok {
		panic("duplicate global: " + name)
	}
	gnum := len(g.names)
	g.name2num[name] = gnum
	g.names = append(g.names, name)
	g.values = append(g.values, val)
	verify.That(len(g.names) == len(g.values))
	return gnum
}

// TestDef sets a global for tests
func (g *globals) TestDef(name string, val Value) {
	g.values[g.Num(name)] = val
}

// Num returns the global number for a name
// adding it if it doesn't exist.
func (g *globals) Num(name string) Gnum {
	if gn, ok := g.check(name); ok {
		return gn
	}
	return g.Add(name, nil)
}

func (g *globals) check(name string) (Gnum, bool) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	gn, ok := g.name2num[name]
	return gn, ok
}

// Name returns the name for a global number
func (g *globals) Name(gnum Gnum) string {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.names[gnum]
}

// Get returns the value for a global
func (g *globals) Get(gnum Gnum) Value {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.values[gnum]
}

// Exists returns whether the name exists - for tests
func (g *globals) Exists(name string) bool {
	_, ok := g.name2num[name]
	return ok
}
