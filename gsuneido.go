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
	"os/exec"
	"runtime"
	"runtime/debug"
	"strconv"
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
var dbmsLocal IDbms
var mainThread *Thread

func main() {
	options.BuiltDate = builtDate
	options.Mode = mode
	options.Parse(os.Args[1:])
	if options.Action == "client" {
		options.Errlog = builtin.ErrlogDir() + "suneido" + options.Port + ".err"
	}
	if options.Mode == "gui" {
		relaunchWithRedirect()
	}
	if options.Action == "" && options.Mode != "gui" {
		options.Action = "repl"
	}

	suneido := new(SuneidoObject)
	suneido.SetConcurrent()
	Global.Builtin("Suneido", suneido)

	switch options.Action {
	case "":
		break
	case "server":
		startServer()
		os.Exit(0)
	case "dump":
		t := time.Now()
		if options.Arg == "" {
			ntables, err := tools.DumpDatabase("suneido.db", "database.su")
			ck(err)
			fmt.Println("dumped", ntables, "tables in",
				time.Since(t).Round(time.Millisecond))
		} else {
			table := strings.TrimSuffix(options.Arg, ".su")
			nrecs, err := tools.DumpTable("suneido.db", table, table+".su")
			ck(err)
			fmt.Println("dumped", nrecs, "records from", table,
				"in", time.Since(t).Round(time.Millisecond))
		}
		os.Exit(0)
	case "load":
		t := time.Now()
		if options.Arg == "" {
			n := tools.LoadDatabase("database.su", "suneido.db")
			fmt.Println("loaded", n, "tables in",
				time.Since(t).Round(time.Millisecond))
		} else {
			table := strings.TrimSuffix(options.Arg, ".su")
			n := tools.LoadTable(table, "suneido.db")
			fmt.Println("loaded", n, "records to", table,
				"in", time.Since(t).Round(time.Millisecond))
		}
		os.Exit(0)
	case "compact":
		t := time.Now()
		ntables, err := tools.Compact("suneido.db")
		ck(err)
		fmt.Println("compacted", ntables, "tables in",
			time.Since(t).Round(time.Millisecond))
		os.Exit(0)
	case "check":
		t := time.Now()
		ck(db19.CheckDatabase("suneido.db"))
		fmt.Println("checked database in", time.Since(t).Round(time.Millisecond))
		os.Exit(0)
	case "repair":
		t := time.Now()
		err := db19.CheckDatabase("suneido.db")
		if err == nil {
			fmt.Println("database ok")
		} else {
			ck(db19.Repair("suneido.db", err))
			fmt.Println("repaired database in", time.Since(t).Round(time.Millisecond))
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
	mainThread = NewThread()
	mainThread.UIThread = true
	MainThread = mainThread
	builtin.UIThread = mainThread
	defer mainThread.Close()
	// dependency injection of GetDbms
	if options.Action == "client" {
		GetDbms = func() IDbms {
			return dbms.NewDbmsClient(options.Arg, options.Port)
		}
		clientErrorLog()
	} else {
		openDbms()
	}
	if options.Action == "repl" ||
		(options.Action == "client" && options.Mode != "gui") {
		run("Init.Repl()")
		repl()
		closeDbms()
	} else {
		run("Init()")
		builtin.Run()
	}
}

func run(src string) {
	defer func() {
		if e := recover(); e != nil {
			printStack(e)
			Fatal("ERROR from "+src+" ", e)
		}
	}()
	compile.EvalString(mainThread, src)
}

func printStack(e interface{}) {
	if builtin.InternalError(e) {
		debug.PrintStack()
		fmt.Println("---")
		PrintStack(mainThread.Callstack())
	} else if se, ok := e.(*SuExcept); ok {
		PrintStack(se.Callstack)
	} else {
		PrintStack(mainThread.Callstack())
	}
}

func ck(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func relaunchWithRedirect() {
	// This is the only way I found to redirect stdout/stderr
	// for built-in output e.g. crashes
	if options.NoRelaunch || options.Redirected() {
		return // to avoid infinite loop
	}
	f, err := os.OpenFile(options.Errlog,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		Fatal(err.Error())
	}
	path, _ := os.Executable()
	cmd := exec.Command(path, os.Args[1:]...)
	cmd.Stdout = f
	cmd.Stderr = f
	err = cmd.Start()
	if err != nil {
		Fatal(err.Error())
	}
	os.Exit(0)
}

func clientErrorLog() {
	dbms := mainThread.Dbms()

	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)
	log.SetPrefix(dbms.SessionId("") + " ")

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
	openDbms()
	//TODO
	closeDbms()
}

var db *db19.Database

func openDbms() {
	startHttpStatus()
	var err error
	db, err = db19.OpenDatabase("suneido.db")
	if err != nil {
		log.Println("ERROR:", err)
		err := db19.Repair("suneido.db", err)
		if err != nil {
			log.Fatalln(err)
		}
		db, err = db19.OpenDatabase("suneido.db")
		if err != nil {
			log.Fatalln(err)
		}
	}
	db19.StartTimestamps()
	db19.StartConcur(db, 10*time.Second) //1*time.Minute) //FIXME
	dbmsLocal = dbms.NewDbmsLocal(db)
	GetDbms = func() IDbms { return dbmsLocal }
	exit.Add(dbmsLocal.Close)
}

func closeDbms() {
	if db != nil {
		db.Close()
	}
}

// HTTP status ------------------------------------------------------

func startHttpStatus() {
	http.HandleFunc("/", httpStatus)
	go func() {
		log.Fatalln(http.ListenAndServe(":3148", nil))
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
				<p>Database: `+mb(GetDbms().Size())+`
				`+threads()+`
				`+trans()+`
			</body>
		</html>`)
}

func mb(n uint64) string {
	return strconv.FormatUint(((n+512*1024)/(1024*1024)), 10) + "mb"
}

func threads() string {
	list := builtin.ThreadList()
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
	list := GetDbms().Transactions()
	n := list.Size()
	var sb strings.Builder
	fmt.Fprintf(&sb, "<p>Transactions: (%d) ", n)
	sep := ""
	for i := 0; i < n; i++ {
		sb.WriteString(sep)
		sb.WriteString(ToStr(list.ListGet(i)))
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
			printStack(e)
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
	result := mainThread.Invoke(fn, nil)
	if result != nil {
		fmt.Println(WithType(result)) // NOTE: doesn't use ToString
	}
}

//-------------------------------------------------------------------

// libload loads a name from the dbms
func libload(t *Thread, gn Gnum, name string) (result Value) {
	defer func() {
		if e := recover(); e != nil {
			// fmt.Println("INFO: error loading", name, e)
			// debug.PrintStack()
			Global.Set(gn, nil)
			panic("error loading " + name + " " + fmt.Sprint(e))
		}
	}()
	defs := t.Dbms().LibGet(name)
	if len(defs) == 0 {
		// fmt.Println("LOAD", name, "MISSING")
		return nil
	}
	for i := 0; i < len(defs); i += 2 {
		lib := defs[i]
		src := defs[i+1]
		if s, ok := LibraryOverrides[lib+":"+name]; ok {
			src = s
		}
		if mode == "gui" && strings.HasSuffix(lib, "webgui") {
			continue
		}
		// want to pass the name from the start (rather than adding after)
		// so it propagates to nested Named values
		result = compile.NamedConstant(lib, name, src)
		Global.Set(gn, result) // required for overload inheritance
		// fmt.Println("LOAD", name, "SUCCEEDED")
	}
	return
}
