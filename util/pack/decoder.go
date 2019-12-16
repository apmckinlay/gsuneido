// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package pack

// Decoder is used to read values from a binary string (created with Encoder)
// It is somewhat similar to strings.Reader
type Decoder struct {
	s string
}

func NewDecoder(s string) *Decoder {
	return &Decoder{s}
}

func (d *Decoder) Get1() byte {
	c := d.s[0]
	d.s = d.s[1:]
	return c
}

func (d *Decoder) Get(n int) string {
	s := d.s[:n]
	d.s = d.s[n:]
	return s
}

func (d *Decoder) Remaining() int {
	return len(d.s)
}
