// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !win32

package console

func OutputToConsole() {
}

func ConsoleAttached() bool {
	return true
}

func RedirOutput() {
}
