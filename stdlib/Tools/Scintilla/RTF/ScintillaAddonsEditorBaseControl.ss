// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonsControl
	{
	New(@args)
		{
		super(@args)
		if .readonly = args.GetDefault('readonly', false)
			.SetBackground(GetSysColor(COLOR.BTNFACE))
		}

	SetReadOnly(readOnly)
		{
		if .readonly
			return

		super.SetReadOnly(readOnly)
		.SetBackground(readOnly is true ? GetSysColor(COLOR.BTNFACE) : CLR.WHITE)
		}

	Hasfocus?: false
	HasFocus?()
		{
		return .Hasfocus? or super.HasFocus?()
		}

	SetValid(valid? = true)
		{
		if (GetFocus() is .Hwnd)
			valid? = true
		// have to check .readonly as well because if we are in a block running from the
		// ignoring_readonly method (like from Set), then the readonly flag will actually
		// be 0 when this is called resulting in the background color incorrectly
		// switching to white
		.SetBackground(.GETREADONLY() is 1 or .readonly
			? GetSysColor(COLOR.BTNFACE)
			: valid? is false ? CLR.ErrorColor : CLR.WHITE)
		}

	SetBackground(bgnd)
		{
		.StyleSetBack(0, bgnd)
		.StyleSetBack(SC.STYLE_DEFAULT, bgnd)
		}

	SCEN_KILLFOCUS()
		{
		if (.Send("Dialog?") isnt true and not .Valid?() and GetFocus() isnt .Hwnd)
			{
			.SetValid(false)
			Beep()
			}
		return super.SCEN_KILLFOCUS()
		}

	SCEN_SETFOCUS()
		{
		super.SCEN_SETFOCUS()
		.SetValid() // don't color invalid when focused
		return 0
		}
	}