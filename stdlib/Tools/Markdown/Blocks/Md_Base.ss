// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Continue(line/*unused*/)
		{
		return false
		}

	NextOpenBlock()
		{
		return false
		}

	Match(line/*unused*/)
		{
		throw 'NOT IMPLEMENTED'
		}

	Add(line/*unused*/)
		{
		throw 'NOT IMPLEMENTED'
		}

	Closed?: false
	Close()
		{
		.Closed? = true
		}

	// helper
	CountLeadingChar(line, char)
		{
		return line.FindRx('[^' $ char $ ']')
		}

	IgnoreLeadingSpaces(line)
		{
		n = .CountLeadingChar(line, ' ')
		if n <= 3 /*=allowed spaces*/
			return n
		return false
		}

	BlankLine?(line)
		{
		return line.Blank?()
		}

	NotBlankLine?(line)
		{
		return not .BlankLine?(line)
		}
	}