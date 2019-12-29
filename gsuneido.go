// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main // import "github.com/apmckinlay/gsuneido"

/*
CheckLibrary('stdlib')

TestRunner.Run(libs: #(stdlib), skipTags: #(gui, windows), quit_on_failure:);;
TestRunner.Run(TestObserverPrint(), libs: #(stdlib), skipTags: #(gui, windows));;

dlv debug -- -c ...
*/

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/apmckinlay/gsuneido/aaainitfirst"
	"github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/database/dbms"
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var builtDate = "Dec 16 2019" // set by: go build -ldflags "-X main.builtDate=..."

var help = `options:
    -c[lient] [ipaddress]
    -p[ort] #
    -r[epl]
    -v[ersion]`

// dbmsLocal is set if running with a local/standalone database.
var dbmsLocal IDbms
var mainThread *Thread

func main() {
	suneido := new(SuObject)
	suneido.SetConcurrent()
	Global.Builtin("Suneido", suneido)

	options.BuiltDate = builtDate
	args := options.Parse(os.Args[1:])
	options.CmdLine = remainder(args)
	if options.Client == "" {
		options.Repl = true
	}

	if options.Version {
		println("gSuneido " + builtDate + " (" + runtime.Version() + " " +
			runtime.GOARCH + " " + runtime.GOOS + ")")
		os.Exit(0)
	}
	if options.Help {
		println(help)
		os.Exit(0)
	}
	Libload = libload // dependency injection
	mainThread = NewThread()
	builtin.UIThread = mainThread
	defer mainThread.Close()
	// dependency injection of GetDbms
	if options.Client != "" {
		addr := options.Client + ":" + options.Port
		GetDbms = func() IDbms { return dbms.NewDbmsClient(addr) }
		clientErrorLog()
	} else {
		dbmsLocal = dbms.NewDbmsLocal()
		GetDbms = func() IDbms { return dbmsLocal }
		eval("Suneido.Print = PrintStdout;;")
	}
	if options.Repl {
		log.SetFlags(0) // no date/time
		aaainitfirst.InputFromConsole()
		repl()
	} else {
		// initLogger()
		eval("Init()")
		builtin.Run()
	}
}

func remainder(args []string) string {
	var sb strings.Builder
	sep := ""
	for _, arg := range args {
		sb.WriteString(sep)
		sep = " "
		if strings.ContainsAny(arg, " '\"") {
			arg = SuStr(arg).String()
		}
		sb.WriteString(arg)
	}
	return sb.String()
}

func clientErrorLog() {
	// unlike cSuneido, client error.log is still in current directory
	// this is partly because stderr has already been redirected
	f, err := os.Open("error.log")
	if err != nil {
		return
	}
	dbms := mainThread.Dbms()
	defer func() {
		f.Close()
		os.Truncate("error.log", 0) // can't remove since open as stderr
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

// REPL -------------------------------------------------------------

var prompt = func(s string) { fmt.Fprintln(os.Stderr, s) }

func repl() {
	if !isTerminal() {
		prompt = func(string) {}
	}

	builtin.Def()
	builtin.Concat()

	built := builtin.Built()
	if options.Client != "" {
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

func isTerminal() bool {
	fm, _ := os.Stdout.Stat()
	mode := fm.Mode()
	return mode&os.ModeDevice == os.ModeDevice &&
		mode&os.ModeCharDevice == os.ModeCharDevice
}

func showOptions() {
	if options.HeapDebug {
		prompt("- HeapDebug enabled")
	}
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
