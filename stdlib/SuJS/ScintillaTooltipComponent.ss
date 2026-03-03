// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
InfoWindowComponent
	{
	New(@args)
		{
		super(@.buildArgs(args))
		}

	buildArgs(args)
		{
		editor = SuRender().GetRegisteredComponent(args.editorId)
		pos = editor.CM.CharCoords(editor.CM.PosFromIndex(args.pos))
		args.x = pos.left
		args.y = pos.bottom
		return args
		}
	}