// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import (
	"math"
	"strings"

	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"golang.org/x/exp/slices"
)

/*
regular expression grammar and compiled form:

regex		:	seq					... DoneSave1
			|	seq (| seq)+		SplitFirst seq (Jump SplitFirst seq)+

seq   		:	element+

element		:	^					opLineStart
			|	$					opLineEnd
			|	\A					opStrStart
			|	\Z					opStrEnd
			|	\<					opWordStart
			|	\>					opWordEnd
			|	(?i)				ignore case (handled by compile)
			|	(?-i)				(handled by compile)
			|	(?q)				quote (handled by compile)
			|	(?-q)				(handled by compile)
			|	(?m)				multi-line mode (handled by compile)
			|	(?-m)				(handled by compile)
			|	simple
			|	simple ?			SplitFirst simple
			|	simple ??			SplitLast simple
			|	simple +			simple SplitFirst
			|	simple +?			simple SplitLast
			|	simple *			SplitFirst simple Jump
			|	simple *?			SplitLast simple Jump

simple		:	.					opAny
			|	char	 			Char c
			|	( regex )			Save # ... Save #+1
			|	[ charmatch+ ]		character class
			|	[^ charmatch+ ]		character class
			|	shortcut			character class

charmatch	:	shortcut			character class
			|	posix				character class
			|	char - char			character class
			|	char				character class

shortcut	:	\d					character class
			|	\D					character class
			|	\w					character class
			|	\W					character class
			|	\s					character class
			|	\S					character class

posix		|	[:alnum:]			character class
			|	[:alpha:]			character class
			|	[:blank:]			character class
			|	[:cntrl:]			character class
			|	[:digit:]			character class
			|	[:graph:]			character class
			|	[:lower:]			character class
			|	[:print:]			character class
			|	[:punct:]			character class
			|	[:space:]			character class
			|	[:upper:]			character class
			|	[:xdigit:]			character class

If a pattern does not start with \A it will be prefixed with .*?

If a pattern does not contain any splits, it will start with OnePass

If the entire pattern is literal characters it is compiled to:
	[opUnanchored] opLiteral characters

If the pattern has a literal prefix it will be compiled to:
	opLitPre len characters <remainder of pattern>
*/

// Compile converts a regular expression string to a Pattern
func Compile(rx string) Pattern {
	co := compiler{src: rx, sn: len(rx), prog: make([]byte, 0, 10+2*len(rx)),
		onePass: true, firstTarget: math.MaxInt}
	return co.compile()
}

type compiler struct {
	src         string
	si          int
	sn          int
	prog        []byte
	ignoreCase  bool
	multiLine   bool
	leftCount   int
	onePass     bool
	firstTarget int
	leftAnchor  bool
}

var uaString = string(rune(opUnanchored))

func (co *compiler) compile() Pattern {
	co.regex()
	if co.si < co.sn {
		panic("regex: closing ) without opening (")
	}
	co.emit(opDoneSave1)

	literal, allLiteral := co.literalPrefix()
	if allLiteral {
		// replace prog with literal
		co.prog = slices.Insert(literal, 0, byte(opLiteral))
	} else {
		if co.onePass {
			co.prog = slices.Insert(co.prog, 0, byte(opOnePass))
		}
		if len(literal) > 0 {
			co.prog = slices.Insert(co.prog, 0, byte(opLitPrefix),
				byte(len(literal)))
			co.prog = slices.Insert(co.prog, 2, literal...)
		}
	}
	if !co.leftAnchor {
		co.prog = slices.Insert(co.prog, 0, byte(opUnanchored))
	}
	//TODO LitPrefix
	return Pattern(hacks.BStoS(co.prog))
}

func (co *compiler) literalPrefix() ([]byte, bool) {
	prefix := []byte{}
	end := ord.Min(len(co.prog), co.firstTarget)
	for i := 0; i < end; i++ {
		if opType(co.prog[i]) == opDoneSave1 {
			return prefix, true
		}
		if opType(co.prog[i]) != opChar {
			return prefix, false
		}
		prefix = append(prefix, co.prog[i+1])
		i++
	}
	return prefix, false
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
	if co.match(`\A`) || (!co.multiLine && co.match("^")) {
		if len(co.prog) == 0 {
			co.leftAnchor = true
		} else {
			co.emit(opStrStart)
		}
	} else if co.match(`\Z`) || (!co.multiLine && co.match("$")) {
		co.emit(opStrEnd)
	} else if co.match("^") {
		co.emit(opLineStart)
	} else if co.match("$") {
		co.emit(opLineEnd)
	} else if co.match("\\<") {
		co.emit(opWordStart)
	} else if co.match("\\>") {
		co.emit(opWordEnd)
	} else if co.match("(?i)") {
		co.ignoreCase = true
	} else if co.match("(?-i)") {
		co.ignoreCase = false
	} else if co.match("(?m)") {
		co.multiLine = true
	} else if co.match("(?-m)") {
		co.multiLine = false
	} else if co.match("(?q)") {
		co.quoted()
	} else if co.match("(?-q)") {
		// handled by quoted
	} else {
		start := len(co.prog)
		co.simple()
		pn := len(co.prog) - start
		// need to match longer first
		if co.match("??") {
			co.insert(start, opSplitFirst, pn+3)
		} else if co.match("?") {
			co.insert(start, opSplitLast, pn+3)
			co.firstTarget = ord.Min(co.firstTarget, start)
		} else if co.match("+?") {
			co.emitOff(opSplitLast, -pn)
			co.firstTarget = ord.Min(co.firstTarget, start)
		} else if co.match("+") {
			co.emitOff(opSplitFirst, -pn)
		} else if co.match("*?") {
			co.emitOff(opJump, -pn-3)
			co.insert(start, opSplitFirst, pn+6)
		} else if co.match("*") {
			co.emitOff(opJump, -pn-3)
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

func (co *compiler) simple() {
	if co.match(".") {
		co.emit(opAnyNotNL)
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
	if co.ignoreCase {
		co.emit(opCharIgnoreCase, ascii.ToLower(c))
	} else {
		co.emit(opChar, c)
	}
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
	if !negate && cc.listLen() == 1 {
		// optimization - treat single character class as just a character
		co.emitChar(cc.list()[0])
		return
	}
	if co.ignoreCase {
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
	co.onePass = false
	co.emit(op, byte(n>>8), byte(n))
}

func (co *compiler) emitCC(b *builder) {
	var data []byte
	setLen := b.setLen()
	if b.listLen() < setLen {
		data = b.list()
		co.emit(opListSet, byte(len(data)))
	} else if setLen == 16 {
		co.emit(opHalfSet)
		data = b[:16]
	} else {
		co.emit(opFullSet)
		data = b[:]
	}
	co.prog = append(co.prog, data...)
}

func (co *compiler) insert(i int, op opType, n int) {
	co.onePass = false
	co.prog = append(co.prog, 0, 0, 0)
	copy(co.prog[i+3:], co.prog[i:])
	co.prog[i] = byte(op)
	co.prog[i+1] = byte(n >> 8)
	co.prog[i+2] = byte(n)
}
