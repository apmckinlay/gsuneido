// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !windows

package aaainitfirst

func init() {
	options.Parse(os.Args[1:])
	console.LogFileAlso()
}
