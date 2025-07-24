// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
StaticComponent
	{
	ContextMenu: true
	New(text, tip = false, color = 'BLUE', bgndcolor = "",
		tabstop = true, font = "", size = "", weight = "")
		{
		super(text, underline:, :color, :tip, :tabstop, :bgndcolor, :font, :size, :weight)
		.El.AddEventListener('click', .click)
		.El.AddEventListener('keydown', .keydown)
		.El.SetStyle('cursor', 'pointer')
		}

	click(event/*unused*/)
		{
		.Event('CLICK')
		}

	keydown(event)
		{
		if event.key is " "
			.Event('CLICK')
		}
	}
