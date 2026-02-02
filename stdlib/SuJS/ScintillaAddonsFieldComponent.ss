// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
ScintillaAddonsComponent
	{
	Xmin: 100
	Xstretch: false
	DefaultFontSize: 9
	New(@args)
		{
		super(@args)
		.AddEventListenerToCM('keydown', .HandleEnterKey)
		}

	HandleEnterKey(unused, event)
		{
		if event.key is "Enter"
			event.PreventDefault()
		}
	}
