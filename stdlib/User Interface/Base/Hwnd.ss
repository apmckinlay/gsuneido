// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Control
	{
	Map: ()
	hfont:	0

	CreateWindow(className, windowName, style, exStyle = 0,
		x = 9999, y = 9999, w = 0, h = 0, id = 0)
		{
		.Hwnd = CreateWindowEx(exStyle, className, windowName, WS.CHILD | style,
			x, y, w, h, .WndProc.Hwnd, id, Instance(), NULL/*param*/)
		if .Hwnd is 0
			{
			if Suneido.User is 'default'
				throw "can't CreateWindow " $ className
			else
				{
				Alert("Can't create window.

Your system may be low on Windows resources.

You may need to restart your computer.", "Error", flags: MB.ICONERROR)
				SuneidoLog("ERROR: can't CreateWindow " $ className)
				throw ''
				}
			}
		.AddHwnd(.Hwnd)
		}

	SendMessage(msg, wParam = 0, lParam = 0)
		{
		return SendMessage(.Hwnd, msg, wParam, lParam)
		}
	SendMessageText(msg, wParam = 0, lParam = 0)
		{
		return SendMessageText(.Hwnd, msg, wParam, lParam)
		}
	SendMessageTextIn(msg, wParam = 0, lParam = 0)
		{
		return SendMessageTextIn(.Hwnd, msg, wParam, lParam)
		}
	SendMessageTextOut(msg, wParam = 0, lParam = 0)
		{
		return SendMessageTextOut(.Hwnd, msg, wParam, lParam)
		}
	PostMessage(msg, wParam = 0, lParam = 0)
		{
		return PostMessage(.Hwnd, msg, wParam, lParam)
		}

	AddHwnd(hwnd)
		{
		.Window.HwndMap[hwnd] = this
		}
	DelHwnd(hwnd)
		{
		.Window.HwndMap.Delete(hwnd)
		}
	EditHwnd() // default
		{
		return .Hwnd
		}
	Mapcmd(cmd)
		{
		return .Window.Mapcmd(cmd)
		}

	SetWinPos(x, y)
		{
		SetWindowPos(.Hwnd, 0, x, y, 0, 0, SWP.NOSIZE | SWP.NOZORDER | SWP.NOACTIVATE)
		}
	SetWinSize(w, h)
		{
		SetWindowPos(.Hwnd, 0, 0, 0, w, h, SWP.NOMOVE | SWP.NOZORDER | SWP.NOACTIVATE)
		}

	Resize(x, y, w, h)
		{
		if not .Member?('Hwnd')
			return
		// NOCOPYBITS prevents redraw problems with order of moving controls like EDIT
		SetWindowPos(.Hwnd, 0, x, y, w, h,
			SWP.NOCOPYBITS | SWP.NOZORDER | SWP.NOACTIVATE)
		}
	Update()
		{
		UpdateWindow(.Hwnd)
		}
	// receive reflected notifications
	Notify(msg, lParam) /*internal*/
		{
		if (.Map.Member?(msg) and .Method?(method = .Map[msg]))
			return this[method](:lParam)
		else
			return 'callsuper'
		}
	ContextMenu(x /*unused*/, y /*unused*/)
		{
		// x and y are screen coordinates
		return 'callsuper'
		}
	GetFont()
		{
		return .hfont
		}
	SetFont(font = "", size = "", weight = "", text = "",
		underline = false, italic = false, strikeout = false, orientation = 0)
		{
		lf = .LogFont(:font, :size, :weight, :underline, :italic, :strikeout,
			:orientation)
		.SetLogFont(lf, :text)
		if Suneido.User is 'default' and
			lf.lfFaceName isnt used = .GetTextFace()
			{
			Print("SetFont used " $ used $ " for " $ lf.lfFaceName)
			StackTrace(10 /* = max trace levels */)
			}
		}
	LogFont(font = "", size = "", weight = "",
		underline = false, italic = false, strikeout = false, orientation = 0)
		{
		factor = .getFontSizeFactor()
		font = StdFonts.Font(font)
		weight = StdFonts.Weight(weight)
		size = StdFonts.FontSize(size)
		if font is ""
			lf = .initLogFont(factor, :size, :weight, :underline, :italic,
				:strikeout, :orientation)
		else
			{
			lf = Object(
				lfFaceName: font,
				lfHeight: -size * factor,
				lfWeight: weight,
				lfUnderline: underline,
				lfItalic: italic,
				lfStrikeOut: strikeout,
				lfOrientation: orientation,
				lfEscapement: orientation,
				lfCharSet: CHARSET[GetLanguage().charset])
			}
		lf.lfHeight = lf.lfHeight.Round(0)
		return lf
		}

	getFontSizeFactor()
		{
		.WithDC()
			{ |hdc|
			return GetDeviceCaps(hdc, GDC.LOGPIXELSY) / PointsPerInch
			}
		}

	initLogFont(factor, size, weight, underline, italic, strikeout, orientation)
		{
		lf = .suneidoLogFont().Copy()
		if not .defaultFont?(:size, :weight, :underline, :italic,
			:strikeout, :orientation)
			{
			if underline isnt false
				lf.lfUnderline = underline
			if weight isnt '' and weight isnt false
				lf.lfWeight = weight
			lf.lfHeight = -size * factor
			if italic isnt false
				lf.lfItalic = italic
			if strikeout isnt false
				lf.lfStrikeOut = strikeout
			if orientation isnt 0
				lf.lfEscapement = lf.lfOrientation = orientation
			}
		return lf
		}

	suneidoLogFont()
		{
		return Suneido.logfont
		}

	defaultFont?(size, weight, underline, italic, strikeout, orientation)
		{
		return size is '' and weight is '' and underline is false and italic is false and
			strikeout is false and orientation is 0
		}

	SetLogFont(logfont, text = "")
		{
		fonts = Suneido.GetInit("HwndFonts", Object().Set_default(Object()))
		cachedhwnd = fonts.Find(logfont)
		if cachedhwnd isnt false
			.hfont = cachedhwnd
		else if logfont is .suneidoLogFont()
			.hfont = Suneido.hfont
		else
			{
			.hfont = CreateFontIndirect(logfont)
			if .hfont isnt 0
				fonts[.hfont] = logfont
			}
		if .hfont is 0
			.hfont = Suneido.hfont
		if text isnt ""
			.MeasureSizes(text)
		SendMessage(.Hwnd, WM.SETFONT, .hfont, true)
		}
	MeasureSizes(text)
		{
		.WithSelectObject(.hfont)
			{|hdc|
			GetTextExtentPoint32(hdc, text, text.Size(), ex = Object())
			.Xmin = ex.x
			.Ymin = ex.y
			GetTextMetrics(hdc, tm = Object())
			.Top = .Ymin - tm.Descent
			.AveCharWidth = tm.AveCharWidth
			}
		}
	GetTextFace()
		{
		.WithSelectObject(.hfont)
			{|hdc|
			return GetTextFace(hdc)
			}
		}
	TextExtent(s)
		{
		.WithSelectObject(.hfont)
			{|hdc|
			GetTextExtentPoint32(hdc, s, s.Size(), ex = Object())
			}
		return ex
		}
	GetStyle(extended = false)
		{
		return GetWindowLong(.Hwnd, (extended) ? GWL.EXSTYLE : GWL.STYLE)
		}
	GetExStyle()
		{
		return GetWindowLong(.Hwnd, GWL.EXSTYLE)
		}
	SetStyle(style, extended = false)
		{
		SetWindowLong(.Hwnd, (extended) ? GWL.EXSTYLE : GWL.STYLE, style)
		}
	AddStyle(style, extended = false)
		{
		.SetStyle(.GetStyle(extended) | style, extended)
		}
	RemStyle(style, extended = false)
		{
		thisStyle = .GetStyle(extended)
		if ((thisStyle & style) is style)
			thisStyle ^= style
		.SetStyle(thisStyle, extended)
		}
	HasStyle?(style, extended = false)
		{
		return (.GetStyle(extended) & style) is style
		}
	HasFocus?()
		{
		return .Hwnd isnt 0 and GetFocus() is .Hwnd
		}
	Repaint(erase = true)
		// TODO: 'Repaint' is a terrible name, since all this does is add the
		//       client rectangle into the update region for future painting.
		//       It should be renamed to InvalidateClientArea()
		{
		.InvalidateClientArea(erase)
		}
	InvalidateClientArea(erase? = true)
		{
		if .Hwnd isnt false    // avoid repainting entire screen
			InvalidateRect(.Hwnd, NULL, erase?)
		}
	GetEnabled()
		{
		return IsWindowEnabled(.Hwnd)
		}
	SetEnabled(enabled)
		{
		if (.GetEnabled() isnt enabled)
			{
			EnableWindow(.Hwnd, enabled)
			.Repaint()
			}
		}
	hidden: false
	GetHidden()
		{
		return .hidden
		}
	SetHidden(hidden)
		{
		.hidden = hidden
		}
	SetVisible(visible)
		{
		Assert(Boolean?(visible))
		visible = not .GetHidden() and visible
		ShowWindow(.Hwnd, visible ? SW.SHOW : SW.HIDE)
		}
	SetReadOnly(readOnly)
		{
		Assert(Boolean?(readOnly))
		.SetEnabled(not readOnly)
		}
	SetFocus()
		{
		SetFocus(.Hwnd)
		}
	GetReadOnly()
		{
		return not .GetEnabled()
		}
	GetRect()
		{
		// returns a Rect object containing the coordinates of
		// the control (relative to parent Hwnd)
		rc = GetWindowRect(.Hwnd)
		pt = Object(x: 0, y: 0)
		if (GetParent(.Hwnd) isnt NULL)
			ScreenToClient(GetParent(.Hwnd), pt)
		return Rect(rc.left + pt.x, rc.top + pt.y, rc.right - rc.left,
			rc.bottom - rc.top)
		}
	GetClientRect()
		{
		GetClientRect(.Hwnd, rc = Object())
		return Rect(rc.left, rc.top, rc.right - rc.left, rc.bottom - rc.top)
		}
	DisableTheme()
		{
		SetWindowTheme(.Hwnd, "\x00\x00", "\x00\x00")
		}
	tip?: false
	ToolTip(tip)
		{
		.Window.ToolTip(.Hwnd, tip)
		.tip? = true
		}
	WithDC(block)
		{
		WithDC(.Hwnd, block)
		}
	WithSelectObject(object, block)
		{
		.WithDC()
			{|hdc|
			orig = SelectObject(hdc, object)
			try
				block(:hdc)
			catch (e)
				{
				SelectObject(hdc, orig)
				throw e
				}
			SelectObject(hdc, orig)
			}
		}

	DestroyWindow()
		{
		if not .Member?('Hwnd')
			return
		if (.Hwnd isnt false)
			DestroyWindow(.Hwnd)
		.DelHwnd(.Hwnd)
		.Hwnd = false
		}
	Destroy()
		{
		if .tip?
			.Window.RemoveTip(.Hwnd)
		.DestroyWindow()
		super.Destroy()
		}
	}
