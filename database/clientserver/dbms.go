package clientserver

// Dbms is the interface to the database
// It has two implementation - DbmsLocal and DbmsClient
type Dbms interface {
	// LibGet returns a list of definitions for name
	// alternating library name and definition
	// in Libraries() order
	LibGet(name string) []string
}

// helloSize is the size of the initial connection message from the server
// the size must match cSuneido and jSuneido
const helloSize = 50
