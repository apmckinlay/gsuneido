// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"golang.org/x/term"

	. "github.com/apmckinlay/gsuneido/core"
	utilexit "github.com/apmckinlay/gsuneido/util/exit"
)

// DEPRECATED use WriteStdout
var _ = builtin(PrintStdout, "(string :string) :void")

func PrintStdout(_ *Thread, args []Value) Value {
	s := ToStr(args[0])
	io.WriteString(os.Stdout, s)
	return nil
}

type suStdio struct {
	staticClass[suStdio]
}

func init() {
	Global.Builtin("Stdio", &suStdio{})
}

func (*suStdio) String() string {
	return "Stdio /* builtin class */"
}

func (sfs *suStdio) Equal(other any) bool {
	return sfs == other
}

func (*suStdio) Lookup(_ *Thread, method string) Value {
	return stdioMethods[method]
}

var stdioMethods = methods("stdio")

var _ = staticMethod(stdio_Write, "(s :string) :void")

func stdio_Write(_ *Thread, args []Value) Value {
	s := ToStr(args[0])
	io.WriteString(os.Stdout, s)
	return nil
}

var _ = staticMethod(stdio_WriteErr, "(s :string) :void")

func stdio_WriteErr(_ *Thread, args []Value) Value {
	s := ToStr(args[0])
	io.WriteString(os.Stderr, s)
	return nil
}

var stdinMu sync.Mutex

// ReadStdinLine reads one line without its newline.
// It returns an unterminated final line before reporting io.EOF.
func ReadStdinLine() (string, error) {
	stdinMu.Lock()
	defer stdinMu.Unlock()
	line, n, err := readLineBytes(os.Stdin, 0)
	if n > 0 && err == io.EOF {
		err = nil
	}
	return line, err
}

var _ = staticMethod(stdio_Readline, "() :false|string")

func stdio_Readline() Value {
	line, err := ReadStdinLine()
	if err == io.EOF {
		return False
	}
	if err != nil {
		panic("ReadLine failed: " + err.Error())
	}
	return StrVal(strings.TrimRight(line, "\r"))
}

var _ = staticMethod(stdio_ReadPassword, "() :false|string")

func stdio_ReadPassword() Value {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return stdio_Readline()
	}
	bytes, err := readPasswordMasked()
	if err == errPasswordInterrupted {
		utilexit.Exit(130) // 128 + SIGINT, what Unix shells report for Ctrl-C
		return False
	}
	if err == io.EOF {
		return False
	}
	if err != nil {
		panic("ReadPassword failed: " + err.Error())
	}
	return StrVal(string(bytes))
}

var errPasswordInterrupted = errors.New("password input interrupted")

func readPasswordMasked() ([]byte, error) {
	stdinMu.Lock()
	defer stdinMu.Unlock()
	fd := int(os.Stdin.Fd())
	old, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	defer term.Restore(fd, old)
	defer fmt.Fprint(os.Stderr, "\r\n") // raw mode needs \r\n

	var pw []byte
	buf := make([]byte, 1)
	for {
		if _, err := os.Stdin.Read(buf); err != nil {
			return nil, err
		}
		switch b := buf[0]; {
		case b == '\r' || b == '\n':
			return pw, nil
		case b == 3: // Ctrl-C arrives as a byte in raw mode, not a signal
			fmt.Fprint(os.Stderr, "^C")
			return nil, errPasswordInterrupted
		case b == 4 && len(pw) == 0: // Ctrl-D on empty input
			return nil, io.EOF
		case b == 127 || b == 8: // backspace / delete
			if len(pw) > 0 {
				// drop one whole UTF-8 character
				for len(pw) > 0 {
					last := pw[len(pw)-1]
					pw = pw[:len(pw)-1]
					if last&0xC0 != 0x80 { // not a continuation byte
						break
					}
				}
				fmt.Fprint(os.Stderr, "\b \b")
			}
		case b >= 32:
			pw = append(pw, b)
			if b&0xC0 != 0x80 { // one mask char per character, not per byte
				fmt.Fprint(os.Stderr, "*")
			}
		}
	}
}

var _ = staticMethod(stdio_Members, "() :object")

func stdio_Members() Value {
	return methodList(stdioMethods)
}
