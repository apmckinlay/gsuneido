// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !windows || !gui

package trace

import "os"

func consolePrint(s string) {
	os.Stdout.WriteString(s)
}
