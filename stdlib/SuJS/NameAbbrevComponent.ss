// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	Recalc()
		{
		super.Recalc()
		firstChild = .Horz.GetChildren()[0]
		.Left = firstChild.Left
		}
	}
