// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Debug'
	Xmin: 500
	CallClass(hwnd, err, calls = false)
		{
		.Dialog(hwnd, err, calls)
		}
	Dialog(hwnd, err, calls = false)
		{
		ToolDialog(hwnd, Object(this, hwnd, err, calls), border: 0)
		}
	Window(hwnd, err, calls = false, onDestroy = false)
		{
		if Suneido.Member?(#Debugger)
			Suneido.Debugger.Destroy()
		Suneido.Debugger = Window(Object(this, hwnd, err, calls), keep_placement:,
			:onDestroy)
		}
	New(hwnd, err, calls)
		{
		.list = .Vert.VertSplit.HorzSplit.ListBox
		.inspect = .Vert.VertSplit.HorzSplit.Inspect
		.debugView = .FindControl(#DebugView)
		.source = .debugView.Editor
		.source.SETCARETLINEVISIBLEALWAYS(true)
		.source.SETYCARETPOLICY(SC.CARET_STRICT | SC.CARET_EVEN)

		.Init(hwnd, err, calls)
		if hwnd is 0 or .Window.Base?(Dialog)
			.FindControl('Destroy_Window').SetEnabled(false)
		}
	Init(.hwnd, err, calls)
		{
		if calls is false
			try calls = err.Callstack()
			catch calls = #()
		RemoveAssertsFromCallStack(calls)
		for c in calls
			.list.AddItem(Unprivatize(Display(c.fn)))
		.calls = HandleSetConcurrentError(calls)
		.SetupError(err)
		}
	SetupError(err, color = false)
		{
		if color is false
			color = CLR.ErrorColor
		errctrl = .FindControl('Err')
		errctrl.Set(err)
		errctrl.SetSelect(0)
		errctrl.SetMarginWidthN(1, 0)
		errctrl.StyleSetBack(0, color)
		errctrl.StyleSetBack(SC.STYLE_DEFAULT, color)
		}
	Startup()
		{
		if .calls.Empty?()
			return
		pos = 0
		while Name(.calls[pos].fn).Has?('Assert')
			++pos
		// need adding delay to make code and call stack list scroll to the proper
		// position when debugger is shown as a dialog
		.Defer()
			{
			.list.SetCurSel(pos)
			.ListBoxSelect(pos, .list)
			}
		}

	Controls()
		{
		buttons = Object(#Horz
			#Skip
			#(Button "Destroy Window")
			#Skip #Skip
			#(Button "Exit")
			#Fill
			#(Button "Go To Definition")
			#Skip)
		buttons.Add(@.ExtraButtons())
		return Object(#Vert
			Object(#VertSplit
				Object(#VertSplit
					Object(.errControl, wrap:, readonly:, ystretch: 0, name: #Err)
					#(CodeView, addons: #(Addon_status: false, Addon_overview_bar: false),
						readonly:, name: DebugView)
				)
				#(HorzSplit
					(Inspect #() themed?: true, xmin: 250, xstretch: 1, ymin: 150)
					(ListBox themed?: true, xmin: 250, xstretch: 1, ymin: 150)
					ystretch: 1
					)
				)
			#(Skip 3)
			buttons
			#(Skip 3)
			)
		}

	errControl: ScintillaAddonsControl { IDE: }

	ClassOutline_SkipHierarchy?()
		{ return true }

	ClassOutline_SelectTreeItem()
		{ return false }

	ExtraButtons()
		{ return #(#(Button Cancel xmin: 80) #Skip) }

	Commands()
		{ return #(#(Toggle_Disasm, "Ctrl+Alt+Shift+/", "Toggle source/disassembly")) }

	ListBoxSelect(i, source)
		{
		if source isnt .list or not .calls.Member?(i)
			return
		.cur_call = i
		try .inspect.Reset(.calls[i].locals)
		.showCurCall()
		}
	ListBoxDoubleClick(i /*unused*/)
		{
		.On_Go_To_Definition()
		}
	On_Go_To_Definition()
		{
		// TODO make this go to the right library
		GotoLibView(.name(class?:), line: .source.LineFromPosition() + 1)
		}
	disassembled: false
	On_Toggle_Disasm()
		{
		// When user types CTRL+ALT+SHIFT+/, toggles between disassembly view
		// and ordinary source code view.
		.disassembled = not .disassembled
		.showCurCall()
		}
	On_Exit()
		{
		Exit(true)
		}
	On_Destroy_Window()
		{
		DestroyWindow(.hwnd)
		}
	On_Cancel() // need when in a Window
		{
		if .Window.Base?(Window)
			.Window.Destroy()
		else
			super.On_Cancel()
		}
	Destroy()
		{
		Suneido.Delete(#Debugger)
		super.Destroy()
		}
	showCurCall()
		{
		frame = .calls[.cur_call]
		.ShowSource(frame, .disassembled)
		}
	ShowSource(frame, disassembled = false)
		{
		fn = frame.fn
		if disassembled
			source = .disassembledText(fn)
		else if false is source = SourceCode(fn)
			source = "SOURCE CODE NOT AVAILABLE"
		.SetDebugView(fn, source)
		if frame.Member?("srcpos")
			.source.SetSelect(frame.srcpos)
		}
	SetDebugView(fn, source)
		{
		.debugView.Table = .lib(fn)
		.debugView.RecName = .name(class?:)
		.debugView.Set([text: source])
		.debugView.Outline_Highlight(.debugView.RecName, .name())
		}
	lib(fn)
		{
		return Display(fn).AfterFirst('/* ').BeforeFirst(' ')
		}
	name(class? = false)
		{
		name = .list.GetText(.list.GetCurSel()).BeforeFirst(' /*')
		return class?
			? name.BeforeFirst('.')
			: name.AfterFirst('.')
		}
	disassembledText(fn)
		{
		disasm = ''
		try
			if not String?(disasm = Disasm(fn))
				throw 'disassembled is not a string'
		catch (e)
			disasm = 'DISASSEMBLY NOT AVAILABLE' $ Opt(' - ', e)
		return disasm
		}
	Getter_Source() { return .source }
	cur_call: false   // Index of current call frame
	Getter_CurCall() { return .cur_call}
	}
