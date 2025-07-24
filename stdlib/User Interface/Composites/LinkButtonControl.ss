// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
StaticControl
	{
	New(name, command = false, tip = false, color = 'BLUE', bgndcolor = "",
		tabstop = true, font = "", size = "", weight = "")
		{
		super(name, underline:, :color, :tip, :tabstop,
			:bgndcolor, :font, :size, :weight)
		.Name = ToIdentifier(name.Trim())
		if command is false
			command = name
		.command = "On_" $ ToIdentifier(command.Trim())
		.SubClass()
		}
	LBUTTONDOWN() // this is not called by StaticControl, only triggered in StaticControl
		{
		.Send(.command)
		return 0
		}
	STN_CLICKED()
		{
		.Send(.command)
		return 0
		}
	ContextMenu(x, y)
		{
		if x is 0 and y is 0 // keyboard
			{
			pt = Object(x: 10, y: 20)
			ClientToScreen(.Hwnd, pt)
			x = pt.x
			y = pt.y
			}
		.Send(.command $ "_ContextMenu", x, y)
		return 0
		}
	MOUSEMOVE()
		{
		SetCursor(LoadCursor(ResourceModule(), IDC.HAND))
		return 0
		}
	SETFOCUS()
		{
		.DrawFocusRect(true)
		return 0
		}
	KILLFOCUS()
		{
		.DrawFocusRect(false)
		return 0
		}
	GETDLGCODE()
		{
		return DLGC.WANTCHARS
		}
	CHAR(wParam)
		{
		if wParam is VK.SPACE
			.STN_CLICKED()
		return 0
		}
	}