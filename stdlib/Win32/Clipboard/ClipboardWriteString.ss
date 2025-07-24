// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// This function writes a nul terminated string to the clipboard
function (str, format = false, add? = false)
	{
	if format is false
		format = CF.TEXT
	if not OpenClipboard(NULL)
		return NULL
	hm = GlobalAllocString(str)
	if not add?
		EmptyClipboard()
	x = SetClipboardData(format, hm)
	CloseClipboard()
	return x
	}
