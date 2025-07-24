// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// abstract base for Window and Dialog
WndProc
	{
	New(border = 0)
		{
		.border = ScaleWithDpiFactor(border)
		.cmdmap = Object().Set_default("")
		.rcmdmap = Object().Set_default(false)
		.mapcmd("On_OK", ID.OK)
		.mapcmd("On_Cancel", ID.CANCEL)
		.Window = .Controller = this
		.HwndMap = Object()
		.move_observers = Object()
		Suneido.OpenWindows = Suneido.GetDefault('OpenWindows', 0) + 1
		}
	New2()
		{ }

	// title can be set in multiple ways, in order of precedence:
	// - argument to Window or Dialog e.g. Dialog(... title: ...)
	// - title or Title member on the control specification Dialog(0, [Editor title: ...])
	// - Title on the constructed control
	FindTitle(title, control)
		{
		if title isnt false
			return title
		if Object?(control)
			{
			if control.Member?(#Title)
				return control.Title
			if control.Member?(#title)
				return control.title
			}
		return .Ctrl.Title
		}

	GetParent() // used by Dialog and ModalWindow
		{
		try
			parentHwnd = _hwnd
		catch
			parentHwnd = 0
		if parentHwnd is 0 or not IsWindow(parentHwnd) or parentHwnd is GetDesktopWindow()
			{
			parentHwnd = GetActiveWindow()
			if ((GetWindowLong(parentHwnd, GWL.STYLE) & WS.POPUP) isnt 0)
				{
				// use GetAncestor (instead of GetParent) to exclude OWNER windows
				// since OWNER WINDOW is an overlapped or pop-up window,
				// and they are destroyed when another window takes focus
				parentHwnd = GetAncestor(parentHwnd, GA.PARENT)
				}
			}
		else
			parentHwnd = GetAncestor(parentHwnd, GA.ROOT)
		if not IsWindow(parentHwnd)
			{
			Print("ERROR: GetParent: failed to get parentHwnd")
			return 0
			}
		return parentHwnd
		}

	// this will also cause the size to be saved on destroy
	DialogCenterSize(parentHwnd, title, control, keep_size, useDefaultSize = false)
		{
		if useDefaultSize is true
			.SetWinSize(@.DefaultSize())
		else
			.SetWinSize(@.RequiredWindowSize())
		if .Ctrl.Xstretch > 0 or .Ctrl.Ystretch > 0
			.restoreKeepSize(parentHwnd, title, control, keep_size)
		.center(parentHwnd)
		}
	DefaultSize()
		{
		dpiFactor = GetDpiFactor()
		return Object(w: 1024 * dpiFactor, h: 768 * dpiFactor)
		}
	center(parentHwnd) // used by Dialog and ModalWindow
		{
		parentRect = .parentRect(parentHwnd)
		rc = GetWindowRect(.Hwnd)
		width = rc.right - rc.left
		height = rc.bottom - rc.top
		// center on parent
		left = parentRect.left + (parentRect.right - parentRect.left - width) / 2
		top = parentRect.top + (parentRect.bottom - parentRect.top - height) / 2
		// ensure it's all visible
		// order is significant if bigger than work area
		wr = GetWorkArea(parentRect)
		if (left + width > wr.right)
			left = wr.right - width
		if (left < wr.left)
			left = wr.left
		if (top + height > wr.bottom)
			top = wr.bottom - height
		if (top < wr.top)
			top = wr.top
		.SetWinPos(left, top)
		return Object(:left, :top)
		}
	parentRect(parentHwnd)
		{
		return parentHwnd is 0
			? GetWorkArea()
			: GetWindowRect(parentHwnd)
		}

	nextcmdnum: 102
	Mapcmd(cmd)
		{
		method = .cmd_to_method(cmd)
		// if a number has already been assigned, return it
		if (.rcmdmap[method] isnt false)
			return .rcmdmap[method]
		cmdnum = .nextcmdnum++
		.mapcmd(method, cmdnum)
		return cmdnum
		}
	mapcmd(method, cmdnum)
		{
		.rcmdmap[method] = cmdnum
		.cmdmap[cmdnum] = method
		}
	Cmdmap(num)
		{
		return .cmdmap[num]
		}
	cmd_to_method(cmd)
		{
		cmd = cmd.BeforeFirst("\t")
		return "On_" $ ToIdentifier(cmd)
		}
	COMMAND(wParam, lParam) /*internal*/
		{
		id = LOWORD(wParam)
		if (lParam is NULL or
			(HIWORD(wParam) is 0 and id isnt 0))
			{
			if (.cmdmap.Member?(id))
				{
				cmd = .cmdmap[id]
				if (cmd.Suffix?("_button"))
					return super.COMMAND(wParam, lParam)
				// menu or accelerator command
				if IsWindowEnabled(.Hwnd)
					.Send(.cmdmap[id])
				}
			return 0
			}
		return super.COMMAND(wParam, lParam)
		}
	Send(@args)
		{
		args.source = this
		if (not .Member?('Ctrl'))
			return 0
		return .Ctrl.Msg(args)
		}
	Msg(args) /*internal*/
		{
		msg = args[0]
		if (.Method?(msg))
			return this[msg](@+1 args)
		else
			return 0
		}

	ModalClose(parentHwnd) // used by Dialog and ModalWindow
		{
		if parentHwnd isnt 0 and not IsWindow(parentHwnd)
			SuneidoLog("INFO: WindowBase.ModalClose: parentHwnd invalid")
		else
			// this is needed to prevent the wrong window being activated
			// when you close a modal that has had another modal
			SetActiveWindow(parentHwnd)

		.disableClose(parentHwnd)
		}

	disableCloseDelay: 1000 // one second
	DisableCloseTimer: false
	disableClose(parentHwnd)
		{
		if not IsWindow(parentHwnd)
			return

		if Object?(openBooks = Suneido.GetDefault('OpenBooks', false))
			if false isnt w = openBooks.FindOne({ .parentBook?(parentHwnd, it) })
				{
				// Disable parent's Close for 1 second to avoid accidental close
				menu = GetSystemMenu(parentHwnd, false)
				EnableMenuItem(menu, SC_CLOSE, MF.BYCOMMAND | MF.GRAYED)
				// Check for false because it will be false after the first kill
				if w.GetDefault('DisableCloseTimer', false) isnt false
					w.DisableCloseTimer.Kill()
				w.DisableCloseTimer = Delay(.disableCloseDelay)
					{ EnableMenuItem(menu, SC_CLOSE, MF.BYCOMMAND) }
				}
		}

	parentBook?(parentHwnd, book)
		{
		try
			return parentHwnd is book.Browser.Window.Hwnd
		return false
		}

	SetupCommands(control)
		{
		if Object?(control) and control.Member?(0)
			control = control[0]
		if String?(control)
			control = Global(control $ "Control")
		.commands = Object()
		if not control.Member?("Commands")
			return
		cmds = control.Commands
		if Function?(cmds)
			cmds = control.Commands()
		.AddCommands(cmds)
		}

	commands: false
	AddCommands(cmds)
		{
		if .commands is false
			.commands = Object()
		for cmd in cmds
			.commands[cmd[0]] = Object(
				accel: cmd.GetDefault(1, ""),
				help: cmd.GetDefault(2, ""),
				bitmap: cmd.GetDefault(3, cmd[0]),
				id: .Mapcmd(cmd[0]))
		return .SetupAccels(cmds)
		}
	Commands()
		{ return .commands }

	// accelerator stuff =============================================
	haccel: 0
	accels: ""
	getter_accelGroups()
		{
		return .accelGroups = Object()
		}
	SetupAccels(cmds)
		{
		newGroup = Object()
		for cmd in cmds
			if cmd.Member?(1)
				{
				id = .Mapcmd(cmd[0])
				accelInfo = .make_accel(cmd[1], id)
				newGroup.Add(accelInfo)
				}
		.accelGroups.Add(newGroup)
		.SetAccels()
		return newGroup
		}
	AddAccel(name, s)
		{
		id = .Mapcmd(name)
		.commands[name] = Object(help: '', bitmap: '', accel: s, :id)
		accelInfo = .make_accel(s, id)
		.accelGroups.Add(Object(accelInfo))
		.accels $= accelInfo.accel
		}
	make_accel(s, cmd)
		{
		ctrl = alt = shift = false
		ac = Object(fVirt: 0, :cmd)
		if (s =~ "Ctrl[+]")
			{
			ac.fVirt |= FCONTROL
			s = s.Replace("Ctrl[+]", "")
			ctrl = true
			}
		if (s =~ "Alt[+]")
			{
			ac.fVirt |= FALT
			s = s.Replace("Alt[+]", "")
			alt = true
			}
		if (s =~ "Shift[+]")
			{
			ac.fVirt |= FSHIFT
			s = s.Replace("Shift[+]", "")
			shift = true
			}
		if (VK.Member?(s.Upper()))
			{
			ac.fVirt |= FVIRTKEY
			ac.key = VK[s.Upper()]
			}
		else if (s.Size() is 1 and ' ' <= s and s <= "~")
			{
			if ((ac.fVirt & FSHIFT) is 0)
				s = s.Lower()
			ac.key = s.Asc()
			if ((ac.fVirt & FCONTROL) isnt 0 and 'a' <= s and s <= 'z')
				ac.key -= 0x60
			}
		else
			return Object(accel: "", id: cmd, :ctrl, :alt, :shift, key: "")
		return Object(accel: ACCEL(ac), id: cmd, :ctrl, :alt, :shift, key: ac.key)
		}
	RestoreAccels(accels)
		{
		.accelGroups.Remove(accels)
		.SetAccels()
		}
	QueryAccel(key, ctrl, alt, shift)
		{
		for group in .accelGroups
			for item in group
				if item.key is key and item.ctrl is ctrl and item.alt is alt and
					item.shift is shift
					return item.id
		return false
		}
	ResetAccels()
		{
		if .haccel isnt 0
			DestroyAcceleratorTable(.haccel)
		.haccel = 0
		}
	SetAccels()
		{
		.ResetAccels()
		if .destroying is true
			return
		.buildAccel()
		if .accels isnt ""
			{
			.haccel = CreateAcceleratorTable(.accels, .accels.Size() / ACCEL.Size())
			Assert(.haccel isnt 0)
			if .haccel is 1
				Print("WARNING: " $ Display(this) $ "
					got an HACCEL equal to 1, will kill MessageLoop()!")
			SetWindowLongPtr(.Hwnd, GWL.USERDATA, .haccel)
			}
		}
	buildAccel()
		{
		.accels = ""
		for group in .accelGroups
			for item in group
				.accels $= item.accel
		}
	// end of accelerator stuff ======================================

	GetChildren()
		{
		return [.Ctrl]
		}
	SetTitle(text)
		{
		SetWindowText(.Hwnd, text)
		}
	GetTitle()
		{
		return GetWindowText(.Hwnd)
		}
	AddToTitle(text)
		{
		SetWindowText(.Hwnd, text $ " - " $ .Ctrl.Title)
		}

	Refresh()
		{
		// needed to combine multiple screen refreshes
		// into one to reduce screen flickering
		if .refreshTimer isnt false
			return
		.refreshTimer = Defer(.refresh)
		}
	refreshTimer: false
	killRefresh()
		{
		if .refreshTimer is false
			return
		.refreshTimer.Kill()
		.refreshTimer = false
		}
	refresh()
		{
		.refreshTimer = false
		if not .Member?(#Ctrl) or .destroying
			return
		.BottomUp(#Recalc)
		GetClientRect(.Hwnd, rc = Object())
		.SIZE(rc.right | (rc.bottom << 16))
		.resizeToMin()
		}
	// run the delayed refresh before ShowWindow/UpdateWindow to avoid flickering
	RunPendingRefresh()
		{
		if .refreshTimer isnt false
			{
			.killRefresh()
			.refresh()
			return true
			}
		return false
		}

	SetReadOnly(readonly = true)
		{
		.Ctrl.SetReadOnly(readonly)
		}
	GETMINMAXINFO(lParam)
		{
		if not .Member?('Ctrl')
			return 1
		GetClientRect(.Hwnd, cr = Object())
		if cr.right is 0 and cr.bottom is 0
			return 1 // restoring from minimized
		StructModify(MINMAXINFO, lParam)
			{|mmi|
			wr = GetWindowRect(.Hwnd)
			x = .Ctrl.Xmin + 2 * .border + (wr.right - wr.left) - cr.right
			y = .Ctrl.Ymin + 2 * .border + (wr.bottom - wr.top) - cr.bottom
			mmi.minTrackSize = [:x, :y]
			if not .ctrlStretchable()
				mmi.maxTrackSize = [:x, :y]
			else if .Ctrl.Xstretch <= 0
				mmi.maxTrackSize = [:x, y: mmi.maxTrackSize.y]
			else if .Ctrl.Ystretch <= 0
				mmi.maxTrackSize = [x: mmi.maxTrackSize.x, :y]
			}
		return 0
		}
	Resize_Ctrl()
		{
		// kludge because Eta_Orders BeforeSave does SelectTab
		if .destroying
			return
		.Ctrl.Resize(.x, .y, .w, .h)
		}
	SIZE(lParam)
		{
		if (not .Member?('Ctrl'))
			return 0
		x = y = .border
		w = LOWORD(lParam) - 2 * .border
		h = HIWORD(lParam) - 2 * .border
		if (.Ctrl.Xstretch is false)
			w = .Ctrl.Xmin
		if (.Ctrl.Ystretch is false)
			h = .Ctrl.Ymin
		.Ctrl.Resize(.x = x, .y = y, .w = w, .h = h)
		return 0
		}
	EXITSIZEMOVE()
		{
		.RemStyle(WS.CLIPCHILDREN)
		return 1 // or 0, no idea
		}
	ENTERSIZEMOVE()
		{
		.AddStyle(WS.CLIPCHILDREN)
		return 1 // or 0, no idea
		}
	ObserveMove(fn)
		{
		.move_observers.Add(fn)
		}
	ObserveMoveRemove(fn)
		{
		.move_observers.Remove(fn)
		}
	WINDOWPOSCHANGING()
		{
		for observer in .move_observers
			observer()
		return 0
		}
	MOVING(lParam)
		{
		// snap to edges of work area
		StructModify(RECT, lParam)
			{|r|
			snap = ScaleWithDpiFactor(4)
			wa = GetWorkArea(r)
			if ((d = wa.left - r.left).Abs() < snap or
				(d = wa.right - r.right).Abs() < snap)
				{
				r.left += d
				r.right += d
				}
			if ((d = wa.top - r.top).Abs() < snap or
				(d = wa.bottom - r.bottom).Abs() < snap)
				{
				r.top += d
				r.bottom += d
				}
			}
		return 0
		}
	resizeToMin()
		{
		if not .Member?(#Ctrl)
			return

		xmin = .Ctrl.Xmin + 2 * .border
		ymin = .Ctrl.Ymin + 2 * .border
		GetClientRect(.Hwnd, cr = Object())
		cr.w = cr.right - cr.left
		cr.h = cr.bottom - cr.top
		if xmin > cr.w or ymin > cr.h or // control bigger than window
			not .ctrlStretchable()
			{
			if .ctrlStretchable()
				{ // don't shrink
				xmin = Max(xmin, cr.w)
				ymin = Max(ymin, cr.h)
				}
			wr = GetWindowRect(.Hwnd)
			w = xmin + (wr.right - wr.left) - cr.w
			h = ymin + (wr.bottom - wr.top) - cr.h
			.SetWinSize(w, h)
			}
		}
	ctrlStretchable()
		{
		return .Ctrl.Xstretch > 0 or .Ctrl.Ystretch > 0
		}

	validationItems: false
	AddValidationItem(item)
		{
		if (.validationItems is false)
			.validationItems = Object()
		.validationItems.Add(item)
		}
	RemoveValidationItem(item)
		{
		.validationItems.RemoveIf({ Same?(it, item) })
		}
	GetValidationItems()
		{
		return .validationItems
		}

	// default method, ListEditWindow redefines this so that its child windows
	// can refer to the main top level window (since the ListEditWindow is a
	// top level window itself)
	MainHwnd()
		{
		return .Window.Hwnd
		}

	tips?: false
	ToolTip(hwnd, tip, rect = false)
		{
		.tips.AddTool(hwnd, tip, :rect)
		}
	RemoveTip(hwnd)
		{
		.tips.RemoveTool(hwnd)
		}
	Tips()
		{
		return .tips
		}
	getter_tips()
		{
		.tips? = true
		return .tips = .Construct(ToolTipControl) // once only
		}
	NewTips()
		{
		if not .tips?
			return
		.tips.Destroy()
		.tips = .Construct(ToolTipControl)
		}

	keep_size: false
	window_info: false
	restoreKeepSize(parentHwnd, title, ctrlspec, keep_size)
		{
		if keep_size is false
			return false
		.keep_size = String?(keep_size) and keep_size isnt ""
			? keep_size
			: .keepSizeName(ctrlspec, title)
		if .keep_size is false or
			false is info = KeyListViewInfo.Get(.keep_size)
			return false
		.window_info = info.window_info
		if not .window_info.Member?(#w)
			return false

		// Only valid parentHwnd can be passed into getWorkArea()
		if IsWindow(parentHwnd)
			{
			min = .RequiredWindowSize()
			wa = .getWorkArea(parentHwnd)
			width = Min(Max(.window_info.w, min.w), wa.right - wa.left)
			height = Min(Max(.window_info.h, min.h), wa.bottom - wa.top)
			.SetWinSize(width, height)
			}
		}
	keepSizeName(ctrlspec, title)
		{
		result = ''
		if Object?(ctrlspec) and ctrlspec.Member?(0) and String?(ctrlspec[0])
			result = ctrlspec[0]
		if title isnt '' and title isnt false
			result = result isnt '' ? (result $ " - " $ title) : title
		return result isnt '' ? result : false
		}
	RequiredWindowSize()
		{
		xmin = .Ctrl.Xmin + 2 * .border
		ymin = .Ctrl.Ymin + 2 * .border
		GetClientRect(.Hwnd, cr = Object())
		cr.w = cr.right - cr.left
		cr.h = cr.bottom - cr.top
		wr = GetWindowRect(.Hwnd)
		wr.w = wr.right - wr.left
		wr.h = wr.bottom - wr.top
		dw = wr.w - cr.w
		dh = wr.h - cr.h
		return Object(w: xmin + dw, h: ymin + dh)
		}
	getWorkArea(parentHwnd)
		{
		return GetWorkArea(GetWindowRect(parentHwnd))
		}
	saveKeepSize()
		{
		if .keep_size is false
			return
		r = GetWindowRect(.Hwnd)
		window_info = .window_info is false ? Object() : .window_info
		window_info.w = r.right - r.left
		window_info.h = r.bottom - r.top
		KeyListViewInfo.Save(.keep_size, window_info)
		}

	On_Cancel()
		{
		.Result(false) // Result is defined in Dialog and Window
		}

	AllowCloseWindow?()
		{
		if false is .Send("Ok_to_CloseWindow?")
			return false

		// Check validation on the window
		if .Send('ConfirmDestroy') is false
			return false

		if .validationItems is false
			return true

		// Warning: validationItems could be modified during the following loop
		// if there are delayed calls during the window closing procedure.
		// This will cause an "object modified during iteration" error. As of 20181016
		// there are no cases of this, however it is possible it could be re-introduced
		for item in .validationItems
			if not item.ConfirmDestroy() and not .discardUnsavedChanges(item)
				return false
		return true
		}

	discardUnsavedChanges(item)
		{
		return (item.Method?('CloseWindowConfirmation') and
			not item.CloseWindowConfirmation())
			? false
			: CloseWindowConfirmation(.Window.Hwnd)
		}

	CLOSE()
		{
		if not .AllowCloseWindow?()
			return 0 // meaning we handled it so don't do default DestroyWindow
		return 'callsuper'
		}
	destroying: false
	DESTROY()
		{
		.destroying = true
		.killRefresh()
		if .DisableCloseTimer isnt false
			{
			.DisableCloseTimer.Kill()
			.DisableCloseTimer = false
			}
		.saveKeepSize()
		if .tips?
			.tips.Destroy()
		if .Member?('Ctrl')
			.Ctrl.Destroy()
		if --Suneido.OpenWindows is 0
			Image.DestroyAllInMemory()
		return 'callsuper'
		}

	DoActivate(@unused) { throw "Suniedo.js only" }
	UnregisterWindow(@unused) { throw "Suniedo.js only" }
	}
