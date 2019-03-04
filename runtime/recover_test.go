package runtime

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/verify"
)

// confirm the behavior of recover
// i.e. Go call stack is as of panic
// but defer's have been done

func TestRecover(*testing.T) {
	a()
}

func a() {
	defer func() {
		if e := recover(); e != nil {
			verify.That(unwound)
			// debug.PrintStack()
		}
	}()
	b()
}

var unwound = false

func b() {
	defer func() {
		unwound = true
	}()
	c()
}

func c() {
	panic("foo")
}
