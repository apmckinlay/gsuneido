// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// MAYBE allow TAB'ing over inserted characters (like Eclipse)
ScintillaAddon
	{
	autodelimit_open: #('(', '[', '{', '"', "'", "`")
	autodelimit_close: #(')', ']', '}', '"', "'", "`")

	selText: 	false
	selMin: 	false
	selMax:		false
	UpdateUI()
		{
		.resetSelection()
		if '' isnt (selText = .GetSelText())
			{
			.selText = selText
			sel = .GetSelect()
			.selMin = sel.cpMin
			.selMax = sel.cpMax
			}
		}

	styleLevel: 1
	Init()
		{
		.indic_delim = .IndicatorIdx(level: .styleLevel)
		}

	Styling()
		{
		return [[level: .styleLevel, indicator: [INDIC.HIDDEN]]]
		}

	CharAdded(c)
		{
		sel = .GetSelect()
		if .autodelimit_close.Has?(c) and .closingChar?(c, sel)
			.On_Delete()
		else if false isnt (i = .autodelimit_open.Find(c))
			.autoDelimit(c, i, sel)
		.resetSelection()
		}

	closingChar?(c, sel)
		{
		return c is .GetAt(sel.cpMin) and .HasIndicator?(sel.cpMin, .indic_delim)
		}

	autoDelimit(c, idx, sel)
		{
		preText = ''
		curLine = .GetLine()
		if .selText is false and (not .at_end_of_line?(sel.cpMin) or
			curLine.Has?('//'))
			return

		if c is '{'
			preText = .handleIndents(curLine)

		.pasteText(preText $ .autodelimit_close[idx])
		.SetIndicator(.indic_delim, sel.cpMin, 1)
		.updateSelect(sel)
		}

	at_end_of_line?(pos)
		{
		while .HasIndicator?(pos, .indic_delim)
			++pos
		c = .GetAt(pos)
		return c is '\x00' or c is '\r' or c is '\n'
		}

	handleIndents(curLine)
		{
		if curLine.Tr(" \t{\r\n") isnt "" or .selText isnt false
			return ' '

		indent = curLine.Extract("^[ \t]*")
		nextLine = .GetLine(1 +.LineFromPosition())
		return nextLine.Extract("^[ \t]*") isnt indent ? '\r\n' $ indent : ''
		}

	pasteText(text)
		{
		if .selText isnt false
			text = .selText $ text
		.Paste(text)
		}

	updateSelect(select)
		{
		if .selText is false
			.SetSelect(select.cpMax)
		else
			.SetSelect(.selMin, .selMax - .selMin + 2)
		}

	resetSelection()
		{
		.selText = .selMin = .selMax = false
		}

	// when you delete an opening delimiter
	// automatically delete a following automatic closing delimiter
	BeforeDelete(pos, len)
		{
		if len isnt 1
			return
		if false isnt (i = .autodelimit_open.Find(.GetAt(pos))) and
			.autodelimit_close[i] is .GetAt(pos + 1) and
			.HasIndicator?(pos + 1, .indic_delim)
			{
			.Defer(.autodelete)
			}
		}
	autodelete()
		{
		pos = .GetSelect().cpMin
		.SetTargetStart(pos)
		.SetTargetEnd(pos + 1)
		SendMessageTextIn(.Hwnd, SCI.REPLACETARGET, 0, "")
		}
	}