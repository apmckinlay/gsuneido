// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// base class for FieldControl (single line) and EditorControl (multi-line)
// WndProc instead of Hwnd so derived controls can subclass if necessary
WndProc
	{
	Name:			'Edit'
	Status:			''
	DefaultWidth: 20
	DefaultHeight: 2

	New(.mandatory = false, .readonly = false, style = 0, bgndcolor = "",
		.textcolor = "", hidden = false, tabover = false,
		font = "", size = "", weight = "", underline = false,
		width = false, height = false, cue = false, readOnlyBgndColor = false,
		.static = false, status = '')
		{
		.createWindow(static, hidden, readonly, style, tabover)
		.SubClass()
		if not .static
			.Send("Data")
		.Map = Object()
		.Map[EN.KILLFOCUS] = 'EN_KILLFOCUS'
		.Map[EN.CHANGE] = 'EN_CHANGE'

		.Map[EN.SETFOCUS] = 'EN_SETFOCUS' // differentiate from WM_
		.err_brush = CreateSolidBrush(CLR.ErrorColor)
		.SetReadOnlyBrush(readOnlyBgndColor)
		if Object?(bgndcolor)
		   bgndcolor = RGB(@bgndcolor)
		.SetBgndColor(bgndcolor)
		if Object?(textcolor)
		   .textcolor = RGB(@textcolor)

		if cue isnt false
			.SetCue(cue)

		.SetFontAndSize(font, size, weight, underline, width, height)
		.AdjustControlSize()
		.ContextExtra = Object()
		if status > ''
			.SetStatus(status)
		}

	createWindow(static, hidden, readonly, style, tabover)
		{
		exStyle = .InListEdit?() ? 0 : WS_EX.CLIENTEDGE
		.SetHidden(hidden)
		if hidden is false
			style |= WS.VISIBLE
		if static
			{
			readonly = true
			exStyle = 0
			}
		if readonly is true
			style |= ES.READONLY
		if tabover is false and (readonly isnt true or static)
			style |= WS.TABSTOP

		// need w and h or else initial text is scrolled horizontally
		.CreateWindow("edit", NULL, style, :exStyle, w: 1000, h: 100)
		}
	readOnlyBrush: false
	readOnlyBgndColor: false
	SetReadOnlyBrush(color = false)
		{
		if color is false
			return
		if Object?(color)
			color = RGB(color[0], color[1], color[2])
		.readOnlyBgndColor = color
		if .readOnlyBrush isnt false
			DeleteObject(.readOnlyBrush)
		.readOnlyBrush = color is CLR.WHITE
			? GetStockObject(SO.WHITE_BRUSH)
			: CreateSolidBrush(color)
		}
	SetFontAndSize(font, size, weight, underline, width, height, text = "M")
		{ // overridden by NumberControl, which also uses text parameter
		.OrigYmin = ymin = .Ymin
		.OrigXmin = xmin = .Xmin
		.SetFont(font, size, weight, :underline, :text)
		margins = .GetMargins()
		.Xmin = xmin isnt 0 ? xmin
			: .Xmin * (width is false ? .DefaultWidth : width) +
				margins.left + margins.right
		.Ymin = ymin isnt 0 ? ymin
			: .Ymin * (height is false ? .DefaultHeight : height)
		}
	AdjustControlSize()
		{
		ymin = .Ymin
		r = Object(bottom: .Ymin, right: .Xmin)
		style = .GetStyle()
		AdjustWindowRectEx(r, style, false, .GetExStyle())
		if WS.VSCROLL is (style & WS.VSCROLL)
			r.right += GetSystemMetrics(SM.CXVSCROLL)
		.Xmin = r.right - r.left + 2
		.Ymin = r.bottom - r.top + 4 /* = offset */
		if not .static // no border when static
			.Top += (.Ymin - ymin) / 2 - 1 // border size
		}
	HorzAdjustInListEdit: 2
	Resize(x, y, w, h)
		{
		if .InListEdit?()
			{
			// align text in editor with text in list
			// so it doesn't appear to move when editor pops up
			x -= .HorzAdjustInListEdit
			y += 2
			}
		super.Resize(x, y, w, h)
		}
	InListEdit?()
		{
		return .Controller.Base?(ListEditWindow) or
			(.Parent.Member?("Parent") and .Parent.Parent.Base?(ListEditWindow))
		}
	On_Delete()
		{
		SendMessage(.Hwnd, WM.CLEAR, 0, 0)
		return 0
		}
	On_Cut()
		{
		SendMessage(.Hwnd, WM.CUT, 0, 0)
		}
	On_Copy()
		{
		sel = SendMessage(.Hwnd, EM.GETSEL, NULL, NULL)
		if (LOWORD(sel) is HIWORD(sel))
			.SelectAll()
		SendMessage(.Hwnd, WM.COPY, 0, 0)
		}
	On_Paste()
		{
		SendMessage(.Hwnd, WM.PASTE, 0, 0)
		}
	On_Undo()
		{
		SendMessage(.Hwnd, WM.UNDO, 0, 0)
		}
	GetSel()
		{
		x = .SendMessage(EM.GETSEL)
		return [LOWORD(x), HIWORD(x)]
		}
	GetSelText()
		{
		s = GetWindowText(.Hwnd)
		range = .GetSel()
		return s[range[0]..range[1]]
		}
	SetSel(start, end)
		{
		.SendMessage(EM.SETSEL, start, end)
		}
	EN_CHANGE()
		{
		.Send("Edit_Change")
		return 0
		}
	ReplaceSel(text)
		{
		SendMessageTextIn(.Hwnd, EM.REPLACESEL, true, text)
		}
	Valid?() // derived classes should define this
		{
		return .validCheck?(.Get(), .mandatory)
		}

	validCheck?(data, mandatory)
		{
		return not (mandatory and data is "")
		}

	ValidData?(@args)
		{
		return .validCheck?(args[0], args.GetDefault('mandatory', false))
		}

	On_Select_All()
		{
		.SelectAll()
		}
	SelectAll()
		{
		.SendMessage(EM.SETSEL, 0, -1)
		}
	Dirty?(dirty = "")
		{
		Assert(dirty is true or dirty is false or dirty is "")
		if (dirty isnt "")
			SendMessage(.Hwnd, EM.SETMODIFY, dirty, 0)
		return SendMessage(.Hwnd, EM.GETMODIFY, 0, 0) is 1
		}
	SetReadOnly(readOnly)
		// pre:		readOnly is a Boolean value
		// post:	this is read-only iff readOnly is true
		{
		if (.readonly)
			return
		Assert(Boolean?(readOnly))
		SendMessage(.Hwnd, EM.SETREADONLY, readOnly, 0)
		}
	SetDefaultReadOnly(readonly, controllerReadOnly)
		{
		.readonly = readonly
		t = GetWindowLong(.Hwnd, GWL.STYLE)
		SetWindowLong(.Hwnd, GWL.STYLE, readonly ? t ^ WS.TABSTOP : t | WS.TABSTOP)
		SendMessage(.Hwnd, EM.SETREADONLY, controllerReadOnly or readonly, 0)
		}
	GetReadOnly()
		// post:	returns true iff this is read-only
		{
		return .HasStyle?(ES.READONLY)
		}
	ShowBalloonTip(msg, icon = 'NONE')
		{
		bt = Object(cbStruct: EDITBALLOONTIP.Size(),
			pszText: MultiByteToWideChar(msg),
			ttiIcon: TTI[icon])
		SendMessageEDITBALLOONTIP(.Hwnd, EM.SHOWBALLOONTIP, 0, bt)
		}
	HideBalloonTip()
		{
		.SendMessage(EM.HIDEBALLOONTIP)
		}
	PosFromChar(i)
		{
		pos = .SendMessage(EM.POSFROMCHAR, i)
		return [x: LOSWORD(pos), y: HISWORD(pos)]
		}
	SetTextColor(.textcolor)
		{
		if Object?(textcolor)
		   .textcolor = RGB(@textcolor)
		.Repaint()
		}
	CTLCOLOREDIT(wParam)
		{
		if .textcolor isnt ""
			SetTextColor(wParam, .textcolor)
		SetBkColor(wParam, .color)
		return .brush
		}
	CTLCOLORSTATIC(wParam)
		{
		if .GetReadOnly() and .readOnlyBrush isnt false
			{
			SetBkColor(wParam, .readOnlyBgndColor)
			return .readOnlyBrush
			}
		return 'callsuper'
		}
	EN_KILLFOCUS()
		{
		.KillFocus()
		if (.Send("Dialog?") isnt true and GetFocus() isnt .Hwnd)
			{
			if 0 is valid? = .Send('Edit_ParentValid?')
				valid? = .Valid?()
			.SetValid(valid?)
			if not valid? and not .GetReadOnly()
				Beep()
			}
		if .Status > ""
			.Send("Status", "")
		return 0
		}
	KillFocus()
		{
		}
	EN_SETFOCUS()
		{
		.Send('Field_SetFocus')
		if .Status > ""
			.Send("Status", .Status)
		.SetValid() // don't color invalid when focused
		return 0
		}
	SetValid(valid? = true, force = false)
		{
		if GetFocus() is .Hwnd and not force
			valid? = true // don't color when we have focus
		.brush = valid? ? .bgndbrush : .err_brush
		.color = valid? ? .bgndcolor : CLR.ErrorColor
		.Repaint()
		// handle non-client frame
		SetWindowPos(.Hwnd, 0, 0, 0, 0, 0,
			SWP.FRAMECHANGED | SWP.NOSIZE | SWP.NOMOVE | SWP.NOZORDER)
		}
	SetBgndColor(color)
		{
		.bgndcolor = color is "" ? CLR.WHITE : color
		.bgndbrush = color is ""
			? GetStockObject(SO.WHITE_BRUSH)
			: CreateSolidBrush(color)
		.color = .bgndcolor
		.brush = .bgndbrush
		.Repaint()
		}
	LBUTTONDBLCLK()
		{
		return 0 is .Send("Edit_DoubleClick") ? 'callsuper' : 0
		}
	GetMargins()
		{
		margins = .SendMessage(EM.GETMARGINS)
		return Object(left: LOWORD(margins), right: HIWORD(margins))
		}
	SetMargins(@margins)
		{
		if margins.Size() is 1 and Object?(margins[0])
			margins = margins[0]
		.SendMessage(EM.SETMARGINS, EC.LEFTMARGIN | EC.RIGHTMARGIN,
			MAKELONG(margins.left, margins.right))
		}
	Get()
		{
		return GetWindowText(.Hwnd)
		}
	Set(value)
		{
		if (not String?(value))
			value = Display(value)
		.Dirty?(false)
		SetWindowText(.Hwnd, value)
		}
	SetCue(cue)
		{
		.SendMessageTextIn(EM.SETCUEBANNER, 0, MultiByteToWideChar(cue))
		}
	AddContextMenuItem(name, runFunc, enabledFunc = function () { #(addToMenu:) })
		{
		if Object?(name)
			{
			// Handle Cascade menu options
			n = Object()
			for idx in name.Members()
				n.Add(Object(name: name[idx], runFunc: runFunc[idx]))
			.ContextExtra.Add(Object(name: n, :enabledFunc))
			}
		else
			.ContextExtra.Add(Object(:name, :runFunc, :enabledFunc))
		}
	FieldMenu: ('&Undo\tCtrl+Z', '',
		'Cu&t\tCtrl+X', '&Copy\tCtrl+C', '&Paste\tCtrl+V', '&Delete', '',
		'Select &All\tCtrl+A')
	ContextMenu(x, y)
		{
		// If the context menu gets called while navigating off the page then
		// SetFocus call will throw an error. Due to dynamic controls; Losing focus will
		// cause a change in said controls, so we need to check if Destroyed? on both ends
		if .Destroyed?()
			return 0

		.SetFocus()
		if .Destroyed?()
			return 0

		if x is 0 and y is 0 // keyboard
			{
			pt = Object(x: 10, y: 20)
			ClientToScreen(.Hwnd, pt)
			x = pt.x
			y = pt.y
			}

		menu = .buildContextMenu()
		i = ContextMenu(menu).Show(.Hwnd, x, y) - 1
		if i is -1 or .Destroyed?()
			return 0
		.callContext(menu, i)
		return 0
		}

	buildContextMenu()
		{
		menu = .FieldMenu.Copy()
		for item in .ContextExtra
			{
			enabledOb = (item.enabledFunc)()
			if enabledOb.addToMenu is true
				{
				if Object?(item.name)
					menu.Add(item.name)
				else
					menu.Add(Object(name: item.name,
						state: enabledOb.GetDefault('state', MFS.ENABLED)
						runFunc: item.runFunc))
				}
			}
		if Suneido.User is 'default'
			menu.Add(@.DevMenu)
		return menu
		}

	callContext(menu, chosen, j = 0)
		{
		for (m = 0; j <= chosen and m < menu.Size(); ++j, ++m)
			{
			if .isCascadeMenu(menu[m])
				{
				if menu[m].Member?('name')
					j = .callContext(menu[m].name, chosen, j) - 1
				else
					j = .callContext(menu[m], chosen, j) - 1
				}
			else if (j is chosen)
				{
				// if we have runFunc call it, otherwise .send a msg
				if Object?(menu[m])
					{
					(menu[m].runFunc)(source: this)
					continue
					}
				.ContextMenuCall(menu[m])
				}

			}
		return j
		}
	isCascadeMenu(menu)
		{
		return (Object?(menu) and
				(not menu.Member?('name') or Object?(menu.name)))
		}
	SetStatus(status)
		{
		.Status = status
		.ToolTip(.Status)
		}

	Destroy()
		{
		if not .static
			.Send("NoData")
		DeleteObject(.err_brush)
		DeleteObject(.bgndbrush)
		if .readOnlyBrush isnt false
			DeleteObject(.readOnlyBrush)
		super.Destroy()
		}
	}
