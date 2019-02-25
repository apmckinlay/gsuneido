package main // import "github.com/apmckinlay/gsuneido"

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	_ "github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/language"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = AddGlobal("Suneido", new(SuObject))

func main() {
	language.Def()
	if len(os.Args) > 1 {
		eval(os.Args[1])
	} else {
		fmt.Println("Press Enter twice (i.e. blank line) to execute, q to quit")
		r := bufio.NewReader(os.Stdin)
		for {
			src := ""
			for {
				fmt.Print("> ")
				line, err := r.ReadString('\n')
				line = strings.TrimRight(line, " \t\r\n")
				if err != nil || line == "q" {
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
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("ERROR:", e)
			if internal(e) {
				debug.PrintStack()
			} else if se, ok := e.(*SuExcept); ok {
				for i := 0; i < se.Callstack.ListSize(); i++ {
					fmt.Println(se.Callstack.ListGet(i))
				}
			}
		}
	}()
	src = "function () {\n" + src + "\n}"
	fn := compile.Constant(src).(*SuFunc)
	// Disasm(os.Stdout, fn)
	th := NewThread()
	result := th.Call(fn)
	if result != nil {
		fmt.Print(">>> ", result)
		fmt.Printf(" <%s %T>", result.TypeName(), result)
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
