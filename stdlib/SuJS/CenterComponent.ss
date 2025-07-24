// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
HtmlContainer
	{
	Name: "Center"
	New(ctrl, .border = 0)
		{
		.CreateElement('div')
		.SetStyles(Object(
			'display': 'grid',
			'grid-template-columns': border $ 'px auto ' $ border $ 'px',
			'grid-template-rows': border $ 'px auto ' $ border $ 'px'))
		.ctrl = .Construct(ctrl)
		.ctrl.SetStyles(Object(
			'grid-column': '2 / 3',
			'grid-row': '2 / 3',
			'place-self': 'center'))

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
		.SetMinSize()
		}

	GetChildren()
		{
		return Object(.ctrl)
		}
	}
