// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"sync"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var builtinNamesOnce sync.Once
var builtinNames *SuObject

var _ = builtin0("BuiltinNames()", func() Value {
	builtinNamesOnce.Do(func() {
		builtinNames = NewSuObject(BuiltinNames()...)
		builtinNames.Sort(nil, False)
		builtinNames.SetReadOnly()
	})
	return builtinNames
})
