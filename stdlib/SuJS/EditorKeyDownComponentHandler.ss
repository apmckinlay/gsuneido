// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(component, event, pressed, extraCommands = #())
		{
		if event.key is "F6"
			.keyDownEvent(component, VK.F6, :pressed, :event)
		if event.ctrlKey is true and event.key is 'f'
			.keyDownEvent(component, VK.F, :pressed, :event)
		if event.ctrlKey is true and event.key is 'p'
			.keyDownEvent(component, VK.P, :pressed, :event)
		if event.key is 'F3'
			.keyDownEvent(component, VK.F3, :pressed, :event)
		if extraCommands.NotEmpty?()
			{
			key = .getKey(event, pressed)
			command = pressed.FindAll(true).Add(key)
			if extraCommands.Any?({ it.EqualSet?(command) })
				.keyDownEvent(component, VK[key.Upper()], :pressed, :event)
			}
		}

	codeToKey: (
		'Slash': '/')
	getKey(event, pressed)
		{
		key = pressed.shift isnt true
			? event.key.Lower()
			: .codeToKey.GetDefault(event.code, event.key).Lower()
		return key.RemovePrefix('arrow')
		}

	keyDownEvent(component, key, pressed, event)
		{
		component.RunWhenNotFrozen()
			{
			component.EventWithOverlay(#KEYDOWN, key, :pressed)
			}
		event.PreventDefault()
		event.StopPropagation()
		}
	}
