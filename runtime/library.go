// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

var LibraryOverrides = make(map[string]string)

func LibraryOverride(lib, name, text string) {
	key := lib + ":" + name
	if text != "" {
		LibraryOverrides[key] = text
	} else {
		delete(LibraryOverrides, key)
	}
	Global.Unload(name)
}

func LibraryOverrideClear() {
	for name := range LibraryOverrides {
		Global.Unload(name)
	}
	LibraryOverrides = make(map[string]string)
}
