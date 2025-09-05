// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Code
	{
	New(info, .indentN, .symbol, .symbolN)
		{
		super(:info)
		}

	requiredSymbolN: 3
	Match(line)
		{
		if false is indentN = .IgnoreLeadingSpaces(line)
			return false

		line = line[indentN..]
		symbol = false
		if .requiredSymbolN <= symbolN = .CountLeadingChar(line, '`')
			symbol = '`'
		else if .requiredSymbolN <= symbolN = .CountLeadingChar(line, '~')
			symbol = '~'
		else
			return false

		info = line[symbolN..].Trim()
		if symbol is '`' and info.Has?('`')
			return false

		return new this(info, indentN, symbol, symbolN)
		}

	Continue(line)
		{
		if false isnt n = .IgnoreLeadingSpaces(line)
			{
			if .symbolN <= m = .CountLeadingChar(line[n..], .symbol)
				{
				if .BlankLine?(line[n+m..])
					.Close()
				}
			}
		return line
		}

	Add(line)
		{
		if .Closed?
			return

		n = .CountLeadingChar(line, ' ')
		line = line[Min(.indentN, n)..]
		super.Add(line)
		}
	}