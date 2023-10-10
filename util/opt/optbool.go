// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package opt

// Bool is an optional bool stored in a single byte.
// The zero value is valid but not set.
type Bool struct {
	val byte
}

const (
	OptBoolFalse byte = 1
	OptBoolTrue  byte = 2
)

func (b *Bool) Set(v bool) {
	if v {
		b.val = OptBoolTrue
	} else {
		b.val = OptBoolFalse
	}
}

func (b Bool) IsSet() bool {
	return b.val != 0
}

func (b Bool) NotSet() bool {
	return b.val == 0
}

func (b Bool) Get() bool {
	switch b.val {
	case OptBoolTrue:
		return true
	case OptBoolFalse:
		return false
	default:
		panic("opt.Bool Get when not set")
	}
}

func (b Bool) GetOr(v bool) bool {
	switch b.val {
	case OptBoolTrue:
		return true
	case OptBoolFalse:
		return false
	default:
		return v
	}
}
