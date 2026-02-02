// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
EditControl
	{
	Name: 			'AutoChoose'
	ComponentName: 	'AutoChoose'
	New(list = false, width = 20, status = "", readonly = false,
		mandatory = false, allowOther = false, style = 0, height = 1, cue = false,
		font = "", size = "")
		{
		super(:status, :readonly, :mandatory, :style, :width, :height, :cue, :font, :size)
		.Window.ObserveMove(.close_list)
		if String?(list) and list[0].Upper?()
			list = Global(list)
		.list = list
		.mandatory = mandatory is true
		if String?(allowOther) and allowOther[0].Upper?()
			allowOther = Global(allowOther)
		else if not Function?(allowOther)
			allowOther = allowOther is true
		.allowOther = allowOther

		.ComponentArgs = Object(width, readonly, height, font, size)
		}

	KEYDOWN(wParam)
		{
		switch (wParam)
			{
		case VK.DOWN :
			return .down()
		case VK.UP :
			return .up()
		case VK.TAB, VK.RETURN :
			return .enter()
		case VK.ESCAPE :
			return .Escape()
		default:
			}
		}

	down()
		{
		if .GetReadOnly()
			return 0
		if .listctrl isnt false
			.listctrl.Down()
		else if .open_list() isnt false
			.listctrl.SelectFirst()
		return 0
		}

	up()
		{
		if .GetReadOnly()
			return 0
		if .listctrl isnt false
			.listctrl.Up()
		else if .open_list() isnt false
			.listctrl.SelectLast()
		return 0
		}

	enter()
		{
		if .listctrl isnt false
			{
			if "" isnt s = .listctrl.Get()
				.Picked(s)
			.close_list()
			}
		return 0
		}

	InsertChoice(s)
		{
		SetWindowText(.Hwnd, s)
		.Send('NewValue', .Get())
		}

	Escape()
		{
		.close_list()
		return 0
		}

	Delay: 200
	ignore: false
	EN_CHANGE(text)
		{
		super.EN_CHANGE(text)
		if .ignore
			.ignore = false
		else if .listctrl is false
			{
			.killtimer()
			.timer = Delay(.Delay, .open_list)
			}
		else
			.change()
		return 0
		}

	EN_KILLFOCUS()
		{
		dirty? = .Dirty?()
		// has to call super here
		super.EN_KILLFOCUS()
		if dirty?
			.Send("NewValue", .Get())
		return 0
		}

	KillFocus()
		{
		if .GetReadOnly()
			return
		.killtimer()
		.close_list()
		.PickIfOneChoice()
		}

	PickIfOneChoice()
		{
		choices = .choices()
		if choices.Size() is 1
			.Picked(choices[0])
		}

	timer: false
	killtimer()
		{
		if .timer is false
			return
		.timer.Kill()
		.timer = false
		}

	change()
		{
		.close_list()
		.open_list()
		}

	listctrl: false
	open_list(choices = false)
		{
		if .listctrl isnt false
			return .listctrl
		if choices is false
			choices = .choices()
		if not choices.Empty?()
			{
			.listctrl = AutoChooseList(this, choices)
			.syncListStatus(true)
			}
		return .listctrl
		}

	syncListStatus(open)
		{
		.Act('SyncListStatus', open)
		}

	OpenList(choices = false)
		{
		.open_list(choices)
		}

	ListOpen?()
		{
		return .listctrl isnt false
		}

	choices()
		{
		return .Choices(.GetPrefix(), .list)
		}

	Choices(prefix, list)
		{
		if prefix is ""
			return #()
		if String?(list)
			list = .Send(#GetField, list)
		if Function?(list) and list.Params() is "()"
			list = list()
		if Object?(list)
			return list.Filter({ it =~ '\<(?i)(?q)' $ prefix })
		else if Function?(list)
			return list(prefix)
		else
			return #()
		}

	GetPrefix()
		{
		return .Get()
		}

	Valid?()
		{
		value = .Get()
		if .mandatory and value is ""
			return false
		if not .ValidChoice?(value, .allowOther, .list)
			return false
		return super.Valid?()
		}

	ValidChoice?(value, allowOther, list)
		{
		if Function?(allowOther)
			allowOther = allowOther(value)
		return allowOther or .validChoice?(value, list)
		}

	validChoice?(value, list)
		{
		return .Choices(value, list).Has?(value)
		}

	IgnoreAfterPick?: true
	Picked(s) // Sent from AutoChooseList
		{
		if .IgnoreAfterPick?
			.ignore = true
		.InsertChoice(s)
		.close_list()
		}
	GetListPos() // called by AutoChooseList
		{
		return GetWindowRect(.Hwnd)
		}

	ListClosed() // sent by listbox
		{
		.listctrl = false
		.syncListStatus(false)
		}

	close_list()
		{
		.killtimer()
		if .listctrl isnt false
			{
			.listctrl.DESTROY()
			.listctrl = false
			.syncListStatus(false)
			}
		}

	Destroy()
		{
		.Window.ObserveMoveRemove(.close_list)
		.close_list()
		super.Destroy()
		}
	}
