// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	SetScrollYmin(rows)
		{
		filters = .FindControl('conditions')
		scroll = .FindControl('Scroll')
		scroll.Ymin = ((filters.Ymin + 2) * rows).Round(0) + 4 /*= border */
		.BottomUp(#Recalc)
		}
	}
