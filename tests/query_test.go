// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package language

// import (
// 	"testing"

// 	"github.com/apmckinlay/gsuneido/db19"
// 	"github.com/apmckinlay/gsuneido/dbms/query"
// 	"github.com/apmckinlay/gsuneido/runtime"
// )

// func TestQuery(t *testing.T) {
// 	db, err := db19.OpenDatabaseRead("../suneido.db")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	query.MakeSuTran = func(qt query.QueryTran) *runtime.SuTran {
// 		return nil
// 	}
// 	tran := db.NewReadTran()
// 	s :=
// 		`(gl_accounts
//             where glacct_num = "#20220511.141319385_inventory"
//             extend gldept_id = "")
// 	join by(glacct_num, gldept_id)
// 		(gl_transactions
//             where gltran_reference = "Physical Count"
//             where gltran_desc = "02-2737"
//             where gltran_date >= #20220501 and gltran_date <= #20220511)`
// 	q := query.ParseQuery(s, tran, nil)
// 	q, _ = query.Setup(q, query.ReadMode, tran)
// 	if row := q.Get(&runtime.Thread{}, runtime.Next); row == nil {
// 		t.Fail()
// 	}
// }
