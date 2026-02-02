// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_ContainerBlock
	{
	Match(line, start)
		{
		if false is n = .IgnoreLeadingSpaces(line, start)
			return false

		if line[start+n::1] isnt '>'
			return false

		bq = new Md_BlockQuote()
		line, start = .advanceBlockQuoteMarker(line, start+n)
		bq.Add(line, start)
		return bq
		}

	advanceBlockQuoteMarker(line, n)
		{
		line = Md_Helper.Detab(line, n+1/*skip '>'*/)
		if line[n+1::1] is ' '
			n++
		return line, n+1
		}

	Continue(line, start)
		{
		if false is n = .IgnoreLeadingSpaces(line, start)
			return false, start

		if line[start+n::1] isnt '>'
			return false, start

		line, start = .advanceBlockQuoteMarker(line, start+n)
		return line, start
		}
	}
