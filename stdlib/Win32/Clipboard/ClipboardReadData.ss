// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// This function reads binary data from the clipboard
function (format)
	{
	if not OpenClipboard(NULL)
		return NULL
	if NULL is hm = GetClipboardData(format)
		{
		CloseClipboard()
		return NULL
		}
	if NULL is data = GlobalData(hm)
		{
		CloseClipboard()
		return NULL
		}
	CloseClipboard()
	return data
	}
