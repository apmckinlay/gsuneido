// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

// Perms holds the permissions for a connection.
// It is set up one time only by the application's Auth function
// (using Perm.ServerEval) and then moved to serverConn.
type Perms struct {
	serverEval map[string]struct{}
	tablePerms map[string]int
	allowAll   bool
}

const (
	PermRead   = 1 << iota // r
	PermAdd                // a
	PermUpdate             // u
	PermDelete             // d
	PermSchema             // s
)

// AddTable adds a table and its associated rights to the permissions.
// Rights must be an ordered subset of "rauds" and is converted to a bit set.
func (p *Perms) AddTable(table string, rights string) {
	match := func(c byte, p int) int {
		if len(rights) > 0 && rights[0] == c {
			rights = rights[1:]
			return p
		}
		return 0
	}
	var bits int
	if rights == "*" {
		bits = PermRead | PermAdd | PermUpdate | PermDelete | PermSchema
	} else {
		bits = match('r', PermRead) |
			match('a', PermAdd) |
			match('u', PermUpdate) |
			match('d', PermDelete) |
			match('s', PermSchema)
		if len(rights) != 0 {
			panic("invalid table rights")
		}
	}
	if p.tablePerms == nil {
		p.tablePerms = make(map[string]int)
	}
	p.tablePerms[table] = bits
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

// TableActAllowed returns whether the given action is allowed on the table.
func (p *Perms) TableActAllowed(table string, action int) bool {
	bits, ok := p.tablePerms[table]
	if !ok {
		return false
	}
	return bits&action != 0
}
