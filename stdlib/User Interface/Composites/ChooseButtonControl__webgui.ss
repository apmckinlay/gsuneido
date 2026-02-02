// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
MenuButtonControl
	{
	Name: 			'ChooseButton'
	ComponentName:	'ChooseButton'
	New(.text, .list, width = false)
		{
		super(text, list, left:, :width)
		.ChooseButton = this
		.Send(#Data)

		.ComponentArgs = Object(text, list, left:, :width)
		}

	Set(value)
		{
		super.Set(value is "" ? .text : value)
		}

	Send(@args)
		{
		if args[0] isnt .Command
			return super.Send(@args)

		.On_ChooseButton(args[1], args.index)
		}

	On_ChooseButton(value, index)
		{
		.Set(value)
		.Send(#NewValue, value)
		super.Send("On_" $ .Name, value, index)
		return 0
		}

	SetList(.list)
		{
		.SetMenu(list)
		.Act('SetList', list)
		}

	GetReadOnly()
		{
		return .Disabled?()
		}
	SetReadOnly(readonly)
		{
		.Disable(readonly)
		.Grayed?(readonly or .Get() is '')
		}

	Destroy()
		{
		.Send(#NoData)
		super.Destroy()
		}
	}
