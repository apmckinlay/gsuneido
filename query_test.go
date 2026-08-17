// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"testing"

	"github.com/apmckinlay/gsuneido/builtin"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// 	/* Schemas:
// 	cus
// 		(c1, c2, c3, c4, ck)
// 		key (ck)

// 	ivc
// 		(ck, i1, i2, i3, i4, ik)
// 		key (ik)
// 		index (ck) in cus

// 	aln
// 		(a1, a2, a3, a4, ak, ik)
// 		key (ik,ak)

// 	bln
// 	    (b1, b2, b3, b4, bk, ik)
//     	key (ik,bk)
// 	*/

func TestFuzzBug(t *testing.T) {
	assert.TestOnlyIndividually(t)
	openDbms()
	defer db.CloseKeepMapped()

	// seed: 7948325488
	query := `((cus extend ik = c3) join ivc) join (aln union (bln union aln))`
	th := &Thread{}
	x := builtin.QueryHash(th, []Value{SuStr(query), True})
	y := builtin.QueryAltHash(th, []Value{SuStr(query), True})
	assert.This(x).Is(y)
}
