// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: "UpDown"
	ComponentName: "UpDown"
	ComponentArgs: #()
	GetReadOnly() // read-only not applicable to updown
		{
		return true
		}

	pos: 0
	UP()
		{
		.pos = Max(0, .pos - 1)
		.Send('VSCROLL', MAKELPARAM(SB.THUMBPOSITION, .pos))
		}

	DOWN()
		{
		.pos = Min(100/*=max*/, .pos + 1)
		.Send('VSCROLL', MAKELPARAM(SB.THUMBPOSITION, .pos))
		}
	}