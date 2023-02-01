// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/db19"
	. "github.com/apmckinlay/gsuneido/dbms/query"
	. "github.com/apmckinlay/gsuneido/runtime"
)

func TestQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	// Global.TestDef("Rule_c",
	// 	compile.Constant("function() { return .b }"))
	db, err := db19.OpenDatabaseRead("../suneido.db")
	if err != nil {
		panic(err.Error())
	}
	MakeSuTran = func(qt QueryTran) *SuTran {
		return nil
	}
	tran := db.NewReadTran()
	s := `(((((cus extend r0) join ivc) union (cus leftjoin ivc)) join (bln where ik is "67")) union ((cus join (ivc rename i1 to r2)) join aln)) sort c4,b4,b3`
	q := ParseQuery(s, tran, nil)
	// trace.QueryOpt.Set()
	// trace.JoinOpt.Set()
	q, _, _ = Setup(q, ReadMode, tran)

	fmt.Println("----------------")
	fmt.Println(Format(q))
	th := &Thread{}
	fmt.Println(q.Get(th, Next))
}

func TestQuery2(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	db, err := db19.OpenDatabaseRead("../suneido.db")
	if err != nil {
		panic(err.Error())
	}
	MakeSuTran = func(qt QueryTran) *SuTran {
		return nil
	}
	tran := db.NewReadTran()
	s := `aln where ik is "67" sort ik`
	q := ParseQuery(s, tran, nil)
	// trace.QueryOpt.Set()
	q = q.Transform()
	const frac = 100
	var index []string
	fixcost, varcost := Optimize(q, ReadMode, index, frac)
	if fixcost+varcost >= 9999999999 {
		panic("invalid query: " + q.String())
	}
	q = SetApproach(q, index, frac, tran)

	fmt.Println("----------------")
	fmt.Println(Format(q))
	th := &Thread{}
	fmt.Println(q.Get(th, Next))
}
