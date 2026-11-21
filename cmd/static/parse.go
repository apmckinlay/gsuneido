// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/apmckinlay/gsuneido/compile"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/core"
)

// parse parses a named Suneido constant (function, class, object, or literal)
// into an AST using the compile package's astAspects.
// It returns nil if the source cannot be parsed.
func parse(name, src string) (result core.Value) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "error parsing %q: %v\n", name, r)
			result = nil
		}
	}()
	p := compile.AstParserNamed(name, src)
	result = p.Const()
	if p.Token != tok.Eof {
		p.Error("did not parse all input")
	}
	return result
}
