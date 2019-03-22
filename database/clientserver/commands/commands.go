package commands

//go:generate stringer -type=Command

// to make stringer: go generate

type Command byte

// command values must match cSuneido and jSuneido
const (
Abort Command = iota
Admin
Auth
Check
Close
Commit
Connections
Cursor
Cursors
Dump
Erase
Exec
Explain
Final
Get
Get1
Header
Info
Keys
Kill
LibGet
Libraries
Load
Log
Nonce
Order
Output
Query
ReadCount
Request
Rewind
Run
SessionId
Size
Timestamp
Token
Transaction
Transactions
Update
WriteCount
)
