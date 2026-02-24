// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
_ScintillaDiffControl
	{
	ComponentName: #ScintillaDiff

	AddMarginText(line, text)
		{
		SuServerPrint('AddMarginText', line, text)
//		.MarginSetStyle(line, SC.STYLE_LINENUMBER)
//		SendMessageTextIn(.Hwnd, SCI.MARGINSETTEXT, line, text)
		}
	}