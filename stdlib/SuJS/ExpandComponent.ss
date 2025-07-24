// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	Name: "Expand"
	Ystretch: 0
	MaxHeight: 8888 // must be different from Control.MaxHeight
	CalcMaxHeight()
		{ return 8888 }
	Recalc()
		{
		pane = .FindControl('ExpandPane')
		control = pane.Vert.GetChildren()[1]
		.Ystretch = pane.GetVisible() and control isnt false ? control.Ystretch : 0
		super.Recalc()
		}
	}
