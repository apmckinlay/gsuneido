// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

// Perms holds the permissions for a connection.
// It is set up one time only by the application's Auth function
// (using Perm.ServerEval) and then moved to serverConn.
type Perms struct {
	serverEval map[string]struct{}
	allowAll   bool
}

// AddServerEval adds name to the whitelist of function names
// that ServerEval is allowed to execute on the server.
// Use "*" to allow all functions.
func (p *Perms) AddServerEval(name string) {
	if name == "*" {
		p.allowAll = true
		return
	}
	if p.serverEval == nil {
		p.serverEval = make(map[string]struct{})
	}
	p.serverEval[name] = struct{}{}
}

// ServerEvalAllowed returns whether name may be used with ServerEval.
// With no whitelist, nothing is allowed.
func (p *Perms) ServerEvalAllowed(name string) bool {
	if p.allowAll {
		return true
	}
	_, ok := p.serverEval[name]
	return ok
}
