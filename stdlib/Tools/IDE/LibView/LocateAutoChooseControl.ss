// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
AutoChooseControl
	{
	New(@args)
		{
		super(@args)
		.noTruncateValue = args.GetDefault('noTruncateValue', false)
		}
	GETDLGCODE(wParam)
		{
		if wParam is VK.ESCAPE
			return DLGC.WANTALLKEYS
		else
			return super.GETDLGCODE(wParam)
		}
	LBUTTONUP()
		{
		.Defer(.SelectAll)
		return 'callsuper'
		}
	EN_SETFOCUS(@args)
		{
		value = .Get()
		if not .noTruncateValue
			value = value.BeforeFirst(' ')
		.Set(value)
		.Defer(.SelectAll)
		.Send('ControlSetFocus')
		return super.EN_SETFOCUS(@args)
		}
	CHAR(wParam)
		{
		if wParam is VK.ESCAPE
			.Send('FieldEscape')
		return super.CHAR(wParam)
		}
	InsertChoice(s)
		{
		if s is 'No matches'
			return 0
		return super.InsertChoice(s)
		}
	// disable function that automatically picks a value if it is only choice
	PickIfOneChoice()
		{
		}
	}