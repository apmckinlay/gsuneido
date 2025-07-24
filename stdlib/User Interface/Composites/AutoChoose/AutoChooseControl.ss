// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
/*
AutoChooseControl is the main control
it uses AutoChooseList to display the popup
and it uses AutoListBoxControl for the list within the popup
this control keeps focus - it does keyboard handling for popup list
allowOther: can be a function e.g. ValidEmailAddress?
list argument can be:
	object = the list itself
	capitalized string = name of global function
	function() = returns complete list
	function(prefix) = returns list of matches
	uncapitalized string = field/rule name
TODO: for rules, set prefix in .<field>_prefix
used as base for MultiAutoChooseControl

TESTING
- TAB and RETURN should select the current item in the list
- ESCAPE should close the list
- in a Dialog TAB, RETURN, and ESCAPE should apply to the list if it's open
	otherwise they should behave as normal in a dialog
- UP and DOWN should open the list

have to test all the combinations of:
- TAB, RETURN, ESCAPE
- list open, list closed
- AutoChoose, MultiAutoChoose
- in a window or in a dialog with OK, or with a different default button

Window(#(Horz (AutoChoose BuiltinNames) Field))

Window(#(Horz (MultiAutoChoose BuiltinNames width: 20) Field))

Ask(ctrl: #(Horz (AutoChoose BuiltinNames) Field))

Ask(ctrl: #(Horz (MultiAutoChoose BuiltinNames width: 20) Field))

Dialog(0, Controller {
	Controls: (Vert,
		(AutoChoose BuiltinNames),
		(MultiAutoChoose BuiltinNames width: 20),
		(Button Go))
	DefaultButton: Go
	On_Go() { .Window.Result(123) }
	})
*/
EditControl
	{
	Name: AutoChoose
	New(list = false, width = 20, status = "", readonly = false,
		mandatory = false, allowOther = false, style = 0, height = 1, cue = false,
		font = "", size = "", .autoSelect = false)
		{
		super(:status, :readonly, :mandatory,
			style: style | (height > 1
				? ES.MULTILINE | ES.AUTOVSCROLL | WS.VSCROLL
				: ES.AUTOHSCROLL),
			:width, :height, :cue, :font, :size)
		.SubClass()
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
		}
	MOUSEWHEEL(wParam, lParam)
		{
		return .listctrl is false
			? 0
			: .listctrl.ListBox.SendMessage(WM.MOUSEWHEEL, wParam, lParam)
		}
	GETDLGCODE(wParam)
		{
		if .ListOpen?()
			{
			if wParam is VK.TAB or wParam is VK.RETURN or wParam is VK.ESCAPE
				return DLGC.WANTALLKEYS
			}
		else if wParam is VK.TAB or wParam is VK.RETURN or wParam is VK.ESCAPE
			return 0 // prevent MULTILINE from absorbing if list not open
		return 'callsuper'
		}
	KEYDOWN(wParam)
		{
		switch (wParam)
			{
		case VK.DOWN :
			return .down()
		case VK.UP :
			return .up()
		case VK.PRIOR, VK.NEXT :
			if .ListOpen?()
				.listctrl.ListBox.SendMessage(WM.KEYDOWN, wParam)
			return 0
		case VK.TAB, VK.RETURN :
			return 0 // handled by CHAR
		case VK.ESCAPE :
			return 0 // prevent MULTILINE from sending CLOSE
		default:
			}
		return 'callsuper'
		}
	CHAR(wParam)
		{ // can't put these in KEYDOWN because then CHAR causes beep
		if wParam is VK.ESCAPE
			return .Escape()
		if wParam is VK.RETURN or wParam is VK.TAB
			return .enter()
		return 'callsuper'
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
	EN_CHANGE()
		{
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
	EN_SETFOCUS()
		{
		if .autoSelect
			Defer(.SelectAll)
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
			.listctrl = AutoChooseList(this, choices)
		return .listctrl
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
		}
	close_list()
		{
		.killtimer()
		if .listctrl isnt false
			{
			DestroyWindow(.listctrl.Hwnd)
			.listctrl = false
			}
		}
	Destroy()
		{
		.Window.ObserveMoveRemove(.close_list)
		.close_list()
		super.Destroy()
		}
	}