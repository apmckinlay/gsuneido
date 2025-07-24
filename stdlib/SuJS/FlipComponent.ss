// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	Xmin:		0
	Ymin:		0
	Xstretch:	0
	Ystretch:	0
	Name:		"Flip"
	New(@ctrlList)
		{
		.controls = Object()
		size = ctrlList.Size(list:)
		for (i = size - 1; 0 <= i; --i)		// add backward to minimize flicker
			{
			.controls.Add(ctrl = .Construct(ctrlList[i]), at: 0)
			.Xmin = Max(ctrl.Xmin, .Xmin)
			.Ymin = Max(ctrl.Ymin, .Ymin)
			.Xstretch = Max(ctrl.Xstretch, .Xstretch)
			.Ystretch = Max(ctrl.Ystretch, .Ystretch)
			if ctrl.Xstretch > 0
				ctrl.SetStyles(#('align-self': 'stretch', 'width': ''))
			if ctrl.Ystretch > 0
				ctrl.SetStyles(#('flex-grow': '1'))
			}
		.SetMinSize()
		}

	GetChildren()
		{ return .controls.Copy() }
	}
