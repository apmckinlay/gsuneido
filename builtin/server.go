// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("ServerIP()", func() Value {
	return SuStr(options.Client)
})

var _ = builtin0("ServerPort()", func() Value {
	n, _ := strconv.Atoi(options.Port)
	return IntVal(n)
})

var _ = builtin0("Server?()", func() Value {
	return False
})
