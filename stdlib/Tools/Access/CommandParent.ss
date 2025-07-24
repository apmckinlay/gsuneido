// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// Make keyboard accelerators for Commands work inside a book
// base class for AccessControl and Access1Control (via AccessBase)
// this could be used more generally (in which case it should be renamed)
// e.g. NextTab, PrevTab
PassthruController
	{
	Startup()
		{
		if .Send('CommandParent_SkipInitCommands?') is true
			return
		.AddCommands()
		}

	AddCommands()
		{
		.curAccels = .Window.SetupAccels(.Commands)
		.tabs = .FindControl('Tabs') // may be false if no tabs
		if .Window.Member?('Ctrl')
			.redirCmds()
		else
			.Defer(.redirCmds)
		}

	redirCmds()
		{
		for cmd in .Commands
			.Window.Ctrl.Redir('On_' $ cmd[0],
				cmd[0].Has?('Tab') ? .tabs : this)
		}

	ResetCommands()
		{
		if not .Window.Member?('Ctrl') // not initialized yet
			return
		.RemoveCommands()
		.AddCommands()
		}

	tabs: false
	curAccels: false
	Destroy()
		{
		.RemoveCommands()
		super.Destroy()
		}

	RemoveCommands()
		{
		if .curAccels isnt false
			{
			if not .Window.Member?('Ctrl')
				return
			.Window.Ctrl.RemoveRedir(this)
			.Window.Ctrl.RemoveRedir(.tabs)
			.Window.RestoreAccels(.curAccels)
			.curAccels = false
			}
		}
	}
