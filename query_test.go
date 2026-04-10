// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

// import (
// 	"testing"

// 	"github.com/apmckinlay/gsuneido/builtin"

// 	. "github.com/apmckinlay/gsuneido/core"
// )

// func TestFuzzBug(t *testing.T) {
// 	openDbms()
// 	defer db.CloseKeepMapped()

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

// 	query := `(((((((ivc leftjoin bln) union (ivc join (aln union bln))) extend x0 = i1)) where ik is "")) where bk is "")`
// 	th := &Thread{}
// 	builtin.QueryHash(th, []Value{SuStr(query), False})
// }
