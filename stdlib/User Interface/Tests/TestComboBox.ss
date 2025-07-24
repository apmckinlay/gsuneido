// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Controls: (Horz StateProv Skip (Button Get) Skip (Button Set) Skip (Button BadSet))
	On_Get()
		{
		Alert(.Horz.StateProv.Get(), title: 'Get', flags: MB.ICONINFORMATION)
		}
	On_Set()
		{
		.Horz.StateProv.Set('SK')
		}
	On_BadSet()
		{
		.Horz.StateProv.Set('XX')
		}
	NewValue(x)
		{
		Alert("NewValue: " $ x, title: 'New Value', flags: MB.ICONINFORMATION)
		}
	}
