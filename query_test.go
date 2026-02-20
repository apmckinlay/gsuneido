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
// 	*/

// 	query := `	(((cus
// 			union /*NOT DISJOINT*/
// 				cus)
// 			union /*NOT DISJOINT*/
// 				cus)
// 		join /*MANY TO MANY*/ by(ck)
// 			(ivc
// 			extend r0))
// 		where r0 is '3'
// 		where ik is ""`
// 	th := &Thread{}
// 	builtin.QueryHash(th, []Value{SuStr(query), False})
// }
