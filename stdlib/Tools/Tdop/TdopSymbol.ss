// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.Token, lbp = 0, .Children = #(), .Position = -1, .Length = 0)
		{
		.Lbp = lbp
		}
	Getter_(member)
		{
		if Number?(member) and member >= 0 and member < .ChildrenSize()
			return .Children[member]
		if .Member?('Getter_' $ member)
			return this['Getter_' $ member]()
		return
		}
	ToString()
		{
		value = not .Member?(#Value)
			? ''
			: Object?(.Value)
				? .Value.ToString()
				: .Value
		id = .Token $ Opt('(', value, ')')
		if .ChildrenSize() is 0
			return id
		else
			return '[' $ id $ ', ' $ .Children.Values().Map(#ToString).Join(', ') $ ']'
		}
	Match(token)
		{
		return .Token is token
		}
	Tokenize(unused)
		{
		return this
		}
	Break(unused)
		{
		return _isNewline() and _getStmtnest() is 0 and .Method?(#Nud)
		}
	ChildrenSize()
		{
		return .Children.Size()
		}
	ToWrite()
		{
		if not .Member?(#Value)
			return TdopTokenToStringMap.GetDefault(.Token, '')

		if .Token is TDOPTOKEN.DATE
			return '#' $ .Value

		return .Value
		}

	ForEachMatch(pattern, block) // non-overlapping matches
		{
		pos = .Position
		do
			{
			match = TdopSearch(this, pattern, pos)
			if match is false
				return
			try
				block(match)
			catch (e, "block:")
				if e is "block:break"
					break
			pos = match[0][0] + Max(1, match[0][1])
			}
		while (pos < .Position + .Length)
		}
	}
