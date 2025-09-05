// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class /*abstract*/
	{
	Xmin:		0
	Xstretch:	false
	Ymin:		0
	Ystretch:	false
	MaxHeight:	99999
	Top:		0		// Top and Left specify the "alignment" point
	Left:		0
	Name:		""
	Title:		""
	DefaultButton: "OK"
	Custom: 	false

	New()
		{
		.New2()
		}
	New2() /*internal*/
		{
		// this is separate from New so it can be overridden by WindowBase
		.Parent = _parent
		.WndProc = _parent.Base?(WndProc) ? _parent : _parent.WndProc
		.Controller = _parent.Base?(Controller) ? _parent : _parent.Controller
		.Window = _parent.Window

		if .Custom is false and _parent.Member?('Custom')
			.Custom = _parent.Custom

		.Xmin = .scale(_ctrlspec.GetDefault("xmin", .Xmin))
		.Ymin = .scale(_ctrlspec.GetDefault("ymin", .Ymin))
		.Top = .scale(_ctrlspec.GetDefault("top", .Top))
		.Left = .scale(.Left)
		if (_ctrlspec.Member?("xstretch"))
			.Xstretch = _ctrlspec.xstretch
		if (_ctrlspec.Member?("ystretch"))
			.Ystretch = _ctrlspec.ystretch
		if (_ctrlspec.Member?("name"))
			.Name = _ctrlspec.name
		}

	scale(value)
		{
		if not Number?(value) or value is 0
			return value
		return ScaleWithDpiFactor(value)
		}

	CallClass(@args)
		{
		args.Add(this, at: 0)
		return Window(args)
		}

	// send a message to controller
	Send(@args)
		{
		if not .Member?(#Controller)
			return 0 // destroyed
		if not args.Member?('source')
			args.source = this
		return .Controller.Msg(args)
		}
	Msg(args) // ignore messages from Window when top control is not a Controller
		{
		if args[0] is 'On_Cancel' and .Window.Method?(#On_Cancel)
			.Window.On_Cancel()
		return 0
		}
	On_Cancel()
		{
		if .Controller.Method?(#On_Cancel)
			.Controller.On_Cancel()
		else if .Window.Method?(#On_Cancel)
			.Window.On_Cancel()
		}

	Construct(@x)
		{
		x = .build(x)
		_parent = this
		_ctrlspec = x
		ctrl = Construct(x, "Control")
		if String?(ctrl.Name) and ctrl.Name isnt ""
			this[ctrl.Name] = ctrl
		return ctrl
		}

	build(x)
		{
		if x.Size() is 1 and x.Member?(0) and Object?(x[0])
			x = x[0]
		if x[0] is 'NoPrompt' and .isFieldString?(x[1])
			x = .handleNoPromptField(x)
		else if .isFieldString?(x[0])
			x = .handleField(x)
		return x
		}

	getDatadictControl(name, orig)
		{
		x = Datadict(name).Control.Copy().Add(name at: 'name')
		for m in orig.Members()
			if String?(m)
				x[m] = orig[m]
		return x
		}

	isFieldString?(item)
		{
		return String?(item) and item =~ "^[_a-z]"
		}

	handleNoPromptField(x)
		{
		name = x[1]
		x = .getDatadictControl(name, x)
		if .Custom isnt false and .Custom.Member?(name)
			x.Merge(.Custom[name])
		return x
		}

	handleField(x)
		{
		name = x[0]
		x = .getDatadictControl(name, x)
		hidden = false
		if .Custom isnt false and .Custom.Member?(name)
			{
			x.Merge(.Custom[name])
			hidden = .Custom[name].GetDefault('hidden', false)
			}
		if "" isnt prompt = Prompt(name)
			{
			if x[0] is 'CheckBox'
				{ // prompt is part of checkbox
				x.text = prompt
				x.hidden = hidden
				}
			else
				x = Object('Pair', Object('Static', prompt, :hidden), x)
			}
		return x
		}

	Resize(x/*unused*/, y/*unused*/, w/*unused*/, h/*unused*/)
		{ }
	Update()
		{ }
	HasFocus?()
		{
		return false
		}
	GetState()
		{
		return Object()
		}
	SetState(unused)
		{ }
	SetEnabled(unused)
		{ }
	SetVisible(unused)
		{ }
	SetReadOnly(unused)
		{ }
	SetFocus()
		{ }
	ClearFocus()
		{
		SetFocus(NULL)
		}
	GetReadOnly()
		{
		return false
		}

	GetChildren()
		{
		return Object()
		}
	FindControl(name, exclude = false)
		{
		for c in .GetChildren()
			{
			if exclude isnt false and c.Base?(exclude)
				continue
			if c.Name is name or
				false isnt (c = c.FindControl(name, :exclude))
				return c
			}
		return false
		}
	GetRect()
		{
		// returns a Rect containing the coordinates of
		// the control (relative to parent Hwnd)
		return Rect(0, 0, 0, 0)
		}
	GetClientRect()
		{
		// returns a Rect containing the coordinates of
		// the control's client area
		return Rect(0, 0, 0, 0)
		}
	Dirty?(unused = "")
		{
		return false
		}
	Valid?()
		{
		return true
		}
	ValidData?(@unused)
		{
		return true
		}
	HandleTab()
		{
		return false
		}
	AlertHwnd()
		{
		// sometimes the Alert methods are used on classes that don't have
		// a window associated with them.
		return .Member?("Window") ? .Window.Hwnd : 0
		}
	AlertError(title, message)
		{
		Alert(message, title, .AlertHwnd(), MB.ICONERROR)
		}
	AlertWarn(title, message)
		{
		Alert(message, title, .AlertHwnd(), MB.ICONWARNING)
		}
	AlertInfo(title, message)
		{
		Alert(message, title, .AlertHwnd(), MB.ICONINFORMATION)
		}
	AlertQuestion(title, message)
		{
		Alert(message, title, .AlertHwnd(), MB.ICONQUESTION)
		}
	BottomUp(@args)
		{
		for c in .GetChildren()
			c.BottomUp(@args)
		if .Method?(args[0])
			this[args[0]](@+1 args)
		return
		}
	TopDown(@args)
		{
		if .Method?(args[0])
			this[args[0]](@+1 args)
		for c in .GetChildren()
			c.TopDown(@args)
		return
		}
	ForeachChild(ctrl, block)
		{
		for c in ctrl.GetChildren()
			{
			block(c)
			.ForeachChild(c, block)
			}
		}
	WindowRefresh()
		{
		if .Member?(#Window)
			.Window.Refresh()
		}
	WindowHwnd()
		{
		try return .Window.Hwnd
		return 0
		}
	WindowActive?()
		{
		return .WindowHwnd() is GetActiveWindow()
		}
	FocusFirst(parentHwnd, custom = false)
		{
		if custom is false
			custom = .Custom

		if custom isnt false
			for field in custom.Members()
				if custom[field].GetDefault('first_focus', false) is true and
					false isnt ctrl = .FindControl(field)
					{
					ctrl.SetFocus()
					return
					}

		hwnd = first = GetNextDlgTabItem(parentHwnd, NULL, false)
		// look for first editable control
		while .readonly(ctrl = .getControl(hwnd))
			{
			hwnd = GetNextDlgTabItem(parentHwnd, hwnd, false)
			// stop if we don't find one
			if hwnd is first
				break
			}
		if ctrl isnt false and ctrl.Method?('SelectAll')
			ctrl.SelectAll()
		SetFocus(hwnd)
		}
	getControl(hwnd)
		{
		return .Window.HwndMap.GetDefault(hwnd, false)
		}
	readonly(ctrl)
		{
		return ctrl is false ? false : ctrl.GetReadOnly()
		}

	CalcMaxHeight()
		{
		return .Ystretch > 0 ? .MaxHeight : .Ymin
		}

	ZoomReadonly(value)
		{
		ZoomControl(0, value, readonly:)
		}

	CalcXminByControls(@unused)
		{
		}

	DoCalcXminByControls(plusCtrls, minusCtrls = #(), overlaps = 0)
		{
		plus = Instance?(plusCtrls) ? plusCtrls.Xmin : plusCtrls.SumWith({ it.Xmin })
		minus = Instance?(minusCtrls) ? minusCtrls.Xmin : minusCtrls.SumWith({ it.Xmin })
		return  plus - minus - overlaps
		}

	// uniqueID is used to prevent too many stacked delayed calls.
	// uniqueID should only be used in cases where it is "safe" to eliminate
	// intermediate delayed function calls (ie block is not doing critical processing
	// such as modifying data)
	timers: #()
	Delay(delay, block, uniqueID = false)
		{
		return .delayed(Delay, block, delay, uniqueID)
		}
	Defer(block, uniqueID = false)
		{
		return .delayed(Defer, block, 0, uniqueID)
		}
	delayed(fn, block, delay, uniqueID)
		{
		if .destroying?
			return false

		if .timers.Readonly?()
			{
			.timers = Object() // intentionally lazy, don't put in New
			.unique = Object()
			}

		if uniqueID isnt false and .unique.Member?(uniqueID)
			{
			.timers.Remove1(.unique[uniqueID])
			.unique[uniqueID].Kill()
			}

		timer = false // needed so block can reference it
		timer = fn(delayMs: delay) { .timers.Remove1(timer); block() }
		.timers.Add(timer)
		if uniqueID isnt false
			.unique[uniqueID] = timer
		return (.killwrap)(timer, .timers)
		}
	killwrap: class
		{
		New(.killer, .timers)
			{
			}
		Kill()
			{
			.killer.Kill()
			.timers.Remove1(.killer)
			}
		}

	Destroyed?()
		{
		return not .Member?('Window')
		}

	destroying?: false
	Destroy()
		{
		.destroying? = true

		// avoid object-modified-during-iteration, also avoid Copy overhead
		while not Same?(.timers, k = .timers.PopLast())
			k.Kill()

		if .Member?(#Parent) and Instance?(.Parent) and String?(.Name) and
			Same?(this, .Parent.GetDefault(.Name, false))
			.Parent.Delete(.Name)

		.Delete(all:) // to help garbage collection
		}

	Act(@unused) { throw "Suniedo.js only" }
	ActWith(@unused) { throw "Suniedo.js only" }
	CancelAct(@unused) { throw "Suniedo.js only" }
	SuSetHidden(@unused) { throw "Suniedo.js only" }
	}