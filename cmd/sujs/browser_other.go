// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !darwin

package main

// findBrowserApp is a no-op on non-macOS platforms.
func findBrowserApp() string {
	return ""
}
