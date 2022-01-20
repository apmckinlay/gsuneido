// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !windows || portable
// +build !windows portable

package options

func Redirected() bool {
	return false
}
