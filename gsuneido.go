package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"strings"

	. "github.com/apmckinlay/gsuneido/base"
	_ "github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/interp"
	"github.com/apmckinlay/gsuneido/interp/global"
)

var _ = global.Add("Suneido", new(SuObject))

func main() {
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
			debug.PrintStack()
			fmt.Println("ERROR:", e)
		}
	}()
	src = "function () {\n" + src + "\n}"
	fn := compile.Constant(src).(*SuFunc)
	//	interp.Disasm(os.Stdout, fn)
	th := interp.NewThread()
	result := th.Call(fn, nil)
	fmt.Print(">>> ", result)
	if result != nil {
		fmt.Print(" (" + reflect.TypeOf(result).String() + ")")
	}
	fmt.Println()
	fmt.Println()
}
