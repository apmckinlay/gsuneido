// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	New(@args)
		{
		super(@args)
		.SubClass()
		}
	GETDLGCODE(lParam)
		{
		if false isnt (m = MSG(lParam)) and m.wParam is VK.RETURN
			.Send('FieldReturn')
		return 'callsuper'
		}
	}
