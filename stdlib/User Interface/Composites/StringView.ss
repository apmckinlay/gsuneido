// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
function (s, title = 'StringView', hwnd = 0)
	{
	ToolDialog(hwnd,
		Object('ScintillaAddonsEditor' set: s, readonly:, xmin: 300, height: 15),
		:title, border: 0)
	}