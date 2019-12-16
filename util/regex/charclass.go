// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex

import "github.com/apmckinlay/gsuneido/util/ascii"

// Character classes are compiled to either listSet or bitSet instructions
// listSet is used for small numbers of characters
// bitSet is 256 bits in 32 bytes

// predefined character class instructions
var (
	blank    = cc().addChars(" \t").build()
	digit    = cc().addRange('0', '9').build()
	notDigit = cc().add(digit).negate().build()
	lower    = cc().addRange('a', 'z').build()
	upper    = cc().addRange('A', 'Z').build()
	alpha    = cc().add(lower).add(upper).build()
	alnum    = cc().add(digit).add(alpha).build()
	word     = cc().addChars("_").add(alnum).build()
	notWord  = cc().add(word).negate().build()
	punct    = cc().addChars("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~").build()
	graph    = cc().add(alnum).add(punct).build()
	print    = cc().addChars(" ").add(graph).build()
	xdigit   = cc().addChars("0123456789abcdefABCDEF").build()
	space    = cc().addChars(" \t\r\n").build()
	notSpace = cc().add(space).negate().build()
	cntrl    = cc().addRange('\u0000', '\u001f').addRange('\u007f', '\u009f').build()
)

const maxList = 8
const setSize = 256 / 8

// builder

type builder struct {
	isSet bool
	data  []byte
}

func cc() *builder {
	return &builder{}
}

// addRange adds a range of characters to a character class instruction
func (b *builder) addRange(from, to byte) *builder {
	if from > to {
		return b
	}
	if !b.isSet && len(b.data)+int(to-from) <= maxList {
		for c := from; c <= to; c++ {
			b.data = append(b.data, c)
		}
	} else {
		b.toSet()
		for c := from; ; c++ {
			b.data[c>>3] |= (1 << (c & 7))
			if c == to {
				break
			}
		}
	}
	return b
}

// addChars adds characters to a character class instruction
func (b *builder) addChars(s string) *builder {
	if !b.isSet && len(b.data)+len(s) <= maxList {
		b.data = append(b.data, s...)
	} else {
		b.toSet()
		for i := 0; i < len(s); i++ {
			c := s[i]
			b.data[c>>3] |= (1 << (c & 7))
		}
	}
	return b
}

func (b *builder) add(in inst) *builder {
	if !b.isSet && in.op == listSet && len(b.data)+len(in.data) <= maxList {
		b.data = append(b.data, in.data...)
	} else if in.op == listSet {
		b.addChars(in.data)
	} else {
		b.toSet()
		for i := 0; i < setSize; i++ {
			b.data[i] |= in.data[i]
		}
	}
	return b
}

// negate inverts a builder
func (b *builder) negate() *builder {
	b.toSet()
	for i := 0; i < setSize; i++ {
		b.data[i] = ^b.data[i]
	}
	return b
}

func (b *builder) toSet() {
	if !b.isSet {
		var bits [setSize]byte
		for _, c := range b.data {
			bits[c>>3] |= (1 << (c & 7))
		}
		b.data = bits[:]
		b.isSet = true
	}
}

// ignore makes the character class ignore case
func (b *builder) ignore() {
	if b.isSet {
		for lo := byte('a'); lo <= 'z'; lo++ {
			up := ascii.ToUpper(lo)
			if b.data[lo>>3]&(1<<(lo&7)) != 0 {
				b.data[up>>3] |= (1 << (up & 7))
			} else if b.data[up>>3]&(1<<(up&7)) != 0 {
				b.data[lo>>3] |= (1 << (lo & 7))
			}
		}
	} else {
		s := ""
		for _, c := range b.data {
			if ascii.IsLetter(c) {
				s += string(ascii.ToLower(c)) + string(ascii.ToUpper(c))
			}
		}
		b.data = b.data[:0]
		b.addChars(s)
	}
}

func (b *builder) empty() bool {
	return len(b.data) == 0
}

func (b *builder) build() inst {
	if b.isSet {
		return inst{op: bitSet, data: string(b.data)}
	}
	return inst{op: listSet, data: string(b.data)}
}

// matchSet returns whether a character is in a bitset
func matchSet(set string, c byte) bool {
	return set[c>>3]&(1<<(c&7)) != 0
}
