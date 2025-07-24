// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
HtmlContainer
	{
	Name: "GroupBox"
	Xstretch: 0
	Ystretch: 0
	styles: `
	.su-group-box {
		display: flex;
		flex-direction: column;
	}
	.su-group-box-container {
		display: flex;
		flex-direction: column;
		flex-grow: 1;
		outline: 1px dotted black;
		padding: 2px 4px 4px 4px;
	}
	.su-group-box-caption {
	}
	`
	New(text, control)
		{
		LoadCssStyles('su-group-box.css', .styles)
		.CreateElement('div', className: 'su-group-box')
		.textEl = CreateElement('div', .El, className: 'su-group-box-caption')
		.SetCaption(text)
		.TargetEl = .box = CreateElement('div', .El, className: 'su-group-box-container')
		.ctrl = .Construct(control)
		if .ctrl.Xstretch > 0
			.ctrl.SetStyles(#('align-self': 'stretch', 'width': ''))
		if .ctrl.Ystretch > 0
			.ctrl.SetStyles(#('flex-grow': '1'))
		.Recalc()
		}

	SetCaption(caption)
		{
		.text = caption
		.textEl.innerText = caption
		}

	GetChildren()
		{
		return Object(.ctrl)
		}

	boxMinPadding: 14
	horzPadding: 4
	topPadding: 2
	bottomPadding: 4
	Recalc()
		{
		metrics = SuRender().GetTextMetrics(.textEl, .text)
		.Xmin = Max(metrics.width + .boxMinPadding, .ctrl.Xmin + .horzPadding * 2)
		.Ymin = .ctrl.Ymin + metrics.height + .topPadding + .bottomPadding
		.Left = .ctrl.Left + .horzPadding
		}
	}
