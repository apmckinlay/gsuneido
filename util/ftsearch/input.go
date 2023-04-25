// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ftsearch

import (
	"github.com/apmckinlay/gsuneido/util/ascii"
	snowballeng "github.com/kljensen/snowball/english"
)

type input struct {
	text string
	pos  int
}

func newInput(text string) *input {
	return &input{text, 0}
}

const maxWordLength = 32
const maxNumberLength = 16

func (src *input) Next() string {
	for {
		// skip non-letter
		for src.pos < len(src.text) && !ascii.IsLetter(src.text[src.pos]) &&
		    !ascii.IsDigit(src.text[src.pos]) {
			src.pos++
		}
		if src.pos >= len(src.text) {
            return ""
        }
		pos := src.pos
		if ascii.IsLetter(src.text[src.pos]) {
			for src.pos < len(src.text) && ascii.IsLetter(src.text[src.pos]) {
				src.pos++
			}
			tok := src.text[pos:src.pos]
			if len(tok) > 1 && len(tok) < maxWordLength {
				if _, ok := stopWords[tok]; !ok {
					return snowballeng.Stem(tok, true)
				}
			}
		} else if ascii.IsDigit(src.text[src.pos]) {
			for src.pos < len(src.text) && ascii.IsDigit(src.text[src.pos]) {
				src.pos++
			}
			tok := src.text[pos:src.pos]
			if len(tok) >= 2 && len(tok) < maxNumberLength {
				return tok
			}
		}
	}
}

// stopWords is copied from snowball
var stopWords = map[string]struct{}{
	"about":      {},
	"again":      {},
	"against":    {},
	"all":        {},
	"am":         {},
	"an":         {},
	"and":        {},
	"any":        {},
	"are":        {},
	"as":         {},
	"at":         {},
	"be":         {},
	"because":    {},
	"been":       {},
	"being":      {},
	"between":    {},
	"both":       {},
	"but":        {},
	"by":         {},
	"can":        {},
	"did":        {},
	"do":         {},
	"does":       {},
	"doing":      {},
	"don":        {},
	"during":     {},
	"each":       {},
	"few":        {},
	"for":        {},
	"from":       {},
	"further":    {},
	"had":        {},
	"has":        {},
	"have":       {},
	"having":     {},
	"he":         {},
	"her":        {},
	"here":       {},
	"hers":       {},
	"herself":    {},
	"him":        {},
	"himself":    {},
	"his":        {},
	"how":        {},
	"if":         {},
	"in":         {},
	"into":       {},
	"is":         {},
	"it":         {},
	"its":        {},
	"itself":     {},
	"just":       {},
	"me":         {},
	"more":       {},
	"most":       {},
	"my":         {},
	"myself":     {},
	"no":         {},
	"nor":        {},
	"not":        {},
	"now":        {},
	"of":         {},
	"once":       {},
	"only":       {},
	"or":         {},
	"other":      {},
	"our":        {},
	"ours":       {},
	"ourselves":  {},
	"out":        {},
	"own":        {},
	"same":       {},
	"she":        {},
	"should":     {},
	"so":         {},
	"some":       {},
	"such":       {},
	"than":       {},
	"that":       {},
	"the":        {},
	"their":      {},
	"theirs":     {},
	"them":       {},
	"themselves": {},
	"then":       {},
	"there":      {},
	"these":      {},
	"they":       {},
	"this":       {},
	"those":      {},
	"through":    {},
	"to":         {},
	"too":        {},
	"until":      {},
	"very":       {},
	"was":        {},
	"we":         {},
	"were":       {},
	"what":       {},
	"when":       {},
	"where":      {},
	"which":      {},
	"while":      {},
	"who":        {},
	"whom":       {},
	"why":        {},
	"will":       {},
	"with":       {},
	"you":        {},
	"your":       {},
	"yours":      {},
	"yourself":   {},
	"yourselves": {},
}
