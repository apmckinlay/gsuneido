// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !windows

package builtin

import (
	"os/exec"
)

func cmdSetup(*exec.Cmd, string, bool) {
}
