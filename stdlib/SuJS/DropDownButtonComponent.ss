// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
Component
	{
	styles: `
		.su-dropdown-button {
			font-family: suneido;
			font-style: normal;
			font-weight: normal;
			background-color: var(--su-color-buttonface);
			border-width: 1px;
			border-color: #767676;
			border-style: solid;
			color: #767676;
			padding: 0px;
			margin: 0px;
			user-select: none;
		}
		.su-dropdown-button:enabled {
			background-color: white;
		}
		.su-dropdown-button:enabled:hover {
			background-color: azure;
			outline: deepskyblue solid 1px;
			outline-offset: -3px;
			color: black;
		}
		.su-dropdown-button:enabled:active {
			background-color: lightblue;
			outline: deepskyblue solid 1px;
			outline-offset: -3px;
			color: black;
		}
		.su-edit:focus + .su-dropdown-button {
			border-color: deepskyblue;
		}`
	New(hidden = false)
		{
		LoadCssStyles('su-dropdown-button.css', .styles)
		.CreateElement('button', '&nbsp;'/*arrow_down in suneido font*/,
			className: 'su-dropdown-button')
		.SetHidden(hidden)
		.El.tabIndex = "-1"
		.El.AddEventListener('keydown', .keydown)
		.El.AddEventListener('mousedown', .mousedown)
		.Recalc()
		}

	Recalc()
		{
		el = .El
		if false isnt field = .Parent.GetChildren().GetDefault(0, false)
			el = field.El
		metrics = SuRender().GetTextMetrics(el, 'M')
		.Xmin = metrics.width
		.Ymin = metrics.height
		pad = .Ymin + .Ymin % 2 // make it even
		.Xmin += pad
		.Ymin += 10 /* = y padding + border*/
		.SetMinSize()
		}

	CLICK()
		{
		field = .Parent.GetChildren()[0]
		.RunWhenNotFrozen()
			{
			r = SuRender.GetClientRect(field.El)
			r.width += .El.clientWidth
			rcExclude = Object(left: -9999, right: 9999, top: r.top, bottom: r.bottom)
			if .Window.Base?(ListEditWindowComponent)
				// Use overlay to prevent users from closing the ListEditWindow after CLICK
				.EventWithOverlay('CLICK', x: r.left, y: r.bottom, :rcExclude, rect: r)
			else
				// Use freeze to avoid the overlay's focus change that triggers AfterField
				.EventWithFreeze('CLICK', x: r.left, y: r.bottom, :rcExclude, rect: r)
			}
		}

	keydown(event)
		{
		if event.key in (#ArrowUp, #ArrowDown)
			.CLICK()
		}

	mousedown(event)
		{
		.CLICK()
		if .Window.Base?(ListEditWindowComponent)
			event.PreventDefault()
		}

	SetReadOnly(readonly)
		{
		.SetEnabled(not readonly)
		}
	}
