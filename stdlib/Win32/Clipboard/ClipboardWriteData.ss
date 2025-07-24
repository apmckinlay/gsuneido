// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// This function writes binary data to the clipboard
function (data, format, add? = false)
	{
	if not OpenClipboard(NULL)
		return NULL
	hm = GlobalAllocData(data)
	if not add?
		EmptyClipboard()
	x = SetClipboardData(format, hm)
	CloseClipboard()
	return x
	}
