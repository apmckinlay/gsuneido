// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex

import (
	"math/bits"

	"github.com/apmckinlay/gsuneido/util/ascii"
)

// Character classes are compiled to either listSet or bitSet instructions
// listSet is used for small numbers of characters
// bitSet is either 128 bits in 16 bytes or 256 bits in 32 bytes
// Note: this is for ASCII only

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

type cclass [32]byte // 256 / 8

func cc() *cclass {
	return &cclass{}
}

// addRange adds a range of characters to a character class instruction
func (cc *cclass) addRange(from, to byte) *cclass {
	if from > to {
		return cc
	}
	for c := from; ; c++ {
		cc[c>>3] |= (1 << (c & 7))
		if c >= to {
			break
		}
	}
	return cc
}

// addChars adds characters to a character class instruction
func (cc *cclass) addChars(s string) *cclass {
	for i := 0; i < len(s); i++ {
		c := s[i]
		cc[c>>3] |= (1 << (c & 7))
	}
	return cc
}

// addChar add a single character to a character class instruction
func (cc *cclass) addChar(c byte) *cclass {
	cc[c>>3] |= (1 << (c & 7))
	return cc
}

// add or's the argument cclass into the receiver
func (cc *cclass) add(b2 *cclass) *cclass {
	for i := range cc {
		cc[i] |= b2[i]
	}
	return cc
}

// negate inverts a builder
func (cc *cclass) negate() *cclass {
	for i := range cc {
		cc[i] = ^cc[i]
	}
	return cc
}

// ignore makes the character class ignore case
func (cc *cclass) ignore() {
	for lo := byte('a'); lo <= 'z'; lo++ {
		up := ascii.ToUpper(lo)
		if cc[lo>>3]&(1<<(lo&7)) != 0 {
			cc[up>>3] |= (1 << (up & 7))
		} else if cc[up>>3]&(1<<(up&7)) != 0 {
			cc[lo>>3] |= (1 << (lo & 7))
		}
	}
}

// setLen returns the length of the cclass as a bit set (16 or 32)
func (cc *cclass) setLen() int {
	for _, b := range cc[16:] {
		if b != 0 {
			return 32
		}
	}
	return 16
}

// listLen returns the length of the cclass as a list of characters
func (cc *cclass) listLen() int {
	n := 0
	for _, x := range cc {
		n += bits.OnesCount8(x)
	}
	return n
}

// list returns the cclass as a list of characters
func (cc *cclass) list() []byte {
	list := make([]byte, 0, 16)
	for i := 0; i < 256; i++ {
		if cc[i>>3]&(1<<(i&7)) != 0 {
			list = append(list, byte(i))
		}
	}
	return list
}

// matchHalfSet returns whether a character is in a half bit set (16 bytes)
func matchHalfSet(set Pattern, c byte) bool {
	return c < 128 && set[c>>3]&(1<<(c&7)) != 0
}

// matchFullSet returns whether a character is in a full bit set (32 bytes)
func matchFullSet(set Pattern, c byte) bool {
	return set[c>>3]&(1<<(c&7)) != 0
}
