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

func TestFastGet(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	Libload = libload // dependency injection
	mainThread.Name = "main"
	mainThread.SetSviews(&sviews)
	MainThread = &mainThread
	openDbms()
	defer db.CloseKeepMapped()

	compile.EvalString(MainThread, `Database('drop tmp')
		Database('create tmp (a, b, c, d, e) key(a) index(b) index(c) index(d)')
		for i in ..1000
			QueryOutput(#tmp, [a: i, b: i % 64, c: i % 16, d: i % 4])`)

	test := func(s string, expected bool) {
		t.Helper()
		fmt.Println(s)
		result := compile.EvalString(MainThread, s)
		fmt.Println("=>", result != False)
		assert.T(t).Msg(s).This(result != False).Is(expected)
	}
	test2 := func(s string, expected bool) {
		t.Helper()
		for _, which := range []string{"Query1", "QueryExists?"} {
			test(which+"("+s+")", expected)
		}
	}
	trace.QueryOpt.Set()
	test2("#company", true)                                        // no filter
	test2("#company, company_state_prov: 'ON'", true)              // empty key
	test2("#company, company_state_prov: 'X'", false)              // empty key
	test2("#taxes, tax_code: 'PST'", true)                         // just index
	test2("#taxes, tax_code: 'X'", false)                          // just index
	test("QueryExists?(#taxes)", true)                             // no filter
	test2("#stdlib, num: 2", true)                                 // key
	test2("#stdlib, num: 2, name: 'X'", false)                     // key + filter
	test2("#stdlib, num: 2, name: 'Beep'", true)                   // key + filter
	test("QueryExists?(#stdlib)", true)                            // no filter
	test("QueryExists?(#stdlib, name: 'Alert')", true)             // just index
	test("QueryExists?(#stdlib, name: 'Alert', text: 'X')", false) // only index

	test("QueryExists?(#tmp)", true)                         // no filter
	test("QueryExists?(#tmp, b: 59)", true)                  // just index
	test("QueryExists?(#tmp, a: 123)", true)                 // key
	test("QueryExists?(#tmp, b: 59, c: 11, d: 3)", true)     // multi = b
	test("QueryExists?(#tmp, b: 59, c: 11, d: 9999)", false) // multi = d
}

func BenchmarkSlow(b *testing.B) {
	openDbms()
	defer db.CloseKeepMapped()
	MainThread = &mainThread
	qry.MakeSuTran = func(qry.QueryTran) *SuTran {
		return nil
	}
	args := &SuObject{}
	args.Add(SuStr("stdlib where num = 2"))
	for b.Loop() {
		dbmsLocal.Get(MainThread, args, Only)
	}
}

func BenchmarkFast(b *testing.B) {
	openDbms()
	defer db.CloseKeepMapped()
	MainThread = &mainThread
	qry.MakeSuTran = func(qry.QueryTran) *SuTran {
		return nil
	}
	args := &SuObject{}
	args.Add(SuStr("stdlib"))
	args.Set(SuStr("num"), SuInt(2))
	for b.Loop() {
		dbmsLocal.Get(MainThread, args, Only)
	}
}
