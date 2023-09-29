// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(fatal, "(msg)")

func fatal(a Value) Value {
	Fatal(ToStrOrString(a))
	return nil
}
