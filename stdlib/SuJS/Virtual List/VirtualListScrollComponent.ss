// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	Name: 		"VirtualListScroll"
	Xstretch: 	1
	Ystretch: 	1
	Xmin: 		50
	Ymin: 		10

	hdrCornerCtrl: false
	New(control, bar, hdrCornerCtrl = false)
		{
		super(control)
		.El.SetStyle('position', 'relative')
		.bar = .Construct(bar)
		if hdrCornerCtrl isnt false
			{
			buttonHeight = SuRender().GetTextMetrics(.El, 'M').height
			hdrCornerCtrl = hdrCornerCtrl.Copy()
			hdrCornerCtrl.buttonHeight = buttonHeight
			.hdrCornerCtrl = .Construct(hdrCornerCtrl)
			.hdrCornerCtrl.SetStyles([
				position: 'absolute',
				left: '0px',
				top: '0px',
				overflow: 'hidden',
				'box-sizing': 'content-box',
				'padding': '2px 1px 2px 1px',
				'z-index': '3'])
			}
		}

	Destroy()
		{
		.bar.Destroy()
		if .hdrCornerCtrl isnt false
			.hdrCornerCtrl.Destroy()
		super.Destroy()
		}
	}
