// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
HtmlContainer
	{
	Name: "Border"
	Ctrl: false
	New(control, border, borderline)
		{
		.CreateElement('div')
		.border = border
		if borderline isnt 0
			.SetStyles(Object('border': borderline $ 'px solid lightgrey'))
		.SetStyles(Object(display: 'inline-flex', 'flex-direction': 'row',
			'box-sizing': 'border-box', padding: (.border - borderline) $ 'px'))

		.xmin_original = .Xmin
		.ymin_original = .Ymin
		.xstretch_original = .Xstretch
		.ystretch_original = .Ystretch
		.SetCtrl(control)
		}

	Recalc()
		{
		if .Ctrl is false
			return
		if .Ctrl.Xmin isnt 0
			.Xmin = Max(.xmin_original, .Ctrl.Xmin + 2 * .border)
		if .Ctrl.Ymin isnt 0
			.Ymin = Max(.ymin_original, .Ctrl.Ymin + 2 * .border)
		if .xstretch_original is false
			.Xstretch = .Ctrl.Xstretch
		if .ystretch_original is false
			.Ystretch = .Ctrl.Ystretch
		.SetMinSize()
		}

	SetCtrl(control)
		{
		if .Ctrl isnt false
			.Ctrl.Destroy()

		.Ctrl = .Construct(control)
		if .Ctrl.Xstretch >= 0
			.Ctrl.SetStyles(#('flex-grow': '1'))
		if .Ctrl.Ystretch > 0
			.Ctrl.SetStyles(#('align-self': 'stretch'))
		else
			.Ctrl.SetStyles(#('align-self': 'baseline'))
		.Recalc()
		}

	GetChild()
		{
		return .Ctrl
		}

	GetChildren()
		{
		return Object(.Ctrl)
		}
	}
