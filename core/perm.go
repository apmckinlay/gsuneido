// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import "sync"

// Perms holds the permissions for a connection.
// It is set up one time only by the application's Auth function
// and then moved to serverConn.
// Tables created at runtime are added with all rights.
type Perms struct {
	serverEval map[string]struct{}
	table      map[string]TablePerm
	lock       sync.Mutex // guards table
	schema     SchemaPerm
}

type SchemaPerm byte

const (
	PermCreate SchemaPerm = 1 // allows creating tables and columns
	PermUpdate            = 3 // includes create
)

type TablePerm byte

const (
	PermRead  TablePerm = 1
	PermWrite           = 3 // includes read
)

func (p *Perms) SetSchema(rights string) {
	if p == nil {
		return
	}
	var sp SchemaPerm
	switch rights {
	case "":
		sp = 0
	case "create":
		sp = PermCreate
	case "update":
		sp = PermUpdate
	default:
		panic("invalid schema rights")
	}
	p.schema = sp
}

// AddTable adds a table and its associated rights to the permissions.
func (p *Perms) AddTable(table string, rights string) {
	if p == nil {
		return
	}
	var tp TablePerm
	switch rights {
	case "":
		tp = 0
	case "read":
		tp = PermRead
	case "write":
		tp = PermWrite
	default:
		panic("invalid table rights")
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.table == nil {
		p.table = make(map[string]TablePerm)
	}
	p.table[table] = tp
}

// AddServerEval adds name to the whitelist of function names
// that ServerEval is allowed to execute on the server.
// Use "*" to allow all functions.
func (p *Perms) AddServerEval(name string) {
	if p.serverEval == nil {
		p.serverEval = make(map[string]struct{})
	}
	p.serverEval[name] = struct{}{}
}

// ServerEvalAllowed returns whether name may be used with ServerEval.
func (p *Perms) ServerEvalAllowed(name string) bool {
	if p == nil {
		return true
	}
	if _, ok := p.serverEval[name]; ok {
		return true
	}
	_, ok := p.serverEval["*"]
	return ok
}

// TableActAllowed returns whether the given action is allowed on the table.
// If the table has no specific rights, the fallback rights added with "*" apply.
func (p *Perms) TableActAllowed(table string, action TablePerm) bool {
	if p == nil {
		return true
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	bits, ok := p.table[table]
	if !ok {
		bits, ok = p.table["*"]
		if !ok {
			return false
		}
	}
	return bits&action == action
}

// SchemaActAllowed returns whether the given action is allowed on the schema.
func (p *Perms) SchemaActAllowed(action SchemaPerm) bool {
	if p == nil {
		return true
	}
	return p.schema&action == action
}
