// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package intern

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestIntern(t *testing.T) {
	assert := assert.T(t).This
	assert(Index("hello")).Is(1)
	assert(Index("world")).Is(2)

	assert(String(1)).Is("hello")
	assert(String(2)).Is("world")

	assert(Index("hello")).Is(1)

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
