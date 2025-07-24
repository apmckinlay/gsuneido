// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	ContextMenu()
		{ #("Inspect\tF4") }

	On_Inspect()
		{
		sel = .GetSelect()
		if .haveSelection?(sel)
			{
			.inspect(.GetSelText())
			return
			}

		if 0 isnt .Send("CurrentName")
			{
			.inspect(.Get())
			return
			}

		.findMatchedEndingBracket()
			{ |start, length|
			.SetSelect(start, length)
			.inspect(.GetSelText())
			}
		}

	haveSelection?(sel)
		{
		return sel.cpMin < sel.cpMax
		}

	findMatchedEndingBracket(block)
		{
		line = .GetLine().Trim()
		if line.Suffix?(')') or line.Suffix?(']')
			{
			pos = .GetCurrentPos()
			text = .Get()
			lineEndPos = text.Find('\n', pos) - 2
			match = .BraceMatch(lineEndPos)
			if match isnt -1 and match isnt pos
				block(match, lineEndPos - match + 1)
			}
		}

	inspect(text)
		{
		if text.Prefix?('(')
			text = '#' $ text
		try
			ob = text.SafeEval()
		catch (err /*unused*/)
			{ ob = text }
		Inspect.Window(ob)
		}
	}
