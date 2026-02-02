// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Control
	{
	ComponentName: 'InfoWindow'

	CallClass(@args)
		{
		args.Add(this, at: 0)
		return Window(args, border: 0, style: WS.POPUP)
		}

	New(.text = "", .title = "", x = false, y = false, width = 300, height = 300,
		marginSize = 15, titleSize = 20, autoClose = false)
		{
		.ComponentArgs = Object(text, title, x, y, width, height, marginSize, titleSize)
		if Number?(autoClose)
			Delay(autoClose.SecondsInMs(), .Close)
		}

	Msg(args)
		{
		if args[0] is 'Inactivate' or args[0] is 'On_Cancel'
			.Close()
		return 0
		}

	GetReadOnly()
		{ return true }

	closing?: false
	Close()
		{
		if .closing? or .Destroyed?()
			return
		.closing? = true
		.Window.CLOSE()
		}
	}
