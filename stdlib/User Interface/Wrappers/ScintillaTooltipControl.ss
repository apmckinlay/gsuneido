// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
InfoWindowControl
	{
	ComponentName: 'ScintillaTooltip'
	New(@args)
		{
		super(@args)
		.ComponentArgs.editorId = args.editorId
		.ComponentArgs.pos = args.pos
		}
	}