// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// helper methods for find and replace
class
	{
	findPat(find, options)
		{
		if options.expr is true
			return .ast(find)
		else
			return .Regex(find, options)
		}
	NeedAst?(options)
		{
		return options.expr is true
		}
	isAstPattern?(findpat)
		{
		return not String?(findpat)
		}
	ast(find)
		{
		try
			return Tdop(find, type: 'expression')
		catch
			return ""
		}
	// options should be an object with regex, word, or case members
	Regex(find, options)
		{
		if find is ""
			return ""

		rx = find
		if options.regex is true
			rx = rx.Unescape() // allow \t, \n etc.
		else
			rx = "(?q)" $ rx $ "(?-q)"
		if options.word is true
			{
			if find[0].Alpha?() or find[0] is '_' or
				(options.regex is true and find[0] is '(')
				rx = "\<" $ rx
			if find[-1].AlphaNum?() or find[-1] is '_' or
				(options.regex is true and find[-1] is ')')
				rx = rx $ "\>"
			}
		if options.case isnt true
			rx = "(?i)" $ rx
		return rx
		}
	// options should be an object with regex member
	Replacement(replace, options)
		{
		if replace is ""
			return ""
		if options.regex is true or options.expr is true
			replace = replace.Unescape()
		else
			replace = `\=` $ replace
		return replace
		}

	DoFind(text, from, options, prev = false)
		{
		_context = Object(:text, ast: options.GetDefault(#ast, false))
		findpat = .findPat(options.find, options)
		if findpat is ""
			{ .beep(); return false }
		return .findNextPrev(from, findpat, :prev)
		}

	findNextPrev(i, findpat, prev = false)
		{
		match = .findOne(i, findpat, :prev)
		retryIndex = prev ? .searchText(0).Size() : 0
		if match is false
			{
			if i is retryIndex
				{ .beep(); return false }
			return .findNextPrev(retryIndex, findpat, :prev)
			}
		match = match[0]
		return Object(match[0], match[1])
		}

	findOne(i, findpat, prev = false)
		{
		if .isAstPattern?(findpat)
			{
			if false is  _context.ast
				return false
			match = TdopSearch(_context.ast.Root, findpat, i, :prev)
			}
		else
			{
			s = prev is false ? .searchText(i) : .searchText(0)[.. i]
			match = false
			.tryIgnoreRegexError({ match = s.Match(findpat, :prev) })
			if match isnt false and not prev
				match[0][0] += i
			}
		return match
		}

	searchText(i)
		{
		return _context.text[i..]
		}

	tryIgnoreRegexError(block)
		{
		try
			block()
		catch (err)
			if not err.Has?('regex')
				SuneidoLog('ERROR: (CAUGHT) Find failed: ' $ err)
		}

	// extract for testing
	beep()
		{
		Beep()
		}

	FindAll(text, options)
		{
		findpat = .findPat(options.find, options)
		if findpat is ""
			{ .beep(); return false }
		if false is s = .isAstPattern?(findpat) ? options.GetDefault(#ast, false) : text
			return false
		matches = Object()
		.tryIgnoreRegexError()
			{
			s.ForEachMatch(findpat)
				{|m|
				matches.Add(m[0].Project(#(0, 1)))
				}
			}
		return matches
		}

	DoReplace(text, selected, from, to, options)
		{
		_context = Object(:text, ast: options.GetDefault(#ast, false))
		findpat = .findPat(options.find, options)
		if findpat is ""
			{ .beep(); return false }
		replacement = .Replacement(options.replace, options)
		if .isAstPattern?(findpat)
			{
			if false is (match = .findOne(from, findpat)) or
				match[0][0] + match[0][1] > to
				return false
			writer = _context.ast.GetNewWriter()
			s = TdopReplace(writer, findpat, replacement, :from, :to)
			}
		else
			{
			s = selected
			if s !~ findpat
				return false
			s = s.Replace(findpat, replacement)
			}
		return s
		}
	}
