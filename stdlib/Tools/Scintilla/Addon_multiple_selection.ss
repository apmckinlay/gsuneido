// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	Init()
		{
		.SetMultipleSelection(true)
		.SetAdditionalSelectionTyping(true)
		}
	ContextMenu()
		{
		return #('Multi Select Next Occurrence\tCtrl+E')
		}
	On_Multi_Select_Next_Occurrence()
		{
		// similar logic as ScintillaControl find_selected
		sel = .GetSelect()
		.SelectCurrentWord()
		if "" is find = .GetSelText().Trim()
			{
			Beep()
			.SetSel(sel.cpMin, sel.cpMax)
			return
			}
		pos = .GetSelect().cpMax
		word = sel.cpMin is sel.cpMax
		if ((false is match =
			Find.DoFind(.SearchText(), pos, [case:, :word, :find])) or match[0] < pos)
			{
			Beep()
			return
			}
		org = match[0]
		end = match[0] + match[1]
		.AddSelection(end, org)
		.EnsureRangeVisible(org, end)
		.ScrollCaret()
		}
	}