// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/tools"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dbg"
	"github.com/apmckinlay/gsuneido/util/exit"
	"github.com/apmckinlay/gsuneido/util/generic/hamt"
	"github.com/apmckinlay/gsuneido/util/regex"
	"github.com/apmckinlay/gsuneido/util/system"
)

var builtDate = "Dec 29 2020" // set by: go build -ldflags "-X main.builtDate=..."
var mode = ""                 // set by: go build -ldflags "-X main.mode=gui"

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
	if err := system.Service("gSuneido", redirect, stopServer); err != nil {
		Fatal(err)
	}
	if options.Action == "" && mode != "gui" {
		options.Action = "repl"
	}

	switch options.Action {
	case "":
		break
	case "server":
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
	case "error":
		Fatal(options.Error)
	case "repl", "client":
		// handled below
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
		dbClose()
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
			printStack(e)
			Fatal("ERROR from", src, e)
		}
	}()
	compile.EvalString(mainThread, src)
}

func printStack(e any) {
	if InternalError(e) {
		dbg.PrintStack()
		PrintStack(mainThread.Callstack())
	} else if se, ok := e.(*SuExcept); ok {
		PrintStack(se.Callstack)
	} else {
		PrintStack(mainThread.Callstack())
	}
}

func ck(err error) {
	if err != nil {
		Fatal(err)
	}
}

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

func startServer() {
	log.Println("starting server")
	openDbms()
	startHttpStatus()
	Libload = libload // dependency injection
	mainThread = &Thread{}
	mainThread.Name = "main"
	run("Init()")
	dbms.Server(dbmsLocal)
}

func stopServer() {
	httpServer.Close()
	dbms.StopServer()
	dbmsLocal.Close()
	log.Println("server stopped")
}

var db *db19.Database

func openDbms() {
	// startHttpStatus()
	var err error
	db, err = db19.OpenDatabase("suneido.db")
	if err != nil {
		Alert("ERROR:", err)
		msg, err := db19.Repair("suneido.db", err)
		if err != nil {
			Fatal("repair:", err)
		}
		Alert(msg)
		db, err = db19.OpenDatabase("suneido.db")
		if err != nil {
			Fatal("open:", err)
		}
	}
	db19.StartTimestamps()
	db19.StartConcur(db, 10*time.Second) //1*time.Minute) //FIXME
	dbmsLocal = dbms.NewDbmsLocal(db)
	DbmsAuth = options.Action == "server" || mode != "gui" || !db.HaveUsers()
	GetDbms = getDbms
	exit.Add(dbmsLocal.Close)
	// go checkState()
}

func getDbms() IDbms {
	if DbmsAuth {
		return dbmsLocal
	}
	return dbms.Unauth(dbmsLocal)
}

func dbClose() {
	if db != nil {
		db.Close()
	}
}

func checkState() {
	for {
		state := db.GetState()
		cksum := state.Meta.Cksum()
		// read meta to verify checksums
		schemaOff, infoOff := state.Meta.Offsets()
		if schemaOff != 0 {
			hamt.ReadChain[string](db.Store, schemaOff, meta.ReadSchema)
		}
		if schemaOff != 0 {
			hamt.ReadChain[string](db.Store, infoOff, meta.ReadInfo)
		}
		time.Sleep(50 * time.Millisecond)
		// recalculate checksum to verify Meta hasn't been mutated
		assert.That(state.Meta.Cksum() == cksum)
	}
}

// HTTP status ------------------------------------------------------

var httpServer *http.Server

func startHttpStatus() {
	http.HandleFunc("/", httpStatus)
	port, _ := strconv.Atoi(options.Port)
	addr := ":" + strconv.Itoa(port+1)
	go func() {
		httpServer = &http.Server{Addr: addr}
		err := httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Println("Server Monitor:", err)
		}
	}()
}
func httpStatus(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w,
		`<html>
			<head>
			<title>Suneido Server Monitor</title>
			<meta http-equiv="refresh" content="5" />
			</head>
			<body>
				<h1>Suneido Server Monitor</h1>
				<p>Built: `+builtin.Built()+`</p>
				<p>Heap: `+mb(builtin.HeapSys())+`</p>
				<p>Database: `+mb(dbmsLocal.Size())+`
				`+threads()+`
				`+trans()+`
				`+dbms.Conns()+`
			</body>
		</html>`)
}

func mb(n uint64) string {
	return strconv.FormatUint(((n+512*1024)/(1024*1024)), 10) + "mb"
}

func threads() string {
	list := builtin.ThreadList()
	sort.Strings(list)
	var sb strings.Builder
	fmt.Fprintf(&sb, "<p>Threads: (%d) ", len(list))
	sep := ""
	for _, s := range list {
		sb.WriteString(sep)
		sb.WriteString(s)
		sep = ", "
	}
	sb.WriteString("<p>\n")
	return sb.String()
}

func trans() string {
	list := dbmsLocal.Transactions()
	n := list.Size()
	var sb strings.Builder
	fmt.Fprintf(&sb, "<p>Transactions: (%d) ", n)
	sep := ""
	for i := 0; i < n; i++ {
		sb.WriteString(sep)
		sb.WriteString(list.ListGet(i).String())
		sep = ", "
	}
	sb.WriteString("<p>\n")
	return sb.String()
}

// REPL -------------------------------------------------------------

var prompt = func(s string) { fmt.Fprintln(os.Stderr, s) }

func repl() {
	log.SetFlags(log.Ltime)
	log.SetPrefix("")

	if !isTerminal(os.Stdin) || !isTerminal(os.Stdout) {
		prompt = func(string) {}
	}

	builtin.Def()
	builtin.Concat()

	built := builtin.Built()
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
			log.Println("ERROR:", e)
			if !strings.HasSuffix(fmt.Sprint(e), "(from server)") {
				printStack(e)
			}
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
			// debug.PrintStack()
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
