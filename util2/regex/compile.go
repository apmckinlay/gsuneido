// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/generic/cache"
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
	prog        []byte
	si          int
	sn          int
	leftCount   int
	ignoreCase  bool
	multiLine   bool
	rightAnchor bool
}

// Compile converts a regular expression string to a Pattern
func Compile(rx string) Pattern {
	co := compile(rx)
	return co.compile2()
}

func (pat Pattern) Literal() (string, bool) {
	if opType(pat[0]) == opLiteralSubstr {
		return string(pat[1:]), true
	}
	return "", false
}

func compile(rx string) *compiler {
	co := compiler{src: rx, sn: len(rx), prog: make([]byte, 0, 10+2*len(rx)),
		multiLine: true} //TEMP for backward compatibility
	co.regex()
	if co.si < co.sn {
		panic("regex: closing ) without opening (")
	}
	co.emit(opDoneSave1)
	return &co
}

// compile2 handles optimizations for all literal, one pass, and literal prefix
func (co *compiler) compile2() Pattern {
	leftAnchor := opType(co.prog[0]) == opStrStart
	literal, allLiteral := co.literalPrefix()
	if allLiteral {
		op := opLiteralSubstr
		if leftAnchor && co.rightAnchor {
			op = opLiteralEqual
		} else if leftAnchor {
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
		if len(literal) > 0 && !leftAnchor {
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

func (co *compiler) literalPrefix() (prefix []byte, allLiteral bool) {
	allLiteral = true
	for i := 0; i < len(co.prog); i++ {
		switch opType(co.prog[i]) {
		case opStrStart:
			if i != 0 {
				return prefix, false
			}
		case opChar:
			prefix = append(prefix, co.prog[i+1])
			i++
		case opSave:
			allLiteral = false
			i++
		default:
			return prefix, false
		case opStrEnd:
			if i != len(co.prog)-2 { // immediately before final opDoneSave1
				return prefix, false
			}
		case opDoneSave1:
			return prefix, allLiteral
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
		co.insert(start, opSplitNext, pn+6)
		patch = append(patch, len(co.prog))
		co.emitOff(opJump, 0)
		start = len(co.prog)
		co.sequence()
		co.rightAnchor = false
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
		co.emit(opStrStart)
	} else if co.match(`\Z`) || (!co.multiLine && co.match("$")) {
		co.emit(opStrEnd)
		co.rightAnchor = (co.si == len(co.src)) // tentative
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
		co.simple() // RECURSE
		pn := len(co.prog) - start
		// need to match longer first
		if co.match("??") {
			co.insert(start, opSplitJump, pn+3)
			co.rightAnchor = false
		} else if co.match("?") {
			co.insert(start, opSplitNext, pn+3)
			co.rightAnchor = false
		} else if co.match("+?") {
			co.emitOff(opSplitNext, -pn)
			co.rightAnchor = false
		} else if co.match("+") {
			co.emitOff(opSplitJump, -pn)
			co.rightAnchor = false
		} else if co.match("*?") {
			co.emitOff(opJump, -pn-3)
			co.insert(start, opSplitJump, pn+6)
			co.rightAnchor = false
		} else if co.match("*") {
			//  compile x* as (x+)? as per golang.org/issue/46123
			co.emitOff(opSplitJump, -pn)
			co.insert(start, opSplitNext, pn+6)
			co.rightAnchor = false
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
	switch c := co.next(); c {
	case '.':
		co.emit(opAnyNotNL)
	case '\\':
		if co.si >= co.sn {
			co.emitChar('\\')
			return
		}
		switch c := co.next(); c {
		case 'd':
			co.emitCC(digit)
		case 'D':
			co.emitCC(notDigit)
		case 'w':
			co.emitCC(word)
		case 'W':
			co.emitCC(notWord)
		case 's':
			co.emitCC(space)
		case 'S':
			co.emitCC(notSpace)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			panic("regex: backreferences not supported")
		default:
			co.emitChar(c)
		}
	case '[':
		co.charClass()
	case '(':
		if co.match(")") {
			panic("regex: empty parenthesis not allowed")
		}
		co.leftCount++
		leftCount := co.leftCount
		if leftCount < 10 {
			co.emit(opSave, 2*byte(leftCount))
		}
		co.regex() // RECURSE
		if leftCount < 10 {
			co.emit(opSave, 2*byte(leftCount)+1)
		}
		co.mustMatch(")")
	default:
		co.emitChar(c)
	}
}

func (co *compiler) charClass() {
	negate := co.match("^")
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
			cc.addChar(co.src[co.si])
			co.si++
		}
		if co.si >= co.sn || co.src[co.si] == ']' {
			break
		}
	}
	co.mustMatch("]")
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

func (co *compiler) next() byte {
	co.si++
	return co.src[co.si-1]
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
		if len(data) == 1 {
			// optimization - treat single character class as just a character
			co.emitChar(data[0])
			return
		}
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
	leftAnchor := opType(co.prog[0]) == opStrStart
	if leftAnchor && len(co.prog) < 1000 && co.onePass1() && co.onePass2() {
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
		case opSplitJump, opSplitNext:
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
		case opSplitJump, opSplitNext:
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
		case opSplitJump, opSplitNext:
			// although we recurse on both branches,
			// by the time we get here, we know one branch is immediate
			var cc1, cc2 cclass
			c1 := co.chars(pi+3, &cc1, inProgress)          // RECURSE
			c2 := co.chars(co.target(pi), &cc2, inProgress) // RECURSE
			if c1 == nil || c2 == nil || !disjoint(c1, c2) {
				return nil
			}
			return cc.add(c1).add(c2)
		case opStrEnd:
			if opType(co.prog[pi+1]) == opDoneSave1 {
				return cc
			}
			return nil
		case opSave:
			pi += 2
		case opWordStart, opWordEnd, opLineStart, opLineEnd, opStrStart:
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
		case opSplitJump, opSplitNext:
			if co.concrete(pi + 3) {
				co.prog[pi] = byte(opBranchNext)
			} else {
				co.prog[pi] = byte(opBranchJump)
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
}

// Cache ----------------------------------------------------------------------

type Cache struct {
	*cache.Cache[string, Pattern]
}

func (c *Cache) Get(s string) Pattern {
	if c.Cache == nil {
		c.Cache = cache.New(Compile)
	}
	return c.Cache.Get(s)
}
