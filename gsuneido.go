package main // import "github.com/apmckinlay/gsuneido"

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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

var builtDate string // set by: go build -ldflags "-X builtin.builtDate=..."

var prompt = func(s string) { fmt.Print(s); _ = os.Stdout.Sync() }

// dbmsLocal is set if running with a local/standalone database.
var dbmsLocal IDbms
var mainThread *Thread

func main() {
	Global.Builtin("Suneido", new(SuObject))
	options.BuiltDate = builtDate
	flag.BoolVar(&options.Client, "c", false, "run as a client")
	flag.Parse()
	// dependency injection of GetDbms
	if options.Client {
		GetDbms = func() IDbms { return dbms.NewDbmsClient("127.0.0.1:3147") }
		fmt.Println("Running as client")
	} else {
		dbmsLocal = dbms.NewDbmsLocal()
		GetDbms = func() IDbms { return dbmsLocal }
	}
	Libload = libload // dependency injection
	mainThread = NewThread()
	defer mainThread.Close()
	repl()
}

func repl() {
	fm, _ := os.Stdin.Stat()
	if fm.Mode().IsRegular() {
		prompt = func(string) {}
	}

	builtin.Def()
	builtin.Concat()

	prompt("Press Enter twice (i.e. blank line) to execute, q to quit\n")
	if options.Client {
		eval("Init()")
	} else {
		eval("Suneido.Print = PrintStdout;;")
	}
	r := bufio.NewReader(os.Stdin)
	for {
		prompt("> ")
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
		fmt.Println()
	}
}

func eval(src string) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("ERROR:", e)
			if internal(e) {
				debug.PrintStack()
				fmt.Println("---")
				printCallStack(mainThread.CallStack())
			} else if se, ok := e.(*SuExcept); ok {
				printCallStack(se.Callstack)
			} else {
				printCallStack(mainThread.CallStack())
			}
		}
	}()
	src = "function () {\n" + src + "\n}"
	v, results := compile.Checked(src)
	for _, s := range results {
		fmt.Println(s)
	}
	fn := v.(*SuFunc)
	// DisasmMixed(os.Stdout, fn, src)

	mainThread.Reset()
	result := mainThread.Start(fn, nil)
	if result != nil {
		prompt(">>> ")
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

func printCallStack(cs *SuObject) {
	if cs == nil {
		return
	}
	for i := 0; i < cs.ListSize(); i++ {
		frame := cs.ListGet(i)
		fn := frame.Get(nil, SuStr("fn"))
		fmt.Println(fn)
		locals := frame.Get(nil, SuStr("locals"))
		s := locals.String()
		if len(s) > 230 {
			s = s[:230] + "..."
		}
		fmt.Println(s)
	}
}

// libload loads a name from the dbms
func libload(t *Thread, name string) (result Value) {
	defer func() {
		if e := recover(); e != nil {
			panic("error loading " + name + " " + fmt.Sprint(e))
			//result = nil
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
		// fmt.Println("LOAD", name, "SUCCEEDED")
	}
	return
}
