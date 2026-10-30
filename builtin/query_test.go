// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestQueryWhere(t *testing.T) {
	names := []Value{SuStr("master_key"), SuStr("flag")}
	as := &ArgSpec{Nargs: 2, Spec: []byte{0, 1}, Names: names}
	args := []Value{DateFromLiteral("19010101"), False}
	assert.T(t).This(queryWhere(as, args)).
		Is("\nwhere master_key is #19010101\nand flag is false")
}

// permsTestSetup sets up a test database, tables, and permissions.
func permsTestSetup(t *testing.T) *Thread {
	t.Helper()
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	t.Cleanup(func() { db.Close() })
	d := dbms.NewDbmsLocal(db)
	d.AdminTest("create builtin_perm_public (id) key(id)")
	d.AdminTest("create builtin_perm_private (id) key(id)")
	th := &Thread{}
	th.SetDbms(d)
	th.SetPerms(AllPerms)
	th.Call(compile.Constant(`function () {
		Transaction(update:) {|t|
			t.QueryDo('insert {id: 1} into builtin_perm_public')
			t.QueryDo('insert {id: 42} into builtin_perm_private')
			}
		}`))
	p := &Perms{}
	p.AddTable("builtin_perm_public", "write")
	th.SetPerms(p)
	return th
}

func TestBuiltinRulePerms(t *testing.T) {
	th := permsTestSetup(t)
	name := "Rule_builtin_perm_nested"
	old := Global.GetIfPresent(name)
	Global.TestDef(name, compile.Constant(`function () {
		return Transaction(read:) {|t|
			t.QueryFirst('builtin_perm_private sort id').id
			}
		}`))
	defer Global.TestDef(name, old)
	fn := compile.Constant(`function () {
		return QueryFirst('builtin_perm_public sort id').builtin_perm_nested
		}`)
	assert.T(t).This(func() { th.Call(fn) }).Panics("not authorized: builtin_perm_private")
	th.Reset()
	th.SetPerms(AllPerms)
	assert.T(t).This(th.Call(fn)).Is(IntVal(42))
}

func TestBuiltinTriggerPerms(t *testing.T) {
	th := permsTestSetup(t)
	restricted := th.Perms()
	name := "Trigger_builtin_perm_public"
	old := Global.GetIfPresent(name)
	defer Global.TestDef(name, old)

	fn := compile.Constant(`function () {
		Transaction(update:) {|t|
			t.QueryDo('insert {id: 2} into builtin_perm_public')
			}
		return true
		}`)

	// triggers run with all permissions, so reading a private table is allowed
	// even though the invoking thread only has access to builtin_perm_public
	Global.TestDef(name, compile.Constant(`function (t, oldrec, newrec) {
		Transaction(read:) {|nested|
			nested.QueryFirst('builtin_perm_private sort id')
			}
		QueryFirst('builtin_perm_private sort id')
		}`))
	assert.T(t).This(th.Call(fn)).Is(True)

	// but the invoking thread's permissions are restored afterwards
	assert.T(t).This(func() {
		th.Call(compile.Constant(
			`function () { QueryFirst('builtin_perm_private sort id') }`))
	}).Panics("not authorized: builtin_perm_private")
	assert.That(th.Perms() == restricted)
	th.Reset()
	th.SetPerms(restricted)

	// triggers can't grant permissions, i.e. NewPerms is nil while they run
	// so Auth can't be fooled into escalating permissions from a trigger
	builder := &Perms{}
	th.SetNewPerms(builder)
	Global.TestDef(name, compile.Constant(`function (t, oldrec, newrec) {
		Permit.Table('builtin_perm_escalate', 'write')
		}`))
	fn2 := compile.Constant(`function () {
		Transaction(update:) {|t|
			t.QueryDo('insert {id: 3} into builtin_perm_public')
			}
		}`)
	assert.T(t).This(func() { th.Call(fn2) }).
		Panics("Permit.Table can only be called from Auth")

	// NewPerms is restored so Auth can continue building permissions
	assert.That(th.NewPerms() == builder)
	th.Reset()
}

func TestBuiltinDatabasePerms(t *testing.T) {
	th := permsTestSetup(t)
	restricted := th.Perms()
	test := func(source string) {
		t.Helper()
		fn := compile.Constant("function () { " + source + " }")
		th.SetPerms(restricted)
		assert.T(t).This(func() { th.Call(fn) }).Panics("not authorized: builtin_perm_private")
		th.Reset()
		th.SetPerms(AllPerms)
		th.Call(fn)
	}
	test("Transaction(read:) {|t| t.QueryFirst('builtin_perm_private sort id') }")
	test("QueryFirst('builtin_perm_private sort id')")
	test("QueryAlt('builtin_perm_private')")
	test("QueryAltHash('builtin_perm_private')")
	test("QueryHash('builtin_perm_private')")
	test("Database.Top10('builtin_perm_private', 'id')")
	test("Database.Distinct('builtin_perm_private')")
	test("Cursor('builtin_perm_private') {|c| Transaction(read:) {|t| c.Next(t) } }")

	fn := compile.Constant("function () { Database('create builtin_perm_created (id) key(id)') }")
	th.SetPerms(restricted)
	assert.T(t).This(func() { th.Call(fn) }).Panics("not authorized")
	th.Reset()
	th.SetPerms(AllPerms)
	th.Call(fn)
}
