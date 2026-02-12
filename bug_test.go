package main

// import (
// 	"testing"

// 	"github.com/apmckinlay/gsuneido/builtin"
// 	. "github.com/apmckinlay/gsuneido/core"
// 	"github.com/apmckinlay/gsuneido/util/exit"
// )

// func TestBug(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode")
// 	}
// 	Libload = libload // dependency injection
// 	mainThread = &Thread{}
// 	mainThread.Name = "main"
// 	mainThread.UIThread = true
// 	MainThread = mainThread
// 	builtin.UIThread = mainThread
// 	exit.Add(func() { mainThread.Close() })

// 	openDbms()
// 	defer db.CloseKeepMapped()

// 	run(`
// 		Init.Repl()
// 		//Use("axonlib")
// 		//Use("Accountinglib")
// 		//Use("etalib")
// 		//Use("pcmiler")
// 		//Use("ticketlib")
// 		//Use("prlib")
// 		//Use("prcadlib")
// 		//Use("etaprlib")
// 		//Use("invenlib")
// 		//Use("wolib")
// 		//Use("polib")
// 		//Use("configlib")
// 		//Use("demobookoptions")
// 		//Use("Test_lib")
// 		Print("running...")
// 		Timer(secs: 60) {
// 			Qfuzz.MakeQuery()
// 		}
// 		`)
// }
