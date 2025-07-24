// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: 'Fill'
	New(min = 0, fill = 1)
		{
		.CreateElement('div')
		if (not .Parent.Member?("Dir"))
			.Xstretch = .Ystretch = fill
		else if (.Parent.Dir is "vert")
			.Ystretch = fill
		else if (.Parent.Dir is "horz")
			.Xstretch = fill
		if (min > 0)
			{
			if (not .Parent.Member?("Dir"))
				.Xmin = .Ymin = min
			else if (.Parent.Dir is "vert")
				.Ymin = min
			else if (.Parent.Dir is "horz")
				.Xmin = min
			}
		.SetMinSize()
		}
	GetReadOnly() // read-only not applicable to fill
		{
		return true
		}
	}
