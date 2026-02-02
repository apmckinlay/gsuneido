// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	commands: #()
	Set(menuOptions)
		{
		.commands = Object()
		accelerators = Object()
		menuOptions.Each()
			{
			menuOptionOb = it.Split('\t').Map('Trim')
			if '' is command = menuOptionOb.GetDefault(1, '')
				continue
			if accelerators.Has?(command)
				.logDuplicateLevels(it)
			else
				{
				accelerators.Add(command)
				.commands.Add([
					command: .formatCommandOb(command.Split('+')),
					method: .translateCommand(menuOptionOb[0])
					])
				}
			}
		}

	logDuplicateLevels(option)
		{ SuneidoLog('ERROR: (CAUGHT) ' $ option $ ' duplicate commands. Skipping') }

	formatCommandOb(commandOb)
		{ return commandOb.Map({ it.Trim().Lower() }).Sort!() }

	translateCommand(option)
		{ return ContextMenu.MakeItemIntoCommand(option).Replace('On_Context_', 'On_') }

	GetMethod(command)
		{
		commandOb = .commands.FindOne({ it.command.EqualSet?(command) })
		return commandOb is false ? false : commandOb.method
		}

	BuildCommandOb(key, pressed = false)
		{
		commandOb = Object(String(key))
		if .keyPressed?(VK.SHIFT, :pressed)
			commandOb.Add('Shift')
		if .keyPressed?(VK.CONTROL, :pressed)
			commandOb.Add('Ctrl')
		if .keyPressed?(VK.MENU, :pressed)
			commandOb.Add('Alt')
		return .formatCommandOb(commandOb)
		}

	keyPressed?(key, pressed)
		{
		return KeyPressed?(key, :pressed)
		}

	GetCommands()
		{
		return .commands.Map({ it.command })
		}

	Destroy()
		{ .commands = Object() }
	}