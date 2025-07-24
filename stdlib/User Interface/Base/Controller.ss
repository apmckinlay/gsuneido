// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Container
	{
	// can't have Name: "Controller" - it overwrites "Controller" in parent !!!
	New(ctrl = false)
		{
		.xmin0 = .Xmin
		.ymin0 = .Ymin

		.Initialize(ctrl)
		}
	Initialize(ctrl)
		{
		.init_redir()
		if ctrl is false
			ctrl = .Val_or_func('Controls')
		if ctrl is false
			return // derived class constructs & sets size/stretch
		if .ctrl isnt false
			.ctrl.Destroy()
		.ctrl = .Construct(ctrl)
		.Recalc()
		// TODO move to Recalc (have to save original values)
		if .Xstretch is false
			.Xstretch = .ctrl.Xstretch
		if .Ystretch is false
			.Ystretch = .ctrl.Ystretch
		if .MaxHeight is Control.MaxHeight
			.MaxHeight = .ctrl.MaxHeight
		}
	Controls: false
	init_redir()
		{
		.redir = Object(
			On_Copy:		'focus',
			On_Cut:			'focus',
			On_Paste:		'focus',
			On_Delete:		'focus',
			On_Select_All: 	'focus',
			On_Undo:		'focus',
			On_Redo:		'focus')
		}
	redir: () // avoid errors if destroyed

	Msg(args) /* internal, used by Control.Send */
		{
		msg = args[0]
		target = this
		if .redir.Member?(msg)
			{
			target = .redir[msg]
			if target is 'focus' and false is target = .GetFocus()
				return 0
			}
		return .call_target(target, msg, args)
		}
	call_target(target, msg, args)
		{
		if Function?(target)
			return target(@+1 args)
		if Instance?(target) or Class?(target)
			{
			if target.Method?(msg)
				return target[msg](@+1 args)
			else if target.Method?(#Recv)
				return target.Recv(@args)
			}
		return 0
		}

	GetFocus()
		{
		return (NULL is (hwnd = GetFocus()) or
			false is .Window.HwndMap.Member?(hwnd))
			? false : .Window.HwndMap[hwnd]
		}

	Redir(msg, ctrl = 'focus')
		{
		if ctrl isnt false
			.redir[msg] = ctrl
		}
	GetRedir(msg)
		{
		return .redir.GetDefault(msg, false)
		}
	DeleteRedir(msg)
		{
		return .redir.Delete(msg)
		}
	RemoveRedir(target)
		{
		.redir.Remove(target)
		}

	ctrl: false
	Resize(x, y, w, h)
		{
		if .ctrl is false
			return
		.ctrl.Resize(x, y, w, h)
		.Recalc()
		}
	Recalc()
		{
		if .ctrl is false
			return
		.Xmin = Max(.xmin0, .ctrl.Xmin)
		.Ymin = Max(.ymin0, .ctrl.Ymin)
		.MaxHeight = .ctrl.MaxHeight
		}

	GetChild()
		{
		return .ctrl
		}
	GetChildren()
		{
		return .ctrl is false ? #() : Object(.ctrl)
		}
	ChangeControl(control)
		{
		.redir = Object()
		if .ctrl isnt false
			.ctrl.Destroy()
		.ctrl = .Construct(control)
		.Window.Refresh()
		}
	DetachControl()
		{
		ctrl = .ctrl
		.ctrl = false
		return ctrl
		}
	AttachControl(control)
		{
		.redir = Object()
		if .ctrl isnt false
			.ctrl.Destroy()
		.ctrl = control
		.ctrl.Parent = this
		.ctrl.Window = .Window
		.ctrl.WndProc = .WndProc
		.setController(.ctrl)
		.Window.Refresh()
		}
	setController(control)
		{
		control.Controller = this
		if not control.Base?(Controller)
			for c in control.GetChildren()
				.setController(c)
		}

	On_Close()
		{
		.Window.Destroy()
		}
	On_Exit()
		{
		.Window.Destroy()
		}
	Update()
		{
		if .ctrl isnt false
			.ctrl.Update()
		super.Update()
		}
	AccessGoto(field, value, wrapper = false)
		{
		if .ctrl isnt false and .ctrl.Method?('AccessGoto')
			.ctrl.AccessGoto(field, value, :wrapper)
		}
	AccessGoTo_CurrentBookOption()
		{
		return .Send('AccessGoTo_CurrentBookOption')
		}
	BookRefresh()
		{
		.Send('BookRefresh')
		}
	}
