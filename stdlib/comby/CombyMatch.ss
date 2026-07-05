// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(search, s)
		{
		items = CombyTemplate(search)
		if items.Size() is 0
			return #()

		tokens = Object()
		scan = Scanner(s)
		start = 0
		while scan isnt type = scan.Next2()
			{
			tokens.Add(Object(:start, end: scan.Position(), :type,
				value: scan.Value()))
			start = scan.Position()
			}

		i = 0
		results = Object()
		env = Object(:s, :tokens, :items, holes: Object())
		while i < tokens.Size()
			{
			if false is next = .match(env, i, 0)
				i++
			else
				{
				results.Add(Object(pos: tokens[i].start, end: tokens[next - 1].end,
					holes: env.holes))
				env.holes = Object()
				i = next
				}
			}
		return results
		}

	match(env, tokenStartIdx, itemStartIdx)
		{
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
				return .matchLastHole(env, item, tokenStartIdx, tokenIdx)
			else
				return .matchHole(env, item, tokenStartIdx, tokenIdx, itemIdx)
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
	matchLastHole(env, item, tokenStartIdx, tokenIdx)
		{
		blockLevel = 0
		while tokenIdx < env.tokens.Size()
			{
			if env.tokens[tokenIdx].value in ('(', '[', '{')
				blockLevel++
			else if env.tokens[tokenIdx].value in (')', ']', '}')
				{
				blockLevel--
				if blockLevel < 0
					break
				}
			tokenIdx++
			}
		env.holes[item.value] = env.s[env.tokens[tokenStartIdx].start..
			env.tokens[tokenIdx - 1].end]
		return Object(tokenIdx - tokenStartIdx, 1)
		}

	matchHole(env, item, tokenStartIdx, tokenIdx, itemIdx)
		{
		blockLevel = 0
		while tokenIdx < env.tokens.Size()
			{
			if blockLevel is 0 and
				false isnt next = .match(env, tokenIdx, itemIdx + 1)
				{
				env.holes[item.value] =  env.s[env.tokens[tokenStartIdx].start..
					env.tokens[tokenIdx - 1].end]
				return Object(next - tokenStartIdx, env.items.Size() - itemIdx)
				}
			if env.tokens[tokenIdx].value in ('(', '[', '{')
				blockLevel++
			else if env.tokens[tokenIdx].value in (')', ']', '}')
				{
				blockLevel--
				if blockLevel < 0
					break
				}
			tokenIdx++
			}
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
