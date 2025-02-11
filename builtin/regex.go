// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = exportMethods(&RegexMethods, "regex")

var _ = method(regex_Disasm, "()")

func regex_Disasm(this Value) Value {
	return SuStr(this.(SuRegex).Pat.String())
}
