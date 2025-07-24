// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
FieldComponent
	{
	New(@args)
		{
		super(@args)
		.El.AddEventListener('keydown', .KEYDOWN)
		}

	KEYDOWN(event)
		{
		if event.key is 'Enter'
			.Event('ENTER')
		}
	}
