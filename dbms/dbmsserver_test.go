// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestLogWithLimit(t *testing.T) {
	assert := assert.T(t)

	// Create a serverConn to test
	sc := &serverConn{}

	// Test 1: Normal logging under limit
	result := sc.limitLog("Hello World") // 11 bytes + 1 newline = 12
	assert.This(result).Is("Hello World")
	assert.This(sc.logSize.Load()).Is(12)

	// Test 2: Add more logs, still under limit
	result = sc.limitLog("Test message") // 12 bytes + 1 newline = 13, total 25
	assert.This(result).Is("Test message")
	assert.This(sc.logSize.Load()).Is(25)

	// Test 3: Add log that crosses the 10KB limit
	largeMsg := strings.Repeat("x", 10240-25-1) // 10214 + 1 newline = 10215, total = 10240
	result = sc.limitLog(largeMsg)
	assert.This(result).Is(largeMsg)
	assert.This(sc.logSize.Load()).Is(10240)

	// Test 4: Next log should trigger warning
	result = sc.limitLog("x") // 1 + 1 newline = 2, total 10242
	assert.This(result).Is("log size limit exceeded (10KB), ignoring further logs")
	assert.This(sc.logSize.Load()).Is(10242)

	// Test 5: Subsequent logs should be ignored (return empty string)
	oldSize := sc.logSize.Load()
	result = sc.limitLog("More logging") // 12 + 1 = 13
	assert.This(result).Is("")
	assert.This(sc.logSize.Load()).Is(oldSize + 13) // Size still increases even though ignored
}

func TestLogWithLimitBoundary(t *testing.T) {
	assert := assert.T(t)

	sc := &serverConn{}

	// Test exactly at the boundary
	result := sc.limitLog(strings.Repeat("a", 10239)) // 10239 + 1 newline = 10240 exactly
	assert.This(result).Is(strings.Repeat("a", 10239))
	assert.This(sc.logSize.Load()).Is(10240)

	// Next log should trigger the limit warning
	result = sc.limitLog("x") // 1 byte + 1 newline = 2, total 10242
	assert.This(result).Is("log size limit exceeded (10KB), ignoring further logs")
	assert.This(sc.logSize.Load()).Is(10242)

	// Subsequent logs should be ignored
	result = sc.limitLog("ignored")
	assert.This(result).Is("")
}

func TestLogWithLimitEmpty(t *testing.T) {
	assert := assert.T(t)

	sc := &serverConn{}

	// Test empty string - still counts 1 byte for the newline that log.Println adds
	result := sc.limitLog("")
	assert.This(result).Is("")
	assert.This(sc.logSize.Load()).Is(1)

	// Test that we can still log after empty
	result = sc.limitLog("hello")
	assert.This(result).Is("hello")
	assert.This(sc.logSize.Load()).Is(int32(7)) // 1 + (5 + 1) = 7
}
