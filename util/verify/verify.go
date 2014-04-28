/*
Package verify provides a simple assertion facility

For example:

	verify.That(size >= 0)
*/
package verify

// That panics if its argument is false
func That(cond bool) {
	if !cond {
		panic("verify failed")
	}
}
