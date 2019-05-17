package runtime

// IDbms is the interface to the dbms package.
// The two implementations, DbmsLocal and DbmsClient, are in the dbms package
type IDbms interface {
	// Admin executes a database request
	Admin(s string)

	// Auth authorizes the connection with the server
	Auth(string) bool

	// Check checks the database like -check
	// It returns "" or an error message.
	Check() string

	// Connections returns a list of the current server connections
	Connections() Value

	// Cursors returns the current number of cursors
	Cursors() int

	// Dump dumps a table or the entire database like -dump
	// It returns "" or an error message.
	Dump(table string) string

	// Exec is used by the new style ServerEval(...)
	Exec(t *Thread, args Value) Value

	// Final returns the current number of final transactions
	Final() int

	// Get returns a single record for Query1, QueryFirst, QueryLast
	Get(tn int, query string, prev, single bool) (Row, *Header)

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
	Libraries() *SuObject

	// Load loads a table or the entire database like -load
	// It returns the number of records loaded.
	Load(table string) int

	// Log writes to the server's error.log
	Log(string)

	// Nonce returns a random string from the server
	Nonce() string

	// Run is used by the old style string.ServerEval()
	Run(code string) Value

	// SessionId sets and/or returns the session id for the current connection
	SessionId(id string) string

	// Size returns the current database size
	Size() int64

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

	Close()
}

// ITran is the interface to a database transaction,
// either local (not implemented yet) or TranClient.
type ITran interface {
	// Abort rolls back the transaction
	Abort()

	// Complete commits the transaction
	Complete() string

	// Get returns a single row for Query1, QueryFirst, QueryLast
	Get(query string, prev, single bool) (Row, *Header)

	// Erase deletes a record
	Erase(adr int)

	// Update modifies a record
	Update(adr int, rec Record) int

	// Request executes an insert, update, or delete
	// and returns the number of records processed
	Request(request string) int

	// Query starts a query
	Query(query string) IQuery

	String() string
}

type Dir byte
const (
	Next Dir = '+'
	Prev Dir = '-'
)

// IQuery is the interface to a database query,
// either local (not implemented yet) or QueryClient.
type IQuery interface {
	// Get returns the next or previous row from a query
	Get(dir Dir) Row

	// Close ends a query
	Close()

	// Header returns the header (columns and fields) for the query
	Header() *Header

	// Keys returns the keys for the query (a list of comma separated strings)
	Keys() *SuObject

	// Order returns the order for the query (a list of columns)
	Order() *SuObject

	// Rewind resets the query to the beginning/end
	Rewind()

	// Strategy returns a description of the optimized query
	Strategy() string
}
