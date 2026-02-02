// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
WndPaneControl
	{
	Name: 'ExpandPane'
	ComponentName: 'ExpandPane'
	New(open, control)
		{
		super(control, "SuBtnfaceArrow")
		.ComponentArgs.Add(open)
		}
	SetOpen(open)
		{
		.Act(#SetOpen, open)
		}
	}