/*
Package regex implements Suneido regular expressions
*/
package regex

import (
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/cmatch"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/verify"
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
 *	simple		:	.					any
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
 *
 * handling ignore case:
 * - compile Chars and CharClass to lower case
 * - match has to convert to lower case
 * - also handled by Backref
 * NOTE: assumes that ignore case state is in sync between compile and match
 * this won't be the case for e.g. (abc(?i)def)+
 *
 * Element.nextPossible is used to optimize match
 * if amatch fails at a certain position
 * nextPossible skips ahead
 * so it doesn't just try amatch at every position
 * This makes match almost as fast as indexOf or contains
 */

// Compile converts a regular expression string to a Pattern
func Compile(rx string) Pattern {
	co := Compiler{src: rx, sn: len(rx)}
	return co.compile()
}

// Result ----------------------------------------------------------------------

// maxResult is the maximum number of elements in Result
const maxResult = 10

type Result [maxResult]part

// part holds the results of a match
type part struct {
	// pos1 is the index of the match + 1 (so zero is invalid)
	pos1 int
	// end is the index after the match i.e. non-inclusive
	end int
}

// Range returns the start and end of part of a match, pos is -1 for no match.
// end is after the match i.e. non-inclusive
func (p part) Range() (pos, end int) {
	return p.pos1 - 1, p.end
}

// Part returns the substring of part of a match, "" for no match
func (p part) Part(s string) string {
	if p.pos1 == 0 {
		return ""
	}
	return s[p.pos1-1 : p.end]
}

func (r *Result) String() string {
	s := ""
	for _, p := range r {
		s += "(" + strconv.Itoa(p.pos1-1) + ", " + strconv.Itoa(p.end) + ") "
	}
	return s
}

// Pattern ---------------------------------------------------------------------

const maxBranch = 1000

// Pattern is a compiled regular expression
type Pattern struct {
	pat []Element
}

// Matches returns whether or not a pattern matches a string
func (p Pattern) Matches(s string) bool {
	var result Result
	return p.FirstMatch(s, 0, &result)
}

// FirstMatch finds the first match in the string at or after pos.
// Returns true if a match is found, else false.
func (p Pattern) FirstMatch(s string, pos int, result *Result) bool {
	var w work
	sn := len(s)
	verify.That(0 <= pos && pos <= sn)
	e := p.pat[1] // skip LEFT0
	for si := pos; si <= sn; si = e.nextPossible(s, si, sn) {
		if p.Amatch(s, si, result, &w) {
			return true
		}
	}
	return false
}

// LastMatch finds the last match in the string before pos.
// Returns true if a match is found, else false.
// Does not use the nextPossible optimization so may be slower;
func (p Pattern) LastMatch(s string, pos int, result *Result) bool {
	var w work
	sn := len(s)
	verify.That(0 <= pos && pos <= sn)
	for si := pos; si >= 0; si-- {
		if p.Amatch(s, si, result, &w) {
			return true
		}
	}
	return false
}

// ForEachMatch calls action for each non-overlapping match in the string.
// The action should return true to continue, false to stop.
func (p Pattern) ForEachMatch(s string, action func(*Result) bool) {
	var w work
	var result Result
	sn := len(s)
	e := p.pat[1] // skip LEFT0
	for si := 0; si <= sn; si = e.nextPossible(s, si, sn) {
		if p.Amatch(s, si, &result, &w) {
			if !action(&result) {
				break
			}
			si = ints.Max(result[0].pos1-1, result[0].end-1)
			// -1 since nextPossible will at least increment
		}
	}
}

// work is temporary values used by Amatch
// declared externally and passed in to avoid zeroing on every Amatch call
type work struct {
	alt_si [maxBranch]int
	alt_pi [maxBranch]int
	tmp [maxResult]int
}

// Amatch tries to match at a specific position.
// Returns true if a match is found, else false.
func (p Pattern) Amatch(s string, si int, result *Result, w *work) bool {
	na := 0
	for pi := 0; pi < len(p.pat); {
		e := p.pat[pi]
		if b, ok := e.(Branch); ok {
			w.alt_pi[na] = pi + b.alt
			w.alt_si[na] = si
			na++
			pi += b.main
		} else if j, ok := e.(Jump); ok {
			pi += j.offset
		} else if left, ok := e.(Left); ok {
			i := left.idx
			if i < maxResult {
				w.tmp[i] = si
			}
			pi++
		} else if right, ok := e.(Right); ok {
			i := right.idx
			if i < maxResult {
				result[i].pos1 = w.tmp[i] + 1
				result[i].end = si
			}
			pi++
		} else {
			si = e.omatch(s, si, result)
			if si >= 0 {
				pi++
			} else if na > 0 {
				// backtrack
				na--
				si = w.alt_si[na]
				pi = w.alt_pi[na]
			} else {
				return false
			}
		}
	}
	return true
}

func (p Pattern) String() string {
	s := ""
	for _, e := range p.pat {
		s += e.String() + " "
	}
	return s
}

// compile -----------------------------------------------------------------

type Compiler struct {
	src                 string
	si                  int
	sn                  int
	pat                 []Element
	ignoringCase        bool
	leftCount           int
	inChars             bool
	inCharsIgnoringCase bool
}

func (co *Compiler) compile() Pattern {
	co.emit(LEFT0)
	co.regex()
	co.emit(RIGHT0)
	if co.si < co.sn {
		panic("regex: closing ) without opening (")
	}
	return Pattern{co.pat}
}

func (co *Compiler) regex() {
	start := len(co.pat)
	co.sequence()
	if co.match("|") {
		pn := len(co.pat) - start
		co.insert(start, Branch{main: 1, alt: pn + 2})
		for {
			start = len(co.pat)
			co.sequence()
			pn = len(co.pat) - start
			if co.match("|") {
				co.insert(start, Branch{main: 1, alt: pn + 2})
				co.insert(start, Jump{offset: pn + 2})
			} else {
				break
			}
		}
		co.insert(start, Jump{offset: pn + 1})
	}
}

func (co *Compiler) sequence() {
	for co.si < co.sn && co.src[co.si] != '|' && co.src[co.si] != ')' {
		co.element()
	}
}

func (co *Compiler) element() {
	if co.match("^") {
		co.emit(startOfLine)
	} else if co.match("$") {
		co.emit(endOfLine)
	} else if co.match("\\A") {
		co.emit(startOfString)
	} else if co.match("\\Z") {
		co.emit(endOfString)
	} else if co.match("\\<") {
		co.emit(startOfWord)
	} else if co.match("\\>") {
		co.emit(endOfWord)
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
			co.insert(start, Branch{main: pn + 1, alt: 1})
		} else if co.match("?") {
			co.insert(start, Branch{main: 1, alt: pn + 1})
		} else if co.match("+?") {
			co.emit(Branch{main: 1, alt: -pn})
		} else if co.match("+") {
			co.emit(Branch{main: -pn, alt: 1})
		} else if co.match("*?") {
			co.emit(Branch{main: 1, alt: -pn})
			co.insert(start, Branch{main: pn + 2, alt: 1})
		} else if co.match("*") {
			co.emit(Branch{main: -pn, alt: 1})
			co.insert(start, Branch{main: 1, alt: pn + 2})
		}
	}
}

func (co *Compiler) quoted() {
	start := co.si
	i := strings.Index(co.src[co.si:], "(?-q)")
	if i == -1 {
		co.si = co.sn
	} else {
		co.si += i
	}
	co.emitChars(co.src[start:co.si])
}

func (co *Compiler) simple() {
	if co.match(".") {
		co.emit(any)
	} else if co.match("\\d") {
		co.emit(CharClass{cm: digit})
	} else if co.match("\\D") {
		co.emit(CharClass{cm: notDigit})
	} else if co.match("\\w") {
		co.emit(CharClass{cm: word})
	} else if co.match("\\W") {
		co.emit(CharClass{cm: notWord})
	} else if co.match("\\s") {
		co.emit(CharClass{cm: space})
	} else if co.match("\\S") {
		co.emit(CharClass{cm: notSpace})
	} else if co.matchBackref() {
		i := int(co.src[co.si-1] - '0')
		co.emit(Backref{idx: i, ignoringCase: co.ignoringCase})
	} else if co.match("[") {
		co.charClass()
		co.mustMatch("]")
	} else if co.match("(") {
		co.leftCount++
		i := co.leftCount
		co.emit(Left{idx: i})
		co.regex() // recurse
		co.emit(Right{idx: i})
		co.mustMatch(")")
	} else {
		if co.si+1 < co.sn {
			co.match("\\")
		}
		co.si++
		co.emitChars(co.src[co.si-1 : co.si])
	}
}

func (co *Compiler) emitChars(s string) {
	if co.inChars && co.inCharsIgnoringCase == co.ignoringCase &&
		!co.next1of("?*+") {
		e := co.pat[len(co.pat)-1].(addable)
		co.pat[len(co.pat)-1] = e.add(s)
	} else {
		if co.ignoringCase {
			co.emit(CharsIgnoreCase{s})
		} else {
			co.emit(Chars{s})
		}
		co.inChars = true
		co.inCharsIgnoringCase = co.ignoringCase
	}
}

func (co *Compiler) next1of(set string) bool {
	return co.si < co.sn && strings.IndexByte(set, co.src[co.si]) != -1
}

func (co *Compiler) charClass() {
	negate := co.match("^")
	chars := ""
	if co.match("]") {
		chars += "]"
	}
	var cm cmatch.CharMatch
	for co.si < co.sn && co.src[co.si] != ']' {
		var elem cmatch.CharMatch
		if co.matchRange() {
			elem = cmatch.InRange(co.src[co.si-3], co.src[co.si-1])
		} else if co.match("\\d") {
			elem = digit
		} else if co.match("\\D") {
			elem = notDigit
		} else if co.match("\\w") {
			elem = word
		} else if co.match("\\W") {
			elem = notWord
		} else if co.match("\\s") {
			elem = space
		} else if co.match("\\S") {
			elem = notSpace
		} else if co.match("[:") {
			elem = co.posixClass()
		} else {
			if co.si+1 < co.sn {
				co.match("\\")
			}
			chars += string(co.src[co.si])
			co.si++
			continue
		}
		cm = cm.Or(elem)
	}
	if !negate && cm == nil && len(chars) == 1 {
		// optimization for class with only one character
		co.emitChars(chars)
		return
	}
	if len(chars) > 0 {
		cm = cm.Or(cmatch.AnyOf(chars))
	}
	if cm == nil {
		panic("empty character class")
	}
	if negate {
		cm = cm.Negate()
	}
	if co.ignoringCase {
		co.emit(CharClassIgnoreCase{cm})
	} else {
		co.emit(CharClass{cm: cm})
	}
}

func (co *Compiler) matchRange() bool {
	if co.src[co.si+1] == '-' &&
		co.si+2 < co.sn && co.src[co.si+2] != ']' {
		co.si += 3
		return true
	}
	return false
}

var blank = cmatch.AnyOf(" \t")
var digit = cmatch.InRange('0', '9')
var notDigit = digit.Negate()
var lower = cmatch.InRange('a', 'z')
var upper = cmatch.InRange('A', 'Z')
var alpha = lower.Or(upper)
var alnum = digit.Or(alpha)
var word = alnum.Or(cmatch.Is('_'))
var notWord = word.Negate()
var punct = cmatch.AnyOf("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")
var graph = alnum.Or(punct)
var print = graph.Or(cmatch.Is(' '))
var xdigit = cmatch.AnyOf("0123456789abcdefABCDEF")
var space = cmatch.AnyOf(" \t\r\n")
var notSpace = space.Negate()
var cntrl = cmatch.InRange('\u0000', '\u001f').Or(cmatch.InRange('\u007f', '\u009f'))

func (co *Compiler) posixClass() cmatch.CharMatch {
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

func (co *Compiler) match(s string) bool {
	if strings.HasPrefix(co.src[co.si:], s) {
		co.si += len(s)
		return true
	}
	return false
}

func (co *Compiler) mustMatch(s string) {
	if !co.match(s) {
		panic("regex: missing '" + s + "'")
	}
}

func (co *Compiler) matchBackref() bool {
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

func (co *Compiler) emit(e Element) {
	co.pat = append(co.pat, e)
	co.inChars = false
}

func (co *Compiler) insert(i int, e Element) {
	co.pat = append(co.pat, nil)
	copy(co.pat[i+1:], co.pat[i:])
	co.pat[i] = e
	co.inChars = false
}

// elements of compiled regex --------------------------------------------------

var FAIL = -1

type Element interface {
	/* returns FAIL or the position after the match */
	omatch(s string, si int, res *Result) int

	// nextPossible is an optional optimization
	nextPossible(s string, si int, sn int) int

	String() string
}

type elemDefaults struct {
}

func (elemDefaults) omatch(string, int, *Result) int {
	panic("should not be called")
}
func (elemDefaults) nextPossible(_ string, si int, _ int) int {
	return si + 1
}

type StartOfLine struct {
	elemDefaults
}

func (e StartOfLine) omatch(s string, si int, _ *Result) int {
	if si == 0 || s[si-1] == '\r' || s[si-1] == '\n' {
		return si
	}
	return FAIL
}

func (e StartOfLine) nextPossible(s string, si int, sn int) int {
	if si == sn {
		return si + 1
	}
	j := strings.IndexByte(s[si+1:], '\n')
	if j == -1 {
		return sn
	}
	return si + 1 + j + 1
}

func (e StartOfLine) String() string {
	return "^"
}

var startOfLine StartOfLine

type EndOfLine struct {
	elemDefaults
}

func (e EndOfLine) omatch(s string, si int, _ *Result) int {
	if si >= len(s) || s[si] == '\r' || s[si] == '\n' {
		return si
	}
	return FAIL
}

func (e EndOfLine) String() string {
	return "$"
}

var endOfLine EndOfLine

type StartOfString struct {
	elemDefaults
}

func (e StartOfString) omatch(_ string, si int, _ *Result) int {
	if si == 0 {
		return si
	}
	return FAIL
}

func (e StartOfString) nextPossible(_ string, _ int, sn int) int {
	return sn + 1 // only the initial position is possible
}

func (e StartOfString) String() string {
	return "\\A"
}

var startOfString StartOfString

type EndOfString struct {
	elemDefaults
}

func (e EndOfString) omatch(s string, si int, _ *Result) int {
	if si >= len(s) {
		return si
	}
	return FAIL
}

func (e EndOfString) String() string {
	return "\\Z"
}

var endOfString EndOfString

type StartOfWord struct {
	elemDefaults
}

func (e StartOfWord) omatch(s string, si int, _ *Result) int {
	if si == 0 || !word.Match(s[si-1]) {
		return si
	}
	return FAIL
}
func (e StartOfWord) String() string {
	return "\\<"
}

var startOfWord StartOfWord

type EndOfWord struct {
	elemDefaults
}

func (e EndOfWord) omatch(s string, si int, _ *Result) int {
	if si >= len(s) || !word.Match(s[si]) {
		return si
	}
	return FAIL
}
func (e EndOfWord) String() string {
	return "\\>"
}

var endOfWord EndOfWord

type Backref struct {
	elemDefaults
	idx          int
	ignoringCase bool
}

func (e Backref) omatch(s string, si int, res *Result) int {
	if res[e.idx].end == -1 {
		return FAIL
	}
	b := s[res[e.idx].pos1-1 : res[e.idx].end]
	bn := len(b)
	if e.ignoringCase {
		if si+bn > len(s) {
			return FAIL
		}
		for i := 0; i < bn; i++ {
			if ascii.ToLower(s[si+i]) != ascii.ToLower(b[i]) {
				return FAIL
			}
		}
	} else if !strings.HasPrefix(s[si:], b) {
		return FAIL
	}
	return si + bn
}
func (e Backref) String() string {
	s := ""
	if e.ignoringCase {
		s = "i"
	}
	return s + "\\" + string('0'+e.idx)
}

type addable interface {
	Element
	add(s string) Element
}

type Chars struct {
	chars string
}

func (e Chars) omatch(s string, si int, _ *Result) int {
	if !strings.HasPrefix(s[si:], e.chars) {
		return FAIL
	}
	return si + len(e.chars)
}

func (e Chars) nextPossible(s string, si int, sn int) int {
	j := strings.Index(s[si+1:], e.chars)
	if j == -1 {
		return sn + 1
	}
	return si + 1 + j
}

func (e Chars) add(s string) Element {
	e.chars += s
	return e
}

func (e Chars) String() string {
	return "'" + e.chars + "'"
}

// extend Chars so compile (simple) can treat them the same
type CharsIgnoreCase struct {
	chars string
}

func (e CharsIgnoreCase) omatch(s string, si int, _ *Result) int {
	cn := len(e.chars)
	if si+cn > len(s) {
		return FAIL
	}
	for i := 0; i < cn; i++ {
		if ascii.ToLower(s[si+i]) != ascii.ToLower(e.chars[i]) {
			return FAIL
		}
	}
	return si + len(e.chars)
}

func (e CharsIgnoreCase) nextPossible(s string, si int, sn int) int {
	cn := len(e.chars)
	for si++; si <= sn-cn; si++ {
		for i := 0; ; i++ {
			if i == cn {
				return si
			} else if ascii.ToLower(s[si+i]) != ascii.ToLower(e.chars[i]) {
				break
			}
		}
	}
	return sn + 1 // no possible match
}

func (e CharsIgnoreCase) add(s string) Element {
	e.chars += s
	return e
}

func (e CharsIgnoreCase) String() string {
	return "i'" + e.chars + "'"
}

type CharClass struct {
	cm cmatch.CharMatch
}

func (e CharClass) omatch(s string, si int, _ *Result) int {
	if si >= len(s) {
		return FAIL
	}
	if e.cm.Match(s[si]) {
		return si + 1
	}
	return FAIL
}

func (e CharClass) nextPossible(s string, si int, sn int) int {
	if si >= sn {
		return si + 1
	}
	j := e.cm.IndexIn(s[si+1:])
	if j == -1 {
		return sn + 1
	}
	return si + 1 + j
}

func (e CharClass) String() string {
	return "[...]"
}

type CharClassIgnoreCase struct {
	cm cmatch.CharMatch
}

func (e CharClassIgnoreCase) omatch(s string, si int, _ *Result) int {
	if si >= len(s) {
		return FAIL
	}
	if e.cm.Match(ascii.ToLower(s[si])) ||
		e.cm.Match(ascii.ToUpper(s[si])) {
		return si + 1
	}
	return FAIL
}

func (e CharClassIgnoreCase) nextPossible(s string, si int, sn int) int {
	for si++; si < sn; si++ {
		if e.cm.Match(ascii.ToLower(s[si])) ||
			e.cm.Match(ascii.ToUpper(s[si])) {
			return si
		}
	}
	return sn + 1 // no possible match
}

func (e CharClassIgnoreCase) String() string {
	return "i[...]"
}

type Any struct {
	elemDefaults
	CharClass
}

func (Any) String() string {
	return "."
}

var any = CharClass{cmatch.AnyOf("\r\n").Negate()}

/*
 * Implemented by amatch.
 * Tries to jump to main first
 * after setting up fallback alternative to jump to alt.
 * main and alt are relative offsets
 */
type Branch struct {
	elemDefaults
	main int
	alt  int
}

func (e Branch) String() string {
	return "Branch(" + strconv.Itoa(e.main) + ", " + strconv.Itoa(e.alt) + ")"
}

/* Implemented by amatch. */
type Jump struct {
	elemDefaults
	offset int
}

func (e Jump) String() string {
	return "Jump(" + strconv.Itoa(e.offset) + ")"
}

/* Implemented by amatch. */
type Left struct {
	elemDefaults
	idx int
}

func (e Left) String() string {
	return "Left" + strconv.Itoa(e.idx)
}

var LEFT0 = Left{idx: 0}

/* Implemented by amatch. */
type Right struct {
	elemDefaults
	idx int
}

func (e Right) String() string {
	return "Right" + strconv.Itoa(e.idx)
}

var RIGHT0 = Right{idx: 0}
