package main // import "github.com/apmckinlay/gsuneido"

/*
CheckLibrary('stdlib')

TestRunner.Run(libs: #(stdlib), skipTags: #(gui, windows), quit_on_failure:);;
TestRunner.Run(TestObserverPrint(), libs: #(stdlib), skipTags: #(gui, windows));;

dlv debug -- -c "WorkSpaceControl();MessageLoop()"
*/

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"

	"strings"

	"github.com/apmckinlay/gsuneido/builtin"
	_ "github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/database/dbms"
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var builtDate = "Oct 19 2019" // set by: go build -ldflags "-X builtin.builtDate=..."

var prompt = func(s string) { fmt.Println(s) }

// dbmsLocal is set if running with a local/standalone database.
var dbmsLocal IDbms
var mainThread *Thread

func main() {
	builtin.Init()
	suneido := new(SuObject)
	suneido.SetConcurrent()
	Global.Builtin("Suneido", suneido)
	options.BuiltDate = builtDate
	flag.BoolVar(&options.Client, "c", false, "run as a client")
	flag.StringVar(&options.NetAddr, "p", "127.0.0.1:3147", "network address and/or port")
	flag.BoolVar(&options.Repl, "r", false, "run REPL (not GUI message loop)")
	flag.Parse()
	if !options.Client {
		options.Repl = true
		options.NetAddr = "" // for ServerIP
	}
	Libload = libload // dependency injection
	mainThread = NewThread()
	builtin.UIThread = mainThread
	defer mainThread.Close()
	// dependency injection of GetDbms
	if options.Client {
		GetDbms = func() IDbms { return dbms.NewDbmsClient(options.NetAddr) }
	} else {
		dbmsLocal = dbms.NewDbmsLocal()
		GetDbms = func() IDbms { return dbmsLocal }
		eval("Suneido.Print = PrintStdout;;")
	}
	if options.Repl {
		log.SetFlags(0)
		repl()
	} else {
		initLogger()
		eval("Init()")
		builtin.Run()
	}
}

func initLogger() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	f, err := os.OpenFile("error.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(io.MultiWriter(f, os.Stderr))
}

// REPL -------------------------------------------------------------

func repl() {
	fm, _ := os.Stdin.Stat()
	if fm.Mode().IsRegular() {
		prompt = func(string) {}
	}

	builtin.Def()
	builtin.Concat()

	prompt(builtin.Built())
	if options.Client {
		prompt("Running as client")
	}
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

func eval(src string) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("ERROR:", e)
			if internal(e) {
				debug.PrintStack()
				fmt.Println("---")
				PrintStack(mainThread.Callstack())
			} else if se, ok := e.(*SuExcept); ok {
				PrintStack(se.Callstack)
			} else {
				PrintStack(mainThread.Callstack())
			}
		}
	}()
	src = "function () {\n" + src + "\n}"
	v, results := compile.Checked(mainThread, src)
	for _, s := range results {
		fmt.Println(s)
	}
	fn := v.(*SuFunc)
	// DisasmMixed(os.Stdout, fn, src)

	mainThread.Reset()
	result := mainThread.Start(fn, nil)
	if result != nil {
		fmt.Println(WithType(result)) // NOTE: doesn't use ToString
	}
}

type internalError interface {
	RuntimeError()
}

func internal(e interface{}) bool {
	_, ok := e.(internalError)
	return ok
}

// libload loads a name from the dbms
func libload(t *Thread, gn Gnum, name string) (result Value) {
	defer func() {
		if e := recover(); e != nil {
			// fmt.Println("INFO: error loading", name, e)
			// debug.PrintStack()
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
		// want to pass the name from the start (rather than adding after)
		// so it propogates to nested Named values
		result = compile.NamedConstant(lib, name, src)
		Global.Set(gn, result) // required for overload inheritance
		// fmt.Println("LOAD", name, "SUCCEEDED")
	}
	return
}
