// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New()
		{
		.tab = .Vert.Tab
		}
	Controls: (Vert
		Tab
		(Horz
			(Button Insert)
			(Button Remove)
			)
		)
	i: 0
	On_Insert()
		{
		.tab.Insert(0, "tab " $ .i,
			Object(tooltip: "this is tab " $ .i))
		++.i
		}
	On_Remove()
		{
		.tab.Remove(0)
		}
	}