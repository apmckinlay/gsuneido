// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// modal/blocking dialog that runs a nested MessageLoop
// disables all other windows while open
// does NOT use the Windows Dialog system (e.g. DefDialogProc)
// (other than the IsDialogMessage that is standard in MessageLoop)
// NOTE: recommended usage is to pass 0 as parentHwnd
Window
	{
	CallClass(parentHwnd, control, style = 0, exStyle = 0, border = 5, title = false,
		posRect = false, keep_size = false, closeButton? = false, useDefaultSize = false)
		{
//		if parentHwnd isnt 0
//			Print('Dialog: parentHwnd is not used, should pass 0')
		parentHwnd = .GetParent()
//		if style isnt 0 or exStyle isnt 0
//			Print("Dialog: does not take style or exStyle, only closeButton?")
		style = WS.CAPTION | WS.SIZEBOX
		// WS.CAPTION is required with WS.SYSMENU based on msdn
		exStyle = .getExStyle(parentHwnd)
		if closeButton?
			{
			style |= WS.SYSMENU
			exStyle |= WS_EX.STATICEDGE | WS_EX.OVERLAPPEDWINDOW
			}
		w = super.CallClass(control, :title, :style, :exStyle, :border, :parentHwnd,
			show: false, skipStartup?:)
		return w.InternalRun(parentHwnd, control, keep_size, posRect, useDefaultSize)
		}
	getExStyle(parentHwnd)
		{
		curPID = GetCurrentProcessId()
		GetWindowThreadProcessId(parentHwnd, lpdwProcessId = Object())
		parentPID = lpdwProcessId.x
		// if the dialog starts from non-suneido window, like command line
		// should not specify WS_EX.TOOLWINDOW so it shows in the taskbar
		return curPID isnt parentPID ? 0 : WS_EX.TOOLWINDOW
		}
	// not meant for public use
	destroy?: true
	InternalRun(.parentHwnd, control, keep_size, posRect, useDefaultSize = false)
		{
		CreateContributionSysMenus('DialogMenus', this)
		.setPos(parentHwnd, control, keep_size, posRect, useDefaultSize)
		SetWindowPos(.Hwnd, HWND.TOP, 0, 0, 0, 0, SWP.NOSIZE | SWP.NOMOVE)
		.setDefaultButtonStyle()
		return .ActivateDialog()
		}
	inLoop: false
	ActivateDialog()
		{
		previouslyDisabled = .disableWindows(.Hwnd)
		.Show(SW.SHOWNORMAL)
		DoStartup(.Ctrl)
		.inLoop = true
		MessageLoop(.Hwnd)
		.inLoop = false
		.reenableWindows(previouslyDisabled)
		.ModalClose(.parentHwnd)

		if .destroy?
			DestroyWindow(.Hwnd)
		else
			.destroy? = true
		return .resultValue
		}
	Unblock()
		{
		if not .inLoop
			return false

		.destroy? = false
		.Result(false)
		}

	setDefaultButtonStyle()
		{
		if .Ctrl.DefaultButton isnt "" and
			false isnt (ctrl = .FindControl(ToIdentifier(.Ctrl.DefaultButton))) and
			(DLGC.BUTTON | DLGC.UNDEFPUSHBUTTON) is ctrl.SendMessage(WM.GETDLGCODE, 0, 0)
			ctrl.SendMessage(BM.SETSTYLE, BS.DEFPUSHBUTTON, true)
		}
	disableWindows(hdlg) // returns a list of the windows that were disabled
		{
		list = Object()
		disableWindowBlock = {|hwnd, unused|
			if hwnd isnt hdlg and IsWindowVisible(hwnd) and
				not EnableWindow(hwnd, false)
				list.Add(hwnd)
			true // continue enumerating
			}
		EnumThreadWindows(GetCurrentThreadId(), disableWindowBlock, NULL)
		ClearCallback(disableWindowBlock)
		return list
		}
	reenableWindows(list)
		{
		for hwnd in list
			EnableWindow(hwnd, true)
		}

	setPos(parentHwnd, control, keep_size, posRect, useDefaultSize)
		{
		.DialogCenterSize(parentHwnd, .GetWindowTitle(), control, keep_size,
			useDefaultSize)
		if posRect isnt false
			.alignToPosRect(posRect)
		}
	alignToPosRect(posRect)
		{
		r = GetWindowRect(.Hwnd)
		wa = GetWorkArea(posRect)
		.SetWinPos(@.calcPos(posRect, r.right - r.left, r.bottom - r.top, wa))
		}
	calcPos(posRect, width, height, wa)
		{
		x = posRect.left
		y = posRect.bottom
		// shrink to fit screen
		width = Min(width, wa.right - wa.left)
		height = Min(height, wa.bottom - wa.top)
		if x + width > wa.right
			x -= ((x + width) - wa.right) // push left to fit on screen
		if y + height > wa.bottom
			{ // if the dialog won't fit below
			jump_height = posRect.bottom - posRect.top
			jump_width = posRect.right - posRect.left
			if y - (height + jump_height) > wa.top
				y -= height + jump_height // above
			else if ((x + jump_width + width) < wa.right)
				{ // to the right
				x += jump_width
				y = .move_into_workarea(y - height / 2, wa, height)
				}
			else if ((x - width) > wa.left)
				{ // to the left
				x -= width
				y = .move_into_workarea(y - height / 2, wa, height)
				}
			}
		return Object(:x, :y)
		}
	move_into_workarea(y, wr, height)
		{
		if y < wr.top
			y = wr.top
		if y + height > wr.bottom
			y = wr.bottom - height
		return y.Int()
		}

	Call(hwnd, msg, wParam, lParam)
		{
		// make the default button work
		if msg is DM.GETDEFID
			return MAKELONG(.Mapcmd(.Ctrl.DefaultButton $ "_button"), DC_HASDEFID)
		super.Call(hwnd, msg, wParam, lParam)
		}

	ACTIVATE(wParam)
		{
		// TEMPORARILY accessing private member of Window, this code should be in Window
		if LOWORD(wParam) isnt WA.INACTIVE and .Window_focus is 0
			{
			hwnd = GetNextDlgTabItem(.Hwnd, 0, false) // first control with TABSTOP
			if hwnd is 0
				hwnd = GetNextDlgGroupItem(.Hwnd, 0, false)
			if IsWindow(hwnd)
				SetFocus(hwnd)
			}
		super.ACTIVATE(wParam)
		}

	resultValue: false
	Result(value)
		{
		.resultValue = value
		.end_dialog()
		}
	end_dialog()
		{
		PostMessage(.Hwnd, WM.NULL, END_MESSAGE_LOOP, END_MESSAGE_LOOP)
		}

	CLOSE()
		{
		if not .AllowCloseWindow?()
			return 0
		.closeDialog()
		return 0 // meaning we handled it so don't do default DestroyWindow
		}
	Destroy() // called by Controller.On_Close
		{
		if not .AllowCloseWindow?()
			return
		.closeDialog()
		}
	closeDialog()
		{
		if .inLoop
			{
			.destroy? = true
			.Result(false)
			}
		else
			DestroyWindow(.Hwnd)
		}
	DoWithWindowsDisabled(block, exclude = 0)
		{
		hwndList = .disableWindows(exclude)
		return Finally(block, { .reenableWindows(hwndList) })
		}
	}
