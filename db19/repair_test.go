// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	rand "math/rand/v2"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestRepair(*testing.T) {
	// go test -tags portable -run ^TestRepair$ ./db19 -v -vet=off -count=1
	if testing.Short() {
		return
	}
	Repair("../suneido.db", nil)
}

func TestSearch(*testing.T) {
	// go test -tags portable -run ^TestSearch$ ./db19 -v -vet=off -count=1
	if testing.Short() {
		return
	}
	store, err := stor.MmapStor("../suneido.db", stor.Read)
	if err != nil {
		fmt.Println(err)
	}
	defer store.Close(true)
	r := repair{store: store}
	i, off, _ := r.search()
	if off == 0 {
		fmt.Println("no valid states found")
	}
	fmt.Println("result", i, off)
}

func TestPrintStates(t *testing.T) {
	// go test -tags portable -run ^TestPrintStates$ ./db19 -v -vet=off -count=1
	if testing.Short() {
		return
	}
	PrintStates("../suneido.db", true)
}

func TestStates(t *testing.T) {
	// go test -tags portable -run ^TestStates$ ./db19 -v -vet=off -count=1
	if testing.Short() {
		return
	}
	store, err := stor.MmapStor("../suneido.db", stor.Read)
	if err != nil {
		t.Fatal(err)
	}

	var offsets []uint64
	off := store.Size()
	for {
		off = store.LastOffset(off, magic1, nil)
		if off == 0 { // no more
			break
		}
		offsets = append(offsets, off)
	}
	ec := &errCorrupt{}
	check := func(i int) {
		off := offsets[i]
		state := getState(store, off)
		if state == nil {
			fmt.Println(i, off, "read state failed")
			return
		}
		ec = checkState(state, checkTable, ec.Table(), ec.Ixcols())
		if ec != nil {
			fmt.Println(i, off, ec)
			return
		}
		fmt.Println(i, off, "good")
	}
	check(0)
	check(len(offsets) - 1)
	for range 10 {
		i := rand.IntN(len(offsets))
		check(i)
	}
}

func TestBisect(t *testing.T) {
	if testing.Short() {
		return
	}
	const n = 10001
	// lastGood := rand.IntN(n)
	for lastGood := 0; lastGood < n; lastGood++ {
		// fmt.Println("last good:", lastGood)
		check := func(i int) bool {
			// fmt.Println("check", i)
			return i >= lastGood
		}
		good := -1
		i := 0
		var prev int
		for skip := 1; ; skip *= 2 {
			if i > n {
				i = n
			}
			if check(i) {
				good = i
				break
			}
			if i == n {
				// fmt.Println("no good value found")
				return
			}
			prev = i
			i += skip
		}
		// fmt.Println("good:", good, "so >", prev, "<=", good)
		lo := prev
		hi := good
		for lo < hi-1 {
			mid := lo + (hi-lo)/2
			if check(mid) {
				hi = mid
			} else {
				lo = mid
			}
		}
		// fmt.Println("=>", hi)
		assert.This(hi).Is(lastGood)
	}
}

func TestScanner(t *testing.T) {
	if testing.Short() {
		return
	}
	// go test -tags portable -run ^TestScanner$ ./db19 -v -vet=off -count=1
	store, err := stor.MmapStor("../suneido.db", stor.Read)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close(true)

	scnr := newScanner(store)
	for i := 0; ; i++ {
		off := scnr.get(i)
		if off == 0 {
			break
		}
		fmt.Println(i, off)
	}
}
