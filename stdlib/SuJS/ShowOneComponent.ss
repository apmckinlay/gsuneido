// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
HtmlContainer
	{
	New(@ctrls)
		{
		super()
		.CreateElement('div')
		.SetStyles(#('display': 'flex', 'flex-direction': 'column'))

		.ctrls = Object()
		for c in ctrls.Values(list:)
			.ctrls.Add(.Construct(c))
		.Recalc()
		}

	show: false
	UpdateShow(.show)
		{
		.Recalc()
		}

	Recalc()
		{
		if not .ctrls.Member?(.show)
			return
		ctrl = .ctrls[.show]
		.Xmin = ctrl.Xmin
		.Ymin = ctrl.Ymin
		.Xstretch = ctrl.Xstretch
		.Ystretch = ctrl.Ystretch
		.SetMinSize()
		}

	GetChildren()
		{
		return .ctrls
		}
	}
