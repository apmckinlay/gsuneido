// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/ascii"
)

func Match(pat Pattern, s string) bool {
	var cur []int
	var next []int
	cur = addstate(cur, 0)
	for si := 0; len(cur) > 0; si++ {
		if si < len(s) {
			// fmt.Println("c:", s[si:si+1])
		} else {
			// fmt.Println("at end of string")
		}
		for ci := 0; ci < len(cur); ci++ { // for each state
			pi := cur[ci]
		loop:
			// fmt.Printf("[%v] %v: %v\n", ci, pi, pat.opstr1(pi))
			switch opType(pat[pi]) {
			case opChar:
				if si < len(s) && s[si] == pat[pi+1] {
					next = addstate(next, pi+2)
				}
			case opCharIgnoreCase:
				if si < len(s) && ascii.ToLower(s[si]) == pat[pi+1] {
					next = addstate(next, pi+2)
				}
			case opAny:
				if si < len(s) && s[si] != '\r' && s[si] != '\n' {
					next = addstate(next, pi+1)
				}
			case opListSet:
				if si < len(s) {
					n := int(pat[pi+1])
					if -1 != strings.IndexByte(string(pat[pi+2:pi+2+n]), s[si]) {
						next = addstate(next, pi+2+n)
					}
				}
			case opHalfSet:
				if si < len(s) && matchHalfSet(pat[pi+1:], s[si]) {
					next = addstate(next, pi+1+16)
				}
			case opFullSet:
				if si < len(s) && matchFullSet(pat[pi+1:], s[si]) {
					next = addstate(next, pi+1+32)
				}

			// zero width

			case opSave:
				//TODO
				pi += 2
				goto loop
			case opStrStart:
				if si == 0 {
					pi++
					goto loop
				}
			case opStrEnd:
				if si == len(s) {
					pi++
					goto loop
				}
			case opLineStart:
				if si == 0 || s[si-1] == '\n' {
					pi++
					goto loop
				}
			case opLineEnd:
				if si >= len(s) || s[si] == '\r' || s[si] == '\n' {
					pi++
					goto loop
				}
			case opWordStart:
				if si == 0 || (si <= len(s) && !matchHalfSet(wordSet, s[si-1])) {
					pi++
					goto loop
				}
			case opWordEnd:
				if si >= len(s) || !matchHalfSet(wordSet, s[si]) {
					pi++
					goto loop
				}
			case opJump:
				jmp := int16(pat[pi+1])<<8 | int16(pat[pi+2])
				pi += int(jmp)
				goto loop
			case opSplitFirst:
				jmp := int16(pat[pi+1])<<8 | int16(pat[pi+2])
				cur = addstate(cur, pi+int(jmp))
				pi += 3
				goto loop
			case opSplitLast:
				jmp := int16(pat[pi+1])<<8 | int16(pat[pi+2])
				cur = addstate(cur, pi+3)
				pi += int(jmp)
				goto loop
			case opStop:
				if si >= len(s) {
					return true
				}
			default:
				panic("not implemented")
			}
		}
		cur, next = next, cur // swap
		next = next[:0]       // clear
	}
	// fmt.Println("reached end of string")
	return false
}

func addstate(states []int, pi int) []int {
	for _, j := range states {
		if j == pi {
			return states
		}
	}
	return append(states, pi)
}
