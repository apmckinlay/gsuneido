// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
_OpenFileName
	{
	// not support multi for now
	CallClass(filter = "", hwnd = 0, flags/*unused*/ = false,
		multi = false, title = "Open", file/*unused*/ = "",
		initialDir/*unused*/ = "", hDrop = false, attachment? = false)
		{
		if false is filename = ToolDialog(hwnd,
			Object('OpenFileName', filter, hDrop, multi, :attachment?), :title)
			return multi ? [] : ''
		return multi ? filename : filename[0]
		}
	}