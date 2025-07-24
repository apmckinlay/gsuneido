// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
WndPaneComponent
	{
	Name: 'ExpandPane'
	New(control, bgColor, open)
		{
		super(control, bgColor)
		.SetOpen(open)
		}
	Recalc()
		{
		open = .open?() // have to get this before recalc
		super.Recalc()
		if open
			.y0min = .Ymin
		else
			.Ymin = 0
		.SetMinSize()
		}

	SetOpen(open)
		{
		if open is .open?()
			return
		if open is false
			{
			.y0min = .Ymin
			.Ymin = 0
			}
		else
			{
			.Ymin = .y0min
			}
		super.SetVisible(open)
		}

	open?()
		{
		return .Ymin isnt 0
		}
	}
