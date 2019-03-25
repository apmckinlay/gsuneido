package main // import "github.com/apmckinlay/gsuneido"

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"runtime/debug"
	"strings"

	_ "github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/database/clientserver"
	"github.com/apmckinlay/gsuneido/language"
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var builtDate string // set by: go build -ldflags "-X builtin.builtDate=..."

var _ = Global.Add("Suneido", new(SuObject))

var prompt = func(s string) { fmt.Print(s); os.Stdout.Sync() }

var dbms clientserver.Dbms

func main() {
	options.BuiltDate = builtDate
	flag.BoolVar(&options.Client, "c", false, "run as a client")
	flag.Parse()
	if options.Client {
		dbms = clientserver.NewDbmsClient("127.0.0.1:3147")
		fmt.Println("Running as client")
	} else {
		dbms = clientserver.NewDbmsLocal()
	}
	Libload = libload
	repl()
}

func repl() {
	fm, _ := os.Stdin.Stat()
	if fm.Mode().IsRegular() {
		prompt = func(string) {}
	}

	language.Def()
	language.Concat()
	if len(flag.Args()) > 1 {
		eval(flag.Arg(1))
	} else {
		prompt("Press Enter twice (i.e. blank line) to execute, q to quit\n")
		r := bufio.NewReader(os.Stdin)
		for {
			src := ""
			for {
				prompt("> ")
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
}

func eval(src string) {
	th := NewThread()
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("ERROR:", e)
			if internal(e) {
				debug.PrintStack()
				fmt.Println("---")
				printCallStack(th.CallStack())
			} else if se, ok := e.(*SuExcept); ok {
				printCallStack(se.Callstack)
			} else {
				printCallStack(th.CallStack())
			}
		}
	}()
	src = "function () {\n" + src + "\n}"
	fn := compile.Constant(src).(*SuFunc)
	// Disasm(os.Stdout, fn)
	result := th.Call(fn)
	if result != nil {
		prompt(">>> ")
		fmt.Print(result)
		if _, ok := result.(SuStr); !ok {
			fmt.Printf(" <%s %T>", result.Type(), result)
		}
		fmt.Println()
	}
	fmt.Println()
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
		fmt.Println(cs.ListGet(i))
	}
}

// libload loads a name from the dbms
// TODO handle multiple libraries
func libload(name string) (result Value) {
	defer func() {
		if e := recover(); e != nil {
			panic("error loading " + name + " " + fmt.Sprint(e))
			result = nil
		}
	}()
	defs := dbms.LibGet(name)
	if len(defs) == 0 {
		return nil
	}
	result = compile.NamedConstant(name, string(defs[1]))
	// fmt.Println("LOAD", name, "SUCCEEDED")
	return
}
