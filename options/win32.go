// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package options

func init() {
	BuiltExtra += " win32"
}
