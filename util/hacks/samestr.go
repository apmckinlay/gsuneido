// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hacks

import "unsafe"

// SameString returns true if the two strings have the same StringData
func SameString(s, t string) bool {
	return len(s) == len(t) &&
		unsafe.StringData(s) == unsafe.StringData(t)
}
