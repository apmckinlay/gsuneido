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

// Matches is a shortcut for Match(s, nil)
func (pat Pattern) Matches(s string) bool {
	return pat.Match(s, nil)
}

// LastMatch finds the last match in s.
// It is inefficient, but rarely used.
func (pat Pattern) LastMatch(s string, cap *Captures) bool {
	for i := len(s); i >= 0; i-- {
		if pat.Match(s[:i], cap) {
			return true
		}
	}
	return false
}

type state struct {
	pi  int16
	cap *Captures
}

// match looks for a match anywhere in the string i.e. not anchored.
// If toEnd is true, the match must go to the end of the string.
// It returns true if a match was found, and false otherwise.
// If it returns true, the captures are updated.
func (pat Pattern) Match(s string, cap *Captures) bool {
	_ = t && trace.Println(pat)
	piStart := int16(0)
	prefix := ""
	switch op := opType(pat[piStart]); op {
	case opPrefix:
		n := int16(pat[piStart+1])
		prefix = string(pat[piStart+2 : piStart+2+n])
		piStart += 2 + n
	case opOnePass:
		return Pattern(pat[piStart+1:]).onePass(s, cap)
	case opLiteralSubstr, opLiteralPrefix, opLiteralSuffix, opLiteralEqual:
		return pat.literalMatch(op, s, cap)
	}
	anchored := opType(pat[piStart]) == opStrStart
	cap2 := dup(cap)
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
		case opSplitJump:
			jmp := int16(pat[pi+1])<<8 | int16(pat[pi+2])
			states = pat.addstate(s, si, live, states, pi+jmp, cap) // RECURSE
			pi += 3
		case opSplitNext:
			jmp := int16(pat[pi+1])<<8 | int16(pat[pi+2])
			states = pat.addstate(s, si, live, states, pi+3, cap) // RECURSE
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

var wordSet = Pattern(word[:])

func dup(cap *Captures) *Captures {
	if cap == nil {
		return nil
	}
	cp := *cap
	return &cp
}

// ------------------------------------------------------------------

func (pat Pattern) onePass(s string, cap *Captures) bool {
	_ = t && trace.Println(">>> one pass")
	cap2 := dup(cap)
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
			if si >= len(s) || !matchFullSet(pat[pi+1:], s[si]) {
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
				cap2[c] = int32(si)
			}
			pi++
		case opDoneSave1:
			if cap != nil {
				*cap = *cap2
				cap[0] = 0
				cap[1] = int32(si)
			}
			return true
		case opJump:
			jmp := int(int16(pat[pi+1])<<8 | int16(pat[pi+2]))
			pi += jmp - 1 // -1 because loop increments
		case opBranchNext:
			if onePass1(pat, pi+3, s, si) {
				pi += 2
			} else {
				jmp := int(int16(pat[pi+1])<<8 | int16(pat[pi+2]))
				pi += jmp - 1 // -1 because loop increments
			}
		case opBranchJump:
			jmp := int(int16(pat[pi+1])<<8 | int16(pat[pi+2]))
			if onePass1(pat, pi+jmp, s, si) {
				pi += jmp - 1 // -1 because loop increments
			} else {
				pi += 2
			}
		default:
			panic(assert.ShouldNotReachHere())
		}
	}
	return false
}

func onePass1(pat Pattern, pi int, s string, si int) bool {
	for ; opType(pat[pi]) == opSave; pi += 2 {
	}
	switch opType(pat[pi]) {
	case opChar:
		return si < len(s) && s[si] == pat[pi+1]
	case opCharIgnoreCase:
		return si < len(s) && ascii.ToLower(s[si]) == pat[pi+1]
	case opListSet:
		n := int(pat[pi+1])
		return si < len(s) &&
			-1 != strings.IndexByte(string(pat[pi+2:pi+2+n]), s[si])
	case opHalfSet:
		return si < len(s) && matchHalfSet(pat[pi+1:], s[si])
	case opFullSet:
		return si < len(s) && matchFullSet(pat[pi+1:], s[si])
	case opDoneSave1:
		return si >= len(s)
	}
	panic(assert.ShouldNotReachHere())
}

// ------------------------------------------------------------------

func (pat Pattern) literalMatch(op opType, s string, cap *Captures) bool {
	_ = t && trace.Println(">>>", op)
	lit := string(pat[1:])
	i := 0
	switch op {
	case opLiteralEqual:
		if s != lit {
			return false
		}
	case opLiteralSubstr:
		i = strings.Index(s, lit)
		if i < 0 {
			return false
		}
	case opLiteralPrefix:
		if !strings.HasPrefix(s, lit) {
			return false
		}
	case opLiteralSuffix:
		if !strings.HasSuffix(s, lit) {
			return false
		}
		i = len(s) - len(lit)
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
