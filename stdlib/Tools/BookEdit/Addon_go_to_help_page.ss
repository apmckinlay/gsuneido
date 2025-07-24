// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonForChanges
	{
	ContextMenu()
		{ #("Go To Help") }
	styleLevel: 20
	urlPattern: `\"(\/[a-zA-Z /#.?-]+?)\"|\'(\w+?\.png)\'|GetHelpPage\(\"(\w+?)\"\)`
	WordChars: "-_.!~*'();?:@&=+$,%#/<>/\" " $ // dash must be first for regex
		"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Init()
		{
		.indic_word = .IndicatorIdx(level: .styleLevel)
		}

	Styling()
		{
		return [[level: .styleLevel, indicator: [INDIC.PLAIN, fore: CLR.Highlight]]]
		}

	ProcessChunk(text, pos)
		{
		.ClearIndicator(.indic_word, pos, text.Size())
		text.ForEachMatch(.urlPattern)
			{|m|
			m = m.Values()[1]
			.mark_word(pos + m[0], m[1])
			}
		}

	mark_word(begin, end)
		{
		.SetIndicator(.indic_word, begin, end)
		}

	On_Go_To_Help()
		{
		if '' is link = .Send('FindLinkedHelpPage')
			return
		.Send('GotoHelp', link)
		}

	DoubleClick()
		{
		preSel = .GetSelect()
		if '' is link = .Send('FindLinkedHelpPage', select?:)
			return

		GetCursorPos(pt = Object())
		res = ContextMenu(#('Go To Help', 'Copy')).Show(.Hwnd, pt.x, pt.y)
		if res is 1
			.Send('GotoHelp', link)
		else if res is 2
			ClipboardWriteString(link)
		else
			.SetSelect(preSel.cpMin, preSel.cpMax - preSel.cpMin)
		}
	}
