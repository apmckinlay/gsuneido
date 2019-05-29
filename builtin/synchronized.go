package builtin

import (
	"sync"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var lock sync.Mutex

var _ = builtin("Synchronized(block)",
	func(t *Thread, args []Value) Value {
		lock.Lock()
		defer lock.Unlock()
		return t.Call(args[0])
	})
