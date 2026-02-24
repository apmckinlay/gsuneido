// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
_MergeDiffControl
	{
	SetProcs()
		{
		.Act('SetProcs', .listCurrent.Hwnd, .listMerge.Hwnd, .listMaster.Hwnd)

		.listMerge.SetupMargin()
		.listMerge.SetFocus()
		}

	ClearCallBacks()
		{
		}
	}