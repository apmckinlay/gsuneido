// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/stor"
)

const dataOffset = 1 * 1024 * 1024 * 1024  // 1 gb
const indexOffset = 2 * 1024 * 1024 * 1024 // 2 gb
const recordSize = 1024
const nRecords = 100 * 1024
const indexFanout = 100
const nNodes = nRecords / indexFanout

var Sum int
var Key string

func BenchmarkAccess1(b *testing.B) {
	read := func(buf []byte, i, n int) int {
		sum := 0
		for _, b := range buf[i:i+n] {
			sum += int(b)
		}
		return sum
	}
	var slotSize = 320
	st, err := stor.MmapStor("../../../big.db", stor.Read)
	if err != nil {
		b.Fatal(err)
	}
	i := 0
	for b.Loop() { // per node
		nodeOffset := indexOffset + uint64(i*slotSize)
		node := st.Data(nodeOffset)
		for j := range indexFanout { // for each slot
			read(node, j*slotSize, slotSize)
		}
		i = (i + 1) % nNodes
	}
}

func BenchmarkAccess2(b *testing.B) {
	getLeafKey := func(st *stor.Stor, off uint64) string {
		data := st.Data(off)
		var sb strings.Builder
		sb.Write(data[100:110])
		sb.Write(data[500:510])
		sb.Write(data[900:910])
		return sb.String()
	}
	read := func(buf []byte, i, n int) int {
		sum := 0
		for _, b := range buf[i:i+n] {
			sum += int(b)
		}
		return sum
	}
	var slotSize = 10
	st, err := stor.MmapStor("../../../big.db", stor.Read)
	if err != nil {
		b.Fatal(err)
	}
	i := 0
	for b.Loop() { // per node
		nodeOffset := indexOffset + uint64(i*slotSize)
		node := st.Data(nodeOffset)
		for j := range indexFanout { // for each slot
			read(node, j*slotSize, slotSize)
			r := rand.IntN(nRecords)
			roff := dataOffset + uint64(r)*recordSize
			Key = getLeafKey(st, roff)
		}
		i = (i + 1) % nNodes
	}
}
