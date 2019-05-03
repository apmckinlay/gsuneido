package intern

import (
	"fmt"
	"runtime"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestIntern(t *testing.T) {
	Assert(t).That(Index("hello"), Equals(1))
	Assert(t).That(Index("world"), Equals(2))

	Assert(t).That(String(1), Equals("hello"))
	Assert(t).That(String(2), Equals("world"))

	Assert(t).That(Index("hello"), Equals(1))

	for c := 'a'; c <= 'z'; c++ {
		Index(string(c))
	}
}

var List []string
var S = "hello"
var T = func() string { return "world" }

func TestMemory(*testing.T) {
	for n := 0; n < 10000; n++ {
		// List = append(List, S+T())
		List = append(List, String(Index(S+T())))
	}
	runtime.GC()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Println("allocated ", ms.Alloc)
}
