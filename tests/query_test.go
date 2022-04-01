// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package language

// func TestQuery(*testing.T) {
// 	db, err := db19.OpenDatabaseRead("../suneido.db")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	tran := db.NewReadTran()
// 	s := `(eta_equip_stmt_history
// 				project etaequipstmt_num, etaequip_num)
// 		leftjoin by(etaequipstmt_num)
// 			(eta_equip_orders
// 				rename etaequip_num to etaeo_equip_num
// 				project etaequipstmt_num, etaeo_equip_num)`
// 	q := query.ParseQuery(s, tran, nil)
// 	q, _ = query.Setup(q, query.ReadMode, tran)
// 	fmt.Println(q)
// 	q.Get(nil, runtime.Next)
// }
