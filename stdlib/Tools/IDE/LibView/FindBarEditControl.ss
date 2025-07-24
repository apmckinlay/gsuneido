// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Xmin: 200
	New()
		{
		.Top = .Horz.Top
		}
	Controls: #(Horz
		(FieldHistory, font: '@mono', size: '+2', width: 10, xstretch: .001,
			trim: false, name: "find")
		Skip
		(Static name: "occurrence"))
	}