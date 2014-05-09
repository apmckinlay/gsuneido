package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/interp"
	"github.com/apmckinlay/gsuneido/value"
)

func main() {
	if len(os.Args) > 1 {
		eval(os.Args[1])
	} else {
		r := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			line, err := r.ReadString('\n')
			if err != nil || line == "q\n" {
				break
			}
			eval(line)
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
	src = "function () { " + strings.TrimSpace(src) + "}"
	fn := compile.Constant(src).(*value.SuFunc)
	interp.Disasm(os.Stdout, fn)
	th := interp.Thread{}
	result := th.Call(fn, interp.SimpleArgSpecs[0])
	fmt.Println(">>>", result)
}
