package runtime

import (
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
}

var g = globals{
	name2num: make(map[string]Gnum),
	// put ""/nil in first slot so we never use gnum of zero
	names:   []string{""},
	values:  []Value{nil},
	missing: &SuExcept{}, // type doesn't matter, just has to be unique
}

var Libload = func(string) Value { return nil }

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

// Get returns the value for a global
func (typeGlobal) Get(gnum Gnum) Value {
	g.lock.RLock()
	x := g.values[gnum]
	g.lock.RUnlock()
	if x == nil {
		// NOTE: can't hold lock during Libload
		// since compile may need to access Global.
		// That means two threads could both load
		// but they should both get the same value.
		x = Libload(Global.Name(gnum))
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
