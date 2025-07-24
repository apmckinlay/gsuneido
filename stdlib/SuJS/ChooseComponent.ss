// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	New(field, button, buttonBefore, hidden)
		{
		.CreateElement('div')
		.SetStyles(#('display': 'inline-flex', 'flex-direction': 'row'))
		.SetHidden(hidden)

		if buttonBefore
			{
			.button = .Construct(button)
			.button.El.SetStyle('border-right', '0px')
			.field = .Construct(field)
			.field.El.SetStyle('border-left', '0px')
			}
		else
			{
			.field = .Construct(field)
			.field.El.SetStyle('border-right', '0px')
			.button = .Construct(button)
			.button.El.SetStyle('border-left', '0px')
			}

		.field.SetStyles(#(flex: '1 1 auto'))
		.children = Object(.field, .button)

		.field.El.AddEventListener('keydown', .keydown)
		.Recalc()
		}

	Recalc()
		{
		.button.Recalc()
		.Xmin = .field.Xmin + .button.Xmin
		.Ymin = .field.Ymin
		.SetMinSize()
		}

	Initialize(@unused) { }

	keydown(event)
		{
		if .field.GetReadOnly()
			return

		if event.key in (#ArrowUp, #ArrowDown) or
			event.key is 'z' and event.altKey is true
			.button.CLICK()
		else if event.key is #Enter
			{
			.field.Event('KEYDOWN', VK.RETURN)
			return
			}
		else if event.key is #Escape
			.field.EventWithOverlay('KEYDOWN', VK.ESCAPE)
		else
			return

		event.PreventDefault()
		event.StopPropagation()
		}

	children: false
	GetChildren()
		{
		return .children isnt false ? .children : Object()
		}

	SetFocus()
		{
		.field.SetFocus()
		}

	Resize(w, h)
		{
		super.Resize(w, h)
		.field.SetStyles(#('min-width': '', 'width': '0'))
		.field.El.RemoveAttribute(#size)
		}
	}
