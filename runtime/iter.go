package runtime

// Iter is the internal type for iterators.
// builtin.SuIter wraps Iter and implements Value and methods
type Iter interface {
	Next() Value
	Infinite() bool
	Dup() Iter
}
