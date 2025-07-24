// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// remove parenthesis, square brackets, curly braces, or quotes
ScintillaAddon
	{
	ContextMenu()
		{
		return #('Unwrap\tCtrl+U')
		}
	On_Unwrap()
		{
		code = .Get()
		pos = .GetSelectionStart()
		end = .GetSelectionEnd()
		if end is pos and false isnt range = CodeNest(code, pos)
			.unwrap(range)
		else
			Beep()
		}
	unwrap(range)
		{
		.SetSel(range[0], range[1] + 1)
		text = .GetSelText()

		unwrappedText = .getUnwrapped(text)
		.ReplaceSel(unwrappedText)
		.SetSel(range[0], range[0] + unwrappedText.Size())
		}
	getUnwrapped(text)
		{
		inner = text[1..-1]

		return text[0] is '{'
			? inner.Trim()
			: inner
		}
	}