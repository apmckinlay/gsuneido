// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
FieldComponent
	{
	Name: 'Link'
	New(@args)
		{
		super(@args)
		.El.AddEventListener('dblclick', .doubleClicked)
		}

	doubleClicked(event)
		{
		if event.button isnt 0
			return
		.RunWhenNotFrozen({ .EventWithFreeze('LBUTTONDBLCLK') })
		}

	SetTextColor(color)
		{
		if color is "" or color is false
			color = CLR.BLACK
		.SetStyles(Object('color': ToCssColor(color)))
		}
	}
