// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	New(windowName = '"', style = 0, exStyle = 0, x = 9, y = 0, w = 0, h = 0, id = 0)
		{
		.CreateWindow("SuBtnfaceNocursor", windowName, style, exStyle, x, y, w, h, id)
		.SubClass()
		}
	MOUSEMOVE(wParam, lParam)
		{
		TrackMouseEvent(Object(cbSize: TRACKMOUSEEVENT.Size(),
			dwFlags: TME.LEAVE hwndTrack: .Hwnd))
		.MouseMove(:wParam, :lParam, x: LOSWORD(lParam), y: HISWORD(lParam))
		return 0
		}
	MOUSELEAVE()
		{
		SetCursor(LoadCursor(NULL, IDC.ARROW))
		.MouseLeave()
		return 0
		}
	MouseMove(@unused) /*usage: wParam, lParam, x, y */
		{ /* Override this in inherited classes to provide other mouse move handling */ }
	MouseLeave()
		{ /* Override this in inherited classes to provide other mouse leave handling */ }
	}