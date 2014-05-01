package intern

import (
	"fmt"
	"runtime"
	"testing"
)

var List []string
var S = "hello"
var T = func() string { return "world" }

func TestIntern(t *testing.T) {
	for n := 0; n < 100000; n++ {
		// List = append(List, S+T())
		List = append(List, Intern(S+T()))
	}
	runtime.GC()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Println("allocated ", ms.Alloc)
}
