// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// Purpose: make ENTER key do validation/formatting
FieldControl
	{
	New(@args)
		{
		super(@args)
		.SubClass()
		}
	GETDLGCODE(lParam)
		{
		if (false isnt (m = MSG(lParam)) and
			m.wParam is VK.RETURN and
			not .Window.Base?(Dialog) and
			(m.message is WM.CHAR or m.message is WM.KEYDOWN))
			return DLGC.WANTALLKEYS
		return 'callsuper'
		}
	CHAR(wParam)
		{
		if wParam is VK.RETURN // Enter
			{
			dirty? = .Dirty?()
			.KillFocus()
			if dirty?
				.Send("NewValue", .Get())
			return 0
			}
		// need this for UomControl in browse in dialog (!)
		if wParam is VK.TAB
			return 0
		return 'callsuper'
		}
	}