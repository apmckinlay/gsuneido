// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/core/trace"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"

	. "github.com/apmckinlay/gsuneido/core"
)

func TestQueryBug(t *testing.T) {
	assert.TestOnlyIndividually(t)
	Libload = libload // dependency injection
	dbmsLocal := openDbms()
	defer dbmsLocal.Close()
	th := NewThread(dbmsLocal, AllPerms)

	query := `ap_checklines summarize apchk_num, apivc_invoice, 
		apchklin_amount_paid = total apchklin_amount_paid`
	args := SuObjectOf(SuStr(query))
	dbmsLocal.Get(th, args, Any)
}

func TestFastGet(t *testing.T) {
	assert.TestOnlyIndividually(t)
	Libload = libload // dependency injection
	dbmsLocal := openDbms()
	defer dbmsLocal.Close()
	th := NewThread(dbmsLocal, AllPerms)

	compile.EvalString(th, `try Database('drop tmp')
		Database('create tmp (a, b, c, d, e) key(a) index(b) index(c) index(d)')
		for i in ..1000
			QueryOutput(#tmp, [a: i, b: i % 64, c: i % 16, d: i % 4])`)

	test := func(s string, expected bool) {
		t.Helper()
		fmt.Println(s)
		result := compile.EvalString(th, s)
		fmt.Println("=>", result != False)
		assert.T(t).Msg(s).This(result != False).Is(expected)
	}
	test2 := func(s string, expected bool) {
		t.Helper()
		for _, which := range []string{"Query1", "not QueryEmpty?"} {
			test(which+"("+s+")", expected)
		}
	}
	trace.QueryOpt.Set()
	test2("#company", true)                                      // no filter
	test2("#company, company_state_prov: 'ON'", true)            // empty key
	test2("#company, company_state_prov: 'X'", false)            // empty key
	test2("#taxes, tax_code: 'PST'", true)                       // just index
	test2("#taxes, tax_code: 'X'", false)                        // just index
	test("QueryEmpty?(#taxes)", false)                           // no filter
	test2("#stdlib, num: 2", true)                               // key
	test2("#stdlib, num: 2, name: 'X'", false)                   // key + filter
	test2("#stdlib, num: 2, name: 'Beep'", true)                 // key + filter
	test("QueryEmpty?(#stdlib)", false)                          // no filter
	test("QueryEmpty?(#stdlib, name: 'Alert')", false)           // just index
	test("QueryEmpty?(#stdlib, name: 'Alert', text: 'X')", true) // only index
	test("QueryEmpty?(#tmp)", false)                             // no filter
	test("QueryEmpty?(#tmp, b: 59)", false)                      // just index
	test("QueryEmpty?(#tmp, a: 123)", false)                     // key
	test("QueryEmpty?(#tmp, b: 59, c: 11, d: 3)", false)         // multi = b
	test("QueryEmpty?(#tmp, b: 59, c: 11, d: 9999)", true)       // multi = d
}

func BenchmarkSlow(b *testing.B) {
	dbmsLocal := openDbms()
	defer dbmsLocal.Close()
	th := NewThread(dbmsLocal, AllPerms)
	qry.MakeSuTran = func(qry.QueryTran) *SuTran {
		return nil
	}
	args := &SuObject{}
	args.Add(SuStr("stdlib where num = 2"))
	for b.Loop() {
		dbmsLocal.Get(th, args, Only)
	}
}

func BenchmarkFast(b *testing.B) {
	dbmsLocal := openDbms()
	defer dbmsLocal.Close()
	th := NewThread(dbmsLocal, AllPerms)
	qry.MakeSuTran = func(qry.QueryTran) *SuTran {
		return nil
	}
	args := &SuObject{}
	args.Add(SuStr("stdlib"))
	args.Set(SuStr("num"), SuInt16(2))
	for b.Loop() {
		dbmsLocal.Get(th, args, Only)
	}
}
