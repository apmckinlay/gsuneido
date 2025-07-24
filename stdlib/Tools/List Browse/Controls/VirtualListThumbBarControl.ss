// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: 'VirtualListThumbBar'
	New(.disableSelectFilter = false)
		{
		super(.layout())

		.thumb = .Horz.Vert.VirtualListThumb
		}

	layout()
		{
		.sysWidth = GetSystemMetrics(SM.CXVSCROLL)
		return Object('Horz',
			Object(.colorLine)
			Object('Vert',
				.arrowButton('home'),
				.arrowButton('up', allowHolding:),
				Object('VirtualListThumb' width: .sysWidth - 1),
				.arrowButton('down', allowHolding:),
				.arrowButton('end'),
				not .disableSelectFilter ? .selectButton() : "",
				xmin: ScaleWithDpiFactor.Reverse(.sysWidth)))
		}

	colorLine: WndProc
		{
		Ystretch: 1
		New()
			{
			.CreateWindow("SuWhiteArrow", "", WS.VISIBLE)
			.Xmin = 1
			}
		}

	arrowButton(name, allowHolding = false)
		{
		ctrl = allowHolding ? 'VirtualListThumbImageButton' : 'EnhancedButton'
		// button theme border is included in the imagePadding
		// the border width is 2 pixel on regular dpi
		// imagePadding should not be smaller than 2 + 1 (offset)
		// otherwise the bottom border will be overlapped by image
		return Object(ctrl $ 'Control',
			command: 'VirtualListThumb_Arrow' $ name.Capitalize()
			image: 'arrow_' $ name $ '.emf', mouseEffect:, imagePadding: .15,
			buttonWidth: .sysWidth, buttonHeight: .sysWidth, name: name.Capitalize())
		}

	selectButton()
		{
		return Object('EnhancedButton', 'Select', command: 'VirtualListThumb_ArrowSelect'
			image: 'zoom.emf', imagePadding: .15, name: 'Select'
			enlargeOnHover: #(imagePadding: .2, direction: 'top-left'),
			buttonWidth: .sysWidth, buttonHeight: .sysWidth)
		}

	SetSelectPressed(pressed = false)
		{
		if .disableSelectFilter
			return
		select = .Horz.Vert.Select
		select.SetPressed?(pressed)
		select.Pushed?(pressed)
		select.Repaint()
		}

	GetWidth()
		{
		return .sysWidth
		}

	Default(@args)
		{
		event = args[0]
		.thumb[event](@+1 args)
		}
	}
