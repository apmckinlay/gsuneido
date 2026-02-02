// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
PassthruController
	{
	Name: "VirtualListScroll"
	ComponentName: "VirtualListScroll"
	hdrCornerCtrl: false
	New(control, bar, floating, thinBorder/*unused*/, hdrCornerCtrl = false,
		expandExtra = false)
		{
		super(control)
		.ctrl = .GetChild()
		.bar = .Construct(bar)
		.floating = .Construct(floating)
		.ComponentArgs.bar = .bar.GetLayout()
		if hdrCornerCtrl isnt false
			{
			.hdrCornerCtrl = .Construct(hdrCornerCtrl)
			.ComponentArgs.hdrCornerCtrl = .hdrCornerCtrl.GetLayout()
			}

		if expandExtra isnt false and expandExtra isnt ''
			.expandExtra = .Construct(expandExtra)
		}

	GetChildren()
		{
		children = Object(.ctrl, .bar, .floating)
		if .hdrCornerCtrl isnt false
			children.Add(.hdrCornerCtrl)
		if .expandExtra isnt false
			children.Add(.expandExtra)
		return children
		}

	GetHdrCornerCtrl()
		{
		return .hdrCornerCtrl
		}

	Default(@unused) { }

	Destroy()
		{
		.ctrl.Destroy()
		.floating.Destroy()
		.bar.Destroy()
		if .hdrCornerCtrl isnt false
			.hdrCornerCtrl.Destroy()
		if .expandExtra isnt false
			.expandExtra.Destroy()
		super.Destroy()
		}
	}
