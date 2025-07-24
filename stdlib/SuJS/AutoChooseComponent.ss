// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
EditComponent
	{
	Name: AutoChoose
	Xstretch: 1
	New(width = 20, readonly = false, height = 1, font = "", size = "")
		{
		super(:readonly, :width, :height, :font, :size)
		.El.AddEventListener('keydown', .keydown)
		}

	getter_keymap()
		{
		.keymap = Object(
			ArrowUp: VK.UP,
			ArrowDown: VK.DOWN,
			Tab: VK.TAB,
			Enter: VK.RETURN,
			Escape: VK.ESCAPE)
		}

	keydown(event)
		{
		if .listOpen? and .keymap.Member?(event.key)
			{
			.Event(#KEYDOWN, .keymap[event.key])
			event.PreventDefault()
			event.StopPropagation()
			}
		}

	GetListPos() // called by AutoChooseListComponent
		{
		return SuRender.GetClientRect(.El)
		}

	listOpen?: false
	SyncListStatus(.listOpen?) {}
	}
