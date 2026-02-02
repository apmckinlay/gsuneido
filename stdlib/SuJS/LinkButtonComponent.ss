// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
StaticComponent
	{
	ContextMenu: true
	New(text, tip = false, .color = 'BLUE', bgndcolor = "",
		tabstop = true, font = "", size = "", weight = "")
		{
		super(text, underline:, :color, :tip, :tabstop, :bgndcolor, :font, :size, :weight)
		.El.AddEventListener('click', .click)
		.El.AddEventListener('keydown', .keydown)
		.El.SetStyle('cursor', 'pointer')
		}

	click(event/*unused*/)
		{
		.RunWhenNotFrozen({ .EventWithOverlay('CLICK') })
		}

	Recalc()
		{
		super.Recalc()
		.Ymin += 8 /*=padding*/
		.SetMinSize()
		.SetStyles(Object('line-height': (.Ymin/*=padding*/) $ 'px'))
		}

	keydown(event)
		{
		if event.key is " "
			.RunWhenNotFrozen({ .EventWithOverlay('CLICK') })
		}

	SetEnabled(.enabled)
		{
		.El.SetStyle('pointer-events', enabled is false ? 'none' : '')
		.El.SetStyle('color',
			ToCssColor(enabled is false ? 'darkgray' : .color))
		}

	GetEnabled()
		{
		return .enabled
		}
	}
