// Governed by the MIT license found in the LICENSE file.

package main

import (
	"testing"

	qry "github.com/apmckinlay/gsuneido/dbms/query"

	. "github.com/apmckinlay/gsuneido/core"
)

// func TestFull(t *testing.T) {
// 	options.BuiltDate = builtDate
// 	Libload = libload // dependency injection
// 	mainThread.Name = "main"
// 	mainThread.SetSviews(&sviews)
// 	MainThread = &mainThread
// 	openDbms()
// 	defer db.CloseKeepMapped()
// 	run("Init.Repl()")
// 	run(`x = QueryAll('taxes')[..1]
// 		Print(Unpack(Zlib.Uncompress(Pack(x, zip:))))`)
// }

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
