// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function (filter = "", hwnd = false, flags = false,
	title = "Save", ext = '', file = "", initialDir = "",
	block = false, alert = '')
	{
	if "" is filename = SaveFileName(filter, hwnd, flags, title, ext, file, initialDir)
		return false

	if block isnt false
		try
			block(filename)
		catch (err)
			{
			Alert((alert is '' ? 'Unable to save file' : alert) $ ': ' $ err,
				title, hwnd, MB.ICONWARNING)
			return false
			}

	return true
	}