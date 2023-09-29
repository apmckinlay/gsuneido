// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

// These methods are specific for int values

var _ = exportMethods(&IntMethods)

var _ = method(int_Int, "()")

func int_Int(this Value) Value {
	return this
}

var _ = method(int_Frac, "()")

func int_Frac(Value) Value {
	return Zero
}
