package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/apmckinlay/gsuneido/compile/codegen"
	"github.com/apmckinlay/gsuneido/compile/parse"
	"github.com/apmckinlay/gsuneido/interp"
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

func eval(line string) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			fmt.Println("ERROR:", e)
		}
	}()
	line = "function () { " + strings.TrimSpace(line) + "}"
	ast := parse.Parse(line).(parse.AstNode)
	fmt.Println("ast", ast.String())
	fn := codegen.Codegen(ast)
	interp.Disasm(os.Stdout, fn)
	th := interp.Thread{}
	result := th.Call(fn, interp.SimpleArgSpecs[0])
	fmt.Println(">>>", result)
}
