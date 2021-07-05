// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !windows portable

package trace

import "os"

func consolePrintln(s string) {
	os.Stdout.WriteString(s)
}
