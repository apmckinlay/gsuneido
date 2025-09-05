// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Code
	{
	requiredSpaces: 4
	Match(line)
		{
		if .CountLeadingChar(line, ' ') >= .requiredSpaces
			return new this(line[.requiredSpaces..])
		return false
		}

	Continue(line)
		{
		if .BlankLine?(line)
			return line[.requiredSpaces..]

		if .CountLeadingChar(line, ' ') >= .requiredSpaces
			return line[.requiredSpaces..]

		return false
		}

	Close()
		{
		if false is start = .Codes.FindIf(.NotBlankLine?)
			.Codes = Object()
		else
			.Codes = .Codes[start...Codes.FindLastIf(.NotBlankLine?)+1]
		}
	}