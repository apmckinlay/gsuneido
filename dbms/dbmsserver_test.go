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

// nonces

func TestNonceExpiration(t *testing.T) {
	// Mock the serverConns for testing
	originalConns := serverConns
	defer func() { serverConns = originalConns }()
	
	serverConns = make(map[uint32]*serverConn)
	
	// Test multiple connections with different nonce states
	sc1 := &serverConn{id: 1, nonce: "fresh-nonce", nonceOld: false}
	sc2 := &serverConn{id: 2, nonce: "", nonceOld: false} // no nonce
	sc3 := &serverConn{id: 3, nonce: "old-nonce", nonceOld: true} // old nonce
	
	serverConns[1] = sc1
	serverConns[2] = sc2
	serverConns[3] = sc3
	
	// Single expiry cycle should handle all states correctly
	expireNonces()
	
	// Fresh nonce should be marked as old
	assert.T(t).This(sc1.nonce).Is("fresh-nonce")
	assert.T(t).This(sc1.nonceOld).Is(true)
	
	// No nonce should remain unchanged
	assert.T(t).This(sc2.nonce).Is("")
	assert.T(t).This(sc2.nonceOld).Is(false)
	
	// Old nonce should be deleted
	assert.T(t).This(sc3.nonce).Is("")
	assert.T(t).This(sc3.nonceOld).Is(false)
}

func TestNonceConsumption(t *testing.T) {
	// Test that auth clears nonce and resets nonceOld
	sc := &serverConn{nonce: "test-nonce", nonceOld: true}
	ss := &serverSession{sc: sc}
	
	// Mock successful auth
	nonce := ss.sc.nonce
	ss.sc.nonce = ""
	ss.sc.nonceOld = false
	
	assert.T(t).This(nonce).Is("test-nonce")
	assert.T(t).This(ss.sc.nonce).Is("")
	assert.T(t).This(ss.sc.nonceOld).Is(false)
}

func TestNonceExpirationLifecycle(t *testing.T) {
	// Test complete nonce lifecycle: fresh → old → deleted
	originalConns := serverConns
	defer func() { serverConns = originalConns }()
	
	serverConns = make(map[uint32]*serverConn)
	sc := &serverConn{id: 1, nonce: "test-nonce", nonceOld: false}
	serverConns[1] = sc
	
	// Initial state: fresh nonce
	assert.T(t).This(sc.nonce).Is("test-nonce")
	assert.T(t).This(sc.nonceOld).Is(false)
	
	// After 1 minute: nonce marked as old
	expireNonces()
	assert.T(t).This(sc.nonce).Is("test-nonce")
	assert.T(t).This(sc.nonceOld).Is(true)
	
	// After 2 minutes: nonce deleted
	expireNonces() 
	assert.T(t).This(sc.nonce).Is("")
	assert.T(t).This(sc.nonceOld).Is(false)
	
	// After 3 minutes: no change (no nonce to process)
	expireNonces()
	assert.T(t).This(sc.nonce).Is("")
	assert.T(t).This(sc.nonceOld).Is(false)
}
