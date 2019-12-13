package builtin

import (
	"log"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("ErrorLog(string)",
	func(arg Value) Value {
		log.Println(arg)
		return nil
	})
