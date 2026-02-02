// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Continue(line/*unused*/, start)
		{
		return false, start
		}

	NextOpenBlock()
		{
		return false
		}

	Match(line/*unused*/, start/*unused*/)
		{
		throw 'NOT IMPLEMENTED'
		}

	Add(line/*unused*/, start/*unused*/)
		{
		throw 'NOT IMPLEMENTED'
		}

	HasEndingBlankLine?()
		{
		return false
		}

	Closed?: false
	Close()
		{
		.Closed? = true
		}

	// helper
	CountLeadingChar(line, start, char)
		{
		if line.Size() <= start
			return 0
		return line.Find1of('^' $ char, start) - start
		}

	IgnoreLeadingSpaces(line, start, limit = 3)
		{
		n = .CountLeadingChar(line, start, ' ')
		if limit is false or n <= limit
			return n
		return false
		}

	BlankLine?(line, start)
		{
		return line[start..].Blank?()
		}

	NotBlankLine?(line, start = 0)
		{
		return not .BlankLine?(line, start)
		}
	}
