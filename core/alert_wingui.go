// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable && !com

package core

import "github.com/apmckinlay/gsuneido/builtin/goc"

func Alert(args ...any) {
	goc.Alert(args...)
}

func AlertCancel(args ...any) bool {
	return goc.AlertCancel(args...)
}

func Fatal2(s string) {
	goc.Fatal(s)
}
