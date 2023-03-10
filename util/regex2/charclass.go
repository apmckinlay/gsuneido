// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import (
	"math/bits"

	"github.com/apmckinlay/gsuneido/util/ascii"
)

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

var wordSet = Pattern(word[:])

// builder

type builder [32]byte // 256 / 8

func cc() *builder {
	return &builder{}
}

// addRange adds a range of characters to a character class instruction
func (b *builder) addRange(from, to byte) *builder {
	if from > to {
		return b
	}
	for c := from; ; c++ {
		b[c>>3] |= (1 << (c & 7))
		if c >= to {
			break
		}
	}
	return b
}

// addChars adds characters to a character class instruction
func (b *builder) addChars(s string) *builder {
	for i := 0; i < len(s); i++ {
		c := s[i]
		b[c>>3] |= (1 << (c & 7))
	}
	return b
}

func (b *builder) add(b2 *builder) *builder {
	for i := range b {
		b[i] |= b2[i]
	}
	return b
}

// negate inverts a builder
func (b *builder) negate() *builder {
	for i := range b {
		b[i] = ^b[i]
	}
	return b
}

// ignore makes the character class ignore case
func (b *builder) ignore() {
	for lo := byte('a'); lo <= 'z'; lo++ {
		up := ascii.ToUpper(lo)
		if b[lo>>3]&(1<<(lo&7)) != 0 {
			b[up>>3] |= (1 << (up & 7))
		} else if b[up>>3]&(1<<(up&7)) != 0 {
			b[lo>>3] |= (1 << (lo & 7))
		}
	}
}

func (b *builder) listLen() int {
	n := 0
	for _, x := range b {
		n += bits.OnesCount8(x)
	}
	return n
}

func (b *builder) setLen() int {
	for _, b := range b[16:] {
		if b != 0 {
			return 32
		}
	}
	return 16
}

func (b *builder) list() []byte {
	list := make([]byte, 0, 16)
	for i := 0; i < 256; i++ {
		if b[i>>3]&(1<<(i&7)) != 0 {
			list = append(list, byte(i))
		}
	}
	return list
}

// matchHalfSet returns whether a character is in a bitset
func matchHalfSet(set Pattern, c byte) bool {
	return c < 128 && set[c>>3]&(1<<(c&7)) != 0
}

// matchSet returns whether a character is in a bitset
func matchFullSet(set Pattern, c byte) bool {
	return set[c>>3]&(1<<(c&7)) != 0
}
