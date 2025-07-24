// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(children = false, match = false, mustMatch = false, token = false,
		at = false, implicit = false)
		{
		if children is false
			children = Object()

		if token isnt false
			t = token
		else if match is false
			{
			t = _token()
			_advance()
			}
		else if false is t = .isMatch(match, implicit)
			{
			if mustMatch is true
				throw 'expected ' $ match $ ', but got ' $ Display(_token())
			t = TdopCreateNode(Object?(match) ? match[0] : match)
			}

		if at isnt false
			children.Add(t, :at)
		else
			children.Add(t)

		return children
		}
	isMatch(match, implicit)
		{
		if Object?(match)
			{
			for m in match
				if false isnt t = .match(m, implicit)
					return t
			return false
			}
		else
			return .match(match, implicit)
		}
	match(match, implicit)
		{
		newline = _isNewline()
		token = _token()
		end = _end

		if implicit is true and (match is TDOPTOKEN.SEMICOLON and
			(token is end or token.Match(TDOPTOKEN.RCURLY) or
				(not token.Match(TDOPTOKEN.SEMICOLON) and newline)))
			return TdopCreateNode(TDOPTOKEN.SEMICOLON)// implicit semicolon insertion
		if token.Match(match)
			{
			_advance()
			return token.Tokenize(match)
			}
		return false
		}
	}