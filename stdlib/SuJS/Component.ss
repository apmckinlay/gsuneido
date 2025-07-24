// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
class
	{
	Xmin:		0
	Xstretch:	false
	Ymin:		0
	Ystretch:	false
	Left:		0
	Name:		""
	Custom: 	false
	UniqueId:	false
	MaxHeight:	99999
	ContextMenu: false

	New()
		{
		.New2()
		}
	New2()
		{
		.Parent = _parent
		.ParentEl = _parent.Member?(#TargetEl) ? _parent.TargetEl : _parent.El
		.Controller = _parent.Base?(HtmlDivComponent)
			? _parent
			: _parent.Controller
		.Window = _parent.Window

		.Xmin = _ctrlspec.GetDefault("xmin", .Xmin)
		.Ymin = _ctrlspec.GetDefault("ymin", .Ymin)
		.Left = .Left
		if (_ctrlspec.Member?("xstretch"))
			.Xstretch = _ctrlspec.xstretch
		if (_ctrlspec.Member?("ystretch"))
			.Ystretch = _ctrlspec.ystretch
		if (_ctrlspec.Member?("name"))
			.Name = _ctrlspec.name
		.InitUniqueId()
		}

	InitUniqueId()
		{
		if (_ctrlspec.Member?("uniqueId"))
			.UniqueId = _ctrlspec.uniqueId
		if .UniqueId isnt false
			SuRender().Register(.UniqueId, this)
		}

	Event(@args)
		{
		.sendEvent(args)
		}

	EventWithOverlay(@args)
		{
		.sendEvent(args, showOverlay?:)
		}

	EventWithFreeze(@args)
		{
		.sendEvent(args, freeze?:)
		}

	sendEvent(args, showOverlay? = false, freeze? = false)
		{
		if .UniqueId is false
			{
			SuRender().Event(false, 'SuneidoLog', Object(
				'WARNING: Component sent an event when .UniqueId is false',
				params: [destroyed?: .Destroyed?(), this: Display(this)].Merge(args)))
			return
			}
		event = args[0]
		args.Delete(0)
		SuRender().Event(.UniqueId, event, args, :showOverlay?, :freeze?)
		}

	UnsupportedFeature(message)
		{
		.Event('AlertInfo', 'Unsupported Browser Feature', message)
		}

	delayedTaskId: false
	RunWhenNotFrozen(block)
		{
		.delayedTaskId = SuRender().RunWhenNotFrozen(block)
		}

	El: false
	CreateElement(tag, text = false, className = false, namespace = false)
		{
		.SetEl(CreateElement(tag, .ParentEl, className, :namespace))
		if text isnt false
			.El.innerHTML = text
		}

	SetEl(el)
		{
		.El = el
		.El.Control(this)
		.El.Window(.Window)
		.El.SetAttribute('su-control', Display(this))
		.El.SetAttribute('su-name', .Name)
		.El.SetAttribute('su-unique-id', Display(.UniqueId))
		.El.AddEventListener('focus', .OnFocus)
		.El.AddEventListener('blur', .OnBlur)
		if .ContextMenu is true
			.El.AddEventListener('contextmenu', .OnContextMenu)
		}

	OnContextMenu(event)
		{
		.RunWhenNotFrozen({
			.EventWithOverlay('ContextMenu', event.clientX, event.clientY) })
		event.StopPropagation()
		event.PreventDefault()
		}

	OnFocus()
		{
		SuRender().UpdateStatus(#Focus, .UniqueId)
		}

	OnBlur(event = false)
		{
		id = false
		if event isnt false
			try
				{
				focusEl = event.relatedTarget
				id = focusEl.Control().UniqueId
				}
		SuRender().UpdateStatus(#Focus, id)
		}

	SetMinSize()
		{
		if .El is false
			return
		if Number?(.Xmin)
			.El.SetStyle('min-width', .Xmin $ 'px')
		if Number?(.Ymin)
			.El.SetStyle('min-height', .Ymin $ 'px')
		if Number?(.MaxHeight)
			.El.SetStyle('max-height', .MaxHeight $ 'px')
		}
	SetStyles(styles, el = false)
		{
		if el is false
			el = .El
		if el is false
			return
		for style, value in styles
			el.SetStyle(style, value)
		}
	SetFont(font = "", size = "", weight = "",
		underline = false, italic = false, strikeout = false, el = false)
		{
		if el is false
			el = .El
		if font isnt ""
			el.SetStyle("font-family", .FontFamily(font))
		if size isnt ""
			{
			el.SetStyle("font-size", .ConvertSize(size))
			}
		if weight isnt ""
			el.SetStyle("font-weight", weight)
		if italic is true
			el.SetStyle("font-style", "italic")
		decoration = ""
		if underline is true
			decoration $= "underline "
		if strikeout is true
			decoration $= "line-through"
		if decoration isnt ""
			el.SetStyle("text-decoration", decoration)
		}
	// Based on StdFonts
	fontMap: #(Ui: "Arial", Mono: "monospace", Serif: "Georgia", Sans: "Verdana")
	FontFamily(font)
		{
		if font.Prefix?('@')
			font = .fontMap.GetDefault(font[1..].Capitalize(), font[1..])
		return '' $ font $ '' $ Opt(', ', SuFontOb.GetDefault(font, ''))
		}
	ConvertSize(size)
		{
		return String?(size)
			? .calcSize(size) / 10 /*=10%*/ + 1 $ "em"
			: -StdFontsSize.LfSize(size, WinDefaultDpi) $ "px"
		}
	calcSize(sizeStr)
		{
		sign = 1
		if sizeStr[0] in ('-', '+')
			{
			sign = sizeStr[0] is '-' ? -1 : 1
			sizeStr = sizeStr[1..]
			}
		size = 0
		for i in ..sizeStr.Size()
			size = size * 10/*=10*/ + sizeStr[i].Asc() - '0'.Asc()
		return sign * size
		}

	CalcXminByControls(@unused) { }
	DoCalcXminByControls(plusCtrls, minusCtrls)
		{
		return .calcXminByControls(plusCtrls) - .calcXminByControls(minusCtrls)
		}

	calcXminByControls(ctrls)
		{
		res = 0
		for id in ctrls
			{
			if false is component = SuRender().GetRegisteredComponent(id)
				continue
			res += component.Xmin
			}
		return res
		}

	hidden: false
	GetHidden()
		{
		return .hidden
		}
	SetHidden(hidden)
		{
		.hidden = hidden
		.SetVisible(.GetVisible())
		}
	prevDisplay: 'initial'
	SetVisible(visible)
		{
		visible = not .GetHidden() and visible
		el = .Member?(#GetContainerEl) ? .GetContainerEl() : .El
		if el.GetStyle("display") isnt 'none'
			if visible
				return
			else
				.prevDisplay = el.GetStyle("display")
		el.SetStyle("display", visible ? .prevDisplay : 'none')
		}
	GetVisible()
		{
		el = .Member?(#GetContainerEl) ? .GetContainerEl() : .El
		return el.GetStyle("display") isnt 'none'
		}
	SkipSetFocus: false
	SetFocus()
		{
		.El.Focus()
		}
	ClearFocus()
		{
		.El.Blur()
		}
	GetEnabled()
		{
		try
			return .El.disabled is false
		catch
			return true
		}
	SetEnabled(enabled)
		{
		if .GetEnabled() isnt enabled
			.El.disabled = not enabled
		}
	GetReadOnly()
		{
		return false
		}
	SetReadOnly(unused)
		{ }

	Construct(@x)
		{
		x = .build(x)
		_parent = this
		_ctrlspec = x
		ctrl = Construct(x, "Component")
		if String?(ctrl.Name) and ctrl.Name isnt ""
			this[ctrl.Name] = ctrl
		return ctrl
		}

	build(x)
		{
		if x.Size() is 1 and x.Member?(0) and Object?(x[0])
			x = x[0]
		return x
		}

	CalcMaxHeight()
		{
		if .Ystretch > 0
			return .MaxHeight
		return .Ymin
		}

	AddToolTip(tip, el = false)
		{
		if tip in (false, '')
			return
		if el is false
			el = .El
		el.SetAttribute('title', tip)
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

	Msg(args)
		{
		msg = args[0]
		target = this
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

	On_Cancel()
		{
		if .Controller.Method?(#On_Cancel)
			.Controller.On_Cancel()
		else if .Window.Method?(#On_Cancel)
			.Window.On_Cancel()
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

	Dirty?(unused = "")
		{
		return false
		}

	Get_Hwnd()
		{
		return .El
		}

	FocusFirst()
		{
		first = false
		list = .El.QuerySelectorAll("input, .su-code-wrapper")
		for i in .. list.length
			{
			el = list.Item(i)

			try
				{
				ctrl = el.Control()
				if ctrl.TabOver?()
					continue

				if first is false
					first = el

				if ctrl.GetReadOnly() isnt true
					{
					ctrl.SetFocus()
					return true
					}
				}
			}
		if first isnt false
			{
			first.Focus()
			return true
			}
		return false
		}

	GetControlFromEl(el)
		{
		while (el isnt false)
			{
			id = false
			try
				{
				control = el.Control()
				id = control.UniqueId
				}
			if id isnt false
				return control
			parentEl = false
			try
				parentEl = el.parentElement
			el = parentEl
			}
		return false
		}

	TabOver?()
		{
		return .El.tabIndex is -1
		}

	WindowRefresh()
		{
		if .Member?(#Window)
			.Window.Refresh()
		}

	SubClass()
		{
		}

	AddContextMenuItem(@unused)
		{
		}

	HasFocus?()
		{
		return SuUI.GetCurrentDocument().activeElement is .El
		}

	HandleTab()
		{
		return false
		}

	// Overriden by ChooseComponent, TextPlusComponent
	// Called by ListEditWindowComponent
	Resize(w, h)
		{
		.Xmin = w
		.Ymin = h
		.SetMinSize()
		.SetStyles(Object(width: w $ 'px', height: h $ 'px'))
		}

	mouseTracking?: false
	StartMouseTracking(mouseupCB, mousemoveCB = false)
		{
		SuRender().SetMouseMoveCB(mousemoveCB)
		SuRender().SetMouseUpCB(mouseupCB)
		.mouseTracking? = true
		}

	StopMouseTracking()
		{
		SuRender().ClearMouseMoveCB()
		SuRender().ClearMouseUpCB()
		.mouseTracking? = false
		}

	Destroyed?()
		{
		return not .Member?('Window')
		}

	Destroy()
		{
		// avoid object-modified-during-iteration, also avoid Copy overhead
//		for (n = .delayed.Size(); n > 0; --n)
//			.delayed[0].Kill()
		if .mouseTracking? is true
			.StopMouseTracking()

		if .delayedTaskId isnt false
			{
			SuRender().CancelDelayedTask(.delayedTaskId)
			.delayedTaskId = false
			}

		if .El isnt false
			{
			.El.Control(false)
			.El.Window(false)
			.El.Remove()
			}

		if .Member?(#Parent) and Instance?(.Parent) and String?(.Name) and
			.Parent.GetDefault(.Name, false) is this
			.Parent.Delete(.Name)

		if .UniqueId isnt false
			SuRender().UnRegister(.UniqueId)
		.Delete(all:) // to help garbage collection
		}
	}
