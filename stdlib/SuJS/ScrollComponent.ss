// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Xstretch: 	1
	Ystretch: 	1
	Xmin: 		50
	Ymin: 		50
	Name: 		"Scroll"

	styles: `
		.su-scroll-container {
			position: relative;
			overflow: auto;
		}
		.su-scroll-container:focus {
			outline: none;
		}
		.su-scroll {
			position: absolute;
			display: inline-flex;
			flex-direction: column;
			top: 0px;
			left: 0px;
			width: 100%;
			height: 100%;
		}
		`
	New(control, .trim)
		{
		LoadCssStyles('su-scroll.css', .styles)
		.CreateElement('div', className: 'su-scroll-container')
		.TargetEl = CreateElement('div', .El, 'su-scroll')
		.ctrl = .Construct(control)
		.Recalc()
		}

	handleTrim()
		{
		if .trim is false
			return

		w = .ctrl.Xmin + 20/*=estimated scroll bar width*/
		if w < .Xmin
			.Xmin = w
		h = .ctrl.Ymin + 20/*=estimated scroll bar height*/
		if h < .Ymin
			.Ymin = h
		if .trim is 'auto' and .ctrl.Ystretch > 0
			.MaxHeight = 99999
		else
			.MaxHeight = .ctrl.Ymin + 20/*=estimated scroll bar height*/
		.SetMinSize()
		}

	Recalc()
		{
		.handleTrim()
		if .ctrl.Xstretch > 0
			.ctrl.SetStyles(#('align-self': 'stretch', 'width': ''))
		if .ctrl.Ystretch > 0
			.ctrl.SetStyles(#('flex-grow': '1'))
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
		}

	SetVisible(visible)
		{
		Assert(Boolean?(visible))
		.ctrl.SetVisible(visible)
		super.SetVisible(visible)
		}

	SetReadOnly(readonly)
		{
		.ctrl.SetReadOnly(readonly)
		}

	VSCROLL(param)
		{
		if param is SB.TOP
			.El.scrollTop = 0
		else if param is SB.BOTTOM
			.El.scrollTop = .El.scrollHeight
		}

	Destroy()
		{
		.ctrl.Destroy()
		super.Destroy()
		}
	}
