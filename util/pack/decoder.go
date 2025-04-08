// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package pack

// Decoder is used to read values from a binary string (created with Encoder)
// It is somewhat similar to strings.Reader
type Decoder struct {
	s string
	i int
}

func MakeDecoder(s string) Decoder {
	return Decoder{s: s}
}

func (d Decoder) Peek() byte {
	return d.s[d.i]
}

func (d *Decoder) Skip(n int) {
	d.i += n
}

func (d *Decoder) Get1() byte {
	c := d.s[d.i]
	d.i++
	return c
}

func (d *Decoder) Get(n int) string {
	s := d.s[d.i : d.i+n]
	d.i += n
	return s
}

func (d Decoder) Remaining() int {
	return len(d.s) - d.i
}

func (d Decoder) Rest() string {
	return d.s[d.i:]
}

func (d Decoder) Slice(n int) Decoder {
	return Decoder{s: d.s[:d.i+n], i: d.i}
}

func (d *Decoder) Prev(offset int) Decoder {
	return Decoder{s: d.s, i: d.i - offset}
}

func (d *Decoder) Pos() int {
	return d.i
}