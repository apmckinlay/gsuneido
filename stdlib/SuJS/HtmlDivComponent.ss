// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
HtmlContainer
	{
	ctrl: false
	Controls: false
	Name: 'Controller'

	New(@args)
		{
		.xmin0 = .Xmin
		.ymin0 = .Ymin
		ctrl = args.GetDefault(0, false)
		.Initialize(ctrl)
		}

	Initialize(ctrl)
		{
		if .ctrl isnt false
			.ctrl.Destroy()

		.CreateElement('div')
		.SetStyles(#('display': 'inline-flex', 'flex-direction': 'column'))

		if ctrl is false
			return

		.ctrl = .Construct(ctrl)
		.SetupCtrl(.ctrl)
		}

	SetupCtrl(ctrl)
		{
		if ctrl.Xstretch > 0
			ctrl.SetStyles(#('align-self': 'stretch', 'width': ''))
		if ctrl.Ystretch > 0
			ctrl.SetStyles(#('flex-grow': '1'))
		.recalcCtrl(ctrl)

		if .Xstretch is false
			.Xstretch = ctrl.Xstretch
		if .Ystretch is false
			.Ystretch = ctrl.Ystretch
		if .MaxHeight is Component.MaxHeight
			.MaxHeight = ctrl.MaxHeight
		}

	Recalc(ctrl = false)
		{
		if ctrl is false
			ctrl = .ctrl
		.recalcCtrl(ctrl)
		}

	recalcCtrl(ctrl)
		{
		if ctrl is false
			return
		.Xmin = Max(.xmin0, ctrl.Xmin)
		.Ymin = Max(.ymin0, ctrl.Ymin)
		.MaxHeight = ctrl.MaxHeight
		.SetMinSize()
		}

	GetChild()
		{
		return .ctrl
		}
	GetChildren()
		{
		return .ctrl is false ? #() : Object(.ctrl)
		}
	}
