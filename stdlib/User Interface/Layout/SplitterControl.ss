// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// primarily for use by Split
// a container that uses this must have a Movesplit(pos) function
CursorWndProc
	{
	Xmin: 0
	Xstretch: 0
	Ymin: 0
	Ystretch: 0
	HandleLength: 50

	New()
		{
		super(style: WS.CHILD | WS.VISIBLE | WS.CLIPSIBLINGS)
		.Dir = .Parent.Dir

		d = ScaleWithDpiFactor(6) /*= size of splitter */
		if .Dir is "vert"
			{ .Xstretch = 1; .Ymin = d }
		else
			{ .Ystretch = 1; .Xmin = d }

		// Coordinates of split "placeholder" rectangle to erase before drawing new one...
		.posrect = Object(left: -1 top: -1 right: -1 bottom: -1)
		.HandleLength = ScaleWithDpiFactor(.HandleLength)
		}
	SIZE()
		{
		InvalidateRect(.Hwnd, NULL, true)
		return 0
		}
	ERASEBKGND()
		{ return 1 }
	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		GetClientRect(.Hwnd, r = Object())
		FillRect(hdc, r, GetSysColorBrush(COLOR.BTNFACE))
		len = .HandleLength
		if .Dir is 'vert'
			{
			r.left = ((r.right - len) / 2).Int()
			r.right = r.left + len
			r.top = (r.bottom / 2).Int()
			r.bottom = r.top + 1
			FillRect(hdc, r, GetSysColorBrush(COLOR.BTNSHADOW))
			r.top--
			r.bottom--
			FillRect(hdc, r, GetSysColorBrush(COLOR.BTNHIGHLIGHT))
			}
		else
			{
			r.left = (r.right / 2).Int()
			r.right = r.left + 1
			r.top = ((r.bottom - len) / 2).Int()
			r.bottom = r.top + len
			FillRect(hdc, r, GetSysColorBrush(COLOR.BTNSHADOW))
			r.left--
			r.right--
			FillRect(hdc, r, GetSysColorBrush(COLOR.BTNHIGHLIGHT))
			}
		EndPaint(.Hwnd, ps)
		return 0
		}
	dragging: false
	LBUTTONDOWN(lParam)
		{
		.dragging = true
		.x0 = LOSWORD(lParam)
		.y0 = HISWORD(lParam)
		SetCapture(.Hwnd)
		SetCursor(LoadCursor(ResourceModule(),
			.Dir is "vert" ? IDC.VSPLITBAR : IDC.HSPLITBAR))
		return 0
		}
	MouseMove(x, y)
		{
		SetCursor(LoadCursor(ResourceModule(),
			.Dir is "vert" ? IDC.VSPLITBAR : IDC.HSPLITBAR))
		if not .dragging
			return 0

		if .Dir is "vert"
			{ if not .Parent.CanMovesplit?(.y + y - .y0) return false }
		else
			{ if not .Parent.CanMovesplit?(.x + x - .x0) return false }

		hdc = GetDC(.Window.Hwnd)
		GetClientRect(.Hwnd, rcClient = Object())
		ClientToScreen(.Hwnd, mouse_pt = Object(:x, :y))
		ScreenToClient(.Window.Hwnd, mouse_pt)

		// Erase old placeholder rectangle:
		DrawFocusRect(hdc, .posrect)

		ClientToScreen(.Hwnd, pt = Object(x: rcClient.left y: rcClient.top))
		ScreenToClient(.Window.Hwnd, pt)

		if .Dir is "vert"
			{
			.posrect.top = .Parent.Group_y
			.posrect.bottom = mouse_pt.y
			.posrect.left = pt.x
			.posrect.right = pt.x + (rcClient.right - rcClient.left)
			}
		else
			{
			.posrect.left = .Parent.Group_x
			.posrect.right = mouse_pt.x
			.posrect.top = pt.y
			.posrect.bottom = pt.y + (rcClient.bottom - rcClient.top)
			}
		DrawFocusRect(hdc, .posrect)
		ReleaseDC(.Window.Hwnd, hdc)

		return 0
		}
	LBUTTONUP(lParam)
		{
		ReleaseCapture()
		if .dragging
			{
			// End dragging
			.dragging = false
			.eraseFocusRect()
			.Parent.Movesplit(.Dir is "vert"
				? .y + HISWORD(lParam) - .y0
				: .x + LOSWORD(lParam) - .x0)
			}
		return 0
		}
	Resize(.x, .y, .w, .h)
		{
		MoveWindow(.Hwnd, x, y, w, h, true)
		}
	eraseFocusRect()
		{
		hdc = GetDC(.Window.Hwnd)
		DrawFocusRect(hdc, .posrect)
		ReleaseDC(.Hwnd, hdc)
		.posrect.left = .posrect.top = .posrect.right = .posrect.bottom = -1
		}
	}
