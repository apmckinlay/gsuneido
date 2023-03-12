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

If the entire pattern is literal characters it is compiled to:
	opLiteralEqual/Prefix/Suffix/Substr characters

If the pattern can be executed by onePass, it is compiled to:
	opOnePass ...

If the pattern has a literal prefix it will be compiled to:
	opPrefix len characters ...
*/

type compiler struct {
	src         string
	si          int
	sn          int
	prog        []byte
	ignoreCase  bool
	multiLine   bool
	leftCount   int
	firstTarget int
	leftAnchor  bool
	rightAnchor bool
}

// Compile converts a regular expression string to a Pattern
func Compile(rx string) Pattern {
	co := compile(rx)
	return co.compile2()
}

func compile(rx string) *compiler {
	co := compiler{src: rx, sn: len(rx), prog: make([]byte, 0, 10+2*len(rx)),
		firstTarget: math.MaxInt}
	co.regex()
	if co.si < co.sn {
		panic("regex: closing ) without opening (")
	}
	co.emit(opDoneSave1)
	return &co
}

func (co *compiler) compile2() Pattern {
	literal, allLiteral := co.literalPrefix()
	if allLiteral {
		op := opLiteralSubstr
		if co.leftAnchor && co.rightAnchor {
			op = opLiteralEqual
		} else if co.leftAnchor {
			op = opLiteralPrefix
		} else if co.rightAnchor {
			op = opLiteralSuffix
		}
		// replace prog with literal
		co.prog = slices.Insert(literal, 0, byte(op))
	} else {
		if co.onePass() {
			co.prog = slices.Insert(co.prog, 0, byte(opOnePass))
		}
		if len(literal) > 0 && !co.leftAnchor {
			if len(literal) > 255 {
				literal = literal[:255]
			}
			co.prog = slices.Insert(co.prog, 0, byte(opPrefix),
				byte(len(literal)))
			co.prog = slices.Insert(co.prog, 2, literal...)
		}
	}
	return Pattern(hacks.BStoS(co.prog))
}

func (co *compiler) literalPrefix() ([]byte, bool) {
	prefix := []byte{}
	end := ord.Min(len(co.prog), co.firstTarget)
	for i := 0; i < end; i++ {
		switch opType(co.prog[i]) {
		case opChar:
			prefix = append(prefix, co.prog[i+1])
			i++
		case opDoneSave1:
			return prefix, true
		case opStrStart, opStrEnd:
			// ignore
		default:
			return prefix, false
		}
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
		}
		co.emit(opStrStart)
	} else if co.match(`\Z`) || (!co.multiLine && co.match("$")) {
		if co.si >= len(co.src) {
			co.rightAnchor = true
		}
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
		if co.leftCount >= 10 {
			panic("regex: too many parenthesized groups")
		}
		co.emit(opSave, 2*byte(co.leftCount))
		co.regex() // recurse
		co.emit(opSave, 2*byte(co.leftCount)+1)
		co.mustMatch(")")
	} else {
		if co.si+1 < co.sn {
			co.match("\\")
		}
		co.emitChar(co.src[co.si])
		co.si++
	}
}

func (co *compiler) charClass() {
	negate := co.match("^")
	chars := ""
	var cc = cclass{}
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

func (co *compiler) posixClass() *cclass {
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

// matching

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

// emit

func (co *compiler) emitChar(c byte) {
	if co.ignoreCase {
		co.emit(opCharIgnoreCase, ascii.ToLower(c))
	} else {
		co.emit(opChar, c)
	}
}

func (co *compiler) emit(op opType, arg ...byte) {
	co.prog = append(append(co.prog, byte(op)), arg...)
}

func (co *compiler) emitOff(op opType, n int) {
	co.emit(op, byte(n>>8), byte(n))
}

func (co *compiler) emitCC(b *cclass) {
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
	co.prog = append(co.prog, 0, 0, 0)
	copy(co.prog[i+3:], co.prog[i:])
	co.prog[i] = byte(op)
	co.prog[i+1] = byte(n >> 8)
	co.prog[i+2] = byte(n)
}

//-------------------------------------------------------------------

// onePass determines if the pattern can be executed in one pass.
// To avoid rewriting the pattern, it requires that each split
// has one branch that leads immediately to a concrete operation (e.g. opChar)
func (co *compiler) onePass() bool {
	if co.leftAnchor && len(co.prog) < 1000 && co.onePass1() && co.onePass2() {
		co.onePass3()
		return true
	}
	return false
}

// onePass1 checks if all splits have an immediate concrete branch
func (co *compiler) onePass1() bool {
	for i := 0; i < len(co.prog); i++ {
		switch opType(co.prog[i]) {
		case opChar, opCharIgnoreCase, opSave:
			i += 1
		case opJump:
			i += 2
		case opSplitFirst, opSplitLast:
			if !co.concrete(i+3) && !co.concrete(co.target(i)) {
				return false
			}
			i += 2
		case opHalfSet:
			i += 16
		case opFullSet:
			i += 32
		case opListSet:
			i += 1 + int(co.prog[i+1])
		}
	}
	return true
}

func (co *compiler) target(i int) int {
	return i + int(int16(co.prog[i+1])<<8|int16(co.prog[i+2]))
}

func (co *compiler) concrete(i int) bool {
	for ; i < len(co.prog) && opType(co.prog[i]) == opSave; i += 2 {
	}
	if i >= len(co.prog) {
		return false
	}
	switch opType(co.prog[i]) {
	case opChar, opCharIgnoreCase, opListSet, opHalfSet, opFullSet:
		return true
	}
	return false
}

// onePass2 checks if all splits are concrete and disjoint
func (co *compiler) onePass2() bool {
	var inProgress BitSet
	for pi := 0; pi < len(co.prog); pi++ {
		switch opType(co.prog[pi]) {
		case opChar, opCharIgnoreCase, opSave:
			pi += 1
		case opJump:
			pi += 2
		case opSplitFirst, opSplitLast:
			inProgress.Clear()
			var cc cclass
			if co.chars(pi, &cc, inProgress) == nil {
				return false
			}
			pi += 2
		case opHalfSet:
			pi += 16
		case opFullSet:
			pi += 32
		case opListSet:
			pi += 1 + int(co.prog[pi+1])
		}
	}
	return true
}

func (co *compiler) chars(pi int, cc *cclass, inProgress BitSet) *cclass {
	if !inProgress.AddNew(int16(pi)) {
		return nil
	}
	pat := co.prog
	for {
		switch opType(pat[pi]) {
		case opChar:
			return cc.addChars(string(rune(pat[pi+1])))
		case opCharIgnoreCase:
			return cc.addChars(string(rune(pat[pi+1]))).
				addChars(string(rune(ascii.ToUpper(pat[pi+1]))))
		case opListSet:
			n := int(pat[pi+1])
			return cc.addChars(string(pat[pi+2 : pi+2+n]))
		case opHalfSet:
			copy(cc[:16], pat[pi+1:])
			return cc
		case opFullSet:
			copy(cc[:], pat[pi+1:])
			return cc
		case opJump:
			jmp := int(int16(pat[pi+1])<<8 | int16(pat[pi+2]))
			pi += jmp // follow jump
		case opSplitFirst, opSplitLast:
			// although we recurse on both branches,
			// by the time we get here, we know one branch is immediate
			var cc1, cc2 cclass
			c1 := co.chars(pi+3, &cc1, inProgress)          // RECURSE
			c2 := co.chars(co.target(pi), &cc2, inProgress) // RECURSE
			if c1 == nil || c2 == nil || !disjoint(c1, c2) {
				return nil
			}
			return cc.add(c1).add(c2)
		case opDoneSave1:
			if co.rightAnchor {
				return cc
			}
			return nil
		case opSave:
			pi += 2
		case opWordStart, opWordEnd, opLineStart, opLineEnd, opStrStart, opStrEnd:
			pi++
		default:
			return nil
		}
	}
}

func disjoint(x, y *cclass) bool {
	for i := range x {
		if (x[i] & y[i]) != 0 {
			return false
		}
	}
	return true
}

// onePass3 modifies splits to
// opSplitNext if the immediate concrete branch follows,
// and opSplitJump means it's after the jump.
func (co *compiler) onePass3() {
	for pi := 0; pi < len(co.prog); pi++ {
		switch opType(co.prog[pi]) {
		case opChar, opCharIgnoreCase, opSave:
			pi += 1
		case opJump:
			pi += 2
		case opSplitFirst, opSplitLast:
			if co.concrete(pi + 3) {
				co.prog[pi] = byte(opSplitNext)
			} else {
				co.prog[pi] = byte(opSplitJump)
			}
			pi += 2
		case opHalfSet:
			pi += 16
		case opFullSet:
			pi += 32
		case opListSet:
			pi += int(co.prog[pi+1])
		}
	}
}
