// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"github.com/apmckinlay/gsuneido/compile/lexer"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("QueryScanner(string)",
	func(arg Value) Value {
		return &suScanner{name: "QueryScanner",
			lxr: *lexer.NewQueryLexer(ToStr(arg))}
	})
