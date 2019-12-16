// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex

import (
	"strings"
)

/*
 * regular expression grammar and compiled form:
 *
 *	regex	:	sequence				LEFT0 ... RIGHT0
 *			|	sequence (| sequence)+	Branch sequence (Jump Branch sequence)+
 *
 *	sequence	:	element+
 *
 *	element		:	^					startOfLine
 *				|	$					endOfLine
 *				|	\A					startOfString
 *				|	\Z					endOfString
 *				|	(?i)				(only affects compile)
 *				|	(?-i)				(only affects compile)
 *				|	(?q)				(only affects compile)
 *				|	(?-q)				(only affects compile)
 *				|	\<					startOfWord
 *				|	\>					endOfWord
 *				|	\#					Backref(#)
 *				|	simple
 *				|	simple ?			Branch simple
 *				|	simple +			simple Branch
 *				|	simple *			Branch simple Branch
 *				|	simple ??			Branch simple
 *				|	simple +?			simple Branch
 *				|	simple *?			Branch simple Branch
 *
 *	simple		:	.					dot
 *				|	[ charmatch+ ]		CharClass
 *				|	[^ charmatch+ ]		CharClass
 *				|	shortcut			CharClass
 *				|	( regex )			Left(i) ... Right(i)
 *				|	chars				Chars(string) // multiple characters
 *
 *	charmatch	:	shortcut			CharClass
 *				|	posix				CharClass
 *				|	char - char			CharClass
 *				|	char				CharClass
 *
 *	shortcut	:	\d					CharClass
 *				|	\D					CharClass
 *				|	\w					CharClass
 *				|	\W					CharClass
 *				|	\s					CharClass
 *				|	\S					CharClass
 *
 *	posix		|	[:alnum:]			CharClass
 *				|	[:alpha:]			CharClass
 *				|	[:blank:]			CharClass
 *				|	[:cntrl:]			CharClass
 *				|	[:digit:]			CharClass
 *				|	[:graph:]			CharClass
 *				|	[:lower:]			CharClass
 *				|	[:print:]			CharClass
 *				|	[:punct:]			CharClass
 *				|	[:space:]			CharClass
 *				|	[:upper:]			CharClass
 *				|	[:xdigit:]			CharClass
 */

// Compile converts a regular expression string to a Pattern
func Compile(rx string) Pattern {
	co := compiler{src: rx, sn: len(rx)}
	return co.compile()
}

type compiler struct {
	src                 string
	si                  int
	sn                  int
	pat                 []inst
	ignoringCase        bool
	leftCount           int
	inChars             bool
	inCharsIgnoringCase bool
}

var (
	left0  = inst{op: left, i: 0}
	right0 = inst{op: right, i: 0}
)

func (co *compiler) compile() Pattern {
	co.emit(left0)
	if co.sn >= 2 && co.startsWithAnything() { //BUG has to be inside grouping
		co.emit(inst{op: startOfLine})
	}
	co.regex()
	co.emit(right0)
	if co.si < co.sn {
		panic("regex: closing ) without opening (")
	}
	return Pattern(co.pat)
}

func (co *compiler) startsWithAnything() bool {
	i := 0
	for co.src[i] == '(' {
		i++
	}
	s := co.src[i:]
	return strings.HasPrefix(s, ".*") || strings.HasPrefix(s, ".+")
}

func (co *compiler) regex() {
	start := len(co.pat)
	co.sequence()
	if co.match("|") {
		pn := len(co.pat) - start
		co.insert(start, inst{op: branch, jump: 1, alt: int16(pn + 2)})
		for {
			start = len(co.pat)
			co.sequence()
			pn = len(co.pat) - start
			if co.match("|") {
				co.insert(start, inst{op: branch, jump: 1, alt: int16(pn + 2)})
				co.insert(start, inst{op: jump, jump: int16(pn + 2)})
			} else {
				break
			}
		}
		co.insert(start, inst{op: jump, jump: int16(pn + 1)})
	}
}

func (co *compiler) sequence() {
	for co.si < co.sn && co.src[co.si] != '|' && co.src[co.si] != ')' {
		co.element()
	}
}

func (co *compiler) element() {
	if co.match("^") {
		co.emit(inst{op: startOfLine})
	} else if co.match("$") {
		co.emit(inst{op: endOfLine})
	} else if co.match("\\A") {
		co.emit(inst{op: startOfString})
	} else if co.match("\\Z") {
		co.emit(inst{op: endOfString})
	} else if co.match("\\<") {
		co.emit(inst{op: startOfWord})
	} else if co.match("\\>") {
		co.emit(inst{op: endOfWord})
	} else if co.match("(?i)") {
		co.ignoringCase = true
	} else if co.match("(?-i)") {
		co.ignoringCase = false
	} else if co.match("(?q)") {
		co.quoted()
	} else if co.match("(?-q)") {
		// handled by quoted
	} else {
		start := len(co.pat)
		co.simple()
		pn := len(co.pat) - start
		if co.match("??") {
			co.insert(start, inst{op: branch, jump: int16(pn + 1), alt: 1})
		} else if co.match("?") {
			co.insert(start, inst{op: branch, jump: 1, alt: int16(pn + 1)})
		} else if co.match("+?") {
			co.emit(inst{op: branch, jump: 1, alt: int16(-pn)})
		} else if co.match("+") {
			co.emit(inst{op: branch, jump: int16(-pn), alt: 1})
		} else if co.match("*?") {
			co.emit(inst{op: branch, jump: 1, alt: int16(-pn)})
			co.insert(start, inst{op: branch, jump: int16(pn + 2), alt: 1})
		} else if co.match("*") {
			co.emit(inst{op: branch, jump: int16(-pn), alt: 1})
			co.insert(start, inst{op: branch, jump: 1, alt: int16(pn + 2)})
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
	co.emitChars(co.src[start:co.si])
}

func (co *compiler) simple() {
	if co.match(".") {
		co.emit(inst{op: dot})
	} else if co.match("\\d") {
		co.emit(digit)
	} else if co.match("\\D") {
		co.emit(notDigit)
	} else if co.match("\\w") {
		co.emit(word)
	} else if co.match("\\W") {
		co.emit(notWord)
	} else if co.match("\\s") {
		co.emit(space)
	} else if co.match("\\S") {
		co.emit(notSpace)
	} else if co.matchBackref() {
		i := int(co.src[co.si-1] - '0')
		if co.ignoringCase {
			co.emit(inst{op: backrefIgnore, i: byte(i)})
		} else {
			co.emit(inst{op: backref, i: byte(i)})
		}
	} else if co.match("[") {
		co.charClass()
		co.mustMatch("]")
	} else if co.match("(") {
		co.leftCount++
		i := co.leftCount
		if i < maxResult {
			co.emit(inst{op: left, i: byte(i)})
		}
		co.regex() // recurse
		if i < maxResult {
			co.emit(inst{op: right, i: byte(i)})
		}
		co.mustMatch(")")
	} else {
		if co.si+1 < co.sn {
			co.match("\\")
		}
		co.si++
		co.emitChars(co.src[co.si-1 : co.si])
	}
}

func (co *compiler) emitChars(s string) {
	if co.inChars && co.inCharsIgnoringCase == co.ignoringCase &&
		!co.next1of("?*+") {
		in := &co.pat[len(co.pat)-1]
		in.data += s
	} else {
		if co.ignoringCase {
			co.emit(inst{op: charsIgnore, data: s})
		} else {
			co.emit(inst{op: chars, data: s})
		}
		co.inChars = true
		co.inCharsIgnoringCase = co.ignoringCase
	}
}

func (co *compiler) next1of(set string) bool {
	return co.si < co.sn && strings.IndexByte(set, co.src[co.si]) != -1
}

func (co *compiler) charClass() {
	negate := co.match("^")
	chars := ""
	if co.match("]") {
		chars += "]"
	}
	var cc = builder{}
	for co.si < co.sn && co.src[co.si] != ']' {
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
	}
	if len(chars) > 0 {
		cc.addChars(chars)
	}
	if negate {
		cc.negate()
	}
	// optimization - treat single character class as just character
	if len(cc.data) == 1 {
		co.emitChars(string(cc.data[0:1]))
		return
	}
	if co.ignoringCase {
		cc.ignore()
	}
	co.emit(cc.build())
}

func (co *compiler) matchRange() bool {
	if co.src[co.si+1] == '-' &&
		co.si+2 < co.sn && co.src[co.si+2] != ']' {
		co.si += 3
		return true
	}
	return false
}

func (co *compiler) posixClass() inst {
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
		panic("bad posix class")
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

func (co *compiler) matchBackref() bool {
	if co.si+2 > co.sn || co.src[co.si] != '\\' {
		return false
	}
	c := co.src[co.si+1]
	if c < '1' || '9' < c {
		return false
	}
	co.si += 2
	return true
}

func (co *compiler) emit(in inst) {
	co.pat = append(co.pat, in)
	co.inChars = false
}

func (co *compiler) insert(i int, in inst) {
	co.pat = append(co.pat, inst{})
	copy(co.pat[i+1:], co.pat[i:])
	co.pat[i] = in
	co.inChars = false
}
