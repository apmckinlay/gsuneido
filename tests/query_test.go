// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

// import (
// 	"testing"

// 	"github.com/apmckinlay/gsuneido/compile"
// 	"github.com/apmckinlay/gsuneido/db19"
// 	"github.com/apmckinlay/gsuneido/dbms/query"
// 	. "github.com/apmckinlay/gsuneido/runtime"
// )

// func TestQuery(t *testing.T) {
// 	Libload = func(t *Thread, name string) (result Value, e any) {
// 		if name != "Rule_rule" {
// 			return nil, nil
// 		}
// 		src := `function ()
// 			{
// 			this['foo' $ Random(10000)]
// 			return .num
// 			}`
// 		return compile.NamedConstant("stdlib", name, src, nil), nil
// 	}
// 	db, err := db19.OpenDatabaseRead("../suneido.db")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	query.MakeSuTran = func(qt query.QueryTran) *SuTran {
// 		return nil
// 	}
// 	tran := db.NewReadTran()
// 	s :=
// 		`etalib
// 		extend rule
// 		where rule isnt 123
// 		sort rule`
// 	for i := 0; i < 16; i++ {
// 		q := query.ParseQuery(s, tran, nil)
// 		q, _ = query.Setup(q, query.ReadMode, tran)
// 		th := &Thread{}
// 		for {
// 			row := q.Get(th, Next)
// 			if row == nil {
// 				break
// 			}
// 		}
// 	}
// }
