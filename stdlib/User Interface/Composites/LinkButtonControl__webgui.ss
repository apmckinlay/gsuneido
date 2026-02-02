// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
StaticControl
	{
	ComponentName: 'LinkButton'
	New(name, command = false, tip = false, color = 'BLUE', bgndcolor = "",
		tabstop = true, font = "", size = "", weight = "")
		{
		super(name, underline:, :color, :tip, :tabstop,
			:bgndcolor, :font, :size, :weight)
		.Name = ToIdentifier(name.Trim())
		if command is false
			command = name
		.command = "On_" $ ToIdentifier(command.Trim())
		.ComponentArgs = Object(name, tip, color, bgndcolor, tabstop, font, size, weight)
		}

	CLICK()
		{
		.Send(.command)
		return 0
		}

	ContextMenu(x, y)
		{
		.Send(.command $ "_ContextMenu", x, y)
		return 0
		}
	}