// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core


// IDbms is the interface to the dbms package.
// The two implementations, DbmsLocal and DbmsClient, are in the dbms package
type IDbms interface {
	// Admin executes a schema change (create, alter, drop)
	Admin(string, *Sviews)

	// Auth authorizes the connection with the server
	Auth(*Thread, string) bool

	// Check checks the database like -check
	// It returns "" or an error message.
	Check() string

	// Close ends a dbms connection
	Close()

	// Connections returns a list of the current server connections
	Connections() Value

	// Cursor is like a query but independent of any one transaction
	Cursor(query string, sv *Sviews) ICursor

	// Cursors returns the current number of cursors
	Cursors() int

	DisableTrigger(table string)
	EnableTrigger(table string)

	// Exec is used by the new style ServerEval(...)
	Exec(th *Thread, args Value) Value

	// Final returns the current number of final transactions
	Final() int

	// Get returns a single record, for Query1 (dir = One),
	// QueryFirst (dir = Next), or QueryLast (dir = Prev)
	Get(th *Thread, query Value, dir Dir) (Row, *Header, string)

	// Info returns an object containing database information
	Info() Value

	// Kill terminates connections with the given session id.
	// It returns the count of connections ended.
	Kill(string) int

	// LibGet returns a list of definitions for name
	// alternating library name and definition
	// in Libraries() order
	LibGet(name string) []string

	// Libraries returns a list of the libraries currently in use
	Libraries() []string

	// Log writes to the server's error.log
	Log(string)

	// Nonce returns a random string from the server
	Nonce(*Thread) string

	// Run is used by the old style string.ServerEval()
	Run(th *Thread, code string) Value

	Schema(table string) string

	// SessionId sets and/or returns the session id for the current connection
	SessionId(th *Thread, id string) string

	// Size returns the current database size
	Size() uint64

	// Timestamp returns a guaranteed unique date/time
	Timestamp() SuDate

	// Token returns data to use with Auth
	Token() string

	// Transaction starts a transaction
	Transaction(update bool) ITran

	// Transactions returns a list of the outstanding transactions
	Transactions() *SuObject

	// Use removes a library from those in use
	Unuse(lib string) bool

	// Use adds a library to those in use
	Use(lib string) bool

	// Unwrap removes DbmsUnauth if DbmsAuth
	Unwrap() IDbms
}

// ITran is the interface to a dbms transaction,
type ITran interface {
	String() string

	// Abort rolls back the transaction
	Abort() string

	// Complete commits the transaction.
	// It returns "" on success, otherwise the conflict.
	Complete() string

	// Delete deletes a record
	Delete(th *Thread, table string, off uint64)

	// Get returns a single record, for Query1 (dir = One),
	// QueryFirst (dir = Next), or QueryLast (dir = Prev)
	Get(th *Thread, query Value, dir Dir) (Row, *Header, string)

	// Query starts a query
	Query(query string, sv *Sviews) IQuery

	// Action executes an insert, update, or delete
	// and returns the number of records processed
	Action(th *Thread, action string) int

	// Update modifies a record
	Update(th *Thread, table string, off uint64, rec Record) uint64

	// ReadCount returns the number of reads done by the transaction
	ReadCount() int

	// WriteCount returns the number of writes done by the transaction
	WriteCount() int

	// Asof sets the date-time for the transaction if asof is non-zero,
	// and returns the date-time for the transaction,
	// with the date-times as unix milli
	Asof(int64) int64

	Num() int
}

type Dir byte

const (
	Only Dir = '1' // Query1
	Next Dir = '+' // QueryFirst
	Prev Dir = '-' // QueryLast
	Any  Dir = '@' // QueryEmpty?
)

func (dir Dir) String() string {
	switch {
    case dir == Next:
        return "QueryFirst"
    case dir == Prev:
        return "QueryLast"
    case dir == Only:
        return "Query1"
    case dir == Any:
        return "QueryEmpty"
    }
    return ""
}

func (dir Dir) Reverse() Dir {
	switch dir {
	case Next:
		return Prev
	case Prev:
		return Next
	}
	return dir
}

// IQuery is the interface to a database query, either local or client
type IQuery interface {
	IQueryCursor

	// Get returns the next or previous row from a query
	// and its table if the query is updateable
	Get(th *Thread, dir Dir) (Row, string)

	// Output outputs a record to a query
	Output(th *Thread, rec Record)
}

// ICursor is the interface to a database query, either local or client
type ICursor interface {
	IQueryCursor

	// Get returns the next or previous row from a cursor
	// and its table if the query is updateable
	Get(th *Thread, tran ITran, dir Dir) (Row, string)
}

type IQueryCursor interface {
	// Close ends a query
	Close()

	// Header returns the header (columns and fields) for the query
	Header() *Header

	// Keys returns the keys for the query (a list of comma separated strings)
	Keys() []string

	// Order returns the order for the query (a list of columns)
	Order() []string

	// Rewind resets the query to the beginning/end
	Rewind()

	// Strategy returns a description of the optimized query
	Strategy(formatted bool) string

	Tree() Value
}

// For timestamps with milliseconds up to TsThreshold,
// a client is allowed to increment the milliseconds TsInitialBatch times
// before requesting another timestamp.
// For timestamps with milliseconds >= TsThreshold,
// a client should use SuTimestamp incrementing extra.
// These values must match jSuneido.
const (
	TsInitialBatch = 5
	TsThreshold    = 500
)
