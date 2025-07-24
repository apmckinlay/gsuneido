// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
// If there is a selection, returns it.
// Otherwise it selects the span of non-blank lines around the cursor.
// The result is trimmed.
// Used by WorkSpaceControl and QueryViewControl
class
	{
	CallClass(editor)
		{
		s = editor.GetSelText()
		if s isnt ""
			return s
		s = editor.Get()
		if s.Blank?()
			return ""

		cPos = editor.GetCurrentPos()
		if .clickBlankLine?(s, cPos)
			return ""

		blankline = "\r?\n[ \t]*\r?\n"
		crlfPos = s[.. cPos].FindRxLast(blankline)
		start = crlfPos is false ? 0 : crlfPos + 2
		len = s[start ..].FindRx(blankline)
		// trim
		while s[start] in ('\r', '\n')
			{
			++start
			--len
			}
		while len > 1 and s[start + len - 1] =~ '\s'
			--len
		s = s[start :: len]
		editor.SetSelect(start, len)
		return s
		}

	clickBlankLine?(s, cPos)
		{
		prevNewLinePos = s.FindLast('\n', cPos - 1)
		nextNewLinePos = s.Find('\n', cPos)
		if prevNewLinePos is false
			prevNewLinePos = 0
		return s[prevNewLinePos..nextNewLinePos].Blank?()
		}
	}