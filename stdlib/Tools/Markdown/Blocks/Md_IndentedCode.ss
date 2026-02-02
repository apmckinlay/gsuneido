// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Code
	{
	requiredSpaces: 4
	Match(line, start)
		{
		if .CountLeadingChar(line, start, ' ') >= .requiredSpaces
			return new this(line[start+.requiredSpaces..])
		return false
		}

	Continue(line, start)
		{
		if .BlankLine?(line, start)
			return line, start+.requiredSpaces

		if .CountLeadingChar(line, start, ' ') >= .requiredSpaces
			return line, start+.requiredSpaces

		return false, start
		}

	lastLineIsBlank?: false
	HasEndingBlankLine?()
		{
		return .lastLineIsBlank?
		}

	Close()
		{
		if .Codes.Size() > 0 and .BlankLine?(.Codes.Last(), 0)
			.lastLineIsBlank? = true

		if false is start = .Codes.FindIf(.NotBlankLine?)
			.Codes = Object()
		else
			.Codes = .Codes[start...Codes.FindLastIf(.NotBlankLine?)+1]
		}
	}
