// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import "sync"

var LibraryOverrides = &libraryOverrides{
	overrides: make(map[string]string),
	originals: make(map[string]Value)}

type libraryOverrides struct {
	overrides map[string]string // by key (lib:name)
	originals map[string]Value  // by name
	lock      sync.Mutex
}

func (lo *libraryOverrides) Put(lib, name, text string) {
	lo.lock.Lock()
	defer lo.lock.Unlock()
	key := lib + ":" + name
	if text != "" {
		if text != lo.overrides[key] {
			if _, ok := lo.overrides[key]; !ok {
				if val := Global.GetIfPresent(name); val != nil {
					lo.originals[name] = val
				}
			}
			lo.overrides[key] = text
			Global.unload(name) // not Unload because it clears original
		}
	} else if _, ok := lo.overrides[key]; ok {
		delete(lo.overrides, key)
		lo.restore(name)
	}
}

func (lo *libraryOverrides) Get(lib, name string) (string, bool) {
	lo.lock.Lock()
	defer lo.lock.Unlock()
	s, ok := lo.overrides[lib+":"+name]
	return s, ok
}

func (lo *libraryOverrides) Unload(name string) {
	lo.lock.Lock()
	defer lo.lock.Unlock()
	delete(lo.originals, name)
}

func (lo *libraryOverrides) ClearOriginals() {
	lo.lock.Lock()
	defer lo.lock.Unlock()
	lo.originals = make(map[string]Value)
}

func (lo *libraryOverrides) Clear() {
	lo.lock.Lock()
	defer lo.lock.Unlock()
	for name := range lo.overrides {
		lo.restore(name)
	}
	lo.overrides = make(map[string]string)
	lo.originals = make(map[string]Value)
}

func (lo *libraryOverrides) restore(name string) {
	if orig, ok := lo.originals[name]; ok {
		Global.SetName(name, orig)
	} else {
		Global.unload(name)
	}
}
