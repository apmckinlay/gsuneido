// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbg

import (
	"bytes"
	"os"
	"runtime"
)

// PrintStack prints the Go call stack to stderr, similar to debug.PrintStack,
// except it limits the size.
func PrintStack() {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	if i := bytes.LastIndexByte(buf[:n], '\n'); i != -1 {
		n = i + 1
	}
	buf = buf[:n]
	buf = bytes.ReplaceAll(buf, []byte("github.com/apmckinlay/"), nil)
	os.Stderr.Write(buf)
}
