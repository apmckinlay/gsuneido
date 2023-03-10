// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import (
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Captures [20]int32 // 2 * 10 (\0 to \9)

// Matches returns whether the pattern is found anywhere in the string
func (pat Pattern) Matches(s string) bool {
	// keep the .* prefix, toEnd = false
	return pat.match(s, nil, false)
}

// Match returns whether the pattern matches the entire string.
func (pat Pattern) Match(s string, cap *Captures) bool {
	// omit the .* prefix, toEnd = true
	return omitUA(pat).match(s, cap, true)
}

func (pat Pattern) FirstMatch(s string, cap *Captures) bool {
	// keep the .* prefix, toEnd = false
	return pat.match(s, cap, false)
}

func (pat Pattern) LastMatch(s string, cap *Captures) bool {
	// inefficient, but rarely used
	// omit the .* prefix, toEnd = false
	pat = omitUA(pat)
	for i := len(s) - 1; i >= 0; i-- {
		if pat.match(s[:i], cap, false) {
			return true
		}
	}
	return false
	// could improve this by figuring out the minimum match length
}

func omitUA(pat Pattern) Pattern {
	return Pattern(strings.TrimPrefix(string(pat), uaString))
}

type state struct {
	pi  int16
	cap *Captures
}

// match looks for a match anywhere in the string i.e. not anchored.
// If toEnd is true, the match must go to the end of the string.
// It returns true if a match was found, and false otherwise.
// If it returns true, the captures are updated.
func (pat Pattern) match(s string, cap *Captures, toEnd bool) bool {
	_ = t && trace.Println(pat)
	cap2 := dup(cap)
	anchored := opType(pat[0]) != opUnanchored
	piStart := int16(0)
	if !anchored {
		piStart = 1
	}
	prefix := ""
	if opType(pat[piStart]) == opLitPrefix {
		n := int16(pat[piStart+1])
		if !anchored {
			prefix = string(pat[piStart+2 : piStart+2+n])
		}
		piStart += 2 + n
	}
	switch opType(pat[piStart]) {
	case opOnePass:
		piStart++
		if anchored {
			return Pattern(pat[piStart:]).onePass(s, cap, toEnd)
		}
	case opLiteral:
		return pat.literalMatch(s, cap, toEnd)
	}
	var cur []state
	var next []state
	var live = &BitSet{}
	matched := false
	for si := 0; si <= len(s); si++ {
		if si < len(s) {
			_ = t && trace.Println("--- si:", si, "c:", s[si:si+1])
		} else {
			_ = t && trace.Println("at end of string")
		}
		if len(cur) == 0 {
			if anchored && si > 0 {
				return matched
			}
			if matched {
				return true // finished exploring alternatives
			}
			if len(prefix) > 0 {
				i := strings.Index(s[si:], prefix)
				if i < 0 {
					return false
				}
				_ = t && trace.Println("skip from", si, "to", si+i)
				si += i
			}
		}
		if !matched {
			if cap != nil {
				cap2[0] = int32(si) // Save 0
			}
			cur = pat.addstate(s, si, live, cur, piStart, cap2)
			live.Clear()
		}
		for i, c := range cur {
			_ = t && trace.Printf("state [%v] %v\n", i, pat.opstr1(c.pi))
		}
		for ci := 0; ci < len(cur); ci++ { // for each state
			pi := cur[ci].pi
			_ = t && trace.Printf("[%v] %v\n", ci, pat.opstr1(pi))
			add := int16(0)
			switch opType(pat[pi]) {
			case opChar:
				if si < len(s) && s[si] == pat[pi+1] {
					_ = t && trace.Println("YES")
					add = pi + 2
				}
			case opCharIgnoreCase:
				if si < len(s) && ascii.ToLower(s[si]) == pat[pi+1] {
					add = pi + 2
				}
			case opAny:
				if si < len(s) {
					add = pi + 1
				}
			case opAnyNotNL:
				if si < len(s) && s[si] != '\r' && s[si] != '\n' {
					add = pi + 1
				}
			case opListSet:
				if si < len(s) {
					n := int16(pat[pi+1])
					if -1 != strings.IndexByte(string(pat[pi+2:pi+2+n]), s[si]) {
						add = pi + 2 + n
					}
				}
			case opHalfSet:
				if si < len(s) && matchHalfSet(pat[pi+1:], s[si]) {
					add = pi + 1 + 16
				}
			case opFullSet:
				if si < len(s) && matchFullSet(pat[pi+1:], s[si]) {
					add = pi + 1 + 32
				}
			case opDoneSave1:
				if toEnd && si < len(s) {
					break
				}
				if cap == nil {
					// if not capturing, any match will do
					return true
				}
				if !matched || int(cap[1]) < si {
					cur[ci].cap[1] = int32(si)
					*cap = *cur[ci].cap
				}
				cur = cur[:0] // cut off lower priority threads
				matched = true
			default:
				panic(assert.ShouldNotReachHere())
			}
			if add > 0 {
				next = pat.addstate(s, si+1, live, next, add, cur[ci].cap)
			}
		}
		cur, next = next, cur // swap
		next = next[:0]       // clear
		live.Clear()
	}
	return matched
}

// addstate adds a state and, recursively, all of its children.
// It processes all zero width instructions
// so the states added will point to character matching instructions.
func (pat Pattern) addstate(s string, si int, live *BitSet, states []state,
	pi int16, cap *Captures) []state {
	for {
		if !live.AddNew(pi) {
			return states
		}
		_ = t && trace.Println("addstate loop", pat.opstr1(pi))
		switch opType(pat[pi]) {
		case opJump:
			jmp := int16(pat[pi+1])<<8 | int16(pat[pi+2])
			pi += jmp
		case opSplitFirst:
			jmp := int16(pat[pi+1])<<8 | int16(pat[pi+2])
			states = pat.addstate(s, si, live, states, pi+jmp, cap)
			pi += 3
		case opSplitLast:
			jmp := int16(pat[pi+1])<<8 | int16(pat[pi+2])
			states = pat.addstate(s, si, live, states, pi+3, cap)
			pi += jmp
		case opSave:
			if cap != nil {
				c := pat[pi+1]
				orig := cap[c]
				cap[c] = int32(si)
				states = pat.addstate(s, si, live, states, pi+2, cap)
				cap[c] = orig
				return states
			}
			pi += 2
		case opStrStart, opStrEnd, opLineStart, opLineEnd, opWordStart, opWordEnd:
			if !boundary(s, si, pat[pi]) {
				return states
			}
			_ = t && trace.Println("YES")
			pi++
		default:
			states = append(states, state{pi: pi, cap: dup(cap)})
			return states
		}
	}
}

func boundary(s string, si int, op byte) bool {
	switch opType(op) {
	case opStrStart:
		return si == 0
	case opStrEnd:
		return si >= len(s)
	case opLineStart:
		return si == 0 || s[si-1] == '\n'
	case opLineEnd:
		return si >= len(s) || s[si] == '\r' ||
			(s[si] == '\n' && (si == 0 || s[si-1] != '\r'))
	case opWordStart:
		return si == 0 || (si <= len(s) && !matchFullSet(wordSet, s[si-1]))
	case opWordEnd:
		return si >= len(s) || !matchFullSet(wordSet, s[si])
	}
	panic(assert.ShouldNotReachHere())
}

func dup(cap *Captures) *Captures {
	if cap == nil {
		return nil
	}
	cp := *cap
	return &cp
}

// ------------------------------------------------------------------

func (pat Pattern) onePass(s string, cap *Captures, toEnd bool) bool {
	_ = t && trace.Println("ONE PASS")
	for si, pi := 0, 0; pi < len(pat); pi++ {
		_ = t && trace.Printf("si %v %q %v\n", si, str.Subn(s, si, 1), pat.opstr1(int16(pi)))
		switch opType(pat[pi]) {
		case opChar:
			if si >= len(s) || s[si] != pat[pi+1] {
				return false
			}
			pi++
			si++
		case opCharIgnoreCase:
			if si >= len(s) || ascii.ToLower(s[si]) != pat[pi+1] {
				return false
			}
			pi++
			si++
		case opAny:
			if si >= len(s) {
				return false
			}
			si++
		case opAnyNotNL:
			if si >= len(s) || s[si] == '\r' || s[si] == '\n' {
				return false
			}
			si++
		case opListSet:
			n := int(pat[pi+1])
			if si >= len(s) ||
				-1 == strings.IndexByte(string(pat[pi+2:pi+2+n]), s[si]) {
				return false
			}
			pi += 1 + n
			si++
		case opHalfSet:
			if si >= len(s) || !matchHalfSet(pat[pi+1:], s[si]) {
				return false
			}
			pi += 16
			si++
		case opFullSet:
			if si >= len(s) && !matchFullSet(pat[pi+1:], s[si]) {
				return false
			}
			pi += 32
			si++
		case opStrStart, opStrEnd, opLineStart, opLineEnd, opWordStart, opWordEnd:
			if !boundary(s, si, pat[pi]) {
				return false
			}
		case opSave:
			if cap != nil {
				c := pat[pi+1]
				cap[c] = int32(si)
			}
			pi++
		case opDoneSave1:
			if toEnd && si < len(s) {
				return false
			}
			if cap != nil {
				cap[1] = int32(si)
			}
			return true
		default:
			panic(assert.ShouldNotReachHere())
		}
	}
	return false
}

// ------------------------------------------------------------------

func (pat Pattern) literalMatch(s string, cap *Captures, toEnd bool) bool {
	_ = t && trace.Println("LITERAL")
	lit := string(pat[1:])
	anchored := true
	if opType(pat[0]) == opUnanchored {
		anchored = false
		lit = lit[1:]
	}
	i := 0
	if anchored {
		if toEnd {
			if s != lit {
				return false
			}
		} else {
			if !strings.HasPrefix(s, lit) {
				return false
			}
		}
	}
	// else not anchored
	i = strings.Index(s, lit)
	if i < 0 {
		return false
	}
	if cap != nil {
		cap[0], cap[1] = int32(i), int32(i+len(lit))
	}
	return true
}

// ------------------------------------------------------------------

const t = false

type tracer struct{}

var trace tracer

func (tracer) Println(args ...any) bool {
	fmt.Println(args...)
	return true
}

func (tracer) Printf(format string, args ...any) bool {
	fmt.Printf(format, args...)
	return true
}
