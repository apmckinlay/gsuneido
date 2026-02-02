// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
TextPlus
	{
	Name: 'CheckBox'
	ComponentName: 'CheckBox'
	CustomizableOptions: #(hidden, readonly, tabover)
	New(@args)
		{
		super(@args)
		.Send('Data')
		if "" isnt set = args.GetDefault(#set, "")
			{
			.Set(set)
			.Send('NewValue', .Get())
			}
		}

	Toggle()
		{
		super.Toggle()
		if not .GetReadOnly()
			.Send('NewValue', .Get())
		}

	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}