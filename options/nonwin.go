// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !windows portable

package options

import "os"

func Redirected() bool {
	return false
}

func Console(s string) {
	os.Stdout.WriteString(s)
}
