// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
FormComponent
	{
	// if a control has stretch, it will take the whole grid cell
	Stretch?(i/*unused*/, c)
		{
		return c.Xstretch isnt false
		}

	// don't automatically span over the following empty group cells
	CalcGridColumn(i, gns, offset/*unused*/)
		{
		start = gns[i] + 1
		next = i + 1 >= gns.Size()
			? gns[i] + 2
			: gns[i + 1] + 1
		return start $ ' / ' $ next
		}
	}