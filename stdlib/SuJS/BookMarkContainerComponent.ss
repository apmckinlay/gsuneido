// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	Name: "BookMarkContainer"
	New(@args)
		{
		super(@args)
		.SetStyles(#(overflow: hidden))
		}
	// to avoid set Xmin, Ymin
	Recalc()
		{
		.GetChild().El.SetStyle('min-width', '')
		}
	}
