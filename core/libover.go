// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import "sync"

var LibraryOverrides = &libraryOverrides{
	overrides: make(map[string]libOver),
	originals: make(map[string]Value)}

type libraryOverrides struct {
	overrides map[string]libOver // by name
	originals map[string]Value   // by name
	lock      sync.Mutex
}

type libOver struct {
	lib  string
	text string
}

func (lo *libraryOverrides) Put(lib, name, text string) {
	lo.lock.Lock()
	defer lo.lock.Unlock()
	if text != "" {
		ov, ok := lo.overrides[name]
		if ok && lib != ov.lib {
			panic("LibraryOverride: override already exists for " +
				ov.lib + ":" + name)
		}
		if !ok || text != ov.text {
			if !ok {
				if val := Global.GetIfPresent(name); val != nil {
					lo.originals[name] = val
				}
			}
			lo.overrides[name] = libOver{lib: lib, text: text}
			Global.unload(name) // not Unload because it clears original
		}
	} else if _, ok := lo.overrides[name]; ok {
		lo.restore(name)
		delete(lo.overrides, name)
		delete(lo.originals, name)
	}
}

func (lo *libraryOverrides) Get(name string) (string, string) {
	lo.lock.Lock()
	defer lo.lock.Unlock()
	ov, _ := lo.overrides[name]
	return ov.lib, ov.text
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
	clear(lo.overrides)
	clear(lo.originals)
}

func (lo *libraryOverrides) restore(name string) {
	if orig, ok := lo.originals[name]; ok {
		Global.SetName(name, orig)
	} else {
		Global.unload(name)
	}
}
