// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
Component
	{
	picker: false
	styles: `
	.su-color-rect {
		align-self: 	center;
		position: 		relative;
	}
	.su-color-rect-input {
		position: 		absolute;
		bottom: 		0px;
		right:			0px;
		visibility: 	hidden;
	}
	`
	New(color, choose? = false)
		{
		LoadCssStyles('su-color-rect.css', .styles)
		.CreateElement('div', className: 'su-color-rect')
		.Set(color)
		.SetMinSize()
		.El.AddEventListener(#click, .click)
		.El.tabIndex = "1"

		if choose? is true
			{
			.picker = CreateElement('input', .El, className: 'su-color-rect-input')
			.picker.type = 'color'
			.picker.value = color
			.picker.AddEventListener(#input, .colorInput)
			.picker.AddEventListener(#change, .colorChange)
			.El.SetStyle('cursor', 'zoom-in')
			SuDelayed(0, .choose)
			}
		}

	Set(color)
		{
		.SetStyles(Object('background-color': color))
		}

	choose()
		{
		if .picker is false
			return
		.picker.Click()
		}

	colorInput(event)
		{
		.El.SetStyle('background-color', event.target.value)
		}

	colorChange(event)
		{
		color = ToCssColor.Reverse(event.target.value)
		.Event(#UpdateColor, color)
		}

	click()
		{
		.Event("LBUTTONDBLCLK")
		.choose()
		}
	}
