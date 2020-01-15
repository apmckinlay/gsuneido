// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("Fatal(msg)", func(a Value) Value {
	goc.Fatal(ToStrOrString(a))
	return nil
})
