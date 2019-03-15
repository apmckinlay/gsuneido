package main // import "github.com/apmckinlay/gsuneido"

import (
	"bufio"
	"fmt"
	"hash/adler32"
	"io"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	_ "github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/language"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = Global.Add("Suneido", new(SuObject))

var prompt = func(s string) { fmt.Print(s); os.Stdout.Sync() }

func main() {
	fm, _ := os.Stdin.Stat()
	if fm.Mode().IsRegular() {
		prompt = func(string) {}
	}

	language.Def()
	Libload = libloadFile
	if len(os.Args) > 1 {
		eval(os.Args[1])
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

func printCallStack(cs *SuObject) {
	if cs == nil {
		return
	}
	for i := 0; i < cs.ListSize(); i++ {
		fmt.Println(cs.ListGet(i))
	}
}

// libload loads a name from the libraries in use
// Currently a temporary version that reads from text files
func libloadFile(name string) (result Value) {
	defer func() {
		if e := recover(); e != nil {
			panic("error loading " + name + " " + fmt.Sprint(e))
			result = nil
		}
	}()
	dir := "../stdlib/"
	hash := adler32.Checksum([]byte(name))
	file := dir + name + "_" + strconv.FormatUint(uint64(hash), 16)
	s, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("LOAD", file, "NOT FOUND")
		return nil
	}
	result = compile.NamedConstant(name, string(s))
	fmt.Println("LOAD", name, "SUCCEEDED")
	return
}
