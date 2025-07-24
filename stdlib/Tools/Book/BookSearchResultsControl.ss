// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
// TODO: make it keyboard accessible up/down/enter
Controller
	{
	New(parent, fieldHwnd, results)
		{
		super(['Border', results, border: 5])
		r = GetWindowRect(fieldHwnd)
		.Window.SetWinPos(r.right - .Xmin - 2, r.bottom)
		.parent = parent
		}
	Commands: ( ("Close", "Escape") )
	Goto(address) // from links
		{
		.parent.Send(#Goto, address)
		}
	Inactivate()
		{
		PostMessage(.Window.Hwnd, WM.CLOSE, 0, 0)
		}
	}