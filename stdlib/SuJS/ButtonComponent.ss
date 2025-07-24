// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Xmin: ""
	Ymin: ""
	Styles: `
		.su-button {
			position: relative;
			border-width: 1px;
			border-style: solid;
			border-radius: 0.3em;
			background-color: lightgrey;
			padding: 0;
			margin: 0;
			user-select: none;
		}
		.su-button:not([data-disable]):enabled:hover {
			background-color: azure;
			border-color: deepskyblue;
		}
		.su-button:not([data-disable]):enabled:active {
			background-color: lightblue;
			border-color: deepskyblue;
		}
		.su-button[data-highlight=true] {
			border-color: deepskyblue;
		}
		.su-button:enabled:focus {
			outline-offset: -3px;
			outline: 1px dashed black;
		}`
	ContextMenu: true
	New(text, font = "", size = "", weight = "", tabover = false, defaultButton = false,
		tip = false, .pad = false, color = false, .width = false,
		underline = false, italic = false, strikeout = false, hidden = false)
		{
		LoadCssStyles('su-button.css', .Styles)
		.text = .RemoveAmpersand(text)
		.CreateElement('button', .text, className: 'su-button')
		.xmin_orig = .Xmin
		.ymin_orig = .Ymin
		.SetHidden(hidden)
		.SetFont(font, size, weight, underline, italic, strikeout)
		.Recalc()
		if color isnt false
			.El.SetStyle(color: ToCssColor(color))
		if defaultButton is true
			.El.Focus()
		if tabover is true
			.El.tabIndex = "-1"

		.El.AddEventListener('click', .CLICKED)
		.El.AddEventListener('focus', .focus)
		.El.AddEventListener('blur', .blur)
		.El.AddEventListener('mousedown', .mousedown)

		.AddToolTip(tip)
		}

	Recalc()
		{
		metrics = SuRender().GetTextMetrics(.El, .text)
		.Xmin = metrics.width
		.Ymin = metrics.height
		if .pad is false
			{
			.pad = .Ymin + .Ymin % 2 // make it even
			}
		if Number?(.width)
			.xmin_orig = SuRender().GetTextMetrics(.El, 'M'.Repeat(.width)).width
		.Xmin += .pad
		.Ymin += 10 /* = y padding + border*/
		if .xmin_orig isnt "" and .xmin_orig > .Xmin
			.Xmin = .xmin_orig
		if .ymin_orig isnt "" and .ymin_orig > .Ymin
			.Ymin = .ymin_orig
		.SetMinSize()
		}

	Highlight(highlight?)
		{
		if .El is false
			return
		.El.dataset.highlight = highlight? is true
		}

	focus()
		{
		.Window.HighlightDefaultButton(false)
		}

	blur()
		{
		if not .Destroyed?()
			.Window.HighlightDefaultButton(true)
		}

	CLICKED()
		{
		.RunWhenNotFrozen({ .EventWithOverlay('CLICKED') })
		}

	mousedown(event)
		{
		if .Window.Base?(ListEditWindowComponent)
			event.PreventDefault()
		}

	Set(text)
		{
		.text = .RemoveAmpersand(text)
		.El.innerHTML = .text
		if .xmin_orig isnt ''
			return
		.WindowRefresh()
		}

	RemoveAmpersand(s)
		{
		return s.Tr('&')
		}
	}
