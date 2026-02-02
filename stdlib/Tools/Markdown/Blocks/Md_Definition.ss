// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Base
	{
	StartDL: false
	EndDL: false
	New(.term, def)
		{
		.defs = [def]
		.ParsedTerm = false
		.ParsedDefs = Object()
		}

	MatchDefLine(line, start)
		{
		line = line[start..]
		if line !~ '^:\s'
			return false
		return line[start+1..].LeftTrim()
		}

	Continue(line, start)
		{
		if .BlankLine?(line, start)
			return false, start
		if Md_ContainerBlock.MatchParagraphInteruptableBlockItem(line, start) isnt false
			return false, start
		return line, start
		}

	Add(line, start)
		{
		if false isnt def = .MatchDefLine(line, start)
			.defs.Add(def)
		else
			.defs[.defs.Size() - 1] $= '\n' $ line[start..]
		}

	ParseInline()
		{
		.ParsedTerm = MarkdownInlineParser2(.term)
		for def in .defs
			.ParsedDefs.Add(MarkdownInlineParser2(def))
		}
	}