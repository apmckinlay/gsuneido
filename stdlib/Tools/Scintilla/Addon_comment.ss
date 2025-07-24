// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// requires On_ methods be forwarded to Scintilla as LibView does
ScintillaAddon
	{
	ContextMenu()
		{
		return #('Comment Lines\tCtrl+/',
			'Comment Selection\tShift+Ctrl+/')
		}
	// by Johan Samyn 15-02-2004
	On_Comment_Lines()
		{
		selStart = .GetSelectionStart()
		selEnd = .GetSelectionEnd()
		textLengthBefore = .GetTextLength()
		colSelEnd = .GetColumn(selEnd)

		// Return if there's no text at all, or we're positioned after the last line.
		if (textLengthBefore is 0 or
			(selStart is selEnd and selEnd is textLengthBefore and colSelEnd is 0))
			return

		lineEnd = .LineFromPosition(selEnd)
		if (selStart isnt selEnd and colSelEnd is 0)
			lineEnd -= 1
		if (colSelEnd > 0 or selStart is selEnd)
			{
			// extend the end of the selection to the end of the line
			++selEnd
			while (selEnd < textLengthBefore and
				.GetColumn(selEnd) > 0)
				++selEnd
			}
		// extend the start of the selection to the start of the line
		// can't simply subtract the column value from selStart because of tabs
		while (selStart > 0 and
			.GetColumn(selStart) > 0)
			--selStart

		.SetSelect(selStart, selEnd - selStart)
		selText = .GetSelText()

		if selText.Prefix?("//")
			{ // uncomment
			selText2 = selText.Replace("^//", "")
			}
		else
			{ // comment
			selText2 = selText.Replace("^", "//").Replace("//$", "")
			}
		SendMessageTextIn(.Hwnd, SCI.REPLACESEL, 0, selText2)

		textLengthAfter = .GetTextLength()
		selEnd += (textLengthAfter - textLengthBefore)
		.SetSelect(selStart, selEnd - selStart)
		}
	On_Comment_Selection()
		{
		selStart = .GetSelectionStart()
		selEnd = .GetSelectionEnd()
		s = .GetRange(selStart - 2, selEnd + 2)
		if s.Prefix?('/*') and s.Suffix?('*/')
			{
			selStart -= 2
			selEnd += 2
			.SetSelect(selStart, selEnd - selStart)
			}
		else
			s = s[2 .. -2]
		if s.Prefix?('/*') and s.Suffix?('*/')
			{ // uncomment
			.Paste(s[2 .. -2])
			.SetSelect(selStart, selEnd - selStart - 4)
			}
		else
			{ // comment
			.Paste('/*' $ s $ '*/')
			.SetSelect(selStart, selEnd - selStart + 4)
			}
		}
	}