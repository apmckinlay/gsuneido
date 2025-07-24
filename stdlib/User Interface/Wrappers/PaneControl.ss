// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// NOTE: contents is actually sibling (i.e. same hwnd parent)
Hwnd
	{
	Name:		Pane
	Xmin:		false
	Ymin:		false
	Xstretch:	false
	Ystretch:	false

	New(control)
		{
		.CreateWindow("SuWhiteArrow", "", WS.VISIBLE, WS_EX.CLIENTEDGE, w: 1, h: 1)
		EnableWindow(.Hwnd, false)
		.ctrl = .Construct(control)
		.initializePaneProperties()
		}
	initializePaneProperties()
		{
		if .Xmin is false
			.Xmin = .ctrl.Xmin + 2 * GetSystemMetrics(SM.CXEDGE)
		yOffset = 1 + GetSystemMetrics(SM.CYEDGE)
		.Top = .ctrl.Top
		if .Ymin is false
			{
			.Ymin = .ctrl.Ymin
			if yOffset isnt 0
				{
				.Ymin += 2 * yOffset
				.Top += yOffset
				}
			}
		borderOffset = 4
		if .Xmin is false
			.Xmin = .ctrl.Xmin + borderOffset
		if .Xstretch is false
			.Xstretch = .ctrl.Xstretch
		if .Ymin is false
			.Ymin = .ctrl.Ymin + borderOffset
		if .Ystretch is false
			.Ystretch = .ctrl.Ystretch
		}
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		clientRect = .GetClientRect()
		clientW = clientRect.GetWidth()
		clientH = clientRect.GetHeight()
		.ctrl.Resize(x + ((w - clientW) / 2).Int(), y + ((h - clientH) / 2).Int(),
			clientW, clientH)
		}
	SetEnabled(enabled)
		{
		Assert(Boolean?(enabled))
		.ctrl.SetEnabled(enabled)
		}
	GetEnabled()
		{
		return .ctrl.GetEnabled()
		}
	SetReadOnly(readOnly)
		{
		Assert(Boolean?(readOnly))
		.ctrl.SetReadOnly(readOnly)
		}
	GetReadOnly()
		{
		return .ctrl.GetReadOnly()
		}
	SetVisible(visible)
		{
		Assert(Boolean?(visible))
		.ctrl.SetVisible(visible)
		super.SetVisible(visible)
		}
	SetFocus()
		{
		.ctrl.SetFocus()
		}
	Update()
		{
		.ctrl.Update()
		super.Update()
		}
	GetChildren()
		{
		return Object(.ctrl)
		}

	Destroy()
		{
		.ctrl.Destroy()
		super.Destroy()
		}
	}
