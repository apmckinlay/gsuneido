// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Container
	{
	Name: "Border"
	New(control, border = 10)
		{
		if Number?(control) and Object?(border)
			{ tmp = control; control = border; border = tmp }
		.Ctrl = .Construct(control)
		.border = ScaleWithDpiFactor(border)
		.xmin_original = .Xmin
		.ymin_original = .Ymin
		.xstretch_original = .Xstretch
		.ystretch_original = .Ystretch
		.Recalc()
		}
	SetCtrl(control)
		{
		.Ctrl.Destroy()
		.Ctrl = .Construct(control)
		.Window.Refresh()
		}
	GetChild()
		{ return .Ctrl }
	GetChildren()
		{ return Object(.Ctrl) }
	Ctrl: false
	Recalc()
		{
		if .Ctrl is false
			return
		.Xmin = Max(.xmin_original, .Ctrl.Xmin + 2 * .border)
		.Ymin = Max(.ymin_original, .Ctrl.Ymin + 2 * .border)
		if .xstretch_original is false
			.Xstretch = .Ctrl.Xstretch
		if .ystretch_original is false
			.Ystretch = .Ctrl.Ystretch
		}
	Resize(x, y, w, h)
		{
		.Ctrl.Resize(x + .border, y + .border, w - 2 * .border, h - 2 * .border)
		}
	}
