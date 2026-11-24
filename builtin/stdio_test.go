// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"io"
	"os"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestStdinReads(t *testing.T) {
	a := assert.T(t)
	long := strings.Repeat("x", MaxLine+1)
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	saved := os.Stdin
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = saved
		r.Close()
	})
	writeDone := make(chan error, 1)
	go func() {
		_, err := io.WriteString(w, "repl\nline\r\n\npassword\n"+long+"\nlast")
		closeErr := w.Close()
		if err == nil {
			err = closeErr
		}
		writeDone <- err
	}()

	line, err := ReadStdinLine()
	a.This(err).Is(nil)
	a.This(line).Is("repl")
	a.This(stdio_Readline()).Is(SuStr("line"))
	a.This(stdio_Readline()).Is(SuStr(""))
	a.This(stdio_ReadPassword()).Is(SuStr("password"))
	a.This(stdio_Readline()).Is(SuStr(long))
	a.This(stdio_Readline()).Is(SuStr("last"))
	a.This(stdio_Readline()).Is(False)
	a.This(<-writeDone).Is(nil)
}
