// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Xstretch: 	1
	Ystretch: 	1
	Xmin: 		50
	Ymin: 		50
	Name: 		"Scroll"
	ComponentName: "Scroll"

	New(control/*, style = 0, .vmanual = false, wndclass = "SuBtnfaceArrow",
		.vdisable = false, .dyscroll = 21*/, trim = false/*, .noEdge = false*/)
		{
		.ctrl = .Construct(control)
		.ComponentArgs = Object(.ctrl.GetLayout(), trim)
		}

	GetChild()
		{
		return .ctrl
		}

	GetChildren()
		{
		return Object(.ctrl)
		}

	SetEnabled(enabled)
		{
		Assert(Boolean?(enabled))
		.ctrl.SetEnabled(enabled)
		super.SetEnabled(enabled)
		.Act('SetEnabled', enabled)
		}

	SetVisible(visible)
		{
		Assert(Boolean?(visible))
		.ctrl.SetVisible(visible)
		super.SetVisible(visible)
		.Act('SetVisible', visible)
		}

	SetReadOnly(readonly)
		{
		.ctrl.SetReadOnly(readonly)
		.Act('SetReadOnly', readonly)
		}

	VSCROLL(wParam)
		{
		switch (param = LOWORD(wParam))
			{
		case SB.TOP, SB.BOTTOM:
			.Act('VSCROLL', param)
		default:
			}
		return 0
		}

	Destroy()
		{
		.ctrl.Destroy()
		super.Destroy()
		}
	}