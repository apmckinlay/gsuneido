// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"strconv"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/options"
)

var _ = builtin(ServerIP, "() :string")

func ServerIP() Value {
	if options.Action == "client" {
		return SuStr(options.Arg)
	}
	return EmptyStr
}

var _ = builtin(ServerPort, "() :number")

func ServerPort() Value {
	n, _ := strconv.Atoi(options.Port)
	return IntVal(n)
}

var _ = builtin(ServerQ, "() :boolean")

func ServerQ() Value {
	return SuBool(options.Action == "server")
}
