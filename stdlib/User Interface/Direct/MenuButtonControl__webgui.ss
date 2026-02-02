// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
ButtonControl
	{
	ComponentName: 'MenuButton'
	New(text, .menu = false, tip = false, tabover = false,
		.left = false, width = false, command = false, .sendParents? = false)
		{
		super(text, pad: 36, :tip, :tabover, :width, :command)
		if command is false
			command = .Name
		.command = 'On_' $ ToIdentifier(command)

		.ComponentArgs = Object(text, tip, tabover, left, width)
		}

	SetMenu(menu)
		{
		.menu = menu
		}

	MenuButton_PullDown(x, y, rcExclude, buttonRect)
		{
		if .disabled is true or not IsWindowEnabled(.Window.Hwnd)
			return

		menu = .menu isnt false ? .menu : .Send('MenuButton_' $ .Name)
		i = ContextMenu(menu).Show(0, x, y, :rcExclude, :buttonRect)
		if i > 0
			.send(.command, menu, i - 1)
		}

	PopupMenu(menu, hwnd)
		{
		i = ContextMenu(menu).Show(hwnd, false, false)
		return i
		}

	send(prefix, menu, chosen, j = 0, parent = false)
		{
		for (m = 0; j <= chosen and m < menu.Size(); ++j, ++m)
			if (Object?(menu[m]))
				j = .send(prefix $ "_" $ ToIdentifier(menu[m - 1]),
					menu[m], chosen, j, prefix) - 1
			else if (j is chosen)
				{
				if parent isnt false and .sendParents?
					.Send(parent, prefix $ "_" $ ToIdentifier(menu[m]))
				.Send(prefix, menu[m], index: m)
				.Send(prefix $ "_" $ ToIdentifier(menu[m]))
				}
		return j
		}

	grayed: false
	Grayed?(state = -1)
		{
		if state isnt -1 and .grayed isnt state
			{
			.grayed = state
			.Act(#Grayed, .grayed)
			}
		return .grayed
		}
	disabled: false
	Disable(readonly)
		{
		.disabled = readonly
		.Act(#Disable, .disabled)
		}
	Disabled?()
		{
		return .disabled
		}

	Getter_Command()
		{
		return .command
		}
	}