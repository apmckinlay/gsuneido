// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "WorkSpace"
	Name: "WorkSpace"

	New()
		{
		.libs = Libraries()
		.editor = .FindControl('Editor')
		.output = .FindControl('Output')
		.tabs = .FindControl('Tabs')
		.inspect = .FindControl('Variables')
		.output_horzsplit = .FindControl('HorzSplit')
		.output_vertsplit = .FindControl('VertSplit')
		.statusbar = .FindControl('Status')
		.startTimer()
		.Redir('On_Copy')
		.Redir('On_Find')
		.Redir('On_Find_Next')
		.Redir('On_Find_Previous')
		.Redir('On_Replace')
		.Redir('On_NextTab', .tabs)
		.Redir('On_PrevTab', .tabs)
		.Redir('On_Comment_Lines')
		.Redir('On_Comment_Selection')

		Plugins().ForeachContribution('WorkSpace', 'Toolbar')
			{|c|
			.Redir("On_" $ c[2], c.target)
			}
		.subs = [
			PubSub.Subscribe('LibraryTreeChange', .setLibs)
			PubSub.Subscribe('Redir_SendToEditors', .sendToEditor)
			]
		.vars = Object()
		}

	setLibs()
		{
		if not .WindowActive?()
			.libs = Libraries()
		}

	sendToEditor(@args)
		{ .editor.SendToAddons(@args) }

	workspace_tab: 0
	find_tab: 1
	Controls()
		{
		IDESettings.Init()
		selectedTabColor = IDESettings.Get('ide_selected_tab_color', false)
		selectedTabBold = IDESettings.Get('ide_selected_tab_bold', true)
		Object(#Vert
			Object(#Horz,
				.toolbar(),
				#(Vert
					(EtchedLine before: 0 after: 0)
					(Horz
						Fill
						LibLocate
						Skip
						(WebLink, "suneido.com", "suneido.com")
						Skip)
					)
				)
			Object('Tabs'
				#(VertSplit
					(WorkSpaceCode ystretch: 4)
					(HorzSplit ystretch: 1
						(WorkSpaceOutput readonly:, xstretch: 5, name: 'Output')
						(Inspect (), name: 'Variables', themed?:, xstretch: 1,
							xmin: 100, ymin: 100)
						)
					Tab: ' Workspace ')
				Object('FindInLibraries' Tab: ' Find '),
				:selectedTabColor, :selectedTabBold, orientation: 'right')
			#Statusbar
			)
		}

	toolbar()
		{
		toolbar = Object('Toolbar', 'Cut', 'Copy', 'Paste', '',
			'Undo', 'Redo', "",
			'Comment_Lines', 'Comment_Selection', "",
			'Clear_Results', 'Run', '')
		Plugins().ForeachContribution('WorkSpace', 'Toolbar')
			{|c|
			toolbar.Add(c[2])
			}
		return toolbar
		}

	cmds: (
		(Find,				"Ctrl+F",	"Find text")
		(Find_Next,			"F3",		"Find the next occurrence")
		(Find_Previous,		"Shift+F3",	"Find the previous occurrence")
		(Replace,			"Ctrl+H",	"Find and replace text in the current item")
		(Run,				"F9",		"Execute the current selection", '!')
		(Run_on_Server,		"Alt+F9",	"ServerEval the current selection")
		(Disassemble,		"Shift+F9",	"Compile and disassemble the current selection")
		(Clear_and_Run,		"Ctrl+F9",	"Clear and then execute the current selection")
		(Clear_Results,		"", 		"", "delete")
		(Access_a_Query)
		(Browse_a_Query)
		(Find_in_Libraries,
			"Ctrl+Shift+F", "Search libraries for a string", Find_in_Folders)
		(Exit				""			"Exit from Suneido")
		(Users_Manual,		"F1")
		(Comment_Lines,
			"Ctrl+/", "Comment/Uncomment selected lines", "CommentLine")
		(Comment_Selection,
			"Shift+Ctrl+/", "Comment/Uncomment selected text", "CommentSpan")
		(NextTab,			"Ctrl+Tab")
		(PrevTab,			"Shift+Ctrl+Tab")
		(Locate,			"Ctrl+L")
		(Copy, 				"Ctrl+C", "Copy the selected text to the clipboard")
		)

	Commands()
		{
		cmds = .cmds.Copy()
		Plugins().ForeachContribution('WorkSpace', 'Toolbar')
			{|c|
			id = 2
			accel = 4
			bitMap = 3
			cmds.Add([c[id], c.GetDefault(accel, ''), '', c[bitMap]])
			}
		return cmds
		}

	Menu:
		(
		("&File",
			"E&xit")
		("&Edit",
			"&Undo", "&Redo", "", "Cu&t", "&Copy", "&Paste", "&Delete", "Select &All", "",
			"&Find...", "Find &Next", "Find &Previous", "R&eplace...",
			"Comment &Lines", "Comment &Selection", "",
			"&Find in Libraries...", /*"&Replace in Libraries...",*/ "",
			"R&un", "&Benchmark", "Profile", "Clear and Run", "Disassemble",
			"Run on Server", "Clear Results")
		)

	Startup()
		{
		Suneido.Print = .Print
		errors = Object()
		for source, call in GetContributions('WorkSpaceStartup')
			try
				call()
			catch (e)
				errors[source] = e
		if errors.NotEmpty?()
			SuneidoLog('ERROR: (CAUGHT) issues encountered during workspace startup',
				params: errors, caughtMsg: 'IDE level error')
		}

	Activate()
		{
		Suneido.Print = .Print
		}
	Print(s)
		{
		// In case Suneido.Print still points here after the WorkSpace has been destroyed
		if .Destroyed?()
			return
		if Sys.MainThread?()
			.print(s)
		else
			TracePrint(s)
		}
	print(s)
		{
		.tabs.Select(.workspace_tab, keepFocus:)
		.output.AppendText(s)
		.output.Update()
		}
	On_Clear_and_Run()
		{
		.clear_output()
		.On_Run()
		}
	On_Run()
		{
		if false is s = .getSelected()
			return

		result = .evalOrRun(s)
		if String?(result) and result.Prefix?("ERROR")
			{
			Alert(result.RemovePrefix("ERROR"), "Run Error", flags: MB.ICONERROR)
			return
			}
		.printResult(result)
		}
	evalOrRun(s)
		{
		if (s =~ "^[A-Z][a-zA-Z_]*[?!]?$")
			try
				return s.Eval2()
			catch (x, "can't find")
				return "ERROR" $ x
		try
			("function () {\n" $ s $ "\n}").Compile()
		catch (e)
			return "ERROR" $ e
		return .run(s)
		}
	On_Benchmark()
		{
		if false is s = .getSelected()
			return
		fn = ("function () {\n" $ s $ "\n}").Compile()
		Print(Bench(fn))
		}
	On_Profile()
		{
		if false is s = .getSelected()
			return
		fn = ("function () {\n" $ s $ "\n}").Compile()
		RunWithProfile(fn)
		}
	On_Run_on_Server()
		{
		if false is s = .getSelected()
			return
		x = ServerEval('WorkSpaceControl.Run_on_Server', s)
		.printResult(x)
		}
	Run_on_Server(code)
		{
		return code.Eval2() // requires Eval
		}
	getSelected()
		{
		if .tabs.GetSelected() isnt .workspace_tab
			return false

		if KeyPressed?(VK.CONTROL)
			.clear_output()
		s = GetRunText(.editor)
		if s is ""
			{
			Beep()
			return false
			}
		return s
		}
	printResult(x)
		{
		// have to check before outputting the result in case the code evaluated
		// destroys the WorkSpace (creating persistent sets for example)
		if .Destroyed?()
			return

		if x.Empty?()
			.Print("ok\r\n")
		else
			{
			x = x[0]
			if String?(x) and x.Has?('\n') and x.Tr(" -~\t\r\n") is ""
				.Print(x $ "\r\n")
			else
				.Print(Display(x) $ "\r\n")
			}
		// need to SetFocus because something is clearing focus
		// if you Run code with Print in it
		if GetActiveWindow() is .Window.Hwnd
			.editor.SetFocus()
		}
	run(s)
		{
		_wsvars = .vars
		// z_ is required on gSuneido to make the blocks closures
		code = .vars.Members().
			Filter({ it.Identifier?() }). // exclude block parameters
			Map({ it $ " = _wsvars." $ it $ "\n" }).Join() $
				"z_ = Type(0); Finally({ z_;;\n" $
				s $ "\n" $
				"}, { z_\n" $
				"_wsvars.Merge(HandleSetConcurrentError(Locals(0))) })"
		x = code.Eval2()
		if .Destroyed?()
			return x
		.vars = _wsvars
		.vars.Delete(#this, #_wsvars)
		if x.Empty?()
			.vars.Delete(#z_)
		else
			try
				.vars.z_ = x[0]
			catch (unused, "*cannot be set to concur")
				.vars.z_ = Display(x[0])
		.inspect.Reset(.vars)
		return x
		}

	On_Disassemble()
		{
		if false is s = .getSelected()
			return
		f = "function () {\n" $ s $ "\n}"
		Print(f.Compile().Disasm(source: f))
		}

	On_Clear_Results()
		{
		ctrl = .tabs.GetSelected() is .workspace_tab ? this : .tabs.GetControl(.find_tab)
		ctrl.ClearResults()
		}

	ClearResults()
		{
		.clear_output()
		.clear_variables()
		}
	clear_output()
		{
		.output.SetReadOnly(false)
		.output.Set("")
		.output.SetReadOnly(true)
		}
	clear_variables()
		{
		for val in .vars
			switch (Type(val))
				{
			case 'Transaction' :
				try val.Rollback()
			case 'COMobject' :
				val.Release()
			default :
				try val.Close()
				}
		.inspect.Reset(.vars = Object())
		}

	TabsControl_SelectTab()
		{
		if .tabs.GetSelected() is 1
			.Defer(.tabs.GetControl(.find_tab).SetDefaultFocus)
		}

	On_Find_in_Libraries()
		{
		.tabs.Select(.find_tab, keepFocus:)
		.tabs.GetControl(.find_tab).SetFocus()
		}

	On_Locate()
		{
		.FindControl('LibLocate').SetFocus()
		}
	Locate(name)
		{
		GotoPersistentWindow('LibViewControl', LibViewControl).Locate(name)
		}

	GetState()
		{
		return Object(text: .editor.Get(),
			editor_state: .editor.GetState(),
			vert_split: .output_vertsplit.GetSplit(),
			horz_split: .output_horzsplit.GetSplit(),
			color_scheme: IDE_ColorScheme.GetCurrent())
		}
	SetState(state)
		{
		if not Object?(state)
			return
		if state.Member?('text')
			.editor.Set(state.text)
		if state.Member?('editor_state')
			.editor.SetState(state.editor_state)
		if state.Member?('vert_split')
			.output_horzsplit.SetSplit(state.vert_split)
		if state.Member?('horz_split')
			.output_horzsplit.SetSplit(state.horz_split)
		.Window.ResetStyle()
		}

	startTimer()
		{
		.timer = SetTimer(NULL, 0, 2.SecondsInMs(), .timerFunc)
		}
	timerFunc(@unused) // args from timer
		{
		.SetStatus()
		if .libs isnt Libraries()
			{
			.libs = Libraries().Copy()
			SvcTable.Publish('TreeChange', type: 'lib', force:)
			// This should only run when two client workspaces are running.
			// This is to handle if ONE of the work enviroments changes the used libs
			if Sys.Client?()
				Unload()
			}
		.printErrorLog()
		}
	killTimer()
		{
		KillTimer(NULL, .timer)
		.timer = false
		ClearCallback(.timerFunc)
		}

	SetStatus()
		{
		.statusbar.Set('  ' $ .built() $ '\t\t' $ WorkSpaceStatus() $ '      ')
		.statusbar.ToolTip(Built() $ '\n' $ WorkSpaceStatus.ResourceDetails())
		}

	built()
		{
		return (Sys.Client?() ? "CLIENT " : "") $
			'Built: ' $ Built().BeforeFirst('(')
		}

	errlog_size: false
	printErrorLog()
		{
		try
			{
			old_size = .errlog_size
			.errlog_size = FileSize(.errfile)
			if old_size is false or .errlog_size is old_size
				return
			File(.errfile)
				{|f|
				f.Seek(old_size)
				while false isnt line = f.Readline()
					Print(ERRORLOG: line)
				}
			}
		}
	getter_errfile()
		{
		return .errfile = Client?()		// once only
			? Paths.Combine(Getenv("APPDATA"), 'suneido' $ ServerPort() $ '.err')
			: Paths.Combine(GetCurrentDirectory(), 'error.log')
		// need to include current directory because OpenFileName changes current dir
		}

	Destroy()
		{
		.subs.Each(#Unsubscribe)
		.killTimer()
		.clear_variables()
		super.Destroy()
		}
	}
