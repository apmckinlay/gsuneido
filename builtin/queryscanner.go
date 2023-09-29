// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"github.com/apmckinlay/gsuneido/compile/lexer"
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(QueryScanner, "(string)")

func QueryScanner(arg Value) Value {
	return &suScanner{name: "QueryScanner",
		lxr: *lexer.NewQueryLexer(ToStr(arg))}
}
