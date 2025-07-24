// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
ButtonControl
	{
	focusRectOffset: -3
	// REFACTOR: if every where using MenuButton can be set sendParents?,
	// we could remove it later
	New(text, .menu = false, tip = false, tabover = false,
		.left = false, width = false, command = false, .sendParents? = false)
		{
		super(text, style: BS.OWNERDRAW, pad: 36,
			:tip, :tabover, :width, :command)
		if command is false
			command = .Name
		.command = 'On_' $ ToIdentifier(command)
		.init()
		}
	init()
		{
		.arrowWidth = ScaleWithDpiFactor(10/*=width*/)
		.arrowHeight = ScaleWithDpiFactor(7/*=height*/)
		.image = ImageResource('arrow_down.emf')
		.xArrowOffset = .focusRectOffset - .arrowWidth
		.yArrowOffset = .focusRectOffset - .arrowHeight
		}
	SetMenu(menu)
		{
		.menu = menu
		}
	GETDLGCODE()
		{
		return DLGC.WANTARROWS
		}
	KEYDOWN(wParam)
		{
		if .disabled is true
			return 0
		if wParam is ' '.Asc() or wParam is VK.DOWN
			.pulldown()
		return 0
		}
	pressed: false
	command: false
	Hwnd: 0
	pulldown()
		{
		if .Hwnd is 0 // in rare case, .Hwnd is invalid when destroying (sugg. 17046)
			return
		.pressed = true

		menu = .menu isnt false	? .menu : .Send('MenuButton_' $ .Name)

		i = .PopupMenu(menu, .Hwnd)

		.pressed = false
		if (i > 0)
			.send(.command, menu, i - 1)
		}
	PopupMenu(menu, hwnd)
		{
		InvalidateRect(hwnd, NULL, false)
		r = GetWindowRect(hwnd)
		rcExclude = Object(left: -9999, right: 9999, top: r.top, bottom: r.bottom)
		i = ContextMenu(menu).Show(hwnd, r.left, r.bottom, left:,
			:rcExclude, buttonRect: r)
		InvalidateRect(hwnd, NULL, false)
		return i
		}
	LBUTTONDOWN()
		{
		if .disabled is true
			return 0
		SetFocus(.Hwnd)
		.pulldown()
		return 0
		}
	hot: false
	MOUSEMOVE()
		{
		if (.hot or .disabled is true)
			return 0
		TrackMouseEvent(Object(cbSize: TRACKMOUSEEVENT.Size(),
			dwFlags: TME.LEAVE, hwndTrack: .Hwnd))
		.hot = true
		InvalidateRect(.Hwnd, NULL, false)
		return 0
		}
	MOUSELEAVE()
		{
		.hot = false
		InvalidateRect(.Hwnd, NULL, false)
		return 0
		}
	send(prefix, menu, chosen, j = 0, parent = false)
		{
		for (m = 0; j <= chosen and m < menu.Size(); ++j, ++m)
			if (Object?(menu[m]))
				j = .send(prefix $ "_" $ ToIdentifier(menu[m - 1]),
					menu[m], chosen, j, prefix) - 1
			else if (j is chosen)
				{
				if parent isnt false and .sendParents?
					.Send(parent, prefix $ "_" $ ToIdentifier(menu[m]))
				.Send(prefix, menu[m], index: m)
				.Send(prefix $ "_" $ ToIdentifier(menu[m]))
				}
		return j
		}
	pushed: false
	Pushed?(state = -1)
		{
		if state isnt -1 and .pushed isnt state
			{
			.pushed = state
			InvalidateRect(.Hwnd, NULL, false)
			}
		return .pushed
		}
	grayed: false
	Grayed?(state = -1)
		{
		if state isnt -1 and .grayed isnt state
			{
			.grayed = state
			InvalidateRect(.Hwnd, NULL, false)
			}
		return .grayed
		}
	disabled: false
	Disable(readonly)
		{
		.disabled = readonly
		}
	Disabled?()
		{
		return .disabled
		}
	ERASEBKGND()
		{ return 1 }
	DRAWITEM(dis)
		{
		dc = dis.hDC
		rect = dis.rcItem
		FillRect(dc, rect, GetSysColorBrush(COLOR.BTNFACE))
		disabled = (dis.itemState & ODS.DISABLED) isnt 0
		.draw(dc, rect, disabled)
		x = rect.right + .xArrowOffset
		y = rect.bottom + .yArrowOffset
		.image.Draw(dc, x, y, .arrowWidth, .arrowHeight,
			GetStockObject(.hot ? SO.BLACK_BRUSH : SO.GRAY_BRUSH))
		if ((dis.itemState & ODS.FOCUS) isnt 0)
			{
			InflateRect(rect, .focusRectOffset, .focusRectOffset)
			DrawFocusRect(dc, rect)
			}
		return true
		}
	draw(dc, rect, disabled)
		{
		state = .pressed or .pushed ? THEME.PBS.PRESSED :
				.hot ? THEME.PBS.HOT :
				disabled or .grayed ? THEME.PBS.DISABLED : THEME.PBS.NORMAL
		hTheme = GetWindowTheme(.Hwnd)
		DrawThemeBackground(hTheme, dc,
			THEME.BP.PUSHBUTTON, state, rect, NULL)
		text = GetWindowText(.Hwnd)
		if .left
			text = '  ' $ text
		DrawThemeText(hTheme, dc, THEME.BP.PUSHBUTTON, state,
			MultiByteToWideChar(text), -1,
			(.left ? DT.LEFT : DT.CENTER) | DT.VCENTER | DT.SINGLELINE,
			0, rect)
		}
	THEMECHANGED()
		{
		.init()
		return 'callsuper'
		}
	}
