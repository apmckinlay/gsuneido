// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !windows || portable || com

package runtime

import "fmt"

func Alert(args ...any) {
	fmt.Println(args...)
}

func AlertCancel(args ...any) bool {
	Alert(args...)
	return true
}

func Fatal2(string) {
}
