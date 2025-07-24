// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
EditControl
	{
	Name: 		'Static'
	Unsortable: true
	Xstretch:	false
	Ystretch:	false
	DefaultHeight: 1
	Hasfocus?:	false
	constructing?: false

	// TODO: if the bgndcolor is "none", the application should draw the background
	New(.text = "", .font = "", .size = "", .weight = "", .justify = "LEFT",
		.underline = false, color = "", .whitebgnd = false, tip = false,
		.tabstop = false, .bgndcolor = "", .hidden = false, .textStyle = false)
		{
		super(@.args(color))
		.setsize()
		.setBgColor()
		if tip isnt false
			.ToolTip(tip)
		.constructing? = false
		}

	firstTime: true
	fontChanged: false
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		if .firstTime
			{
			.firstTime = false
			.alignVertPadding()
			}
		if .fontChanged
			{
			.fontChanged = false
			.changeFont() // call it again to avoid text bouncing problem with custom font
			}
		}

	alignVertPadding()
		{
		SendMessageRect(.Hwnd, EM.GETRECT, 0, rect = Object())
		rect.top += 1 // does not seem to be scale with dpi
		rect.bottom += 1
		SendMessageRect(.Hwnd, EM.SETRECT, 0, rect)
		}

	CalcXminByControls(@args)
		{
		.Xmin = .OrigXmin = .DoCalcXminByControls(@args)
		}

	args(color)
		{
		.constructing? = true
		if .textStyle isnt false and StaticTextStyles.Member?(.textStyle)
			{
			.size = .size is '' ? StaticTextStyles[.textStyle].size : .size
			.weight = .weight is '' ? StaticTextStyles[.textStyle].weight : .weight
			color = color isnt '' ? color : StaticTextStyles[.textStyle].color
			}
		.setColor(color)
		return Object(mandatory: false, readonly:,
			style: ES.MULTILINE | ES[.justify],
			hidden: .hidden, tabover: not .tabstop,
			font: .font, size: .size, weight: .weight, underline: .underline,
			height: 1,
			width: false,
			static:)
		}


	color: ''
	setColor(color)
		{
		color = TranslateColor(color)
		if Object?(color)
		   color = RGB(@color)
		if color is .color
			return false
		.color = color
		return true
		}

	setsize(refreshRequired? = false)
		{
		.Set(.text, :refreshRequired?)
		.Ymin -= 3 /*= to align with existing control*/
		lines = .CalcLines(.text, .OrigXmin)
		.Ymin *= lines
		.Ymin -= lines - 1 /*= to align with existing control when multi lines */
		++.Top
		if .OrigXmin isnt 0
			.Xmin = .OrigXmin
		if .OrigYmin isnt 0
			.Ymin = .OrigYmin
		.SetMargins(Object(left: 0, right: 0)) // set margin at end, so it works on .Set
		}

	// TODO: handle StaticWrapControl
	CalcLines(text, orig_xmin /*unused*/)
		{
		lines = text is "" ? #(" ") : text.Lines()
		return Max(lines.Size(), 1)
		}

	setBgColor()
		{
		.bgColor = .whitebgnd
			? CLR.white
			: .bgndcolor is ''
				? CLR.ButtonFace
				: .bgndcolor
		if Object?(.bgndcolor)
			.bgColor = RGB(@.bgndcolor)
		.bgBrush = CreateSolidBrush(.bgColor)
		}

	Get()
		{ return .untranslated }

	Set(value, logfont = false, refreshRequired? = false)
		{
		prevLines = .text.Lines().Size()
		value = String(value)
		.untranslated = value
		value = value.Tr('\r').Replace('\n', '\r\n')
		.setTabLength(value)
		.text = TranslateLanguage(value)
		curLines = .text.Lines().Size()
		changed? = .textChanged?()
		if logfont isnt false
			.SetLogFont(logfont, "M")
		else if prevLines isnt curLines
			.MeasureSizes('M')
		super.Set(.text)
		.calcX(.text is '' ? Object(' ') : .text.Lines())
		if logfont isnt false or prevLines isnt curLines // for ymin
			{
			.AdjustControlSize()
			.setsize()
			refreshRequired? = false
			}
		.refresh(refreshRequired?, changed?)
		}

	setTabLength(value)
		{
		if value.Has?('\t') // 8 Dialog units = 4 spaces
			SendMessagePoint(.Hwnd, EM.SETTABSTOPS, 1, Object(x: 8))
		}

	textChanged?()
		{ return .constructing? ? false : GetWindowText(.Hwnd) isnt .text }

	calcX(lines)
		{
		if .OrigXmin isnt 0
			return
		.WithSelectObject(.GetFont())
			{ |hdc|
			.Xmin = lines.Map({
				DrawTextEx(hdc, it, -1, xy = Object(), DT.CALCRECT | DT.NOPREFIX, 0)
				xy.right
				}).Max()
			}
		}

	refresh(refreshRequired?, changed?)
		{
		if not .InListEdit?() and (refreshRequired? or changed?)
			.WindowRefresh()
		}

	SetColor(color)
		{
		if .setColor(color)
			.Repaint()
		}

	GetColor()
		{
		return .color
		}

	CTLCOLORSTATIC(wParam)
		{
		if .color isnt ""
			SetTextColor(wParam, .color)
		SetBkColor(wParam, .bgColor)
		return .bgBrush
		}

	// SetReadOnly already handled by EditControl
	// GetReadOnly already handled by EditControl
	// SetFont should not be required
	// SetBgndColor already handled by EditControl
	// SetBgndBrush non usage found

	KEYDOWN()
		{
		return 0
		}

	EN_SETFOCUS()
		{
		.SETCURSOR()
		return 0
		}

	LBUTTONDOWN()
		{
		if 0 isnt .Send('Static_Click')
			return 0
		return 'callsuper' // for selecting text
		}

	LBUTTONUP()
		{
		.Send('Static_LButtonUp')
		return 'callsuper'
		}

	LBUTTONDBLCLK()
		{
		result = .Send('Static_DoubleClick')
		return result is 0 ? 'callsuper' : result
		}

	getter_cursor()
		{
		.cursor = LoadCursor(NULL, IDC.IBEAM)
		}

	SETCURSOR()
		{
		HideCaret(.Hwnd)
		SetCursor(.cursor)
		return 0
		}

	DevMenu: #()
	FieldMenu: #('Copy', 'Select &All')
	ContextMenu(x, y)
		{
		result = .Send('Static_ContextMenu', :x, :y)
		return result is 0 ? super.ContextMenu(x, y) : result
		}

	focusRect?: false
	DrawFocusRect(state = false) // used by LinkButton
		{
		.focusRect? = state
		.Repaint()
		.SubClass() // to handle PAINT
		}

	focusRect: false
	PAINT(wParam, lParam)
		{
		// NOTE: this is only used if we SubClass
		// either because of .DrawFocusRect or because of tip
		// it cannot use Begin/EndPaint (which clips)
		// because it draws the focus rectangle outside the client area
		// AND because it does Callsuper which would nest Begin/EndPaint
		WithDC(.Hwnd)
			{|dc|
			if .focusRect isnt false
				DrawFocusRect(dc, .focusRect) // remove
			.focusRect = false
			.Callsuper(.Hwnd, WM.PAINT, wParam, lParam)
			if .focusRect?
				DrawFocusRect(dc, .focusRect = .getFocusRect()) // add
			}
		return 0
		}

	getFocusRect()
		{
		rc = GetWindowRect(.Hwnd)
		--rc.left; --rc.top; ++rc.right; ++rc.bottom
		ScreenToClient(.Hwnd, pt = Object(x: rc.left, y: rc.top))
		rc.left = pt.x; rc.top = pt.y
		ScreenToClient(.Hwnd, pt = Object(x: rc.right, y: rc.bottom))
		rc.right = pt.x; rc.bottom = pt.y
		return rc
		}

	SetFont(font = "", size = "", weight = "", underline = false)
		{
		if .constructing?
			{
			super.SetFont(.font, .size, .weight, "M", .underline)
			return
			}
		if font isnt ""
			.font = font
		if size isnt ""
			.size = size
		if weight isnt ""
			.weight = weight
		if underline is true
			.underline = true
		.fontChanged = true
		.changeFont() // call it first to avoid alignment issue with other fields
		.WindowRefresh()
		}

	changeFont()
		{
		super.SetFont(.font, .size, .weight, "M", .underline)
		.SetMargins(Object(left: 0, right: 0)) // lost margin after WM.SETFONT message
		.alignVertPadding() // avoid text jumping up and down
		.AdjustControlSize()
		.setsize()
		}

	MOUSEWHEEL(wParam, lParam)
		{
		return .WndProc.Callsuper(.Hwnd, WM.MOUSEWHEEL, wParam, lParam)
		}

	Destroy()
		{
		DeleteObject(.bgBrush)
		super.Destroy()
		}
	}
