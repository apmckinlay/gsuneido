package runtime

// IDbms is the interface to the dbms package.
// It has two implementation - DbmsLocal and DbmsClient
type IDbms interface {
	// LibGet returns a list of definitions for name
	// alternating library name and definition
	// in Libraries() order
	LibGet(name string) []string

	// Libraries returns a list of the libraries currently in use
	Libraries() *SuObject

	// Timestamp returns a guaranteed unique date/time
	Timestamp() SuDate
}
