// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Container
	{
	ComponentName: 'HtmlDiv'
	ctrl: false
	New(ctrl = false)
		{
		.Initialize(ctrl)
		}
	Initialize(ctrl)
		{
		.init_redir()
		.ComponentArgs = Object()
		if ctrl is false
			ctrl = .Val_or_func('Controls')
		if ctrl is false
			return // derived class constructs & sets size/stretch
		if .ctrl isnt false
			.ctrl.Destroy()
		.ctrl = .Construct(ctrl)

		.ComponentArgs.Add(.ctrl.GetLayout())
		}
	Controls: false

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

	GetChild()
		{
		return .ctrl
		}
	GetChildren()
		{
		return .ctrl is false ? #() : Object(.ctrl)
		}

	GetFocus()
		{
		if false is uniqueId = GetFocus()
			return false
		return SuRenderBackend().GetRegisteredControl(uniqueId)
		}

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
	On_Close()
		{
		.Window.Destroy()
		}
	On_Exit()
		{
		.Window.Destroy()
		}
	}
