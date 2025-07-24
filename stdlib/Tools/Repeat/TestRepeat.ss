// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Ymin: 400
	New()
		{
		.repeat = .FindControl('Repeat')
		}
	Controls: #(Border (Vert
		(Repeat (Horz name Skip date))
		Skip
		(HorzEqual (Button Get) Skip (Button Set) Skip
			(Button Dirty?) Skip (Button Valid?))
		))
	On_Get()
		{
		Inspect(.repeat.Get())
		}
	On_Set()
		{
		data = []
		for i in .. 3
			data.Add([name: i, date: Date().NoTime().Plus(days: i)])
		.repeat.Set(data)
		}
	On_Dirty()
		{
		Alert(.repeat.Dirty?(), 'Dirty?')
		}
	On_Valid()
		{
		Alert(.repeat.Valid?(), 'Valid?')
		}
	}
