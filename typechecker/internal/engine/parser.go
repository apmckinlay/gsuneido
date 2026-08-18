// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"github.com/apmckinlay/gsuneido/compile"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/core"
)

func ParseClass(cls string) core.Value {
	p := compile.AstParser(cls)
	result := p.Const()
	if p.Token != tok.Eof {
		p.Error("could not parse all the input")
	}
	return result
}
