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

	// Get returns a single record for Query1, QueryFirst, QueryLast
	Get(query string, prev, single bool) (Row, *Header)

	// Erase deletes a record
	Erase(adr int)

	// Update modifies a record
	Update(adr int, rec Record) int

	// Request executes an insert, update, or delete
	// and returns the number of records processed
	Request(request string) int

	String() string
}

type Header struct {
	Fields  [][]string
	Columns []string
	Map     map[string]RowAt
}

// RowAt specifies the position of a field within a Row
type RowAt struct {
	Reci int16
	Fldi int16
}

type DbRec struct {
	Record
	Adr int
}

type Row []DbRec

func (row Row) Get(hdr *Header, fld string) Value {
	at,ok := hdr.Map[fld]
	if !ok || int(at.Reci) >= len(row) {
		return nil
	}
	return row[at.Reci].GetVal(int(at.Fldi))
}

func (row Row) GetRaw(hdr *Header, fld string) string {
	at,ok := hdr.Map[fld]
	if !ok || int(at.Reci) >= len(row) {
		return ""
	}
	return row[at.Reci].GetRaw(int(at.Fldi))
}
