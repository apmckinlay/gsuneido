// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !windows

package aaainitfirst

import (
	"os"

	"github.com/apmckinlay/gsuneido/console"
	"github.com/apmckinlay/gsuneido/options"
)

func init() {
	options.Parse(os.Args[1:])
	console.LogFileAlso()
}
