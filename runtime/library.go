// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

var LibraryOverrides = make(map[string]string) // by key (lib:name)
var LibraryOriginals = make(map[string]Value)  // by name

func LibraryOverride(lib, name, text string) {
	key := lib + ":" + name
	if text != "" {
		if text != LibraryOverrides[key] {
			if _, ok := LibraryOverrides[key]; !ok {
				if val := Global.GetIfPresent(name); val != nil {
					LibraryOriginals[name] = val
				}
			}
			LibraryOverrides[key] = text
			Global.unload(name) // not Unload because it clears original
		}
	} else if _, ok := LibraryOverrides[key]; ok {
		delete(LibraryOverrides, key)
		overrideRestore(name)
	}
}

func LibraryOverrideClear() {
	for name := range LibraryOverrides {
		overrideRestore(name)
	}
	LibraryOverrides = make(map[string]string)
	LibraryOriginals = make(map[string]Value)
}

func overrideRestore(name string) {
	if val, ok := LibraryOriginals[name]; ok {
		Global.SetName(name, val)
	} else {
		Global.Unload(name)
	}
}
