// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = exportMethods(&RegexMethods, "regex")

var _ = method(regex_Disasm, "() :string")

func regex_Disasm(this Value) Value {
	return StrVal(this.(SuRegex).Pat.String())
}
