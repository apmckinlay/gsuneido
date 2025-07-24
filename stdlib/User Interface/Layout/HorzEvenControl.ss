// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
HorzControl
	{
	Name: HorzEven
	Xstretch: 1
	Resize(x, y, w, h)
		{
		children = .GetChildren()
		n = 0
		for child in children
			if not child.Base?(SkipControl)
				++n
		cw = (w / n).Int()
		for child in children
			if not child.Base?(SkipControl)
				child.Xstretch = Max(0, cw - child.Xmin)
		super.Resize(x, y, w, h)
		}
	}