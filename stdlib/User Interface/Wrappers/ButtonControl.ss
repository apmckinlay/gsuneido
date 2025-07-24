// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// NOTE: When making modifications, please ensure that \ buttons still work
//  as well as buttons that are inside controllers.
//  Also make sure buttons wrapped in multiple controllers still work.
WndProc
	{
	Xmin: ""
	Ymin: ""
	New(text, command = false, font = "", size = "",
		weight = "", tabover = false, defaultButton = false, style = 0,
		tip = false, pad = false, color = false, width = false,
		underline = false, italic = false, strikeout = false, hidden = false)
		{
		if .Name is ""
			.Name = ToIdentifier(text)
		if command is false
			command = .Name
		command = command.Trim()
		.command = "On_" $ ToIdentifier(command)
		style = .style(style, hidden, tabover, defaultButton)
		text = TranslateLanguage(text)
		.CreateWindow("button", text, style, x: 10, y: 10,
			id: .Mapcmd(command $ "_button"))
		.SubClass()
		.xmin_orig = .Xmin
		.ymin_orig = .Ymin
		.SetFont(font, size, weight, text, underline, italic, strikeout)
		.initXYmin(pad, width)
		.Map = Object()
		.Map[BN.CLICKED] = 'CLICKED'
		if tip isnt false
			.ToolTip(tip)
		if color isnt false
			.brush = CreateSolidBrush(color)
		}

	style(style, hidden, tabover, defaultButton)
		{
		.SetHidden(hidden)
		if hidden is false
			style |= WS.VISIBLE
		style |= (tabover is true or .tabOver?() ? 0 : WS.TABSTOP)
		if defaultButton is true
			style |= BS.DEFPUSHBUTTON
		return style
		}

	initXYmin(pad, width)
		{
		if pad is false
			pad = .Ymin + (.Ymin % 2) // make it even
		else
			pad = ScaleWithDpiFactor(pad)
		.pad = pad
		if Number?(width)
			.xmin_orig = .CalcWidth(width)
		.Xmin += pad
		.Ymin += 10 /* = y padding*/
		.Top += 5   /* = top padding*/
		if .xmin_orig isnt "" and .xmin_orig > .Xmin
			.Xmin = .xmin_orig
		if .ymin_orig isnt "" and .ymin_orig > .Ymin
			.Ymin = .ymin_orig
		}

	GETDLGCODE() // to make Enter key working when focused
		{
		return DLGC.BUTTON | DLGC.DEFPUSHBUTTON | DLGC.UNDEFPUSHBUTTON
		}

	CalcWidth(width)
		{
		return width * .AveCharWidth + .pad
		}

	tabOver?()
		{
		return Settings.Get('ButtonControl_Tabover?') is true
		}

	target: false
	CLICKED()
		{
		if .target isnt false
			.target[.command]()
		else
			.Send(.command)
		return 0
		}

	SetCommandTarget(.target) { }

	brush: 'callsuper'
	CTLCOLORBTN()
		{
		return .brush
		}

	// NOTE: pushed state is not visible on ButtonControl
	// use EnhancedButtonControl if you want a visible pushed state
	pushed: false
	Pushed?(state = '')
		{
		if state isnt ''
			.pushed = state
		return .pushed
		}
	// stub to override Hwnd: read-only not applicable to button
	SetReadOnly(readOnly /*unused*/)
		{ }
	GetReadOnly()			// stub to override Hwnd: read-only not applicable to button
		{ return true }
	GetCommand()
		{
		return .command
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
	Get()
		{
		return GetWindowText(.Hwnd)
		}
	Set(text)
		{
		if text is .Get()
			return
		SetWindowText(.Hwnd, text)
		.WithSelectObject(.GetFont())
			{|hdc|
			GetTextExtentPoint32(hdc, text, text.Size(), ex = Object())
			}
		.SetXmin(ex.x + .pad)
		}
	SetXmin(xmin)
		{
		if .xmin_orig isnt ''
			return
		.Xmin = xmin
		.WindowRefresh()
		}
	Destroy()
		{
		if 'callsuper' isnt .brush
			if not DeleteObject(.brush)
				throw "can't DeleteObject(.brush => " $ Display(.brush) $ ")"
		.brush = 'callsuper'
		super.Destroy()
		}
	}
