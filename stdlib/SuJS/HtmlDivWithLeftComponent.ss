// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	New(@args)
		{
		super(@args)
		.Recalc()
		}

	Recalc()
		{
		super.Recalc()
		if false is child = .GetChild()
			return
		.Left = child.Left
		}
	}
