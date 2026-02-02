// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
HorzComponent
	{
	Name: 'Flow'
	Xstretch: 1
	Shrinkable: true
	New(@args)
		{
		super(@args)
		.SetStyles(Object('flex-wrap': 'wrap'))
		}
	Recalc()
		{
		super.Recalc()
		.Xmin = .GetChildren().Map({ it.Xmin }).Max()
		.Ymin = .GetChildren().Map({ it.Ymin }).Max()
		.SetMinSize()
		}
	}