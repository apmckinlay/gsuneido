// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
HorzControl
	{
	Name: 'HorzEqual'
	New(@args)
		{
		super(@args)
		maxXmin = 0
		for c in .GetChildren()
			if c.Base?(ButtonControl)
				maxXmin = Max(maxXmin, c.Xmin)
		maxXmin += args.GetDefault('pad', 20)
		.Xmin = 0
		for c in .GetChildren()
			{
			if c.Base?(ButtonControl)
				c.Xmin = maxXmin
			.Xmin += c.Xmin
			}
		}
	}