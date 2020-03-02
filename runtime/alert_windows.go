// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build win32

package runtime

import "github.com/apmckinlay/gsuneido/builtin/goc"

func Alert(args ...interface{}) {
	goc.Alert(args...)
}

func Fatal(args ...interface{}) {
	goc.Fatal(args...)
}
