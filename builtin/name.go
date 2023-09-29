// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/str"
)

var _ = builtin(Name, "(value)")

func Name(arg Value) Value {
	if named, ok := arg.(Named); ok {
		return SuStr(str.AfterFirst(named.GetName(), ":"))
	}
	return EmptyStr
}
