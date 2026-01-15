// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/core/trace"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"

	. "github.com/apmckinlay/gsuneido/core"
)

func TestQueryBug(t *testing.T) {
	TestOnlyIndividually(t)
	Libload = libload // dependency injection
	mainThread.Name = "main"
	mainThread.SetSviews(&sviews)
	MainThread = &mainThread
	openDbms()
	defer db.CloseKeepMapped()

	query := `
		(ap_checklines where apivc_invoice is "12567" 
		join by(apchk_num) 
		(ap_checks where bizpartner_num is #20260113.104240513110 
			and apchk_void? isnt true 
			and apchk_type is "CK" 
		rename 
			etaequipstmt_name to etaequipstmt_name_hdr, 
			etaequipstmt_date to etaequipstmt_date_hdr, 
			etaequip_num to etaequip_num_hdr
		)) 
		summarize max apchk_date, total apchklin_amount_paid`
	args := SuObjectOf(SuStr(query))
	dbmsLocal.Get(MainThread, args, Only)
}

func TestFastGet(t *testing.T) {
	TestOnlyIndividually(t)
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

func TestOnlyIndividually(t *testing.T) {
	for _, arg := range os.Args {
		if strings.Contains(arg, "-test.run") &&
			!strings.Contains(arg, "|") {
			return
		}
	}
	t.SkipNow()
}
