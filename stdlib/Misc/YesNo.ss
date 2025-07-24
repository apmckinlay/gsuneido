// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function (prompt = "", title = "", hwnd = 0, flags = 0)
	{
	return ID.YES is Alert(prompt, title, hwnd, flags | MB.YESNO)
	}
