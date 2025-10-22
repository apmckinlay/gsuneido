// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !windows || !gui

package core

import "log"

func Alert(args ...any) {
	log.Println(args...)
}

func AlertCancel(args ...any) bool {
	Alert(args...)
	return true
}

func Fatal2(string) {
}
