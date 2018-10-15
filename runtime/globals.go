package runtime

import (
	"sync"

	"github.com/apmckinlay/gsuneido/util/verify"
)

var (
	lock     sync.RWMutex
	name2num = make(map[string]int)
	// put nil in first slot so we never use gnum of zero
	names  = []string{""}
	values = []Value{nil}
)

// Add adds a new name and value to globals.
//
// This is used for set up of built-in globals
// The return value is so it can be used like:
// var _ = globals.Add(...)
func AddGlobal(name string, val Value) int {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := name2num[name]; ok {
		panic("duplicate global: " + name)
	}
	gnum := len(names)
	name2num[name] = gnum
	names = append(names, name)
	values = append(values, val)
	verify.That(len(names) == len(values))
	return gnum
}

// Num returns the global number for a name
// adding it if it doesn't exist.
func GlobalNum(name string) int {
	gn, ok := check(name)
	if ok {
		return gn
	}
	return AddGlobal(name, nil)
}

func check(name string) (int, bool) {
	lock.RLock()
	defer lock.RUnlock()
	gn, ok := name2num[name]
	return gn, ok
}

// Name returns the name for a global number
func GlobalName(gnum int) string {
	lock.RLock()
	defer lock.RUnlock()
	return names[gnum]
}

// Get returns the value for a global
func GetGlobal(gnum int) Value {
	lock.RLock()
	defer lock.RUnlock()
	return values[gnum]
}

// Exists returns whether the name exists - for tests
func GlobalExists(name string) bool {
	_, ok := name2num[name]
	return ok
}
