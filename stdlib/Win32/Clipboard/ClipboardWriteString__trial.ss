// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (str, format = false, add? = false)
	{
	if format is false
		format = CF.TEXT
	if not OpenClipboard(NULL)
		return NULL
	hm = GlobalAllocData(str $ '\x00')
	if not add?
		EmptyClipboard()
	x = SetClipboardData(format, hm)
	CloseClipboard()
	return x
	}
