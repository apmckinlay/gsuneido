// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !portable

package options

import (
	"syscall"
)

func EscapeArg(arg string) string {
	return syscall.EscapeArg(arg)
}
