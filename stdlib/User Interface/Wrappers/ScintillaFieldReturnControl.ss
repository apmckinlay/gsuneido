// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonsFieldControl
	{
	ComponentName: 'ScintillaFieldReturn'
	Width: 20
	New(@args)
		{
		super(@.disableWrap(args))
		.SendMessageTextIn(SCI.SETHSCROLLBAR, false)
		.SendMessageTextIn(SCI.SETVSCROLLBAR, false)
		if false isnt status = args.GetDefault('status', false)
			.ToolTip(status)
		}

	disableWrap(args)
		{
		args.wrap = false
		return args
		}

	GETDLGCODE(wParam, lParam)
		{
		if wParam is VK.RETURN
			.FieldReturn()
		return super.GETDLGCODE(lParam)
		}

	FieldReturn()
		{
		.Send('FieldReturn')
		}
	}