// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Xmin: 20
	Ymin: 20
	isReadOnly: false
	New(hidden = false, .allowReadOnlyDropDown = false)
		{
		.SetHidden(hidden)
		.CreateWindow("button", '', (hidden ? 0 : WS.VISIBLE) | BS.OWNERDRAW)
		.SubClass()
		.init()
		}
	init()
		{
		.imageWidth = ScaleWithDpiFactor(9/*=width*/)
		.imageHeight = ScaleWithDpiFactor(9/*=height*/)
		.image = ImageResource('arrow_down.emf')
		}
	pressed: false
	LBUTTONDOWN()
		{
		.pressed = true
		InvalidateRect(.Hwnd, NULL, false)
		.Send('On_DropDown')
		return 0
		}
	LBUTTONUP()
		{
		.pressed = false
		InvalidateRect(.Hwnd, NULL, false)
		return 0
		}
	DRAWITEM(dis)
		{
		disabled = (dis.itemState & ODS.DISABLED) isnt 0
		.draw(dis.hDC, dis.rcItem, disabled)
		return true
		}
	draw(dc, rect, disabled)
		{
		state = .pressed ? THEME.PBS.PRESSED :
				.hot ? THEME.PBS.HOT :
				disabled ? THEME.PBS.DISABLED : THEME.PBS.NORMAL
		if .hot
			DrawThemeBackground(GetWindowTheme(.Hwnd), dc,
				THEME.BP.PUSHBUTTON, state, rect, NULL)
		else
			FillRect(dc, rect,
				disabled ? 1 + COLOR.BTNFACE : GetStockObject(SO.WHITE_BRUSH))
		x = (rect.right - rect.left - .imageWidth) / 2
		y = (rect.bottom - rect.top - .imageHeight) / 2 + ScaleWithDpiFactor(1)// centered
		.image.Draw(dc, x, y, .imageWidth, .imageHeight,
			GetStockObject(.hot ? SO.BLACK_BRUSH : SO.GRAY_BRUSH))
		}
	hot: false
	MOUSEMOVE()
		{
		if .hot
			return 0
		TrackMouseEvent(Object(cbSize: TRACKMOUSEEVENT.Size(),
			dwFlags: TME.LEAVE, hwndTrack: .Hwnd))
		.hot = true
		InvalidateRect(.Hwnd, NULL, false)
		return 0
		}
	MOUSELEAVE()
		{
		.hot = .pressed = false
		InvalidateRect(.Hwnd, NULL, false)
		return 0
		}
	THEMECHANGED()
		{
		.init()
		return 0
		}
	SetReadOnly(readOnly)
		{
		.isReadOnly = readOnly
		if .allowReadOnlyDropDown is true
			return
		super.SetReadOnly(readOnly)
		}

	GetReadOnly()
		{
		if .allowReadOnlyDropDown is true
			return .isReadOnly
		return super.GetReadOnly()
		}

	AddBorder()
		{
		.AddStyle(WS.BORDER)
		}
	}
