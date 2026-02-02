// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Code
	{
	New(info, .indentN, .symbol, .symbolN)
		{
		super(:info)
		}

	requiredSymbolN: 3
	Match(line, start)
		{
		if false is indentN = .IgnoreLeadingSpaces(line, start)
			return false

		symbol = false
		if .requiredSymbolN <= symbolN = .CountLeadingChar(line, start+indentN, '`')
			symbol = '`'
		else if .requiredSymbolN <= symbolN = .CountLeadingChar(line, start+indentN, '~')
			symbol = '~'
		else
			return false

		info = line[start+indentN+symbolN..].Trim()
		if symbol is '`' and info.Has?('`')
			return false

		return new this(info, indentN, symbol, symbolN)
		}

	Continue(line, start)
		{
		if false isnt n = .IgnoreLeadingSpaces(line, start)
			{
			if .symbolN <= m = .CountLeadingChar(line, start+n, .symbol)
				{
				if .BlankLine?(line, start+n+m)
					.Close()
				}
			}
		return line, start
		}

	Add(line, start)
		{
		if .Closed?
			return

		n = .CountLeadingChar(line, start, ' ')
		super.Add(line, start+Min(.indentN, n))
		}
	}
