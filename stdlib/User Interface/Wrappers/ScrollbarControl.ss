// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	/* purpose:	wraps the windows standard scroll bar control */
	// data:
	Xstretch:	0
	Ystretch:	0
	Border:		false
	Name:		"Scrollbar"
	// interface:
	New(vert = true, enabled = true)
		{
		// set stretching and sizing data
		if vert
			{
			.Xmin = GetSystemMetrics(SM.CXVSCROLL)
			.Ystretch = 1
			}
		else
			{
			.Ymin = GetSystemMetrics(SM.CYHSCROLL)
			.Xstretch = 1
			}
		// create and subclass window
		.CreateWindow("SCROLLBAR", "", WS.VISIBLE | (vert ? SBS.VERT : SBS.HORZ))
		.SubClass()
		.SetEnabled(enabled)
		}
	GetRange()
		{
		GetScrollInfo(.Hwnd, SB.CTL, info = Object(cbSize: SCROLLINFO.Size(),
			fMask: SIF.RANGE))
		return Range(info.nMin, info.nMax)
		}
	SetRange(range)
		// pre:	range is a Range of integers
		// post:	sets this' scrolling range to range
		{
		.SendMessage(SBM.SETRANGE, range.GetLow(), range.GetHigh())
		}
	GetPos()
		// post:	returns this' scroll position
		{
		return .SendMessage(SBM.GETPOS, 0, 0)
		}
	SetPos(pos, redraw = true)
		// pre:	pos is an integer AND redraw is a boolean value
		// post:	sets this' scroll position to pos
		{
		.SendMessage(SBM.SETPOS, pos, redraw)
		}
	enabled: true
	SetEnabled(enabled)
		// pre:	enabled is a Boolean value
		// post:	this is enabled iff enabled is true
		{
		.enabled = enabled
		.SendMessage(SBM.ENABLE_ARROWS, enabled ? ESB.ENABLE_BOTH : ESB.DISABLE_BOTH, 0)
		}
	GetEnabled()
		{
		return .enabled
		}
	GetPageSize()
		{
		GetScrollInfo(.Hwnd, SB.CTL, info = Object(cbSize: SCROLLINFO.Size(),
			fMask: SIF.PAGE))
		return info.nPage
		}
	SetPageSize(size)
		{
		SetScrollInfo(.Hwnd, SB.CTL,
			Object(cbSize:	SCROLLINFO.Size(), fMask: SIF.PAGE, nPage: size), true)
		}
	GetReadOnly()			// read-only not applicable to scrollbar
		{
		return true
		}

	// interface (windows messages):
	VSCROLL(wParam)
		{
		switch (LOWORD(wParam))
			{
		case SB.LINEUP:		.Send('ScrollUp')
		case SB.LINEDOWN:		.Send('ScrollDown')
		case SB.PAGEUP:		.Send('PageUp')
		case SB.PAGEDOWN:		.Send('PageDown')
		case SB.THUMBTRACK: 	.Send('Thumbtrack', HIWORD(wParam))
		case SB.THUMBPOSITION:	.Send('Thumbposition', HIWORD(wParam))
		default:
			}
		return 0
		}
	HSCROLL(wParam)
		{
		switch (LOWORD(wParam))
			{
		case SB.LINELEFT:		.Send('ScrollLeft')
		case SB.LINERIGHT:	.Send('ScrollRight')
		case SB.PAGELEFT:		.Send('PageLeft')
		case SB.PAGERIGHT:	.Send('PageRight')
		case SB.THUMBTRACK: 	.Send('Thumbtrack', HIWORD(wParam))
		case SB.THUMBPOSITION:	.Send('Thumbposition', HIWORD(wParam))
		default:
			}
		return 0
		}
	MOUSEWHEEL(wParam)
		{
		clicks = (HISWORD(wParam) / 120).Int()  /*= Standard Mouse Wheel Delta */
		wsl = SPI_GetWheelScrollLines()
		lines = clicks * wsl
		if wsl is -1
			.Send(clicks > 0 ? 'PageUp' : 'PageDown')
		else
			{
			msg = clicks > 0 ? 'ScrollUp' : 'ScrollDown'
			for (i = 0; i < lines.Abs(); ++i)
				.Send(msg)
			}
		return 0
		}
	}
