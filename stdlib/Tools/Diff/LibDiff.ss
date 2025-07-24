// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
// uses DiffControl
Controller
	{
	Title: 'LibDiff'
	CallClass(lib1, name1, lib2, name2 = false)
		{
		Window([this, lib1, name1, lib2, name2])
		}
	New(@args)
		{
		super(.control(@args))
		}
	control(.lib1, .name1, .lib2, .name2 = false)
		{
		if name2 is false
			.name2 = name1
		return ['Vert', .diff(), name: 'Content']
		}
	On_Refresh()
		{
		.Content.Remove(0)
		.Content.Append(.diff())
		}
	diff()
		{
		['Diff2',
			Query1(.lib1, group: -1, name: .name1).lib_current_text,
			Query1(.lib2, group: -1, name: .name2).lib_current_text,
			.lib1, .name1,
			.lib1 $ ' ' $ .name1, .lib2 $ ' ' $ .name2,
			refresh:]
		}
	}
