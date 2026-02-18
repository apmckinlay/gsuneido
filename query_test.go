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
// 	*/

// 	query := `	((cus
// 				extend a1 = c3)
// 			join by(ck,a1)
// 					(ivc
// 				join by(ik)
// 					aln))
// 		union /*NOT DISJOINT*/
// 					((ivc
// 				join /*MANY TO MANY*/ by(ik)
// 					aln)
// 			leftjoin /*MANY TO MANY*/ by(ck,a1)
// 					((cus
// 				union /*NOT DISJOINT*/
// 					(cus
// 					extend a1 = c4))
// 				where ck <= ""))`
// 	th := &Thread{}
// 	builtin.QueryHash(th, []Value{SuStr(query), False})
// }
