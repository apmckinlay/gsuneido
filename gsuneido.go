// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/tools"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/exit"
	"github.com/apmckinlay/gsuneido/util/regex"
	"github.com/apmckinlay/gsuneido/util/system"
)

var builtDate = "Dec 29 2020 12:34" // set by: go build -ldflags "-X main.builtDate=..."
var mode = ""                       // set by: go build -ldflags "-X main.mode=gui"

var help = `options:
	-check
	-c[lient] [ipaddress] (default 127.0.0.1)
	-d[ump] [table]
	-h[elp] or -?
	-l[oad] [table]
	-n[o]r[elaunch]
	-p[ort] # (default 3147)
	-repair
	-r[epl]
	-s[erver]
	-u[nattended]
	-v[ersion]`

// dbmsLocal is set if running with a local/standalone database.
var dbmsLocal *dbms.DbmsLocal
var mainThread *Thread

func main() {
	options.BuiltDate = builtDate
	options.Mode = mode
	options.Parse(os.Args[1:])
	if options.Action == "client" {
		options.Errlog = builtin.ErrlogDir() + "suneido" + options.Port + ".err"
	}
	if mode == "gui" {
		redirect()
	}
	if err := system.Service("gSuneido", redirect, exit.RunFuncs); err != nil {
		Fatal(err)
	}
	if options.Action == "" && mode != "gui" {
		options.Action = "repl"
	}

	switch options.Action {
	case "":
		break
	case "server":
		if mode == "gui" {
			Fatal("Please use gsport for server mode")
		}
		startServer()
		os.Exit(0)
	case "dump":
		t := time.Now()
		if options.Arg == "" {
			nTables, nViews, err := tools.DumpDatabase("suneido.db", "database.su")
			ck(err)
			Alert("dumped", nTables, "tables", nViews, "views in",
				time.Since(t).Round(time.Millisecond))
		} else {
			table := strings.TrimSuffix(options.Arg, ".su")
			nrecs, err := tools.DumpTable("suneido.db", table, table+".su")
			ck(err)
			Alert("dumped", nrecs, "records from", table,
				"in", time.Since(t).Round(time.Millisecond))
		}
		os.Exit(0)
	case "load":
		t := time.Now()
		if options.Arg == "" {
			nTables, nViews, err := tools.LoadDatabase("database.su", "suneido.db")
			ck(err)
			Alert("loaded", nTables, "tables", nViews, "views in",
				time.Since(t).Round(time.Millisecond))
		} else {
			table := strings.TrimSuffix(options.Arg, ".su")
			n, err := tools.LoadTable(table, "suneido.db")
			ck(err)
			Alert("loaded", n, "records to", table,
				"in", time.Since(t).Round(time.Millisecond))
		}
		os.Exit(0)
	case "compact":
		t := time.Now()
		nTables, nViews, err := tools.Compact("suneido.db")
		ck(err)
		Alert("compacted", nTables, "tables", nViews, "views in",
			time.Since(t).Round(time.Millisecond))
		os.Exit(0)
	case "check":
		t := time.Now()
		ck(db19.CheckDatabase("suneido.db"))
		Alert("checked database in", time.Since(t).Round(time.Millisecond))
		os.Exit(0)
	case "repair":
		t := time.Now()
		err := db19.CheckDatabase("suneido.db")
		if err == nil {
			Alert("database ok")
		} else {
			msg, err := db19.Repair("suneido.db", err)
			ck(err)
			Alert(msg,
				"\nrepaired database in", time.Since(t).Round(time.Millisecond))
		}
		os.Exit(0)
	case "version":
		Alert("gSuneido " + builtDate + " (" + runtime.Version() + " " +
			runtime.GOARCH + " " + runtime.GOOS + ")")
		os.Exit(0)
	case "help":
		Alert(help)
		os.Exit(0)
	case "repl", "client":
		// handled below
	case "error":
		Fatal(options.Error)
	default:
		Alert("invalid action:", options.Action)
		os.Exit(1)
	}
	Libload = libload // dependency injection
	mainThread = &Thread{}
	mainThread.Name = "main"
	mainThread.UIThread = true
	MainThread = mainThread
	builtin.UIThread = mainThread
	exit.Add(func() { mainThread.Close() })
	defer func() {
		if e := recover(); e != nil {
			log.Println("ERROR:", e, "(exiting)")
			exit.Exit(1)
		}
		exit.Exit(0)
	}()
	// dependency injection of GetDbms
	if options.Action == "client" {
		conn, jserver := dbms.ConnectClient(options.Arg, options.Port)
		if jserver {
			mainThread.SetDbms(dbms.NewDbmsClient(conn))
			GetDbms = func() IDbms {
				conn, _ := dbms.ConnectClient(options.Arg, options.Port)
				return dbms.NewDbmsClient(conn)
			}
		} else {
			client := dbms.NewMuxClient(conn)
			GetDbms = func() IDbms {
				return client.NewSession()
			}
		}
		clientErrorLog()
	} else {
		openDbms()
	}
	if options.Action == "repl" ||
		(options.Action == "client" && options.Mode != "gui") {
		run("Init.Repl()")
		repl()
	} else {
		run("Init()")
		builtin.Run()
	}
}

func redirect() {
	if err := system.Redirect(options.Errlog); err != nil {
		Fatal("Redirect failed:", err)
	}
}

func run(src string) {
	defer func() {
		if e := recover(); e != nil {
			LogUncaught(mainThread, src, e)
			Fatal("ERROR from", src, e)
		}
	}()
	compile.EvalString(mainThread, src)
}

func ck(err error) {
	if err != nil {
		Fatal(err)
	}
}

// clientErrorLog sends the client's error log to the server.
// This is to record errors that occurred on the client
// when the server was not connected.
func clientErrorLog() {
	dbms := mainThread.Dbms()

	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)
	log.SetPrefix(mainThread.SessionId("") + " ")

	f, err := os.Open(options.Errlog)
	if err != nil {
		return
	}
	defer func() {
		f.Close()
		os.Truncate(options.Errlog, 0) // can't remove since open as stderr
		if e := recover(); e != nil {
			dbms.Log("log previous errors: " + fmt.Sprint(e))
		}
	}()
	// send errors to server
	in := bufio.NewScanner(f)
	in.Buffer(nil, 1024)
	nlines := 0
	for in.Scan() {
		dbms.Log("PREVIOUS: " + in.Text())
		if nlines++; nlines > 1000 {
			dbms.Log("PREVIOUS: too many errors")
			break
		}
	}
}

// startServer does not return
func startServer() {
	log.Println("starting server")
	startHttpStatus()
	openDbms()
	Libload = libload // dependency injection
	mainThread = &Thread{}
	mainThread.Name = "main"
	run("Init()")
	options.DbStatus.Store("")
	exit.Add(stopServer)
	dbms.Server(dbmsLocal)
}

func stopServer() {
	httpServer.Close()
	dbms.StopServer()
	log.Println("server stopped")
}

var db *db19.Database

func openDbms() {
	var err error
	db, err = db19.OpenDatabase("suneido.db")
	if errors.Is(err, fs.ErrNotExist) || errors.Is(err, fs.ErrPermission) {
		Fatal(err)
	}
	if err != nil {
		if !AlertCancel("ERROR:", err, "\nwill try to repair") {
			Fatal("database corrupt, not repaired")
		}
		options.DbStatus.Store("repairing")
		msg, err := db19.Repair("suneido.db", err)
		if err != nil {
			Fatal("repair:", err)
		}
		Alert(msg)
		db, err = db19.OpenDatabase("suneido.db")
		if err != nil {
			Fatal("open:", err)
		}
		options.DbStatus.Store("starting")
	}
	db19.StartTimestamps()
	db19.StartConcur(db, persistInterval())
	dbmsLocal = dbms.NewDbmsLocal(db)
	DbmsAuth = options.Action == "server" || mode != "gui" || !db.HaveUsers()
	GetDbms = getDbms
	exit.Add(db.CloseKeepMapped) // keep mapped to avoid errors during shutdown
	// go checkState()
}

func persistInterval() time.Duration {
	if options.Action == "server" {
		return 30 * time.Second // FIXME should probably be bigger
	}
	// else standalone
	return 10 * time.Second
}

func getDbms() IDbms {
	if DbmsAuth {
		return dbmsLocal
	}
	return dbms.Unauth(dbmsLocal)
}

// func checkState() {
// 	for {
// 		state := db.GetState()
// 		cksum := state.Meta.Cksum()
// 		// read meta to verify checksums
// 		schemaOff, infoOff := state.Meta.Offsets()
// 		if schemaOff != 0 {
// 			hamt.ReadChain[string](db.Store, schemaOff, meta.ReadSchema)
// 		}
// 		if infoOff != 0 {
// 			hamt.ReadChain[string](db.Store, infoOff, meta.ReadInfo)
// 		}
// 		time.Sleep(50 * time.Millisecond) // adjust for overhead
// 		// recalculate checksum to verify Meta hasn't been mutated
// 		assert.That(state.Meta.Cksum() == cksum)
// 	}
// }

// REPL -------------------------------------------------------------

var prompt = func(s string) { fmt.Fprintln(os.Stderr, s) }

func repl() {
	builtin.InheritHandles = true
	log.SetFlags(log.Ltime)
	log.SetPrefix("")

	if !isTerminal(os.Stdin) || !isTerminal(os.Stdout) {
		prompt = func(string) {}
	}

	builtin.DefDef()
	builtin.DefConcat()

	built := builtin.BuiltStr()
	if options.Action == "client" {
		built += " - client"
	}
	prompt(built)
	showOptions()
	prompt("Press Enter twice (i.e. blank line) to execute, q to quit")
	r := bufio.NewReader(os.Stdin)
	for {
		prompt("~~~")
		src := ""
		for {
			line, err := r.ReadString('\n')
			line = strings.TrimRight(line, " \t\r\n")
			if line == "q" || (err != nil && (err != io.EOF || src == "")) {
				return
			}
			if line == "" {
				break
			}
			src += line + "\n"
		}
		eval(src)
	}
}

func isTerminal(f *os.File) bool {
	fm, err := f.Stat()
	if err != nil {
		return true // ???
	}
	mode := fm.Mode()
	return mode&os.ModeDevice == os.ModeDevice &&
		mode&os.ModeCharDevice == os.ModeCharDevice
}

func showOptions() {
	if options.ThreadDisabled {
		prompt("- Thread disabled")
	}
	if options.TimersDisabled {
		prompt("- Timers disabled")
	}
	if options.ClearCallbackDisabled {
		prompt("- ClearCallback disabled")
	}
}

func eval(src string) {
	defer func() {
		if e := recover(); e != nil {
			LogUncaught(mainThread, "repl", e)
		}
	}()
	src = "function () {\n" + src + "\n}"
	v, results := compile.Checked(mainThread, src)
	for _, s := range results {
		fmt.Println("(" + s + ")")
	}
	fn := v.(*SuFunc)
	// DisasmMixed(os.Stdout, fn, src)

	mainThread.Reset()
	result := mainThread.Call(fn)
	if result != nil {
		fmt.Println(WithType(result)) // NOTE: doesn't use ToString
	}
}

//-------------------------------------------------------------------

// libload loads a name from the dbms
func libload(t *Thread, name string) (result Value, e any) {
	defer func() {
		if e = recover(); e != nil {
			// fmt.Println("INFO: error loading", name, e)
			// dbg.PrintStack()
			result = nil
		}
	}()
	defs := t.Dbms().LibGet(name)
	if len(defs) == 0 {
		// fmt.Println("LOAD", name, "MISSING")
		return nil, nil
	}
	for i := 0; i < len(defs); i += 2 {
		lib := defs[i]
		src := defs[i+1]
		if s, ok := LibraryOverrides.Get(lib, name); ok {
			src = s
		}
		if mode == "gui" && strings.HasSuffix(lib, "webgui") {
			continue
		}
		result = llcompile(lib, name, src, result)
		// fmt.Println("LOAD", name, "SUCCEEDED")
	}
	return result, nil
}

var winErr = regex.Compile("gSuneido does not implement (dll|struct|callback)")

func llcompile(lib, name, src string, prevDef Value) Value {
	defer func() {
		if e := recover(); e != nil {
			es := fmt.Sprint(e)
			if !winErr.Matches(es) {
				panic(e)
			}
		}
	}()
	// want to pass the name from the start (rather than adding after)
	// so it propagates to nested Named values
	return compile.NamedConstant(lib, name, src, prevDef)
}
