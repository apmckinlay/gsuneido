// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Container
	{
	Name: "Center"
	New(control, .border = 0)
		{
		.ctrl = .Construct(control)
		.Xstretch = .ctrl.Xstretch
		.Ystretch = .ctrl.Ystretch

		if not .Parent.Member?("Dir")
			.Xstretch = .Ystretch = 1
		else if .Parent.Dir is "vert"
			.Ystretch = 1
		else if .Parent.Dir is "horz"
			.Xstretch = 1

		.Xmin = .ctrl.Xmin + 2 * border
		.Ymin = .ctrl.Ymin + 2 * border
		}
	Resize(x, y, w, h)
		{
		// in case control's size has changed
		.Xmin = .ctrl.Xmin + 2 * .border
		.Ymin = .ctrl.Ymin + 2 * .border

		xs = Max(0, w - .Xmin)
		ys = Max(0, h - .Ymin)
		.ctrl.Resize(x + .border + xs / 2, y + .border + ys / 2,
			.ctrl.Xmin, .ctrl.Ymin)
		}
	GetChildren()
		{
		return Object(.ctrl)
		}
	}
