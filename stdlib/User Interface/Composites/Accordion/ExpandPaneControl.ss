// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
WndPaneControl
	{
	Name: ExpandPane
	New(open, control)
		{
		super(control, "SuBtnfaceArrow")
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
