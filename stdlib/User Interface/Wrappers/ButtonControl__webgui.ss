// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Control
	{
	ComponentName: 'Button'
	New(.text, command = false, font = "", size = "",
		weight = "", tabover = false, defaultButton = false,
		style/*unused*/ = 0, tip = false, pad = false,
		color = false, width = false, underline = false,
		italic = false, strikeout = false, hidden = false)
		{
		if .Name is ""
			.Name = ToIdentifier(text)
		if (command is false)
			command = .Name
		command = command.Trim()
		.command = "On_" $ ToIdentifier(command)

		.SuSetHidden(hidden)
		.ComponentArgs = Object(text, font, size, weight, tabover, defaultButton,
			tip, pad, color, width, underline, italic, strikeout, hidden)
		}

	target: false
	CLICKED()
		{
		// ensure focus to trigger kill focus on virtual list from server
		// for example, clicking outside button when there is a list edit window open
		.SetFocus()
		.BeforeCallCommand()
		if .target isnt false
			.target[.command]()
		else
			.Send(.command)
		}
	BeforeCallCommand()
		{
		}

	SetCommandTarget(.target) { }

	Get()
		{
		return .text
		}

	Set(.text)
		{
		.Act(#Set, .text)
		}

	// NOTE: pushed state is not visible on ButtonControl
	// use EnhancedButtonControl if you want a visible pushed state
	pushed: false
	Pushed?(state = '')
		{
		if state isnt ''
			.pushed = state
		return .pushed
		}

	SetHidden(hidden)
		{
		.SuSetHidden(hidden)
		}

	ContextMenu(x, y)
		{
		.Send(.command $ "_ContextMenu", x, y)
		return 0
		}
	}
