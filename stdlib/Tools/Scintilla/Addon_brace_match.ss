// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	ContextMenu()
		{
		return #("Go To Matching Brace\tCtrl+M",
			"Select To Matching Brace\tCtrl+Shift+M")
		}

	Init()
		{
		.DefineStyle(SC.STYLE_BRACELIGHT, .GetSchemeColor('braceGoodFore'),
			back: .GetSchemeColor('braceGoodBack'), bold:)
		.DefineStyle(SC.STYLE_BRACEBAD, .GetSchemeColor('braceBadFore'),
			back: .GetSchemeColor('braceBadBack'), bold:)
		}

	getter_additionalBraces()
		{
		if 0 is braces = .Send(#BraceMatch_AdditionalBraces)
			braces = ''
		return .additionalBraces = braces
		}

	UpdateUI()
		{
		.clear()
		pos = .GetCurrentPos() - 1
		if .curCharIsBrace(pos)
			{
			.StyleToEnd() // need to finish styling or BraceMatch doesn't work
			other = .BraceMatch(pos)
			if other is -1
				.BraceBadLight(pos)
			else if Abs(other - pos) > 1
				.highlight(pos, other)
			}
		}

	On_Go_To_Matching_Brace()
		{
		if false is pos = .getBracesPositions()
			return

		.GoToPos(pos.match + 1)
		}

	On_Select_To_Matching_Brace()
		{
		if false is (pos = .getBracesPositions())
			return

		if pos.pos < pos.match
			{
			start = pos.pos
			end = pos.match
			}
		else
			{
			start = pos.match
			end = pos.pos
			}

		.SetSelectionStart(start)
		.SetSelectionEnd(end+1)
		}

	getBracesPositions()
		{
		pos = .GetCurrentPos() - 1
		if not .curCharIsBrace(pos)
			return false

		match = .BraceMatch(pos)
		return match is -1
			? false
			: Object(:pos, :match)
		}

	clear()
		{
		.BraceHighlight(-1, -1)
		.SetHighlightGuide(0)
		}

	highlight(pos, other)
		{
		.BraceHighlight(pos, other)
		col = .GetColumn(pos)
		.SetHighlightGuide(col)
		}

	curCharIsBrace(pos)
		{
		c = .GetAt(pos)
		braces = "[](){}" $ .additionalBraces
		return braces.Has?(c)
		}
	}