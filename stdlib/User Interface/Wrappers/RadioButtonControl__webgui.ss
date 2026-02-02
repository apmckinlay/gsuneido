// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
TextPlus
	{
	Name: 'RadioButton'
	ComponentName: 'RadioButton'
	Toggle()
		{
		if .Get() isnt true
			.Send('Picked', .GetText())
		}
	}