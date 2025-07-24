// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
FieldComponent
	{
	New(@args)
		{
		super(@args)
		.El.AddEventListener('keydown', .keydown)
		}

	keydown(event)
		{
		if event.key isnt 'Enter'
			return
		event.PreventDefault()
		event.StopPropagation()
		.Event('FieldReturn')
		}
	}
