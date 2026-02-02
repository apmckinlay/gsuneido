// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Control
	{
	New()
		{
		.HwndMap = Object()
		.cmdmap = Object().Set_default("")
		.rcmdmap = Object().Set_default(false)
		.mapcmd("On_OK", ID.OK)
		.mapcmd("On_Cancel", ID.CANCEL)
		.Window = .Controller = this
		.move_observers = Object()
		}

	New2() { }

	GetLayout()
		{
		return .Ctrl.GetLayout()
		}

	GetWindowComponentName(excludeModalWindow = false)
		{
		componentName = .ComponentName
		if componentName is #Window and
			SuRenderBackend().WindowManager.ShowingModalWindow?(excludeModalWindow)
			componentName = #ModalWindow

		if componentName in (#ModalWindow, #Dialog)
			SuRenderBackend().WindowManager.AddModalWindow(this)
		if componentName is 'Window'
			SuRenderBackend().WindowManager.AddTaskbarWindow(this)
		return componentName $ 'Component'
		}

	On_Cancel()
		{
		.Result(false) // Result is defined in Dialog and Window
		}

	DialogCenterSize(title, control, keep_size, useDefaultSize = false)
		{
		if .restoreKeepSize(title, control, keep_size) is false
			if useDefaultSize is true
				.SetWinSize(@.DefaultSize())

		.Act('Center')
		}

	DefaultSize()
		{
		return Object(w: 1024, h: 768)
		}

	keep_size: false
	window_info: false
	orig_window_info: false
	restoreKeepSize(title, ctrlspec, keep_size)
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
		.orig_window_info = info.window_info.Copy()
		if not .window_info.Member?(#w_web)
			return false
		.SetWinSize(.window_info.w_web, .window_info.h_web)
		return true
		}

	SetWinSize(w, h)
		{
		.Act(#SetWinSize, w, h)
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

	WINDOWRESIZE(rect)
		{
		.windowSize = Object(w: rect.width, h: rect.height)
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

	getter_validationItems()
		{
		return .validationItems = Object()
		}
	AddValidationItem(item)
		{
		.validationItems.Add(item)
		}
	RemoveValidationItem(item)
		{
		.validationItems.Remove(item)
		}
	GetValidationItems()
		{
		return .validationItems
		}

	MainHwnd()
		{
		return .Window.Hwnd
		}

	Refresh()
		{
		.Act(#Refresh)
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
	AddCommands(cmds)
		{
		for cmd in cmds
			.commands[cmd[0]] = Object(
				accel: cmd.GetDefault(1, ""),
				help: cmd.GetDefault(2, ""),
				bitmap: cmd.GetDefault(3, cmd[0]),
				id: .Mapcmd(cmd[0]))
		.SetupAccels(cmds)
		}
	Commands()
		{ return .commands }

	// accelerator stuff =============================================
	getter_accels()
		{
		return .accels = Object()
		}
	getter_accelGroups()
		{
		return .accelGroups = Object()
		}
	SetupAccels(cmds)
		{
		.ResetAccels()
		newGroup = Object()
		for cmd in cmds
			if cmd.Member?(1)
				{
				id = .Mapcmd(cmd[0])
				accel = .make_accel(cmd[1], id)
				newGroup.Add(Object(:accel, :id))
				}
		.accelGroups.Add(newGroup)
		.SetAccels()
		return newGroup
		}
	AddAccel(name, s)
		{
		id = .Mapcmd(name)
		.commands[name] = Object(help: '', bitmap: '', accel: s, :id)
		accel = .make_accel(s, id)
		.accelGroups.Add(Object(Object(:accel, :id)))
		}
	convert: (
		space: ' ',
		add: '=', // js keydown's key is '=' instead of '+'
		subtract: '-',
		down: 'arrowdown',
		left: 'arrowleft',
		right: 'arrowright',
		up: 'arrowup')
	make_accel(s, cmd)
		{
		ac = Object(:cmd, idx: 0)
		if s =~ "Ctrl[+]"
			{
			ac.idx |= 0x100
			s = s.Replace("Ctrl[+]", "")
			}
		if s =~ "Alt[+]"
			{
			ac.idx |= 0x10
			s = s.Replace("Alt[+]", "")
			}
		if s =~ "Shift[+]"
			{
			ac.idx |= 0x1
			s = s.Replace("Shift[+]", "")
			}
		ac.key = s.Lower()
		ac.key = .convert.GetDefault(ac.key, ac.key)
		return ac
		}
	RestoreAccels(accels)
		{
		.ResetAccels()
		.accelGroups.Remove(accels)
		.SetAccels()
		}
	ResetAccels()
		{
		.Act(#ResetAccels)
		}
	SetAccels()
		{
		.buildAccel()
		if .accels.NotEmpty?()
			.Act(#SetAccels, .accels)
		}
	buildAccel()
		{
		.accels = Object()
		for group in .accelGroups
			for item in group
				.accels.Add(item.accel)
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
	cmd_to_method(cmd)
		{
		cmd = cmd.BeforeFirst("\t")
		return "On_" $ ToIdentifier(cmd)
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

	COMMAND(id)
		{
		if .cmdmap.Member?(id)
			{
			.Send(.cmdmap[id])
			return
			}
		SuServerPrint("COMMAND id not found: " $ id)
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

	SetTitle(text)
		{
		.Act(#SetTitle, text)
		}

	// used by SetWindowText to set value directly
	SetText(text)
		{
		.SetTitle(text)
		}

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
//		if parentHwnd is 0 or not IsWindow(parentHwnd) or parentHwnd is GetDesktopWindow()
//			{
//			parentHwnd = GetActiveWindow()
//			if ((GetWindowLong(parentHwnd, GWL.STYLE) & WS.POPUP) isnt 0)
//				{
//				// use GetAncestor (instead of GetParent) to exclude OWNER windows
//				// since OWNER WINDOW is an overlapped or pop-up window,
//				// and they are destroyed when another window takes focus
//				parentHwnd = GetAncestor(parentHwnd, GA.PARENT)
//				}
//			}
//		else
//			parentHwnd = GetAncestor(parentHwnd, GA.ROOT)
//		if not IsWindow(parentHwnd)
//			{
//			Print("ERROR: GetParent: failed to get parentHwnd")
//			return 0
//			}
		return parentHwnd
		}

	windowSize: false
	saveKeepSize()
		{
		if .keep_size is false or .windowSize is false
			return
		window_info = .window_info is false ? Object() : .window_info
		for m, v in .windowSize
			window_info[m $ '_web'] = v
		if .orig_window_info isnt window_info
			KeyListViewInfo.Save(.keep_size, window_info)
		}

	GetChildren()
		{
		return [.Ctrl]
		}

	DoActivate()
		{
		SuRenderBackend().WindowManager.ActivateWindow(this)
		}

	UnregisterWindow()
		{
		SuRenderBackend().WindowManager.UnregisterWindow(this)
		}

	IsMinimized?()
		{
		return false
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
			return
		.DESTROY()
		}

	DESTROY()
		{
		if .Destroyed?()
			return
		Finally({
			.saveKeepSize()
			if .Member?('Ctrl')
				.Ctrl.Destroy()
			}, {
			.Act(#Destroy)
			.UnregisterWindow()
			})
		super.Destroy()
		}

	Destroy()
		{
		if .Destroyed?()
			return
		.DESTROY()
		}
	}
