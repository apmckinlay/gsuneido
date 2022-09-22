// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/apmckinlay/gsuneido/db19"
// 	"github.com/apmckinlay/gsuneido/dbms/query"
// 	. "github.com/apmckinlay/gsuneido/runtime"
// )

// func TestQuery(t *testing.T) {
// 	db, err := db19.OpenDatabaseRead("../suneido.db")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	query.MakeSuTran = func(qt query.QueryTran) *SuTran {
// 		return nil
// 	}
// 	tran := db.NewReadTran()
// 	s :=
// 		`(((gl_accounts extend gldept_id = "")
// 		join  by(glacct_num, gldept_id)
// 		(gl_transactions where gltran_date >= #20220526 and gltran_date <= #20220526))
// 			extend depts = false)
// 		union
// 		(((gl_accounts
// 			where glacct_type in ("Revenue", "Expense"))
// 		times gl_departments)
// 		join by(glacct_num, gldept_id)
// 		(gl_transactions
// 			where gltran_date >= #20220526 and gltran_date <= #20220526))
// 		sort glacct_type_order, glacct_abbrev, glacct_name, gldept_id, gltran_month, gltran_date`
// 	q := query.ParseQuery(s, tran, nil)
// 	q, _ = query.Setup(q, query.ReadMode, tran)
// 	hdr := q.Header()
// 	for i, fs := range hdr.Fields {
// 		fmt.Println(i, fs)
// 	}
// 	th := &Thread{}
// 	for {
// 		row := q.Get(th, Next)
// 		if row == nil {
// 			break
// 		}
// 		a := row.GetVal(hdr, "glacct_abbrev", th, nil)
// 		if a.Equal(SuStr("4050")) {
// 			d := row.GetVal(hdr, "gldept_id", th, nil)
// 			s := ""
// 			for _, r := range row {
// 				if r.Record == "" {
// 					s += "nil "
// 				} else {
// 					s += "rec "
// 				}
// 			}
// 			fmt.Println(a, s, d)
// 		}
// 	}
// }
