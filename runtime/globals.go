package runtime

import (
	"sync"

	"github.com/apmckinlay/gsuneido/util/verify"
)

// Global is a reference to a global name/value
// Globals are constant
// 0 is invalid
type Global = int

var (
	lock     sync.RWMutex
	name2num = make(map[string]Global)
	// put nil in first slot so we never use gnum of zero
	names  = []string{""}
	values = []Value{nil}
)

// Add adds a new name and value to globals.
//
// This is used for set up of built-in globals
// The return value is so it can be used like:
// var _ = globals.Add(...)
func AddGlobal(name string, val Value) Global {
	lock.Lock()
	defer lock.Unlock()
	if gn, ok := name2num[name]; ok {
		if values[gn] != nil {
			panic("duplicate global: " + name)
		}
		values[gn] = val
		return gn
	}
	gnum := Global(len(names))
	name2num[name] = gnum
	names = append(names, name)
	values = append(values, val)
	verify.That(len(names) == len(values))
	return gnum
}

// Num returns the global number for a name
// adding it if it doesn't exist.
func GlobalNum(name string) Global {
	gn, ok := check(name)
	if ok {
		return Global(gn)
	}
	return AddGlobal(name, nil)
}

func check(name string) (Global, bool) {
	lock.RLock()
	defer lock.RUnlock()
	gn, ok := name2num[name]
	return gn, ok
}

// Name returns the name for a global number
func GlobalName(gnum Global) string {
	lock.RLock()
	defer lock.RUnlock()
	return names[gnum]
}

// Get returns the value for a global
func GetGlobal(gnum Global) Value {
	lock.RLock()
	defer lock.RUnlock()
	return values[gnum]
}

// Exists returns whether the name exists - for tests
func GlobalExists(name string) bool {
	_, ok := name2num[name]
	return ok
}
