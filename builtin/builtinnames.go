// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"sync"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var builtinNamesOnce sync.Once
var builtinNames *SuObject

var _ = builtin(BuiltinNames, "()")

func BuiltinNames() Value {
	builtinNamesOnce.Do(func() {
		builtinNames = NewSuObject(GetBuiltinNames())
		builtinNames.Sort(nil, False)
		builtinNames.SetReadOnly()
	})
	return builtinNames
}
