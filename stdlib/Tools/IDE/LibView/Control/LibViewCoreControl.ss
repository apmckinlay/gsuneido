// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
/* USAGE:
LibViewCoreControl manages the core functionality of the library code editors.
It can be ran as is, or overridden with extra addons defined as:
LibViewCoreControl
	{
	Addons: (
		Addon_a:,
		Addon_b:,
		Addon_c:
		)
	}
NOTE: These addons should inherit from LibViewAddon (see LibViewControl for examples).

Ultimately, the only addon this class requires to function is: Addon_LibView_Explorer.
Addon_LibView_Explorer provides the core components:
	ExplorerMultiControl / ExplorerMultiTreeControl
		-> LibTreeModel
		-> LibViewViewControl
*/
ExplorerAppsControl
	{
	Title: 	'LibraryView'
	Addons: #(/* Override this object to customize the LibView's addons */)
	New()
		{
		.Clipformat = RegisterClipboardFormat('Suneido_LIBVIEW')

		.addons.Send('Init')
		.PluginTools({|cmd, target| .Redir('On_' $ cmd, .adapt(target, this)) })

		.subs = [
			PubSub.Subscribe('LibraryTreeChange', .reset)
			PubSub.Subscribe('LibraryRecordChange', .refresh)
			]
		}

	/*
	BuildCommands is first called before New() by ExplorerAppsControl.More_Commands().
	This means we cannot set class variables in this function.

	If Addon_LibViewToolbar is included in Addons, we call this twice in order to get the
	commands only provided by the addons.
	*/
	BuildCommands(sequenced? = false)
		{
		inst = .addons isnt false ? .addons : .addonsManager()
		inst.Send('Commands', cmds = Object())
		.PluginTools({|cmd, icon, shortcut| cmds.Add([cmd, shortcut, '', icon]) })
		if not sequenced?
			cmds.Map!({ it.Member?('seq') ? it.Copy().Delete('seq') : it })
		return cmds
		}

	addons: false
	Controls()
		{
		.addons = .addonsManager()
		ctrls = Object('Vert')
		.addons.Collect('Ctrl').
			Sort!({|x, y| x.order < y.order }).
			Each({ ctrls.Add(it.ctrl) })
		return Object('Horz', #(Skip 4), ctrls, #(Skip 4))
		}

	addonsManager()
		{
		return AddonManager(this, Object(Addon_LibView_Explorer:).Merge(.Addons))
		}

	adapt(target, libview)
		{
		return { target(:libview) }
		}

	PluginTools(block)
		{
		Plugins().ForeachContribution('LibView', 'Tools')
			{|c|
			cmd = c[2] /*= command */
			icon = c[3] /*= icon */
			shortcut = c.GetDefault(4, '') /*= shortcut */
			block(:cmd, :icon, :shortcut, target: c.target)
			}
		}

	Collect(@args)
		{
		return .addons.Collect(@args)
		}

	Addon(@args)
		{
		.addons.Addon(@args)
		}

	ResetAddons()
		{
		.addons.Send('Init')
		}

	Default(@args) // used by context menu
		{
		return .send(args)
		}

	send(args)
		{
		if .addons isnt false
			.addons.ConditionalSend({ it.AddonReady?() }, args)
		return 0
		}

	Recv(@args)
		{
		return .send(args)
		}

	Getter_Explorer()
		{ return .Explorer = .FindControl('Explorer') }

	Getter_View()
		{ return .Explorer.View }

	Getter_Editor()
		{ return .View is false ? false : .View.Editor }

	Getter_Libs()
		{ return Libraries() }

	More_commands()
		{ return .BuildCommands() }

	Menu()
		{
		.ToolbarMenu(menu = Object())
		return menu
		}

	Menu_Use_Library()
		{ return LibraryTables().Difference(Libraries()).Sort!() }

	Menu_Unuse_Library()
		{ return Libraries().Remove('stdlib') }

	AllowRootDelete?()
		{ return true }

	CurrentTable()
		{ return .View.CurrentTable() }

	CurrentName()
		{ return .View.CurrentName() }

	CurrentLibView()
		{ return this }

	CloseTab?(data)
		{ return data.group and IDESettings.Get(#ide_move_tab, true) }

	Save()
		{ .Explorer.On_Save() }

	ResetCtrls() // Called via LibView Addons
		{
		state = .GetState()
		.Explorer.Reset()
		.SetState(state, resetting:)
		}

	ForceEntabAll()
		{
		.Explorer.ForeachTab()
			{|tab|
			editor = tab.FindControl('Editor')
			if editor.Get() isnt editor.Get().Entab()
				{
				pos = editor.GetCurrentPos()
				firstLine = editor.GetFirstVisibleLine()
				editor.PasteOverAll(editor.Get().Entab())
				editor.GoToPos(pos)
				editor.SetFirstVisibleLine(firstLine)
				}
			}
		}

	reset(args)
		{
		if TestRunner.RunningTests?()
			return

		if .Explorer.ResetControls(force: args.Any?({ it.force }))
			SvcCommitChecker.ClearPreCheck()
		else
			.Explorer.ForeachTab({ |view| view.Invalidate() })

		if QcIsEnabled()
			Qc_ContinuousChecks.ResetCache()
		}

	AlertTestResult(observer)
		{
		alert = observer.HasError?() ? 'AlertError' : 'AlertInfo'
		if observer.Result isnt ''
			this[alert]('Run Test', observer.Result)
		return observer.HasError?() ? observer.Result.FirstLine() : ''
		}

	Locate(item)
		{
		name = item.AfterFirst(':')
		libs = Object(item.BeforeFirst(':'))
		paths = Gotofind(:name, :libs, exact:)
		// exact match in one library should return, at most, one match
		if paths.Size() is 1
			.Explorer_RestoreTab(paths[0])
		.focusEditor()
		}

	LocateEscape()
		{ .focusEditor() }

	focusEditor() // Extracted for tests
		{
		if .Editor isnt false
			.Editor.SetFocus()
		}

	GetState()
		{
		return Object(splitterpos: .CtrlHorzSplit(.Explorer),
			outline_split: .CtrlHorzSplit(.View),
			tabs: .Explorer.GetTabsPaths(),
			activeTabPath: .Explorer.Getpath(.Explorer.CurItem))
		}

	CtrlHorzSplit(ctrl, splitData = false)
		{
		if ctrl is false or false is split = ctrl.FindControl('HorzSplit')
			return #(0, 0)

		if splitData isnt false
			split.SetSplit(splitData)
		return split.GetSplit()
		}

	SetState(statedata, resetting = false)
		{
		if not resetting
			.Explorer.RestoreState(statedata)
		if statedata.Member?('outline_split')
			.CtrlHorzSplit(.View, statedata.outline_split)
		}

	CanPaste?()
		{ return IsClipboardFormatAvailable(.Clipformat) }

	Goto_GetPath()
		{ return .Explorer.Getpath(.Explorer.GetSelected()).BeforeLast('/') }

	GotoMethodLine(method) // called by GotoLibView
		{
		text = .Editor.Get()
		if false is pos = ClassHelp.FindMethod(text, method)
			{
			if method[0].Upper?() and
				false isnt x = ClassHelp.FindBaseMethod(.CurrentTable(), text, method)
				return .GotoBaseMethod(x.lib, x.name, method)
			if false is pos = ClassHelp.FindDotDeclarations(text, method)
				return false
			}
		.gotoLine(.Editor.LineFromPosition(pos))
		return true
		}

	gotoLine(line)
		{
		.Editor.GotoLine(line, noFocus?:)
		.Defer(.focusEditor)
		}

	GotoBaseMethod(lib, name, method)
		{
		if .Explorer.GotoPath(.formatPath(LibHelp.NamePath(lib, name)), skipFolder?:)
			return .GotoMethodLine(method)
		return true
		}

	formatPath(path)
		{
		lib = path.BeforeFirst('/')
		return not Libraries().Has?(lib)
			? '(' $ lib $ ')/' $ path.AfterFirst('/')
			: path
		}

	GotoPathLine(path, line, skipFolder? = false) // called by GotoLibView
		{
		if .Explorer.GotoPath(.formatPath(path), skipFolder?)
			.gotoLine(line - 1)
		}

	Try_run(which, block = false, quiet? = false, wrapper = function (b){ return b })
		{
		if block is false
			{
			s = .Editor.Get()
			block = { .run_all(s, :quiet?) }
			}

		// LibraryRecordChange runs in the delayed publish subscriber
		.Defer(uniqueID: 'LibView_Try_run')
			{
			result = .runBlock(wrapper(block))
			.setViewStatusWithResults(result, which, .View)
			}
		}
	runBlock(block)
		{
		blockResult = #()
		err = ''
		try
			blockResult = block()
		catch (e)
			err = e
		return Object(:blockResult, :err)
		}
	setViewStatusWithResults(result, which, view)
		{
		if result.err isnt ''
			view.Status("Run: run " $ which $ " failed: " $ result.err, invalid:)
		else
			{
			if not result.blockResult.Empty?()
				Print(Display(result.blockResult[0]))
			view.Status("Run: run " $ which $ " successful", valid:)
			}
		}

	run_all(text, quiet? = false)
		{
		text = text.Trim()
		x = text.Eval() // should be Eval
		if Object?(x) and not Class?(x)
			result = [Window(x)]
		else if Class?(x) and x.Base?(Test)
			result = .run_test(x, :quiet?)
		else if Function?(x) or Class?(x)
			result = (text $ "()").Eval2()
		else
			result = [x]
		return result
		}

	run_test(x /*unused*/, quiet? = false)
		{
		name = .CurrentName()
		lib = .CurrentTable()
		.Editor.SendToAddons('On_BeforeAllTests')
		LibViewRunTest(.Editor, lib, name)
			{
			observer = it.RunTest(name, quiet:)
			}
		.Editor.SendToAddons('On_AfterAllTests')
		if not quiet? and '' isnt result = .AlertTestResult(observer)
			throw result
		else if observer.HasError?()
			throw observer.Result.FirstLine()
		return [] // for run_all
		}

	Print(s)
		// pre:	s is a string
		// post:	if a console for this exists, s is appended to this' console ELSE
		//		a console for this is created and s is appended to it
		{
		r = GetWorkArea()
		if not (Suneido.Member?('Console') and Suneido.Console isnt false)
			Window(#(Console), x: r.left, y: r.top, w: .35, h: .5)
		Suneido.Console.Append(s)
		}

	// warning - if dirty and called prior to save, changes are lost
	refresh(args)
		{
		if .Explorer.Refresh(args)
			SvcCommitChecker.ClearPreCheck()
		}

	CloseFolder(path)
		{
		tree = .Explorer.Tree
		list = tree.GetChildren(TVI.ROOT)
		path = path.Split('/')

		item = false
		for (i = 0; i < path.Size(); i++)
			{
			if list.Empty?()
				return
			item = false
			for child in list
				if path[i] is tree.GetName(child)
				{
				item = child
				break
				}
			if item isnt false
				list = tree.GetChildren(item)
			}
		tree.ExpandItem(item, true)
		}

	LibExportFile(lib = false, name = false, item = false)
		{
		if item is false
			item = .Explorer.GetSelected()
		if .Explorer.Container?(item)
			return
		if lib is false
			lib = .CurrentTable()
		if name is false
			name = .CurrentName()
		DoWithSaveFileName(title: 'Export (append) to', hwnd: .Window.Hwnd,
			flags: OFN.PATHMUSTEXIST | OFN.HIDEREADONLY | OFN.NOCHANGEDIR)
			{ |fileName|
			LibIO.Export(lib, name, fileName, interactive:)
			}
		}

	Destroy()
		{
		.subs.Each(#Unsubscribe)
		.addons.Send('Destroy')
		super.Destroy()
		}
	}
