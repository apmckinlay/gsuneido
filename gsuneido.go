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
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/tools"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/exit"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/regex"
	"github.com/apmckinlay/gsuneido/util/system"
	// sync "github.com/sasha-s/go-deadlock"
)

var builtDate = "Jan 23 2023 12:34" // set by: go build -ldflags "-X main.builtDate=..."
var mode = ""                       // set by: go build -ldflags "-X main.mode=gui"

var help = `options:
	-check
	-c[lient][=ipaddress] (default 127.0.0.1)
	-compact
	-d[ump] [table]
	-h[elp] or -?
	-l[oad] [table]
	-p[ort][=#] (default 3147)
	-repair
	-s[erver]
	-v[ersion]
	-w[eb][=#] (default -port + 1)`

// dbmsLocal is set if running with a local/standalone database.
var dbmsLocal *dbms.DbmsLocal
var mainThread Thread
var sviews Sviews

func main() {
	options.BuiltDate = builtDate
	options.Mode = mode
	options.Parse(os.Args[1:])
	if options.Action == "client" {
		options.Errlog = builtin.ErrlogDir() + "suneido" + options.Port + ".err"
	}
	Exit = exit.Exit
	if mode == "gui" {
		redirect()
	} else {
		svc, err := system.Service("gSuneido", redirect, exit.RunFuncs)
		if err != nil {
			Fatal(err)
		}
		if svc {
			Exit = system.StopService
		}
	}

	Libload = libload // dependency injection
	mainThread.Name = "main"
	mainThread.SetSviews(&sviews)
	MainThread = &mainThread

	switch options.Action {
	case "":
		break
	case "server":
		if mode == "gui" {
			Fatal("Please use gsport for server mode")
		}
		runServer()
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
		nTables, nViews, oldSize, newSize, err := tools.Compact("suneido.db")
		ck(err)
		oldSize /= 1024 * 1024
		newSize /= 1024 * 1024
		Alert("compacted", nTables, "tables", nViews, "views",
			"in", time.Since(t).Round(time.Millisecond),
			oldSize, "-", (oldSize - newSize), "=", newSize, "mb")
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
		Alert("gSuneido " + options.BuiltStr())
		os.Exit(0)
	case "help":
		Alert(help)
		os.Exit(0)
	case "client":
		if options.WebServer {
			options.DbStatus.Store("")
			startHttpStatus()
		}
	case "error":
		Fatal(options.Error)
	default:
		Alert("invalid action:", options.Action)
		os.Exit(1)
	}
	mainThread.UIThread = true
	builtin.UIThread = &mainThread
	defer func() {
		if e := recover(); e != nil {
			log.Println("ERROR:", e, "(exiting)")
			Exit(1)
		}
		Exit(0)
	}()
	// dependency injection of GetDbms
	if options.Action == "client" {
		conn, jserver := dbms.ConnectClient(options.Arg, options.Port)
		if jserver {
			mainThread.SetDbms(dbms.NewJsunClient(conn))
			GetDbms = func() IDbms {
				conn, _ := dbms.ConnectClient(options.Arg, options.Port)
				return dbms.NewJsunClient(conn)
			}
		} else {
			client := dbms.NewDbmsClient(conn)
			GetDbms = func() IDbms {
				return client.NewSession()
			}
		}
		if mode == "gui" {
			clientErrorLog()
		}
	} else {
		openDbms()
		if options.WebServer {
			options.DbStatus.Store("")
			startHttpStatus()
		}
	}
	if mode == "gui" {
		run("Init()")
		builtin.Run()
	} else {
		run("Init.Repl()")
		repl()
	}
}

func redirect() {
	getId := func() string {
		if MainThread == nil {
			return ""
		}
		return MainThread.Session()
	}
	if err := system.Redirect(options.Errlog, getId); err != nil {
		Fatal("Redirect failed:", err)
	}
}

func run(src string) {
	defer func() {
		if e := recover(); e != nil {
			LogUncaught(&mainThread, src, e)
			Fatal("ERROR from", src, e)
		}
	}()
	compile.EvalString(&mainThread, src)
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

	f, err := os.Open(options.Errlog)
	if err != nil {
		return
	}
	defer func() {
		f.Close()
		os.Truncate(options.Errlog, 0) // can't remove since open as stderr
		if e := recover(); e != nil {
			dbms.Log("send previous errors: " + fmt.Sprint(e))
		}
	}()
	// send errors to server
	in := bufio.NewScanner(f)
	in.Buffer(nil, 1024)
	nlines := 0
	for in.Scan() {
		dbms.Log("PREV: " + in.Text())
		if nlines++; nlines > 1000 {
			dbms.Log("PREV: too many errors")
			break
		}
	}
}

// runServer does not return
func runServer() {
	log.Println("starting server")
	openDbms()
	startHttpStatus()
	run("Init()")
	options.DbStatus.Store("")
	exit.Add(stopServer)
	dbms.Server(dbmsLocal)
	log.Fatalln("FATAL server should not return")
}

func stopServer() {
	log.Println("server stopping")
	httpServer.Close()
	dbms.StopServer()
	heap := builtin.HeapSys()
	log.Println("server stopped, heap:", heap/(1024*1024), "mb,",
		"goroutines:", runtime.NumGoroutine())
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
	exit.Add(func() {
		if options.Action == "server" {
			log.Println("database closing")
			defer log.Println("database closed")
		}
		db.CloseKeepMapped()
	}) // keep mapped to avoid errors during shutdown
	// go checkState()
}

func persistInterval() time.Duration {
	if options.Action == "server" {
		return 60 * time.Second
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

	built := options.BuiltStr()
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
			LogUncaught(&mainThread, "repl", e)
		}
	}()
	src = "function () {\n" + src + "\n}"
	v, results := compile.Checked(&mainThread, src)
	for _, s := range results {
		fmt.Println("(" + s + ")")
	}
	fn := v.(*SuFunc)
	// fmt.Println(DisasmMixed(fn, src))

	mainThread.Reset()
	mainThread.SetSviews(&sviews)
	result := mainThread.Call(fn)
	if result != nil {
		fmt.Println(WithType(result)) // NOTE: doesn't use ToString
	}
}

//-------------------------------------------------------------------

// libload loads a name from the dbms
func libload(th *Thread, name string) (result Value, e any) {
	defer func() {
		if e = recover(); e != nil {
			// fmt.Println("INFO: error loading", name, e)
			// dbg.PrintStack()
			result = nil
		}
	}()
	libs := LibsList.Load()
	if libs == nil {
		libs = th.Dbms().Libraries()
		LibsList.Store(libs)
	}
	defs := th.Dbms().LibGet(name)
	ovLib, ovText := LibraryOverrides.Get(name)
	i := 0
	for _, lib := range libs {
		var src string
		if slc.StartsWith(defs[i:], lib) {
			src = defs[i+1]
			i += 2
		}
		if mode == "gui" && strings.HasSuffix(lib, "webgui") {
			continue
		}
		if lib == ovLib {
			src = ovText
		}
		if src != "" {
			result = llcompile(lib, name, src, result)
		}
	}
	if ovLib == "" && ovText != "" {
		result = llcompile("", name, ovText, result)
	}
	if i < len(defs) {
		Fatal("libraries changed without unload", "("+defs[i]+")")
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
