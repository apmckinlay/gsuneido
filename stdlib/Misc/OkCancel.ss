// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// NOTE: flags only used when prompt is just a string
// This is the dialog, see OkCancelControl for the buttons
function (prompt = "", title = "", hwnd = 0, flags = 0)
	{
	if String?(prompt)
		return 1 is Alert(prompt, title, hwnd, flags | MB.OKCANCEL)

	return ToolDialog(hwnd, [OkCancelWrapper, prompt], :title, closeButton?: false)
	}
