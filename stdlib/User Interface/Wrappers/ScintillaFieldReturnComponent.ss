// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonsFieldComponent
	{
	New(@args)
		{
		super(@args)
		.CM.SetOption("scrollbarStyle", "null")
		}

	HandleEnterKey(cm, event)
		{
		super.HandleEnterKey(cm, event)
		if event.key is 'Enter'
			.Event('FieldReturn')
		}
	}