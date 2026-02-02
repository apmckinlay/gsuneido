// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
_Diff2Control
	{
	SetProcs()
		{
		.Act('SetProcs', .listNew.Hwnd, .listOld.Hwnd)

		if .base isnt false
			.listNew.SetupMargin()
		}

	ClearCallBacks()
		{
		}
	}