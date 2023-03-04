// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type Captures [20]int32 // 10 * 2

const npre = 7

// Matches returns whether the pattern is found anywhere in the string
func (pat Pattern) Matches(s string) bool {
	// keep the .* prefix, toEnd = false
	return optPre(pat, s).prefixMatch(s, nil, false)
}

// Match returns whether the pattern matches the entire string.
func (pat Pattern) Match(s string, cap *Captures) bool {
	// omit the .* prefix, toEnd = true
	return Pattern(pat[npre:]).prefixMatch(s, cap, true)
}

func (pat Pattern) FirstMatch(s string, cap *Captures) bool {
	// keep the .* prefix, toEnd = false
	return optPre(pat, s).prefixMatch(s, cap, false)
}

func (pat Pattern) LastMatch(s string, cap *Captures) bool {
	// omit the .* prefix, toEnd = false
	pat = pat[npre:]
	for i := len(s) - 1; i >= 0; i-- {
		if pat.prefixMatch(s[:i], cap, false) {
			return true
		}
	}
	return false
	// could improve this by figuring out the minimum match length
}

// optPre removes the .* prefix if the pattern starts with \A
// or if it starts with ^ and the string is short and doesn't contain \n
func optPre(pat Pattern, s string) Pattern {
	// +2 is to skip Save 0
	if pat[npre+2] == byte(opStrStart) || (pat[npre+2] == byte(opLineStart) &&
		len(s) < 1000 && !strings.Contains(s, "\n")) {
		pat = pat[npre:]
	}
	return pat
}

type state struct {
	pi  int16
	cap *Captures
}

// prefixMatch looks for a match starting at the beginning of the string.
// If toEnd is true, the match must go to the end of the string.
// It return true if a match was found, and false otherwise.
// If it returns true, the captures are updated.
func (pat Pattern) prefixMatch(s string, cap *Captures, toEnd bool) bool {
	trace.Println(pat)
	var cur []state
	var next []state
	var ss = &SparseSet{}
	cur = pat.addstate(s, 0, ss, cur, 0, dup(cap))
	ss.Clear()
	matched := false
	for si := 0; ; si++ {
		if si < len(s) {
			trace.Println("--- si:", si, "c:", s[si:si+1])
		} else {
			trace.Println("at end of string")
		}
		// for i, c := range cur {
		// 	trace.Printf("state [%v] %v\n", i, pat.opstr1(c.pi))
		// }
		for ci := 0; ci < len(cur); ci++ { // for each state
			pi := cur[ci].pi
			trace.Printf("[%v] %v\n", ci, pat.opstr1(pi))
			add := int16(0)
			switch opType(pat[pi]) {
			case opChar:
				if si < len(s) && s[si] == pat[pi+1] {
					trace.Println("YES")
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
			case opDone:
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
				next = pat.addstate(s, si+1, ss, next, add, cur[ci].cap)
			}
		}
		if len(next) == 0 {
			return matched
		}
		cur, next = next, cur // swap
		next = next[:0]       // clear
		ss.Clear()
	}
}

// addstate adds a state and, recursively, all of its children.
// It processes all zero width instructions
// so the states added will point to character matching instructions.
func (pat Pattern) addstate(s string, si int, ss *SparseSet, states []state,
	pi int16, cap *Captures) []state {
	// trace.Println("addstate ss", ss.dense)
	for {
		if !ss.AddNew(pi) {
			return states
		}
		trace.Println("addstate loop", pat.opstr1(pi))
		switch opType(pat[pi]) {
		case opJump:
			jmp := int16(pat[pi+1])<<8 | int16(pat[pi+2])
			pi += jmp
		case opSplitFirst:
			jmp := int16(pat[pi+1])<<8 | int16(pat[pi+2])
			states = pat.addstate(s, si, ss, states, pi+jmp, cap)
			pi += 3
		case opSplitLast:
			jmp := int16(pat[pi+1])<<8 | int16(pat[pi+2])
			states = pat.addstate(s, si, ss, states, pi+3, cap)
			pi += jmp
		case opSave:
			if cap != nil {
				c := pat[pi+1]
				orig := cap[c]
				cap[c] = int32(si)
				states = pat.addstate(s, si, ss, states, pi+2, cap)
				cap[c] = orig
				return states
			}
			pi += 2
		case opStrStart, opStrEnd, opLineStart, opLineEnd, opWordStart, opWordEnd:
			if !boundary(s, si, pat[pi]) {
				return states
			}
			trace.Println("YES")
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
		return si == 0 || (si <= len(s) && !matchHalfSet(wordSet, s[si-1]))
	case opWordEnd:
		return si >= len(s) || !matchHalfSet(wordSet, s[si])
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

type tracer struct{}

var trace tracer

func (tracer) Println(args ...any) {
	// fmt.Println(args...)
}

func (tracer) Printf(format string, args ...any) {
	// fmt.Printf(format, args...)
}
