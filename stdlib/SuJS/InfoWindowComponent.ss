// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Component
	{
	marginSize:	15
	titleSize:	20

	New(text = "", title = "", x = false, y = false,
		width/*unused*/ = 300, height/*unused*/ = 300,
		marginSize = 15, titleSize = 20)
		{
		.CreateElement('div')
		.SetStyles([
			'padding': marginSize $ 'px',
			'background-color': 'lightgoldenrodyellow',
			'position': 'fixed'])
		if title isnt ""
			{
			.titleEl = CreateElement('div', .El)
			.titleEl.innerText = title
			.titleEl.SetStyle('height', titleSize $ 'px')
			}
		.textEl = CreateElement('div', .El)
		.textEl.innerText = text

		PlaceElement(.El,
			x is false ? SuRender().GetCursorPos().x : x,
			y is false ? SuRender().GetCursorPos().y : y, 0, [])
		SuRender().RegisterKeydownListener(.close)
		}

	Startup()
		{
		if .Window.Method?(#SetBackdropDismiss)
			.Window.SetBackdropDismiss()
		}

	closed: false
	close(event/*unused*/)
		{
		if .closed is true
			return
		.closed = true
		.Event('Close')
		}

	Destroy()
		{
		SuRender().UnregisterKeydownListener(.close)
		super.Destroy()
		}
	}
