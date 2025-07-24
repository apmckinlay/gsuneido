// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: "UpDown"

	styles: `
		.su-updown-button {
			margin: 0px;
			padding: 1px 6px;
			flex-grow: 1;
			line-height: 1px;
			border-width: 1px;
			border-style: solid;
			background-color: lightgrey;
		}
		.su-updown-button:enabled:hover {
			background-color: azure;
			border-color: deepskyblue;
		}
		.su-updown-button:enabled:active {
			background-color: lightblue;
			border-color: deepskyblue;
		}
		.su-updown-down {
			width: 0px;
			height: 0px;
			display: inline-block;
			vertical-align: middle;
			border-top: 4px dashed;
			border-left: 4px solid transparent;
			border-right: 4px solid transparent;
		}
		.su-updown-up {
			width: 0px;
			height: 0px;
			display: inline-block;
			vertical-align: middle;
			border-bottom: 4px dashed;
			border-left: 4px solid transparent;
			border-right: 4px solid transparent;
		}
		`

	New()
		{
		LoadCssStyles('updown-button.css', .styles)
		.CreateElement('div')
		.SetStyles(['display': 'inline-flex',
			'flex-direction': 'column',
			'align-content': 'stretch',
			'align-self': 'stretch'])
		.up = CreateElement('button', .El, className: 'su-updown-button')
		CreateElement('span', .up, className: 'su-updown-up')
		.down = CreateElement('button', .El, className: 'su-updown-button')
		CreateElement('span', .down, className: 'su-updown-down')

		.up.tabIndex = .down.tabIndex = "-1"
		.up.AddEventListener('mousedown', .onUp)
		.down.AddEventListener('mousedown', .onDown)
		}

	onUp(event)
		{
		.Event('UP')
		event.PreventDefault()
		}

	onDown(event)
		{
		.Event('DOWN')
		event.PreventDefault()
		}

	enabled: true
	SetEnabled(enabled)
		{
		if .enabled is enabled
			return
		.enabled = enabled
		.up.disabled = not enabled
		.down.disabled = not enabled
		}
	GetEnabled()
		{
		return .enabled
		}
	}
