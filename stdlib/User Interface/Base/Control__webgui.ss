// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
class
	{
	Xmin:		0
	Xstretch:	false
	Ymin:		0
	Ystretch:	false
	Name:		""
	Title:		""
	DefaultButton: "OK"
	Custom: 	false
	Top:		0		// Top and Left specify the "alignment" point
	Left:		0

	New()
		{
		.New2()
		}
	New2() /*internal*/
		{
		// this is separate from New so it can be overridden by WindowBase
		.Parent = _parent
		.Controller = _parent.Base?(Controller) ? _parent : _parent.Controller
		.Window = _parent.Window

		if .Custom is false and _parent.Member?('Custom')
			.Custom = _parent.Custom
		if (_ctrlspec.Member?("name"))
			.Name = _ctrlspec.name

		.buildExtraSpec(_ctrlspec)
		}

	extraSpec: #()
	buildExtraSpec(ctrlspec)
		{
		.extraSpec = Object()
		for m in #((xmin, 0), (ymin, 0), (xstretch, false), (ystretch, false))
			{
			member = m[0].Capitalize()
			if this[member] isnt m[1]
				.extraSpec[m[0]] = this[member]
			if ctrlspec.Member?(m[0])
				.extraSpec[m[0]] = ctrlspec[m[0]]
			}
		.extraSpec.uniqueId = .UniqueId
		.extraSpec.name = .Name
		}
	AddExtraSpec(member, value)
		{
		.extraSpec[member] = value
		}

	Getter_UniqueId()
		{
		.UniqueId = SuRenderBackend().NextId()
		.AddHwnd(.UniqueId)
		return .UniqueId
		}

	Getter_Hwnd()
		{
		return .UniqueId
		}

	AddHwnd(hwnd)
		{
		SuRenderBackend().Register(hwnd, this)
		.Window.HwndMap[hwnd] = this
		}

	DelHwnd(hwnd)
		{
		SuRenderBackend().UnRegister(hwnd)
		.Window.HwndMap.Delete(hwnd)
		}

	EditHwnd() // default
		{
		return .Hwnd
		}

	GetLayout()
		{
		layout = Object(.ComponentName)
		layout.Append(.ComponentArgs)
		layout.Append(.extraSpec)
		return layout
		}

	CalcXminByControls(plusCtrls, minusCtrls = #())
		{
		.Act(#CalcXminByControls,
			.convertCtrlToUniqueId(plusCtrls),
			.convertCtrlToUniqueId(minusCtrls))
		}

	convertCtrlToUniqueId(ctrls)
		{
		if Instance?(ctrls)
			ctrls = Object(ctrls)
		converted = Object()
		for ctrl in ctrls
			converted.Add(ctrl.UniqueId)
		return converted
		}

	Act(@args)
		{
		.act(args)
		}

	ActWith(block)
		{
		reservation = SuRenderBackend().ReserveAction()
		args = false
		Finally({
			args = block(:reservation)
			}, {
			if args is false
				SuRenderBackend().CancelReserve(reservation.at)
			})
		.act(args, reservation.at)
		}

	act(args, at = false)
		{
		action = args[0]
		args.Delete(0)
		uniqueId = .UniqueId
		if args.Member?(#actUniqueId)
			{
			uniqueId = args.actUniqueId
			args.Delete(#actUniqueId)
			}
		SuRenderBackend().RecordAction(uniqueId, action, args, :at)
		}

	CancelAct(action, block = false)
		{
		SuRenderBackend().CancelAction(.UniqueId, action, block)
		}

	SendMessage(msg, wParam, lParam)
		{
		SendMessage(.Hwnd, msg, wParam, lParam)
		}

	SubClass() { }

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

	GetState()
		{
		return Object()
		}
	SetState(unused)
		{ }
	enabled: true
	GetEnabled()
		{
		return .enabled
		}
	SetEnabled(enabled)
		{
		if .GetEnabled() isnt enabled
			{
			.enabled = enabled
			.Act('SetEnabled', enabled)
			}
		}
	SetVisible(visible)
		{
		Assert(Boolean?(visible))
		.Act('SetVisible', visible)
		}
	readonly: false
	GetReadOnly()
		{
		return .readonly
		}
	SetReadOnly(readonly)
		{
		if .GetReadOnly() isnt readonly
			{
			.readonly = readonly
			.Act('SetReadOnly', readonly)
			}
		}
	SetFocus()
		{
		SetFocus(.UniqueId)
		}
	ClearFocus()
		{
		SetFocus(NULL)
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
		.Act('FocusFirst', actUniqueId: parentHwnd)
		}

	ZoomReadonly(value)
		{
		ZoomControl(0, value, readonly:)
		}

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
	HasFocus?()
		{
		return .Hwnd isnt 0 and GetFocus() is .Hwnd
		}

	ToolTip(tip)
		{
		.Act('AddToolTip', tip)
		}

	hidden: false
	GetHidden()
		{
		return .hidden
		}
	SuSetHidden(hidden)
		{
		.hidden = hidden
		.Act(#SetHidden, hidden)
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

	SetFont(font = "", size = "", weight = "", text = "",
		underline = false, italic = false, strikeout = false)
		{
		.Act(#SetFont, :font, :size, :weight, :text, :underline, :italic, :strikeout)
		.WindowRefresh()
		}

	ContextMenuCall(menu)
		{
		this['On_' $ ToIdentifier(menu.BeforeFirst('\t'))]()
		}
	DevMenu: ('', 'Copy Field Name')
	On_Copy_Field_Name()
		{
		fieldName = .getFieldName()
		ClipboardWriteString(fieldName)
		}
	getFieldName()
		{
		name = .Send('GetFieldName')
		return String?(name)
			? name
			: .Name is 'Value' // for choose fields
				? .Parent.Name
				: .Name
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
			.Parent.GetDefault(.Name, false) is this
			.Parent.Delete(.Name)
		if .Member?(#UniqueId)
			.DelHwnd(.UniqueId)
		.Delete(all:) // to help garbage collection
		}

	Resize(x/*unused*/, y/*unused*/, w/*unused*/, h/*unused*/) { }
	Repaint(@unused) { }
	}
