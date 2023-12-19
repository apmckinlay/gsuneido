// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build gogen

package str

import (
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

// Builder is a simplified version strings.Builder
// that allows inserting.
type Builder struct {
	buf []byte
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b *Builder) Len() int {
	return len(b.buf)
}

// String returns the accumulated string.
// Unlike strings.Builder, this clears the contents.
// This is necessary because we allow inserting.
func (b *Builder) String() string {
	s := hacks.BStoS(b.buf)
	b.buf = nil
	return s
}

// Add appends the contents of s to b's buffer.
func (b *Builder) Add(s string) {
	b.buf = append(b.buf, s...)
}

func (b *Builder) Adds(ss ...string) {
	for _, s := range ss {
		b.buf = append(b.buf, s...)
	}
}

func (b *Builder) Insert(at int, s string) {
	b.buf = slc.Grow(b.buf, len(s))
	copy(b.buf[at+len(s):], b.buf[at:])
	copy(b.buf[at:], s)
}
