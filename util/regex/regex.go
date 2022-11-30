// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package regex implements Suneido regular expressions
package regex

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/generic/cache"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
)

// Pattern is a compiled regular expression
type Pattern []inst

// inst is a single compiled pattern instruction
type inst struct {
	// data is used by chars and charclass
	data string
	// jump is used by jump and branch
	jump int16
	// alt is used by branch
	alt int16
	op  byte
	// i is used by left, right, and backref
	i byte
}

// op codes
const (
	dot = iota
	// chars is a sequence of literal characters (in data)
	chars
	// charIgnore is a sequence of literal characters to be matched ignoring case
	charsIgnore
	// listSet is a character class represented as a list of characters in data
	listSet
	// bitSet is a character class represented as a 32 byte bit set in data
	bitSet
	branch
	jump
	left
	right
	backref
	backrefIgnore
	startOfLine
	endOfLine
	startOfString
	endOfString
	startOfWord
	endOfWord
)

func (pat Pattern) String() string {
	var sb strings.Builder
	sep := ""
	for _, in := range pat {
		sb.WriteString(sep)
		sep = " "
		sb.WriteString(in.String())
	}
	return sb.String()
}

func (pat Pattern) Display() string {
	var sb strings.Builder
	for i, in := range pat {
		if in.op == branch || in.op == jump {
			in.jump += int16(i)
			in.alt += int16(i)
		}
		fmt.Fprintf(&sb, "%d: %s\n", i, in.String())
	}
	return sb.String()
}

var ops = map[byte]string{
	dot:           ".",
	left:          "Left",
	right:         "Right",
	backref:       "\\",
	jump:          "Jump",
	branch:        "Branch",
	startOfLine:   "^",
	endOfLine:     "$",
	startOfString: "\\A",
	endOfString:   "\\Z",
	startOfWord:   "\\<",
	endOfWord:     "\\>",
	bitSet:        "[...]",
}

func (in inst) String() string {
	s := ops[in.op]
	switch in.op {
	case backrefIgnore:
		s = "i\\"
		fallthrough
	case left, right, backref:
		return s + strconv.Itoa(int(in.i))
	case charsIgnore:
		s = "i"
		fallthrough
	case chars:
		return s + "'" + in.data + "'"
	case listSet:
		return "[" + in.data + "]"
	case jump:
		return s + "(" + strconv.Itoa(int(in.jump)) + ")"
	case branch:
		return s + "(" + strconv.Itoa(int(in.jump)) + ", " + strconv.Itoa(int(in.alt)) + ")"
	default:
		return s
	}
}

type alternate struct {
	pi  int
	si  int
	si2 int
}

// Matches returns whether or not a pattern matches a string
func (pat Pattern) Matches(s string) bool {
	var result Result
	return pat.FirstMatch(s, 0, &result) != -1
}

// FirstMatch finds the first match in the string at or after pos.
// Returns the position if a match is found, else -1.
func (pat Pattern) FirstMatch(s string, pos int, result *Result) int {
	return pat.match(s, pos, +1, result)
}

// LastMatch finds the first match in the string before pos.
// Returns the position if a match is found, else -1.
func (pat Pattern) LastMatch(s string, pos int, result *Result) int {
	return pat.match(s, pos, -1, result)
}

// ForEachMatch calls action for each non-overlapping match in the string.
// The action should return true to continue, false to stop.
func (pat Pattern) ForEachMatch(s string, action func(*Result) bool) {
	var result Result
	pos := 0
	for {
		i := pat.match(s, pos, +1, &result)
		if i == -1 || !action(&result) {
			break
		}
		pos = ord.Max(result[0].end, i+1)
	}
}

const maxAlt = 100
const loopLimit = 1000

var repeat = inst{op: branch, jump: -1, alt: 1}

// match searches for a match.
// incdec should be +1 to search forward, -1 to search backward,
// or 0 to only try at the given position.
// It returns the position of the first match, or -1 if no match found.
func (pat Pattern) match(s string, pos, incdec int, result *Result) int {
	var alts [maxAlt]alternate
	var tmp [maxResult]int
outer:
	for ; 0 <= pos && pos <= len(s); pos += incdec {
		loops := 0
		ai := 0
		si := pos
		first := 0 // used to identify first non-left pattern element
		for pi := 0; pi < len(pat); pi++ {
			m := true
			in := &pat[pi]
			trace("si", si, "pi", pi, in)
			switch in.op {
			case dot:
				if pi+1 < len(pat) && pat[pi+1] == repeat && si < len(s)-1 {
					// for .* or .+ shortcut looping to end of line
					if ai >= len(alts) {
						panic("regex: too many branches")
					}
					trace("shortcut to end")
					alts[ai].pi = pi + 2
					alts[ai].si = si + 1
					j := -1
					if si < len(s) {
						j = strings.IndexAny(s[si:], "\r\n")
					}
					if j == -1 {
						si = len(s)
					} else {
						si += j
					}
					alts[ai].si2 = si
					ai++
					pi++
				} else {
					m = si < len(s) && s[si] != '\r' && s[si] != '\n'
					si++
				}
			case chars:
				if pi == first && incdec == +1 && ai == 0 {
					j := strings.Index(s[si:], in.data)
					if j == -1 {
						return -1 // chars don't exist, give up
					}
					if j > 0 {
						// skip ahead and restart match where chars are
						pos += j
						si = pos
						pi = -1 // -1 because loop increments
						continue
					}
					// else j == 0, m stays true
				} else {
					m = si+len(in.data) <= len(s) && strings.HasPrefix(s[si:], in.data)
				}
				si += len(in.data)
			case charsIgnore:
				m = si+len(in.data) <= len(s) && hasPrefixIgnore(s[si:], in.data)
				si += len(in.data)
			case listSet:
				m = si < len(s) && -1 != strings.IndexByte(in.data, s[si])
				si++
			case bitSet:
				m = false
				if si < len(s) {
					c := s[si]
					m = in.data[c>>3]&(1<<(c&7)) != 0
				}
				si++
			case branch:
				if ai > 0 && alts[ai-1].pi == pi+int(in.alt) && si == alts[ai-1].si2+1 {
					alts[ai-1].si2++ // expand existing entry, avoid stack growth
				} else {
					if ai >= len(alts) {
						panic("regex: too many branches")
					}
					alts[ai].pi = pi + int(in.alt)
					alts[ai].si = si
					alts[ai].si2 = si
					ai++
				}
				fallthrough
			case jump:
				pi += int(in.jump) - 1 // -1 because loop increments
			case left:
				if pi == first {
					first++
				}
				if in.i < maxResult {
					tmp[in.i] = si
				}
			case right:
				if in.i < maxResult {
					result[in.i].pos1 = tmp[in.i] + 1
					result[in.i].end = si
				}
			case backref:
				m, si = backrefMatch(s, si, result[in.i], strings.HasPrefix)
			case backrefIgnore:
				m, si = backrefMatch(s, si, result[in.i], hasPrefixIgnore)
			case startOfLine:
				m = si == 0 || (si <= len(s) && s[si-1] == '\n')
				if !m && pi == first && incdec == +1 && ai == 0 {
					j := strings.IndexByte(s[si:], '\n')
					if j == -1 {
						return -1
					}
					pos = si + j // skip ahead
				}
			case endOfLine:
				m = si >= len(s) || s[si] == '\r' ||
					(s[si] == '\n' && (si < 1 || s[si-1] != '\r'))
			case startOfString:
				m = si == 0
				if !m && pi == first && incdec == +1 && ai == 0 {
					return -1
				}
			case endOfString:
				m = si >= len(s)
			case startOfWord:
				m = si == 0 || (si <= len(s) && !matchSet(word.data, s[si-1]))
			case endOfWord:
				m = si >= len(s) || !matchSet(word.data, s[si])
			default:
				panic("bad regex pattern op code")
			}
			if !m {
				if ai > 0 {
					// backtrack
					prev := pi
					pi = alts[ai-1].pi - 1 // -1 because loop increments
					if alts[ai-1].si2 > alts[ai-1].si {
						trace("backtrack shorten")
						si = alts[ai-1].si2
						alts[ai-1].si2--
					} else {
						ai--
						trace("backtrack pop", "si", si, "alt", alts[ai])
						si = alts[ai].si
					}
					if pi < prev {
						if loops++; loops > loopLimit {
							panic("regex: too many loops")
						}
					}
				} else if incdec != 0 {
					continue outer
				} else {
					break outer
				}
			}
		}
		return pos // matched to end of pattern
	}
	return -1 // didn't match at any position
}

func trace(args ...any) {
	// fmt.Println(args...)
}

// hasPrefixIgnore returns whether s has pre as a prefix
// WARNING: s must be as long as pre
func hasPrefixIgnore(s, pre string) bool {
	for i := 0; i < len(pre); i++ {
		if ascii.ToLower(s[i]) != ascii.ToLower(pre[i]) {
			return false
		}
	}
	return true
}

func backrefMatch(s string, si int, p part, fn func(s, p string) bool) (bool, int) {
	if p.pos1 == 0 || p.end == -1 || p.end < p.pos1-1 {
		return false, si
	}
	pos := p.pos1 - 1
	end := ord.Min(p.end, len(s))
	b := ""
	if pos < end && pos >= 0 && end <= len(s) {
		b = s[pos:end]
	}
	return si+len(b) <= len(s) && fn(s[si:], b), si + len(b)
}

// Literal returns the string,true if the pattern is just a literal string,
// otherwise it returns "",false
func (pat Pattern) Literal() (string, bool) {
	if len(pat) == 2 && pat[0] == left0 && pat[1] == right0 {
		return "", true
	}
	if len(pat) == 3 && pat[0] == left0 && pat[1].op == chars && pat[2] == right0 {
		return pat[1].data, true
	}
	return "", false
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
