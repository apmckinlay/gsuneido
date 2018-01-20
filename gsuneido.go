package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/interp"
)

var _ = interp.AddG("Suneido", new(interp.SuObject))

func main() {
	if len(os.Args) > 1 {
		eval(os.Args[1])
	} else {
		r := bufio.NewReader(os.Stdin)
		for {
			src := ""
			for {
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
	fn := compile.Constant(src).(*interp.SuFunc)
	//	interp.Disasm(os.Stdout, fn)
	th := interp.Thread{}
	result := th.Call(fn, interp.SimpleArgSpecs[0])
	fmt.Print(">>> ", result)
	if result != nil {
		fmt.Print(" (" + reflect.TypeOf(result).String() + ")")
	}
	fmt.Println()
	fmt.Println()
}
