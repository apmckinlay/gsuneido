// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import "github.com/apmckinlay/gsuneido/util/ascii"

// Character classes are compiled to either listSet or bitSet instructions
// listSet is used for small numbers of characters
// bitSet is either 128 bits in 16 bytes or 256 bits in 32 bytes

// predefined character class instructions
var (
	blank    = cc().addChars(" \t")
	digit    = cc().addRange('0', '9')
	notDigit = cc().add(digit).negate()
	lower    = cc().addRange('a', 'z')
	upper    = cc().addRange('A', 'Z')
	alpha    = cc().add(lower).add(upper)
	alnum    = cc().add(digit).add(alpha)
	word     = cc().addChars("_").add(alnum)
	notWord  = cc().add(word).negate()
	punct    = cc().addChars("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")
	graph    = cc().add(alnum).add(punct)
	print    = cc().addChars(" ").add(graph)
	xdigit   = cc().addChars("0123456789abcdefABCDEF")
	space    = cc().addChars(" \t\r\n")
	notSpace = cc().add(space).negate()
	cntrl    = cc().addRange('\u0000', '\u001f').addRange('\u007f', '\u009f')
)

const maxList = 8
const setSize = 256 / 8

var wordSet = Pattern(word.data)

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
		for c := from; ; c++ {
			b.data = append(b.data, c)
			if c >= to { // before increment to avoid overflow
				break
			}
		}
	} else {
		b.toSet()
		for c := from; ; c++ {
			b.data[c>>3] |= (1 << (c & 7))
			if c >= to {
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

func (b *builder) addBytes(buf []byte) *builder {
	if !b.isSet && len(b.data)+len(buf) <= maxList {
		b.data = append(b.data, buf...)
	} else {
		b.toSet()
		for i := 0; i < len(buf); i++ {
			c := buf[i]
			b.data[c>>3] |= (1 << (c & 7))
		}
	}
	return b
}

func (b *builder) add(b2 *builder) *builder {
	if !b.isSet && b2.isSet && len(b.data)+len(b2.data) <= maxList {
		b.data = append(b.data, b2.data...)
	} else if !b2.isSet {
		b.addBytes(b2.data)
	} else {
		b.toSet()
		for i := 0; i < setSize; i++ {
			b.data[i] |= b2.data[i]
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
	} else if b.hasLetter() {
		buf := make([]byte, 0, len(b.data))
		for _, c := range b.data {
			if ascii.IsLetter(c) {
				buf = append(buf, ascii.ToLower(c), ascii.ToUpper(c))
			} else {
				buf = append(buf, c)
			}
		}
		b.data = buf
		if len(b.data) > maxList {
			b.toSet()
		}
	}
}

func (b *builder) hasLetter() bool {
	for _, c := range b.data {
		if ascii.IsLetter(c) {
			return true
		}
	}
	return false
}

// matchHalfSet returns whether a character is in a bitset
func matchHalfSet(set Pattern, c byte) bool {
	return c < 128 && set[c>>3]&(1<<(c&7)) != 0
}

// matchSet returns whether a character is in a bitset
func matchFullSet(set Pattern, c byte) bool {
	return set[c>>3]&(1<<(c&7)) != 0
}
