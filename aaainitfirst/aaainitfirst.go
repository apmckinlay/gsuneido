// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package aaainitfirst

import (
	"log"
	"os"
)

func logFileAlso() {
	log.SetOutput(logWriter{})
}

type logWriter struct{}

// Write outputs to Stderr and error.log
func (lw logWriter) Write(p []byte) (n int, err error) {
	os.Stderr.Write(p)
	f, err := os.OpenFile("error.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(p)
}
