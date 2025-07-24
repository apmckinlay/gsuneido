// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
WindowBase
	{
	parent: 		false
	keyboardHook:	false
	row:			false
	col:			false
	msgCount:		0

	New(control, readonly, .col, .row, .parent, cellRect, custom = false,
		.customFields = false)
		{
		style = WS.BORDER | WS.POPUP
		exStyle = WS_EX.TOOLWINDOW
		cellRect.left--
		cellRect.right++
		AdjustWindowRectEx(cellRect, style, false, exStyle)
		.Hwnd = CreateWindowEx(exStyle, "SuBtnfaceArrow", "", style,
			cellRect.left, cellRect.top, cellRect.right - cellRect.left,
			cellRect.bottom - cellRect.top, parent.Window.Hwnd, 0, Instance(), 0)
		if .Hwnd is 0
			throw "CreateWindow failed"
		.SubClass()
		.createControl(control, readonly, style, exStyle, custom)
		.SetVisible(true)
		}
	createControl(control, readonly, style, exStyle, custom)
		{
		name = .parent.GetCol(.col)
		ctrl = .ctrlProperties(control, name, custom, readonly)

		// needed so the pop up control will inherit .Custom from this window
		// - needed for KeyControl using CustomizableMap
		.Custom = .customFields
		.Ctrl = .Construct(ctrl)
		if (ctrl.Member?('readonly') and
			ctrl.readonly is true and not .Ctrl.GetReadOnly())
			.Ctrl.SetReadOnly(true)
		.setInitialCtrlSize(style, exStyle)

		.Send("Setup")				// call Controller's Setup Method
		dataRow = .parent.GetRow(.row)
		data = dataRow[name]	// get the data
		if not .Ctrl.GetReadOnly() and
			'' isnt invalidVal = ListControl.GetInvalidFieldData(dataRow, name)
			data = invalidVal
		.Ctrl.Set(data)
		}
	ctrlProperties(control, name, custom, readonly)
		{
		field = false
		if 0 is control						// no control set
			{
			field = Datadict(name)
			control = field.Control
			}
		ctrl = Object?(control) ? control.Copy() : Object(control)
		if field isnt false and custom isnt false
			ctrl.Merge(custom)
		ctrl.name = name
		if (readonly)
			ctrl.readonly = true
		return ctrl
		}
	setInitialCtrlSize(style, exStyle)
		{
		r = .GetClientRect().ToWindowsRect()
		if (.Ctrl.Member?("Controller_ctrl") and .Ctrl.Controller_ctrl isnt false)
			{
			r.right = .Ctrl.Xmin
			r.bottom = .Ctrl.Ymin
			.adjustWindow(r, style, exStyle)
			}
		else
			{
			 if (r.bottom < .Ctrl.Ymin)
				{
				r.bottom = .Ctrl.Ymin
				.adjustWindow(r, style, exStyle)
				r = .GetClientRect().ToWindowsRect()
				}
			.Ctrl.Resize(0, 0, r.right - r.left, r.bottom - r.top)
			}
		}
	adjustWindow(r, style, exStyle)
		{
		AdjustWindowRectEx(r, style, false, exStyle)
		w = r.right - r.left
		h = r.bottom - r.top
		ClientToScreen(.Hwnd, pt = Object(x: r.left, y: r.top))
		r = .parent.GetClientRect().ToWindowsRect()
		ClientToScreen(.parent.Hwnd, ppt = Object(x: 0, y: 0))
		x = Min(pt.x, ppt.x + r.right - w)
		y = Min(pt.y, ppt.y + r.bottom - h)
		SetWindowPos(.Hwnd, NULL, x, y, w, h, SWP.NOZORDER | SWP.NOACTIVATE)
		}
	sending?: false
	sendToParent(dir)
		{
		if .sending? is true
			return
		.sending? = true // .Destroy will clear the flag
		if (.msgCount++ is 0 and .parent isnt false)
			{
			.ClearFocus()			// so control gets KillFocus and validates
			.listCommit(dir)
			}
		.Destroy()
		}
	dirty?: false
	Msg(args)
		{
		// TODO: should maybe only pass on specific messages
		// (like Status and GetField)
		// instead of only stopping certain messages
		msg = args[0]
		if msg isnt 'Data' and msg isnt 'NewValue' and .parent.Member?('Controller')
			{
			if ((msg is "GetField" or msg is "SetField") and
				not .parent.Controller.Base?(BrowseControl))
				return .parent[msg](@+1 args)
			if .parent.Controller.Base?(VirtualListEdit)
				return .parent.Controller[msg](@+1 args)
			if .parent.Controller.Method?(msg)
				return .parent.Controller[msg](@+1 args)
			}
		if msg is 'NewValue'
			.dirty? = true
		return 0
		}
	Return()
		{
		.listCommit()
		}
	listCommit(dir = 0)
		{
		if not this.Member?('Ctrl') or .Ctrl.Empty?() // got destroyed somehow already
			return

		// setting _hwnd to parent hwnd in case the application book is no longer the
		// active window (user switches to another application such as a web browser)
		_hwnd = .parent.Window.Hwnd
		val = .Ctrl.Get()
		valid? = .Ctrl.Valid?()
		readonly = .Ctrl.GetReadOnly()
		dirty? = .dirty?
		unvalidated_val = not valid? and .Ctrl.Method?('GetUnvalidated')
			? .Ctrl.GetUnvalidated()
			: ""
		parent = .parent
		col = .col
		row = .row
		.Destroy()
		parent.ListEditWindow_Commit(col, row, dir, val, valid?,
			:unvalidated_val, :readonly, :dirty?)
		if dir is 0
			SetFocus(parent.Hwnd)
		}
	ChildOf?(hwnd)
		{
		isChild? = false
		EnumChildWindows(.Hwnd)
			{ |childHwnd|
			enumNext? = true
			if childHwnd is hwnd
				{
				isChild? = true
				enumNext? = false
				}
			enumNext?
			}
		return isChild? or (hwnd is .Hwnd)
		}

	ENABLE(wParam)
		{
		if not .parent.Member?('Window')
			{
			SuneidoLog('INFO: ListEditWindow does not have parent Window',
				params: Object(ctrl: .Ctrl, parent: .parent))
			return 0
			}

		EnableWindow(.parent.Window.Hwnd, wParam isnt 0)
		return 0
		}

	ClosingListEdit: false
	ACTIVATE(wParam, lParam)
		{
		if (LOWORD(wParam) is WA.INACTIVE)
			{
			.endHook()
			.Send("Inactivate")
			if lParam is NULL
				.sendToParent(0)
			}
		else
			{
			if .ClosingListEdit
				{
				// this also destroys popup window created from ListEditWindow,
				// e.g. ChooseList
				.Return()
				return 0
				}
			.startHook()
			.Send("Activate")
			if (not .ChildOf?(GetFocus()))
				SetFocus(GetNextDlgTabItem(.Hwnd, NULL, false))
			}
		return 0
		}

	startHook()
		{
		if .destroying or .keyboardHook isnt false
			return
		.keyboardHook = SetWindowsHookEx(WH.KEYBOARD,
			.HookFunc, 0, GetCurrentThreadId())
		if (.keyboardHook is NULL)
			{
			.keyboardHook = false
			throw "unable to hook focus messages"
			}
		}
	HookFunc(nCode, wParam, lParam)
		{
		try _hwnd = .parent.Window.Hwnd // parent since the edit window gets destroyed
		if ((lParam & 0x80000000) is 0 and not .Ctrl.Empty?())	// is key being pressed?
			{
			if (wParam is VK.TAB)
				{
				if KeyPressed?(VK.SHIFT)
					.sendToParent(-1)
				else if (false is .Ctrl.HandleTab())	// if control doesn't handle tabs
					.sendToParent(1)
				}
			else if ((wParam is VK.RETURN and not KeyPressed?(VK.CONTROL)) or
				(wParam is VK.F4 and (((lParam >> 29) & 1) is 1)))
				.sendToParent(0)		// return pressed without ctrl, or alt-F4
			}
		return .keyboardHook is false ? 0 :
			CallNextHookEx(.keyboardHook, nCode, wParam, lParam)
		}

	On_Cancel(@unused)
		{
		if .Ctrl.GetReadOnly() is true
			.Destroy()
		}

	endHook()
		{
		if (.keyboardHook isnt false)
			{
			if (not UnhookWindowsHookEx(.keyboardHook))
				throw "unable to unhook focus messages"
			ClearCallback(.HookFunc)
			.keyboardHook = false
			}
		}

	destroying: false
	Destroy()
		{
		if .destroying
			return
		.destroying = true
		.endHook()
		super.Destroy()
		}
	}