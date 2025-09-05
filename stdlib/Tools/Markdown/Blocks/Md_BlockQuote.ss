// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_ContainerBlock
	{
	Match(line)
		{
		if false is n = .IgnoreLeadingSpaces(line)
			return false

		if not line[n..].Prefix?('>')
			return false

		bq = new Md_BlockQuote()
		bq.Add(.removeBlockQuoteMarker(line[n..]))
		return bq
		}

	removeBlockQuoteMarker(line)
		{
		line = Md_Helper.Detab(line[1..]/*remove '>'*/)
		if line.Prefix?(' ')
			line = line[1..]
		return line
		}

	Continue(line)
		{
		if false is n = .IgnoreLeadingSpaces(line)
			return false

		if not line[n..].Prefix?('>')
			return false

		return .removeBlockQuoteMarker(line[n..])
		}
	}