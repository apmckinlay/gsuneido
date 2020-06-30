// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"math/rand"
	"testing"

	"github.com/apmckinlay/gsuneido/util/verify"
)

func TestCheckerStartStop(*testing.T) {
	ck := NewCheck()
	const ntrans = 20
	var trans [ntrans]int
	const ntimes = 5000
	for i := 0; i < ntimes; i++ {
		j := rand.Intn(ntrans)
		if trans[j] == 0 {
			trans[j] = ck.StartTran()
			// fmt.Println("start", trans[j])
		} else {
			if rand.Intn(2) == 1 {
				// fmt.Println("commit", trans[j])
				ck.Commit(trans[j])
			} else {
				// fmt.Println("abort", trans[j])
				ck.Abort(trans[j])
			}
			trans[j] = 0
			// fmt.Println("#", len(ck.trans))
		}
	}
	// fmt.Println("#", len(ck.trans))
	for _, tn := range trans {
		if tn != 0 {
			ck.Commit(tn)
		}
	}
	verify.That(len(ck.trans) == 0)
}
