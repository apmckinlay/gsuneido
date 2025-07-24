// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// This function reads a nul terminated string from the clipboard
function (format = false)
	{
	if format is false
		format = CF.TEXT
	if not OpenClipboard(NULL)
		return NULL
	if NULL is hm = GetClipboardData(format)
		{
		CloseClipboard()
		return NULL
		}
	if NULL is str = GlobalString(hm)
		{
		CloseClipboard()
		return NULL
		}
	CloseClipboard()
	return str
	}
