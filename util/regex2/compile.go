// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

/*
 * regular expression grammar and compiled form:
 *
 *	regex	:	sequence				LEFT0 ... RIGHT0
 *			|	sequence (| sequence)+	Branch sequence (Jump Branch sequence)+
 *
 *	sequence	:	element+
 *
 *	element		:	^					opLineStart
 *				|	$					opLineEnd
 *				|	\A					opStringStart
 *				|	\Z					opStringEnd
 *				|	\<					opWordStart
 *				|	\>					opWordEnd
 *				|	(?i)				(only affects compile)
 *				|	(?-i)				(only affects compile)
 *				|	(?q)				(only affects compile)
 *				|	(?-q)				(only affects compile)
 *				|	simple
 *				|	simple ?			opSplitFirst simple
 *				|	simple ??			opSplitLast simple
 *				|	simple +			simple opSplitFirst
 *				|	simple +?			simple opSplitLast
 *				|	simple *			opSplitFirst simple opJump
 *				|	simple *?			opSplitLast simple opJump
 *
 *	simple		:	.					opAny
 *				|	[ charmatch+ ]		character class
 *				|	[^ charmatch+ ]		character class
 *				|	shortcut			character class
 *				|	( regex )			opSave ... opSave
 *				|	char				opChar
 *
 *	charmatch	:	shortcut			character class
 *				|	posix				character class
 *				|	char - char			character class
 *				|	char				character class
 *
 *	shortcut	:	\d					character class
 *				|	\D					character class
 *				|	\w					character class
 *				|	\W					character class
 *				|	\s					character class
 *				|	\S					character class
 *
 *	posix		|	[:alnum:]			character class
 *				|	[:alpha:]			character class
 *				|	[:blank:]			character class
 *				|	[:cntrl:]			character class
 *				|	[:digit:]			character class
 *				|	[:graph:]			character class
 *				|	[:lower:]			character class
 *				|	[:print:]			character class
 *				|	[:punct:]			character class
 *				|	[:space:]			character class
 *				|	[:upper:]			character class
 *				|	[:xdigit:]			character class
 */

// Compile converts a regular expression string to a Pattern
func Compile(rx string) Pattern {
	co := compiler{src: rx, sn: len(rx)}
	return co.compile()
}

type compiler struct {
	src          string
	si           int
	sn           int
	prog         []byte
	ignoringCase bool
	leftCount    int
}

func (co *compiler) compile() Pattern {
	co.emit(opSave, 0)
	co.regex()
	co.emit(opSave, 1)
	co.emit(opStop)
	if co.si < co.sn {
		panic("regex: closing ) without opening (")
	}
	return Pattern(hacks.BStoS(co.prog))
}

func (co *compiler) regex() {
	var patch []int
	start := len(co.prog)
	co.sequence()
	for co.match("|") {
		pn := len(co.prog) - start
		co.insert(start, opSplitFirst, pn+6)
		patch = append(patch, len(co.prog))
		co.emitOff(opJump, 0)
		start = len(co.prog)
		co.sequence()
	}
	n := len(co.prog)
	for _, i := range patch {
		p := n - i
		co.prog[i+1] = byte(p >> 8)
		co.prog[i+2] = byte(p)
	}
}

func (co *compiler) sequence() {
	for co.si < co.sn && co.src[co.si] != '|' && co.src[co.si] != ')' {
		co.element()
	}
}

func (co *compiler) element() {
	if co.match("^") {
		co.emit(opLineStart)
	} else if co.match("$") {
		co.emit(opLineEnd)
	} else if co.match("\\A") {
		co.emit(opStrStart)
	} else if co.match("\\Z") {
		co.emit(opStrEnd)
	} else if co.match("\\<") {
		co.emit(opWordStart)
	} else if co.match("\\>") {
		co.emit(opWordEnd)
	} else if co.match("(?i)") {
		co.ignoringCase = true
	} else if co.match("(?-i)") {
		co.ignoringCase = false
	} else if co.match("(?q)") {
		co.quoted()
	} else if co.match("(?-q)") {
		// handled by quoted
	} else {
		start := len(co.prog)
		co.simple()
		pn := len(co.prog) - start
		if co.match("?") {
			co.insert(start, opSplitFirst, pn+3)
		} else if co.match("??") {
			co.insert(start, opSplitLast, pn+3)
		} else if co.match("+") {
			co.emitOff(opSplitFirst, -pn)
		} else if co.match("+?") {
			co.emitOff(opSplitLast, -pn)
		} else if co.match("*") {
			co.emitOff(opJump, -pn)
			co.insert(start, opSplitFirst, pn+6)
		} else if co.match("*?") {
			co.emitOff(opJump, -pn)
			co.insert(start, opSplitLast, pn+6)
		}
	}
}

func (co *compiler) quoted() {
	start := co.si
	i := strings.Index(co.src[co.si:], "(?-q)")
	if i == -1 {
		co.si = co.sn
	} else {
		co.si += i
	}
	for _, c := range []byte(co.src[start:co.si]) {
		co.emitChar(c)
	}
}

func (co *compiler) leadingAnything() bool {
	for i := 0; i < co.si; i++ {
		if co.src[i] != '(' {
			return false
		}
	}
	if co.si+2 > co.sn ||
		co.src[co.si] != '.' ||
		(co.src[co.si+1] != '*' && co.src[co.si+1] != '+') {
		return false
	}
	return true
}

func (co *compiler) simple() {
	if co.match(".") {
		co.emit(opAny)
	} else if co.match("\\d") {
		co.emitCC(digit)
	} else if co.match("\\D") {
		co.emitCC(notDigit)
	} else if co.match("\\w") {
		co.emitCC(word)
	} else if co.match("\\W") {
		co.emitCC(notWord)
	} else if co.match("\\s") {
		co.emitCC(space)
	} else if co.match("\\S") {
		co.emitCC(notSpace)
	} else if co.match("[") {
		co.charClass()
	} else if co.match("(") {
		if co.match(")") {
			panic("regex: empty parenthesis not allowed")
		}
		co.leftCount++
		if co.leftCount >= 16 {
			panic("regex: too many parenthesized groups")
		}
		i := byte(co.leftCount)
		co.emit(opSave, 2*i)
		co.regex() // recurse
		co.emit(opSave, 2*i+1)
		co.mustMatch(")")
	} else {
		if co.si+1 < co.sn {
			co.match("\\")
		}
		co.emitChar(co.src[co.si])
		co.si++
	}
}

func (co *compiler) emitChar(c byte) {
	if co.ignoringCase {
		co.emit(opCharIgnoreCase, ascii.ToLower(c))
	} else {
		co.emit(opChar, c)
	}
}

func (co *compiler) next1of(set string) bool {
	return co.si < co.sn && strings.IndexByte(set, co.src[co.si]) != -1
}

func (co *compiler) charClass() {
	negate := co.match("^")
	chars := ""
	var cc = builder{}
	for co.si < co.sn {
		if co.matchRange() {
			cc.addRange(co.src[co.si-3], co.src[co.si-1])
		} else if co.match("\\d") {
			cc.add(digit)
		} else if co.match("\\D") {
			cc.add(notDigit)
		} else if co.match("\\w") {
			cc.add(word)
		} else if co.match("\\W") {
			cc.add(notWord)
		} else if co.match("\\s") {
			cc.add(space)
		} else if co.match("\\S") {
			cc.add(notSpace)
		} else if co.match("[:") {
			cc.add(co.posixClass())
		} else {
			if co.si+1 < co.sn {
				co.match("\\")
			}
			chars += co.src[co.si : co.si+1]
			co.si++
		}
		if co.si >= co.sn || co.src[co.si] == ']' {
			break
		}
	}
	co.mustMatch("]")
	if len(chars) > 0 {
		cc.addChars(chars)
	}
	if !negate && len(cc.data) == 1 {
		// optimization - treat single character class as just character
		co.emitChar(cc.data[0])
		return
	}
	if co.ignoringCase {
		cc.ignore()
	}
	if negate {
		cc.negate()
	}
	co.emitCC(&cc)
}

func (co *compiler) matchRange() bool {
	if co.si+2 < co.sn && co.src[co.si+1] == '-' && co.src[co.si+2] != ']' {
		co.si += 3
		return true
	}
	return false
}

func (co *compiler) posixClass() *builder {
	if co.match("alpha:]") {
		return alpha
	} else if co.match("alnum:]") {
		return alnum
	} else if co.match("blank:]") {
		return blank
	} else if co.match("cntrl:]") {
		return cntrl
	} else if co.match("digit:]") {
		return digit
	} else if co.match("graph:]") {
		return graph
	} else if co.match("lower:]") {
		return lower
	} else if co.match("print:]") {
		return print
	} else if co.match("punct:]") {
		return punct
	} else if co.match("space:]") {
		return space
	} else if co.match("upper:]") {
		return upper
	} else if co.match("xdigit:]") {
		return xdigit
	} else {
		panic("regex: bad posix class")
	}
}

// helpers

func (co *compiler) match(s string) bool {
	if strings.HasPrefix(co.src[co.si:], s) {
		co.si += len(s)
		return true
	}
	return false
}

func (co *compiler) mustMatch(s string) {
	if !co.match(s) {
		panic("regex: missing '" + s + "'")
	}
}

func (co *compiler) emit(op opType, arg ...byte) {
	co.prog = append(append(co.prog, byte(op)), arg...)
}

func (co *compiler) emitOff(op opType, n int) {
	co.emit(op, byte(n>>8), byte(n))
}

func (co *compiler) emitCC(b *builder) {
	data := b.data
	if b.isSet {
		assert.That(len(data) == 32)
		if smallSet(data) {
			co.emit(opHalfSet)
			data = data[:16]
		} else {
			co.emit(opFullSet)
		}
		co.prog = append(co.prog, data...)
	} else {
		co.emit(opListSet, byte(len(data)))
	}
}

func smallSet(data []byte) bool {
	for _, b := range data[16:] {
		if b != 0 {
			return false
		}
	}
	return true
}

func (co *compiler) insert(i int, op opType, n int) {
	co.prog = append(co.prog, 0, 0, 0)
	copy(co.prog[i+3:], co.prog[i:])
	co.prog[i] = byte(op)
	co.prog[i+1] = byte(n >> 8)
	co.prog[i+2] = byte(n)
}
