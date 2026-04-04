// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestDatabase_Top10(t *testing.T) {
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()

	qry.DoAdmin(db, "create tmp (id, c) key(id)", nil)
	ut := db.NewUpdateTran()
	id := 1
	for range 25 {
		qry.DoAction(nil, ut,
			fmt.Sprintf("insert { id: %d, c: 'hot' } into tmp", id))
		id++
	}
	for range 15 {
		qry.DoAction(nil, ut,
			fmt.Sprintf("insert { id: %d, c: 'warm' } into tmp", id))
		id++
	}
	for i := range 30 {
		qry.DoAction(nil, ut,
			fmt.Sprintf("insert { id: %d, c: 'u%02d' } into tmp", id, i))
		id++
	}
	ut.Commit()

	th := &Thread{}
	th.SetDbms(dbms.NewDbmsLocal(db))
	v := db_Top10(th, []Value{SuStr("tmp"), SuStr("c")})
	ob := v.(*SuObject)

	assert.T(t).This(ob.Size()).Is(10)
	assert.T(t).This(ToInt(ob.Get(nil, SuStr("hot")))).Is(25)
	assert.T(t).This(ToInt(ob.Get(nil, SuStr("warm")))).Is(15)
}

func TestDatabase_Distinct(t *testing.T) {
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()

	qry.DoAdmin(db, "create tmp (a, b, c) key(a) index(b,c)", nil)
	ut := db.NewUpdateTran()
	qry.DoAction(nil, ut, "insert { a: 1, b: 'x', c: 10 } into tmp")
	qry.DoAction(nil, ut, "insert { a: 2, b: 'x', c: 20 } into tmp")
	qry.DoAction(nil, ut, "insert { a: 3, b: 'y', c: 20 } into tmp")
	ut.Commit()

	th := &Thread{}
	th.SetDbms(dbms.NewDbmsLocal(db))
	v := db_Distinct(th, []Value{SuStr("tmp")})
	ob := v.(*SuObject)

	assert.T(t).This(ToInt(ob.Get(nil, SuStr("a")))).Is(3)
	assert.T(t).This(ToInt(ob.Get(nil, SuStr("b")))).Is(2)
	assert.T(t).This(ToInt(ob.Get(nil, SuStr("c")))).Is(2)
}
