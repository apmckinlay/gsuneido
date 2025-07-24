// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// User Interface/Wrappers/UpDownControl
// Wraps the windows UpDown common control; implements basic behavior...
/*
This control automatically sends WM_HSCROLL or WM_VSCROLL messages to its
parent window
*/
WndProc
	{
	Name: "UpDown"

	New(style = 0, exStyle = 0)
		{
		.CreateWindow(UPDOWN_CLASS, "UDC", style | WS.VISIBLE, exStyle)

		// Set xmin, ymin
		.vert? = ((style & UDS.HORZ) isnt UDS.HORZ)
		if (.vert?)
			{
			.Xmin = GetSystemMetrics(SM.CXVSCROLL)
			.Ymin = 2 * .Xmin
			}
		else
			{
			.Ymin = GetSystemMetrics(SM.CYHSCROLL)
			.Xmin = 2 * .Ymin
			}
		}

	SetRange(low = 0, high = 1)
		{
		// Message has no return value...
		.SendMessage(UDM.SETRANGE, 0, MAKELONG(high, low))
		}
	SetPos(pos = 0)
		{
		return LOSWORD(SendMessage(.Hwnd, UDM.SETPOS, 0, LOSWORD(pos)))
		}
	GetPos()
		{
		return LOSWORD(SendMessage(.Hwnd, UDM.GETPOS, 0, 0))
		}
	SetBuddy(hwndBuddy, sizetoBuddy = true)
		{
		if (sizetoBuddy)
			{
			rc = GetWindowRect(hwndBuddy)
			if (.vert?)
				.Ymin = rc.bottom - rc.top
			else
				.Xmin = rc.right - rc.left
			GetClientRect(.Window.Hwnd, rc)
			.SendMessage(WM.SIZE, WMSIZE.RESTORED, rc.right | rc.bottom << 16)
			}
		// Returns handle to previous buddy window
		return SendMessage(.Hwnd, UDM.SETBUDDY, hwndBuddy, 0)
		}
	GetReadOnly() // read-only not applicable to updown
		{
		return true
		}
	}
