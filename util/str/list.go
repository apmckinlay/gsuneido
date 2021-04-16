// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package str

// List is a type wrapper for a slice of strings to allow additional methods.
type List []string

// Note: This wrapper approach doesn't fit well with operations
// that change the length of the list.
// That's why we have Without rather than Remove.

// Has returns true if the list contains the string, false otherwise
func (list List) Has(str string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

// Index returns position of the first occurrence of the given string,
// or -1 if not found.
func (list List) Index(str string) int {
	for i, s := range list {
		if s == str {
			return i
		}
	}
	return -1
}

// Without returns a new slice of strings
// with any occurences of a given string removed,
// maintaining the existing order.
func (list List) Without(str string) []string {
	dest := make([]string, 0, len(list))
	for _, s := range list {
		if s != str {
			dest = append(dest, s)
		}
	}
	return dest
}

// Reverse reverses the order of elements.
func (list List) Reverse() {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

// HasPrefix returns true if list2 is a prefix of list
func (list List) HasPrefix(list2 []string) bool {
	if len(list) < len(list2) {
		return false
	}
	for i := range list2 {
		if list[i] != list2[i] {
			return false
		}
	}
	return true
}

// Equal returns true if list2 is equal to list
func (list List) Equal(list2 []string) bool {
	if len(list) != len(list2) {
		return false
	}
	for i := 0; i < len(list); i++ {
		if list[i] != list2[i] {
			return false
		}
	}
	return true
}
