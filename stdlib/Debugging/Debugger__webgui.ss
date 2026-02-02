// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
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
	New(hwnd, .err, calls)
		{
		.list = .Vert.VertSplit.HorzSplit.ListBox
		.inspect = .Vert.VertSplit.HorzSplit.Inspect
		.debugView = .FindControl('DebugView')
		.allLocals = .FindControl(#allLocals)
		.source = .debugView.Editor
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
			.list.AddItem(Unprivatize(String?(c.fn) ? c.fn : Display(c.fn)))
		.calls = HandleSetConcurrentError(calls)
		.SetupError(err)
		}
	SetupError(err, color = false)
		{
		if color is false
			color = CLR.ErrorColor
		errctrl = .FindControl('Err')
		errctrl.Set(err)
		errctrl.SetReadOnlyBrush(color)
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
			#Fill,
			#(CheckBox, 'All Locals', name: 'allLocals'),
			#Skip,
			#(MenuButton, 'Copy Error', ('Full Stack', 'From Selected Stack')),
			#Skip)
		buttons.Add(@.ExtraButtons())
		return Object(#Vert
			Object(#VertSplit
				Object(#VertSplit
					#(Editor, readonly:, height: 2, name: 'Err'),
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

	ExtraButtons()
		{
		return #(#(Button Cancel xmin: 80) #Skip)
		}

	cur_call: false   // Index of current call frame
	Getter_CurCall() { return .cur_call}
	ListBoxSelect(i, source)
		{
		if (source isnt .list)
			return
		.cur_call = i
		try .inspect.Reset(.calls[i].locals)
		.showCurCall()
		}

	showCurCall()
		{
		frame = .calls[.cur_call]
		.ShowSource(frame)
		}

	ShowSource(frame)
		{
		fn = frame.fn
		if false is source = SourceCode(fn)
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
//		.debugView.Outline_Highlight(.debugView.RecName, .name())
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

	On_Exit()
		{
		ExitClient()
		}

	On_Destroy_Window()
		{
		DestroyWindow(.hwnd)
		}

	On_Copy_Error_Full_Stack()
		{
		.copyError(0)
		}

	copyError(start)
		{
		ClipboardWriteString('Error:\r\n\t' $ .err $ '\r\n' $
			'Locals:\r\n\t' $ .localsStr(start) $ '\r\n' $
			'Stack:\r\n' $ .stackStr(start, .allLocals.Get() is true) $ '\r\n')
		}

	localsStr(sel)
		{
		return String(.calls[sel].locals)
		}

	stackStr(sel, allLocals?)
		{
		stack = Object()
		for i in sel .. .list.GetCount()
			stack.Add('\t' $ .list.GetText(i) $
				(allLocals?
					? '\r\n\t\t' $ .localsStr(i)
					: ''))
		return stack.Join('\r\n')
		}

	On_Copy_Error_From_Selected_Stack()
		{
		.copyError(.list.GetCurSel())
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
	}