// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Container
	{
	Name:		"GroupBox"
	Xstretch: 0
	Ystretch: 0
	New(text, control)
		{
		.ctrl = .Construct(control)
		.box = .Construct(Object(GroupBox, text))
		.Recalc()
		}
	horzPadding: 4
	topPadding: 2
	bottomPadding: 4
	Resize(x, y, w, h)
		{
		.box.Resize(x, y, w, h)
		.ctrl.Resize(x + .horzPadding, y + .box.TextHeight + .topPadding,
			w - 2 * .horzPadding, h - .totalVertPadding())
		}
	totalVertPadding()
		{
		return .box.TextHeight + .topPadding + .bottomPadding
		}
	GetChildren()
		{
		return Object(.ctrl, .box)
		}
	SetCaption(caption)
		{
		SetWindowText(.GetChildren()[1].Hwnd, caption)
		}
	GetCaption()
		{
		return GetWindowText(.GetChildren()[1].Hwnd)
		}
	boxMinPadding: 14
	Recalc()
		{
		.Xmin = Max(.box.Xmin + .boxMinPadding, .ctrl.Xmin + .horzPadding * 2)
		.Ymin = .ctrl.Ymin + .totalVertPadding()
		.Left = .ctrl.Left + .horzPadding
		}
	Update()
		{
		.ctrl.Update()
		.box.Update()
		super.Update()
		}
	}
