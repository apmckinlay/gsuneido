// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
/*
CombyMatch is
	- token-based
	- structural pattern matching (Skips non-structural tokens like whitespace, comments
		and newlines)
	- where template holes can (:[name]) lazily capture text between literal token matches
		and respect bracket/delimiter nesting
The idea is based on https://comby.dev/
*/
class
	{
	MAX_HOLE_MATCHES: 100
	CallClass(search, s, pos = 0, prev = false)
		{
		items = CombyTemplate(search)
		if items.Size() is 0
			return false

		tokens = .buildTokens(s)
		env = Object(:s, :tokens, :items, holes: Object())

		return prev is false
			? .findNext(env, pos)
			: .findPrev(env, pos)
		}

	buildTokens(s)
		{
		tokens = Object()
		scan = Scanner(s)
		start = 0
		while scan isnt type = scan.Next2()
			{
			tokens.Add(Object(:start, end: scan.Position(), :type,
				value: scan.Value()))
			start = scan.Position()
			}
		return tokens
		}

	findNext(env, pos)
		{
		tokens = env.tokens
		i = 0
		while i < tokens.Size() and tokens[i].start < pos
			i++
		while i < tokens.Size()
			{
			env.holes = Object()
			if false isnt next = .match(env, i, 0)
				return Object(pos: tokens[i].start, end: tokens[next - 1].end,
					holes: env.holes)
			i++
			}
		return false
		}

	findPrev(env, pos)
		{
		tokens = env.tokens
		i = tokens.Size() - 1
		while i >= 0 and tokens[i].end > pos
			i--
		while i >= 0
			{
			env.holes = Object()
			if false isnt next = .match(env, i, 0)
				{
				if tokens[next - 1].end <= pos
					return Object(pos: tokens[i].start,
						end: tokens[next - 1].end,
						holes: env.holes)
				}
			i--
			}
		return false
		}

	All(search, s)
		{
		items = CombyTemplate(search)
		if items.Size() is 0
			return #()

		tokens = .buildTokens(s)
		i = 0
		results = Object()
		env = Object(:s, :tokens, :items, holes: Object())
		while false isnt next = .findNext(env, i)
			{
			results.Add(next)
			i = next.end
			}
		return results
		}

	match(env, tokenStartIdx, itemStartIdx)
		{
		if itemStartIdx >= env.items.Size()
			return false

		i = tokenStartIdx
		j = itemStartIdx
		while j < env.items.Size()
			{
			if false is advance = .matchItem(env, i, j)
				return false
			i += advance[0]
			j += advance[1]
			}
		return i
		}

	matchItem(env, tokenStartIdx, itemIdx)
		{
		if tokenStartIdx >= env.tokens.Size()
			return false

		item = env.items[itemIdx]

		tokenIdx = tokenStartIdx

		if item.type is #COMMENT
			{
			return env.tokens[tokenIdx].type is #COMMENT ? Object(1, 1) : false
			}

		if item.type is #WHITESPACE
			return .matchWhiteSpace(env, tokenStartIdx, tokenIdx)

		if item.type is #HOLE
			{
			if not env.items.Member?(itemIdx + 1)
				return .matchLastHole(env, item, tokenStartIdx, tokenIdx,
					expr?: item.GetDefault(#expr, false))
			else
				return .matchHole(env, item, tokenStartIdx, tokenIdx, itemIdx,
					expr?: item.GetDefault(#expr, false))
			}
		return item.value is env.tokens[tokenStartIdx].value ? Object(1, 1) : false
		}

	matchWhiteSpace(env, tokenStartIdx, tokenIdx)
		{
		while tokenIdx < env.tokens.Size() and
			env.tokens[tokenIdx].type in (#COMMENT, #WHITESPACE, #NEWLINE)
			tokenIdx++
		return Object(tokenIdx - tokenStartIdx, 1)
		}

	// match until meet the end of source or the close delimiter
	matchLastHole(env, item, tokenStartIdx, tokenIdx, expr? = false)
		{
		blockEnd = .findEndOfBlock(env, tokenIdx, :expr?)
		if blockEnd is tokenIdx
			return false
		env.holes[item.value] = env.s[env.tokens[tokenStartIdx].start..
			env.tokens[blockEnd - 1].end]
		return Object(blockEnd - tokenStartIdx, 1)
		}

	matchHole(env, item, tokenStartIdx, tokenIdx, itemIdx, expr? = false)
		{
		match = 0
		.findEndOfBlock(env, tokenIdx, :expr?)
			{ |blockLevel, ti|
			if match >= .MAX_HOLE_MATCHES
				return false
			if blockLevel is 0
				{
				match++
				if false isnt next = .match(env, ti, itemIdx + 1)
					{
					env.holes[item.value] =  env.s[env.tokens[tokenStartIdx].start..
						env.tokens[ti - 1].end]
					return Object(next - tokenStartIdx, env.items.Size() - itemIdx)
					}
				}
			}
		return false
		}

	findEndOfBlock(env, tokenIdx, expr? = false, block = false)
		{
		blockLevel = 0
		tok = false
		while tokenIdx < env.tokens.Size()
			{
			if block isnt false
				block(blockLevel, tokenIdx)

			prevTok = tok
			tok = env.tokens[tokenIdx]
			if tok.type is ''
				{
				if tok.value in ('(', '[', '{')
					blockLevel++
				else if tok.value in (')', ']', '}')
					{
					blockLevel--
					if blockLevel < 0
						break
					}
				}
			else if .endExpr?(expr?, blockLevel, tok, prevTok)
				break
			tokenIdx++
			}
		return tokenIdx
		}

	endExpr?(expr?, blockLevel, tok, prevTok)
		{
		if expr? isnt true or blockLevel isnt 0
			return false
		if tok.type in (#WHITESPACE, #COMMENT)
			return true
		if tok.type is #NEWLINE
			return not (prevTok isnt false and prevTok.type is '' and
				prevTok.value is '.')
		return false
		}

	GetHint(search)
		{
		hint = ''
		for item in CombyTemplate(search)
			{
			if ((item.type is #IDENTIFIER and item.keyword? is false or
				item.type is #NUMBER) and item.text.Size() > hint.Size())
				hint = item.text
			}
		return hint is '' ? false : hint
		}
	}
