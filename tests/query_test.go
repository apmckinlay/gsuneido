// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/apmckinlay/gsuneido/db19"
// 	"github.com/apmckinlay/gsuneido/dbms/query"
// 	. "github.com/apmckinlay/gsuneido/runtime"
// 	"github.com/apmckinlay/gsuneido/runtime/trace"
// )

// func TestQuery(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode")
// 	}
// 	// Global.TestDef("Rule_c",
// 	// 	compile.Constant("function() { return .b }"))
// 	db, err := db19.OpenDatabaseRead("../suneido.db")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	query.MakeSuTran = func(qt query.QueryTran) *SuTran {
// 		return nil
// 	}
// 	tran := db.NewReadTran()
// 	s := `(cus join ((((ivc join aln) union (ivc join bln)) union ((ivc join bln) union (ivc join bln))) union ((ivc join aln) extend x8763182)))`
// 	// s := `bln where ik = 86 and b2 is 45`
// 	q := query.ParseQuery(s, tran, nil)
// 	trace.QueryOpt.Set()
// 	trace.JoinOpt.Set()
// 	q, _, _ = query.Setup(q, query.ReadMode, tran)
// 	fmt.Println("----------------")
// 	fmt.Println(query.Format(q))
// 	th := &Thread{}
// 	q.Get(th, Next)
// }
