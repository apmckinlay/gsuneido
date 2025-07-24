// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
/* purpose: displays a tool-tip-like information window */
WndProc
	{
	// data:
	Xmin:		100
	Ymin:		100
	marginSize:	15
	titleSize:	20
	isNew:		true

	// interface:
	CallClass(@args)
		{
		return new InfoWindowControl(@args)
		}

	New(.text = "", .title = "", x = false, y = false, width = 300, height = 300,
		.marginSize = 15, .titleSize = 20, autoClose = false)
		{
		width = ScaleWithDpiFactor(width)
		height = ScaleWithDpiFactor(height)
		.marginSize = ScaleWithDpiFactor(.marginSize)
		if x is false or y is false
			{
			GetCursorPos(pt = Object())
			x = pt.x
			y = pt.y
			}
		style = WS.BORDER | WS.POPUP
		exStyle = WS_EX.TOOLWINDOW // to keep off of taskbar
		.Window = .Controller = this
		if NULL is .Hwnd = CreateWindowEx(exStyle, "SuToolArrow", title, style,
			x, y, width, height, 0, 0, Instance(), 0)
			throw "Can't create window: SuToolArrow (InfoWindowControl)"
		.SubClass()
		.HwndMap = Object()
		.paint(hdc = GetDC(.Hwnd))
		ReleaseDC(.Hwnd, hdc)
		SetCapture(.Hwnd)
		.SetFocus()
		if Number?(autoClose)
			Delay(autoClose.SecondsInMs(), .Destroy)
		}
	GETDLGCODE()
		{ return DLGC.WANTALLKEYS }
	New2()
		{ }
	SetTitle(title /*unused*/)
		{
		}
	GetTitle()
		{ return .title }
	SetText()
		{}
	GetText()
		{ return .text }
	getTitleRect(rcClient)
		{
		rcTitle = rcClient.Copy()
		rcTitle.top += .marginSize
		rcTitle.left += .marginSize
		rcTitle.right -= .marginSize
		rcTitle.bottom = rcTitle.top + .titleSize
		return rcTitle
		}
	getTextRect(rcClient)
		{
		rcText = rcClient.Copy()
		if .titleSize > 0
			rcText.top += .titleSize + .marginSize
		rcText.top += .marginSize
		rcText.bottom -= .marginSize
		rcText.left += .marginSize
		rcText.right -= .marginSize
		return rcText
		}
	paint(hdc)
		{
		WithHdcSettings(hdc, .hdcSettings())
			{
			GetClientRect(.Hwnd, rc = Object())
			rcTitle = .getTitleRect(rc)
			flags = DT.TOP | DT.LEFT | DT.NOPREFIX
			DrawTextEx(hdc, .title, .title.Size(), rcTitle,
				flags | DT.SINGLELINE, NULL)
			// The below SelectObject will be restored by the
			// surrounding: WithHdcSettings -> GetStockObject(SO.SYSTEM_FONT)
			SelectObject(hdc, GetStockObject(SO.DEFAULT_GUI_FONT))
			rcText = .getTextRect(rc)
			if .isNew
				{
				DrawTextEx(hdc, .text, .text.Size(), rcText,
					flags | DT.EXPANDTABS | DT.WORDBREAK | DT.CALCRECT, NULL)
				rcWindow = GetWindowRect(.Hwnd)
				rcWindow.right = rcWindow.left + rcText.right + .marginSize
				rcWindow.bottom = rcWindow.top + rcText.bottom + .marginSize
				MakeRectVisible(rcWindow)
				MoveWindow(.Hwnd, rcWindow.left, rcWindow.top,
					rcWindow.right - rcWindow.left,
					rcWindow.bottom - rcWindow.top, true)
				.isNew = false
				ShowWindow(.Hwnd, SW.SHOWNA)
				}
			DrawTextEx(hdc, .text, .text.Size(), rcText,
				flags | DT.EXPANDTABS | DT.WORDBREAK, NULL)
			}
		}
	hdcSettings()
		{
		return [
			GetStockObject(SO.SYSTEM_FONT),
			SetBkMode: TRANSPARENT,
			SetTextColor: GetSysColor(COLOR.INFOTEXT)]
		}
	mouseEvent(flag, lParam)
		{
		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		Mouse_event(flag, x, y, 0, NULL)
		}
	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		.paint(hdc)
		EndPaint(.Hwnd, ps)
		return 0
		}
	LBUTTONDOWN(lParam)
		{
		.mouseEvent(MOUSEEVENTF.LEFTDOWN, lParam)
		.Destroy()
		return 0
		}
	RBUTTONDOWN(lParam)
		{
		.mouseEvent(MOUSEEVENTF.RIGHTDOWN, lParam)
		.Destroy()
		return 0
		}
	KEYDOWN()
		{
		.Destroy()
		return 0
		}
	SYSKEYDOWN()
		{
		.Destroy()
		return 0
		}
	CAPTURECHANGED()
		{
		.Destroy()
		return 0
		}
	GetReadOnly()
		{ return true }

	dead?: false
	Destroy()
		{
		if .dead?
			return
		.dead? = true
		ReleaseCapture()
		super.Destroy()
		}
	}
