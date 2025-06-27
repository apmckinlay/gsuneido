// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"cmp"
	"fmt"
	"maps"
	"slices"

	op "github.com/apmckinlay/gsuneido/core/opcodes"
	"github.com/apmckinlay/gsuneido/util/exit"
)

var _ = exit.Add("opseq stats", printOpSeqStats)

func printOpSeqStats() {
	fmt.Println("total", opSeqCount)
	
	// sort by count in descending order
	sortedSequences := slices.SortedFunc(maps.Keys(opSequences),
		func(a, b [3]op.Opcode) int {
			return cmp.Compare(opSequences[b], opSequences[a])
		})

	limit := min(len(sortedSequences), 20)
	for i := range limit {
		seq := sortedSequences[i]
		count := opSequences[seq]
		fmt.Printf("%2.1f  %-14s %-14s %-14s\n",
			float64((count * 1000) / opSeqCount) / 10, // percentage
			seq[0].String(),
			seq[1].String(),
			seq[2].String())
	}
}
