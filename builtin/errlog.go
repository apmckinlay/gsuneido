// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"log"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("ErrorLog(string)",
	func(arg Value) Value {
		log.Println(ToStrOrString(arg))
		return nil
	})
