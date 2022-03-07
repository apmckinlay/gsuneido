// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"math/rand"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/bits"
	"github.com/apmckinlay/gsuneido/util/generic/hamt"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestMeta(t *testing.T) {
	tbl := InfoHamt{}.Mutable()
	const n = 1000
	data := make([]string, n)
	randStr := str.UniqueRandom(4, 4)
	for i := 0; i < n; i++ {
		data[i] = randStr()
		tbl.Put(&Info{
			Table: data[i],
		})
	}
	tbl = tbl.Freeze()
	meta := &Meta{info: hamt.Chain[string, *Info]{Hamt: tbl}}

	for i := 0; i < 4; i++ {
		m := meta.Mutable()
		for i := 0; i < 5; i++ {
			m.GetRoInfo(data[rand.Intn(100)])
			m.GetRwInfo(data[rand.Intn(100)])
		}
		// end of transaction, merge back to global
		meta = m.LayeredOnto(meta)
	}

	// persist state
	st := stor.HeapStor(32 * 1024)
	meta.Write(st)

	// test that nothing is written if no changes
	size := st.Size()
	meta.Write(st)
	assert.T(t).This(st.Size()).Is(size)
}

func TestChunkingSimulation(t *testing.T) {
	const n = 128
	const maxChain = 5
	clock := 0
	data := make([]int, n)
	// chunks represents what would be stored in the database
	chunks := [][]int{}
	ages := []int{} // parallel with chunks
	updates := rand.Perm(n)
	// run simulation
	for step := 0; step < n; step++ {
		data[updates[step]] = clock
		// fmt.Println("--- update", updates[step])
		// number of previous chunks to merge with
		merge := ord.Min(len(chunks), bits.TrailingOnes(clock))
		if len(chunks) >= maxChain {
			// fmt.Println("MAX")
			merge = maxChain
		}
		oldest := clock
		if merge > 0 {
			oldest = ages[len(ages)-merge]
		}
		// fmt.Println("clock", clock, "merge", merge, "oldest", oldest)
		chunk := []int{}
		for i, p := range data {
			if p >= oldest {
				chunk = append(chunk, i)
			}
		}
		chunks = chunks[:len(chunks)-merge]
		ages = ages[:len(ages)-merge]
		chunks = append(chunks, chunk)
		ages = append(ages, oldest)
		// fmt.Println("chunks", chunks)
		// fmt.Println("ages", ages)
		assert.That(len(chunks) == len(ages))
		clock++
	}
	assert.T(t).This(len(chunks)).Is(1)
	assert.T(t).This(ages[0]).Is(0)
	for i := 0; i < n; i++ {
		assert.That(chunks[0][i] == i)
	}
}
