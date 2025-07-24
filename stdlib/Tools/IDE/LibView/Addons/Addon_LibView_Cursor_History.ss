// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
LibViewAddon
	{
	maxHistory: 	50
	historyPos:		-1
	idleSeconds: 	3
	history: 		false
	lastChecked: 	false
	Initialized: 	false

	Commands(cmds)
		{
		cmds.Add(
			#(Previous_Location, "Alt+Left", "Back", "Back", seq: 3.0)
			#(Next_Location, "Alt+Right", "Forward", "Forward", seq: 3.1)
			)
		}

	New(@args)
		{
		super(@args)
		.initVars()
		}

	initVars()
		{
		.history = Object()
		.historyPos = -1
		.lastChecked = Date().Minus(seconds: .idleSeconds)
		}

	Init()
		{
		if not .Initialized = .Explorer isnt false and .View isnt false and
			.Editor isnt false
			{
			.initVars()
			return
			}

		.Redir('Next_Location')
		.Redir('Previous_Location')
		}

	On_Previous_Location()
		{
		if not .history.Empty?() and .historyPos is .history.Size() - 1
			.save(.View)
		.historyPos = Max(.historyPos - 1, 0)
		.gotoHistoryPos()
		}

	On_Next_Location()
		{
		.incrementHistoryPos()
		.gotoHistoryPos()
		}

	incrementHistoryPos()
		{ .historyPos = Min(.historyPos + 1, .history.Size() - 1) }

	// Prevents saving duplicate locations when already changing locations
	going?: false
	gotoHistoryPos()
		{
		if .history.Empty?()
			return
		name = .View.CurrentName()
		lib = .View.CurrentTable()
		rec = .history[.historyPos]
		.going? = true
		if name is rec.name and lib is rec.lib
			.View.GetChild().GotoLine(rec.pos)
		else
			GotoLibView(rec.lib $ ':' $ rec.name, line: rec.pos + 1, libview: .Parent)
		.going? = false
		}

	Explorer_DeselectTab(oldView)
		{ .save(oldView) }

	TreeView_ItemClicked(@unused)
		{ .newSelection() }

	LibView_Goto()
		{ .newSelection() }

	newSelection()
		{
		bef = .historyPos
		.save(.View)
		if bef is .historyPos
			.incrementHistoryPos()
		}

	EditorChange()
		{
		if .lastChecked.Plus(seconds: .idleSeconds) < Date()
			.save(.View)
		}

	save(view)
		{
		.savePosition(view.CurrentTable(), view.CurrentName(), view.CurrentLine())
		.lastChecked = Date()
		}

	savePosition(lib, name, pos)
		{
		historyOb = Object(:lib, :name, :pos)
		if historyOb is .history.GetDefault(.historyPos, false) or .going?
			return
		for (i = .history.Size() - 1; i > .historyPos; i--)
			.history.Delete(i)

		.history.Add(historyOb)
		if .history.Size() > .maxHistory
			.history.Delete(0)
		else
			.historyPos = Min(.historyPos + 1, .history.Size())
		}
	}
